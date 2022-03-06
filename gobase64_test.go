package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSampleInput(t *testing.T) {
	input := "Many hands make light work."
	expected := "TWFueSBoYW5kcyBtYWtlIGxpZ2h0IHdvcmsu"
	var buffer bytes.Buffer
	btoa(strings.NewReader(input), &buffer)
	actual := buffer.String()
	if actual != expected {
		t.Logf("\nExpected:\n%v\nActual:\n%v", expected, actual)
		t.Fail()
	}
}
