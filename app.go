// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/HuguesGuilleus/go-parsersa"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// A client app, must be check in the provider. Fill all fields exept Forget.
type App struct {
	Key      *rsa.PublicKey // The public key of the auth provider.
	Audience string         // The audience claim in JWT.
	Cookie   string         // JWT cookie name. (by default it's "auth")
	Error    ErrorResponse  // Send error response.
	Mux      http.ServeMux  // Used direcly to handle no identification request
	Forget   func(u *User)  // Forget a user. Can be nil.
}

// Get the public key of the providers, and fill Mux with standars handlers.
func NewApp(id, provider string) (*App, error) {
	u, err := url.ParseRequestURI(provider)
	if err != nil {
		return nil, err
	}
	u.Path = "/"
	u.RawQuery = ""
	urlProvider := u.String()

	rep, err := http.Get(urlProvider + "publickey")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}
	rep.Body.Close()
	k, err := parsersa.Public(body)
	if err != nil {
		return nil, err
	}

	app := &App{
		Key:      k,
		Audience: id,
		Cookie:   "auth",
		Error:    ErrorResponseDefault,
	}
	app.addHandlers(urlProvider)
	return app, nil
}

// Add standards handlers for a standard app. No used Avatar.
func (app *App) addHandlers(urlProvider string) {
	// Avatar URL
	app.Mux.HandleFunc("/avatar", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, urlProvider+"avatar?"+r.URL.RawQuery, http.StatusTemporaryRedirect)
	})
	// Logout
	app.Mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", `auth=none; Max-Age=0`)
		to := "/"
		if r := r.URL.Query().Get("r"); r != "" {
			to = r
		}
		http.Redirect(w, r, to, http.StatusTemporaryRedirect)
	})
	// Login
	app.Mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		to := query.Get("r")
		if to == "" {
			to = "/"
		}
		if jwt := query.Get("jwt"); jwt != "" {
			_, err := UnmarshalJWT(app.Key, app.Audience, jwt)
			if err != nil {
				app.Error(w, r, err, http.StatusBadRequest)
				return
			}
			w.Header().Add("Set-Cookie", (&http.Cookie{
				Name:     app.Cookie,
				Value:    jwt,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   24 * 60 * 60,
				SameSite: http.SameSiteStrictMode,
				Secure:   true,
			}).String())
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			to = strconv.QuoteToASCII(to)
			w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="utf-8"><script>document.location.replace(` + to + `);</script></head><body><a href="` + to + `">Redirect: ` + to + `</a></body></html>`))
		} else {
			u := urlProvider + "auth?" + (&url.Values{
				"app": []string{app.Audience},
				"r":   []string{to},
			}).Encode()
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		}
	})
	// Inforlmation about the user.
	app.Mux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		// Get JWT
		c, _ := r.Cookie(app.Cookie)
		if c == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("null"))
			return
		}
		jwt := c.Value

		// Check JWT
		u, err := UnmarshalJWT(app.Key, app.Audience, jwt)
		if err != nil {
			app.Error(w, r, err, http.StatusBadRequest)
			return
		}

		// Get expiration date
		var date struct {
			Exp int64 `json:"exp"`
		}
		json.Unmarshal([]byte(strings.Split(jwt, ".")[1]), &date)

		// Response
		data, _ := json.Marshal(struct {
			User
			Expiration time.Time
		}{
			User:       *u,
			Expiration: time.Unix(date.Exp, 0),
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	// Forget a user.
	app.HandleFunc("/forget", LevelAdministrator, func(w http.ResponseWriter, r *Request) {
		if r.Method != "DELETE" {
			w.Header().Set("Allow", "DELETE")
			app.Error(w, r.R(), errors.New("Need DELETE method"), http.StatusMethodNotAllowed)
			return
		}

		u, err := UnmarshalJWT(app.Key, app.Audience, r.URL.Query().Get("jwt"))
		if err != nil {
			app.Error(w, r.R(), err, http.StatusBadRequest)
			return
		}

		if f := app.Forget; f != nil {
			f(u)
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

var (
	ErrNotLogged = errors.New("You are not logged")
	ErrLowLevel  = errors.New("Your level is too low")
)

// Wrap the handler with a user level checker: if the level is lower than level,
// return error and the handler is not call.
func (a *App) Handle(pattern string, level UserLevel, handler Handler) {
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		jwt := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
		if jwt == "" {
			cookie, _ := r.Cookie(a.Cookie)
			if cookie == nil {
				if level != LevelNo {
					a.Error(w, r, ErrNotLogged, http.StatusUnauthorized)
					return
				}
				handler.ServeHTTP(w, &Request{
					Request: *r,
					User:    nil,
				})
			}
			jwt = cookie.Value
		}

		u, err := UnmarshalJWT(a.Key, a.Audience, jwt)
		if err != nil {
			a.Error(w, r, err, http.StatusBadRequest)
			return
		} else if u.Level < level {
			a.Error(w, r, ErrLowLevel, http.StatusForbidden)
			return
		}

		handler.ServeHTTP(w, &Request{
			Request: *r,
			User:    u,
		})
	})
}

func (a *App) HandleFunc(pattern string, level UserLevel, handler func(http.ResponseWriter, *Request)) {
	a.Handle(pattern, level, HandlerFunc(handler))
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Mux.ServeHTTP(w, r)
}
