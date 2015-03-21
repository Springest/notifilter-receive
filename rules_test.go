package main

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleKeyDoesNotMatch(t *testing.T) {
	jt := types.JsonText(`{"foo": 1, "bar": 2}`)
	s := Stat{"Mark", jt}

	r := Rule{
		Key: "notmark",
	}

	result := r.Met(&s)
	fmt.Println(result)
	assert.Equal(t, false, result)
}

func TestRuleKeyDoesMatch(t *testing.T) {
	jt := types.JsonText(`{"foo": 1, "bar": 2}`)
	s := Stat{"Mark", jt}

	r := Rule{
		Key: "Mark",
	}

	result := r.Met(&s)
	fmt.Println(result)
	assert.Equal(t, true, result)
}

func TestBool(t *testing.T) {
}
