package utils

import (
	"testing"
)

func TestParseFilter(t *testing.T) {
	t1 := map[string]string{"limit": "10", "offset": "10"}
	l, o, err := ParseQueries(t1)
	if l != 10 {
		t.Errorf("Wrong parse limit")
	}
	if o != 10 {
		t.Errorf("Wrong parse offset")
	}
	if err != nil {
		t.Errorf("Wrong parse %s", err)
	}

	l, o, err = ParseQueries(map[string]string{})
	if l != defaultLimit {
		t.Errorf("Wrong parse limit")
	}
	if o != 0 {
		t.Errorf("Wrong parse offset")
	}
	if err != nil {
		t.Errorf("Wrong parse %s", err)
	}
}
