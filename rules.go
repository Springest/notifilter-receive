package main

type Rule struct {
	Key    string
	Demand string
	Value  string
}

func (r *Rule) Met(s *Stat) bool {
	if r.Key != s.Key {
		return false
	}
	return true
}
