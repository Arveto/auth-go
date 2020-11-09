package main

import (
	"fmt"
	"github.com/Arveto/auth-go"
	"log"
	"net/http"
)

func main() {
	app, err := auth.NewApp("app.example.com", "https://auth.dev.arveto.io/")
	if err != nil {
		fmt.Println("New app fail:", err)
		return
	}

	app.Forget = func(u *auth.User) {
		log.Println("[FORGET]", u.Pseudo)
	}

	// Standard page
	app.HandleFunc("/", auth.LevelNo, func(w http.ResponseWriter, r *auth.Request) {
		log.Println("[REQ]", r.URL)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		if r.User != nil {
			fmt.Fprintln(w, "Your are logged:<br>")
			fmt.Fprintf(w, "%#+v<br>\r\n", r.User)
			fmt.Fprintf(w, "<img src='/avatar?u=%s' alt='avatar'><br>\r\n", r.User.ID)
			fmt.Fprintf(w, `<a href="/logout">Logout</a>`)
		} else {
			fmt.Fprintf(w, "You are not logged.<br>\r\n")
			fmt.Fprintf(w, `<a href="/login">Login</a>`)
		}
	})

	// Examples of level handler filters.
	app.HandleFunc("/visitor", auth.LevelVisitor, func(w http.ResponseWriter, _ *auth.Request) {
		w.Write([]byte("You can access here because you have standard level!"))
	})
	app.HandleFunc("/admin", auth.LevelVisitor, func(w http.ResponseWriter, _ *auth.Request) {
		w.Write([]byte("You can access here because you have administrator level!"))
	})

	fmt.Println("[LISTEN]")
	http.ListenAndServe(":8000", app)
}
