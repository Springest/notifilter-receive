package main

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

var jt = types.JsonText(`{"active": true, "name": "Go"}`)

func TestRuleKeyDoesNotMatch(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key: "notactive",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestRuleKeyDoesMatch(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key: "active",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestBoolFalse(t *testing.T) {
	jt := types.JsonText(`{"active": false}`)
	s := Stat{"Mark", jt}

	r := Rule{
		Key:    "active",
		Demand: "boolean",
		Value:  "true",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestBoolTrue(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:    "active",
		Demand: "boolean",
		Value:  "true",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestStringDoesNotMatch(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:    "name",
		Demand: "string",
		Value:  "NotGo",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestStringDoesMatch(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:    "name",
		Demand: "string",
		Value:  "Go",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}
