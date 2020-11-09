// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJWT(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1048)
	assert.NoError(t, err)
	bob := User{
		ID:     "6751fcc68f",
		Pseudo: "Bob",
		Email:  "bob@arveto.io",
		Bot:    false,
		Level:  LevelStandard,
		Teams: map[string]bool{
			"dev": true,
		},
	}

	jwt, err := MarchalJWT(key, "yolo", &bob)
	assert.NoError(t, err)

	u, err := UnmarshalJWT(&key.PublicKey, "yolo", jwt)
	assert.NoError(t, err)
	assert.Equal(t, &bob, u)
}
