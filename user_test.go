// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTeamsMarshalJSON(t *testing.T) {
	team := Teams{
		"dev":     true,
		"kitchen": true,
	}
	data, err := team.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `["dev","kitchen"]`, string(data))
}

func TestTeamsUnmarshalJSON(t *testing.T) {
	team := Teams{}
	assert.NoError(t, team.UnmarshalJSON([]byte(`["dev","kitchen"]`)))
	assert.Equal(t, Teams{
		"dev":     true,
		"kitchen": true,
	}, team)
}
