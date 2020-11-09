// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/HuguesGuilleus/go-parsersa"
	"net/http"
)

type Provider struct {
	K   *rsa.PrivateKey
	PEM []byte
}

func NewProvier(keyPath string) (*Provider, error) {
	k, err := parsersa.GenPrivKey(keyPath, 4096)
	if err != nil {
		return nil, err
	}

	buff := bytes.Buffer{}
	pem.Encode(&buff, &pem.Block{
		Type:  "BEGIN PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&k.PublicKey),
	})

	return &Provider{
		K:   k,
		PEM: buff.Bytes(),
	}, nil
}

// Return the JWT signed with the provider key.
func (p *Provider) JWT(aud string, u *User) (string, error) {
	return MarchalJWT(p.K, aud, u)
}

func (p *Provider) ServerKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Write(p.PEM)
}
