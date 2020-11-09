// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type jwtBody struct {
	User
	Audience string `json:"aud"`
	Exp      int64  `json:"exp"`
}

const (
	// The header `{"alg": "RS256","typ": "JWT"}` already encoded + dot
	jwtHead = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9."
	// The duration of a day
	day = 24 * time.Hour
)

var (
	JWTNeedUserField   = errors.New("User need an ID, a Pseudo and an email")
	JWTEmpty           = errors.New("This JWT is empty")
	JWTOutDate         = errors.New("This JWT is out date")
	JWTWrongAudience   = errors.New("This JWT is made for an other audience")
	JWTWrongHead       = errors.New("JWT wrong head")
	JWTWrongSyntax     = errors.New("JWT wrong syntax")
	JWTWrongSyntaxHead = errors.New("JWT wrong syntax in head")
)

// Create a new JWT to a specific audience.
func MarchalJWT(key *rsa.PrivateKey, aud string, u *User) (string, error) {
	if u.ID == "" || u.Pseudo == "" || u.Email == "" {
		return "", JWTNeedUserField
	}
	buff := bytes.NewBufferString(jwtHead)
	body, err := json.Marshal(jwtBody{
		User:     *u,
		Audience: aud,
		Exp:      time.Now().Add(day).Unix(),
	})
	if err != nil {
		return "", err
	}
	buff.Write(tob64(body))

	h := sha256.Sum256(buff.Bytes())
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}
	buff.WriteRune('.')
	buff.Write(tob64(sig))

	return buff.String(), nil
}

// Parse the JWT from provider and return the user.
func UnmarshalJWT(key *rsa.PublicKey, aud string, j string) (*User, error) {
	if j == "" {
		return nil, JWTEmpty
	}
	parts := strings.SplitN(j, ".", -1)
	if len(parts) != 3 {
		return nil, JWTWrongSyntax
	}

	// Check the head
	h := struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}{}
	if err := json.Unmarshal([]byte(fromb64(parts[0])), &h); err != nil {
		return nil, err
	} else if h.Alg != "RS256" || h.Typ != "JWT" {
		return nil, JWTWrongHead
	}

	// Get Body
	var body jwtBody
	if err := json.Unmarshal(fromb64(parts[1]), &body); err != nil {
		return nil, err
	} else if body.Audience != aud {
		return nil, JWTWrongAudience
	} else if body.Exp < time.Now().Unix() {
		return nil, JWTOutDate
	}

	// Check the signature
	sig := fromb64(parts[2])
	hash := sha256.Sum256([]byte(j[:len(parts[0])+len(parts[1])+1]))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], sig); err != nil {
		return nil, err
	}

	return &body.User, nil
}

// Return decoding base 64 Raw URL. If error, then output is empty.
func fromb64(in string) []byte {
	out, err := base64.RawURLEncoding.DecodeString(in)
	if err != nil {
		return nil
	}
	return out
}

// Return the string encoding in base64
func tob64(in []byte) []byte {
	return []byte(base64.RawURLEncoding.EncodeToString(in))
}
