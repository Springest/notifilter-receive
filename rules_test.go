package main

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleKeyDoesNotMatch(t *testing.T) {
	jt := types.JsonText(`{"active": false}`)
	s := Stat{"Mark", jt}

	r := Rule{
		Key: "notactive",
	}

	result := r.Met(&s)
	assert.Equal(t, false, result)
}

func TestRuleKeyDoesMatch(t *testing.T) {
	jt := types.JsonText(`{"active": false}`)
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
	jt := types.JsonText(`{"active": true}`)
	s := Stat{"Mark", jt}

	r := Rule{
		Key:    "active",
		Demand: "boolean",
		Value:  "true",
	}

	result := r.Met(&s)
	assert.Equal(t, true, result)
}
