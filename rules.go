package main

import (
	"encoding/json"
	"log"
	"strconv"
)

type Rule struct {
	Key    string
	Demand string
	Value  string
}

func (r *Rule) Met(s *Stat) bool {
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(s.Value), &parsed)
	if err != nil {
		log.Fatal("json.Unmarshal", err)
	}

	// check if key is in the map
	// first value is actual value of key in the map
	if _, ok := parsed[r.Key]; !ok {
		return false
	}

	if r.Demand == "boolean" {
		val := parsed[r.Key]
		needed_val, _ := strconv.ParseBool(r.Value)
		if val.(bool) != needed_val {
			return false
		}
	}

	return true
}
