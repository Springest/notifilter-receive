package main

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

var jt = types.JsonText(`{"active": true, "name": "Go", "number": 12}`)

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
		Key:   "active",
		Type:  "boolean",
		Value: "true",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestBoolTrue(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:   "active",
		Type:  "boolean",
		Value: "true",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestStringDoesNotMatch(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:   "name",
		Type:  "string",
		Value: "NotGo",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestStringDoesMatch(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:   "name",
		Type:  "string",
		Value: "Go",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestNumberDoesNotEqual(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:   "number",
		Type:  "number",
		Value: "11",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestNumberEqual(t *testing.T) {
	s := Stat{"Mark", jt}

	r := Rule{
		Key:   "number",
		Type:  "number",
		Value: "12",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}
