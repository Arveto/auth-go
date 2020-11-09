// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/json"
	"errors"
	"sort"
)

// One user
type User struct {
	ID     string    `json:"id"`
	Pseudo string    `json:"pseudo"`
	Email  string    `json:"email"`
	Level  UserLevel `json:"level"`
	Bot    bool      `json:"bot"`
	Teams  Teams     `json:"teams"`
}

// The user's teams. To save in JSON, it convert to an array.
type Teams map[string]bool

func (t Teams) MarshalJSON() ([]byte, error) {
	tab := make([]string, 0, len(t))
	for k, ok := range t {
		if ok {
			tab = append(tab, k)
		}
	}
	sort.Slice(tab, func(i int, j int) bool { return tab[i] < tab[j] })
	return json.Marshal(tab)
}
func (t *Teams) UnmarshalJSON(data []byte) error {
	tab := make([]string, 0)
	if err := json.Unmarshal(data, &tab); err != nil {
		return err
	}
	*t = make(Teams, len(tab))
	for _, k := range tab {
		(*t)[k] = true
	}
	return nil
}

// The user's or bot's accreditation level.
type UserLevel uint

const (
	LevelNo            = iota
	LevelCandidate     = iota
	LevelVisitor       = iota
	LevelStandard      = iota
	LevelAdministrator = iota
)

var UserLevelUnknown = errors.New("Unknown UserLevel")

func (l UserLevel) String() string {
	switch l {
	case LevelNo:
		return "no"
	case LevelCandidate:
		return "candidate"
	case LevelVisitor:
		return "visitor"
	case LevelStandard:
		return "standard"
	case LevelAdministrator:
		return "administrator"
	}
	return "?"
}
func (l UserLevel) MarshalText() ([]byte, error) {
	s := l.String()
	if s == "?" {
		return nil, errors.New("Unknown UserLevel")
	}
	return []byte(s), nil
}
func (l *UserLevel) UnmarshalText(text []byte) error {
	switch string(text) {
	case "no":
		*l = LevelNo
	case "candidate", "Candidate":
		*l = LevelCandidate
	case "visitor", "Visitor":
		*l = LevelVisitor
	case "standard", "Std":
		*l = LevelStandard
	case "administrator", "Admin":
		*l = LevelAdministrator
	default:
		*l = LevelNo
		return UserLevelUnknown
	}
	return nil
}
