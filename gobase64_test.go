package main

import (
	"bytes"
	"testing"
)

func TestCharacters(t *testing.T) {
	input := []byte{
		0, 16, 131, 16, 81, 135, 32, 146,
		139, 48, 211, 143, 65, 20, 147, 81,
		85, 151, 97, 150, 155, 113, 215, 159,
		130, 24, 163, 146, 89, 167, 162, 154,
		171, 178, 219, 175, 195, 28, 179, 211,
		93, 183, 227, 158, 187, 243, 223, 191}
	output := make([]byte, 256)

	n, err := Encode(&input, len(input), &output)

	expected := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	if n != len(expected) || !bytes.Equal(output[:n], expected) || err != nil {
		t.Logf("\nInput:\n%v\t%v", len(input), input)
		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", len(expected), expected, n, output[:n])
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestPaddingZero(t *testing.T) {
	input := []byte("aaa")
	output := make([]byte, 256)

	n, err := Encode(&input, len(input), &output)

	expected := []byte("YWFh")
	if n != len(expected) || !bytes.Equal(output[:n], expected) || err != nil {
		t.Logf("\nInput:\n%v\t%v", len(input), input)
		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", len(expected), expected, n, output[:n])
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestPaddingOne(t *testing.T) {
	input := []byte("aaaa")
	output := make([]byte, 256)

	n, err := Encode(&input, len(input), &output)

	expected := []byte("YWFhYQ==")
	if n != len(expected) || !bytes.Equal(output[:n], expected) || err != nil {
		t.Logf("\nInput:\n%v\t%v", len(input), input)
		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", len(expected), expected, n, output[:n])
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestPaddingTwo(t *testing.T) {
	input := []byte("aaaaa")
	output := make([]byte, 255)

	n, err := Encode(&input, len(input), &output)

	expected := []byte("YWFhYWE=")
	if n != len(expected) || !bytes.Equal(output[:n], expected) || err != nil {
		t.Logf("\nInput:\n%v\t%v", len(input), input)
		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", len(expected), expected, n, output[:n])
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestEncodePartOfBufferOnly(t *testing.T) {
	input := []byte("aaaaa")
	output := make([]byte, 256)
	octets := 3

	n, err := Encode(&input, octets, &output)

	expected := []byte("YWFh")
	if n != len(expected) || !bytes.Equal(output[:n], expected) || err != nil {
		t.Logf("\nInput:\n%v\t%v", octets, input)
		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", len(expected), expected, n, output[:n])
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func TestSmallOutputBuffer(t *testing.T) {
	input := []byte("aaaaaaaaaa")
	output := make([]byte, 4)

	n, err := Encode(&input, len(input), &output)

	if n != 0 || err == nil {
		t.Logf("\nInput:\n%v\t%v", len(input), input)
		t.Logf("\nExpected:\n%v\t%v\nActual:\n%v\t%v", 0, []byte{}, n, output[:n])
		t.Logf("err: %v", err)
		t.Fail()
	}
}
