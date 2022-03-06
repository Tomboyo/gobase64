package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCharacters(t *testing.T) {
	input := bytes.NewReader([]byte{
		0, 16, 131, 16, 81, 135, 32, 146,
		139, 48, 211, 143, 65, 20, 147, 81,
		85, 151, 97, 150, 155, 113, 215, 159,
		130, 24, 163, 146, 89, 167, 162, 154,
		171, 178, 219, 175, 195, 28, 179, 211,
		93, 183, 227, 158, 187, 243, 223, 191})
	var output bytes.Buffer

	n, err := Encode(input, &output)

	expected := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: len: %v\n%q", output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestPaddingZero(t *testing.T) {
	input := strings.NewReader("aaa")
	var output bytes.Buffer

	n, err := Encode(input, &output)

	expected := "YWFh"
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: len: %v\n%q", output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestPaddingOne(t *testing.T) {
	input := strings.NewReader("aaaaa")
	var output bytes.Buffer

	n, err := Encode(input, &output)

	expected := "YWFhYWE="
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: len: %v\n%q", output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestPaddingTwo(t *testing.T) {
	input := strings.NewReader("aaaa")
	var output bytes.Buffer

	n, err := Encode(input, &output)

	expected := "YWFhYQ=="
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: len: %v\n%q", output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestLargeInput(t *testing.T) {
	// Input is large enough to require more than one Read from the reader.
	input := bytes.NewReader(make([]byte, 3_000))
	var output bytes.Buffer

	n, err := Encode(input, &output)

	expected := strings.Repeat("A", 4_000)
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: len: %v\n%q", output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

// func TestSmallOutputBuffer(t *testing.T) {
// 	input := []byte("aaaaaaaaaa")
// 	output := make([]byte, 4)

// 	n, err := EncodeArray(input, len(input), output)

// 	if n != 0 || err == nil {
// 		t.Logf("\nInput:\n%v\t%v", len(input), input)
// 		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", 0, []byte{}, n, output[:n])
// 		t.Logf("err: %v", err)
// 		t.Fail()
// 	}
// }

// func TestEncodeArrayPartOfBufferOnly(t *testing.T) {
// 	input := []byte("aaaaa")
// 	output := make([]byte, 256)
// 	octets := 3

// 	n, err := EncodeArray(input, octets, output)

// 	expected := []byte("YWFh")
// 	if n != len(expected) || !bytes.Equal(output[:n], expected) || err != nil {
// 		t.Logf("\nInput:\n%v\t%v", octets, input)
// 		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", len(expected), expected, n, output[:n])
// 		t.Logf("err: %v", err)
// 		t.Fail()
// 	}
// }
