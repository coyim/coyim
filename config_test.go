package main

import "testing"

func TestParseYes(t *testing.T) {
	if ok := parseYes("Y"); !ok {
		t.Errorf("parsed Y as %v", ok)
	}

	if ok := parseYes("y"); !ok {
		t.Errorf("parsed y as %v", ok)
	}

	if ok := parseYes("YES"); !ok {
		t.Errorf("parsed YES as %v", ok)
	}

	if ok := parseYes("yes"); !ok {
		t.Errorf("parsed yes as %v", ok)
	}

	if ok := parseYes("Yes"); !ok {
		t.Errorf("parsed yes as %v", ok)
	}

	if ok := parseYes("anything"); ok {
		t.Errorf("parsed something else as %v", ok)
	}
}
