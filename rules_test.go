package main

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

var jt = types.JsonText(`{"active": true, "name": "Go", "number": 12}`)

func TestRuleKeyDoesNotMatch(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key: "notactive",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestRuleKeyDoesMatch(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key: "active",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestBoolFalse(t *testing.T) {
	jt := types.JsonText(`{"active": false}`)
	s := Event{"Mark", jt}

	r := rule{
		Key:   "active",
		Type:  "boolean",
		Value: "true",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestBoolTrue(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:   "active",
		Type:  "boolean",
		Value: "true",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestStringDoesNotMatch(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:   "name",
		Type:  "string",
		Value: "NotGo",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestStringDoesMatch(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:   "name",
		Type:  "string",
		Value: "Go",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestNumberDoesNotEqual(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "eq",
		Value:   "11",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestNumberEqual(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "eq",
		Value:   "12",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestNumberNotGt(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "gt",
		Value:   "13",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestNumberNotGtEqual(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "gt",
		Value:   "12",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestNumberGt(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "gt",
		Value:   "11",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}

func TestNumberNotLt(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "lt",
		Value:   "11",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestNumberNotLtEqual(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "lt",
		Value:   "12",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestNumberLt(t *testing.T) {
	s := Event{"Mark", jt}

	r := rule{
		Key:     "number",
		Type:    "number",
		Setting: "lt",
		Value:   "13",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}
