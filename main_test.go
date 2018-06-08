package main

import (
	"testing"
)

func TestToUpper(t *testing.T) {
	s := toUpper("foo-bar")
	if s != "FOO_BAR" {
		t.Error("hoge: ", s)
	}

	s = toUpper("!@#$%^&*()_+|~")
	if s != "______________" {
		t.Error("hoge: ", s)
	}
}
