package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
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
	input := bytes.NewReader(make([]byte, 6_000))
	var output bytes.Buffer

	n, err := Encode(input, &output)

	expected := strings.Repeat("A", 8_000)
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n", len(expected))
		t.Logf("\nActual: len: %v\n", output.Len())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

// Benchmark runs Encode on two files. File IO is slow, so this will emphasize our buffering strategy.
// cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
// SERIAL
// Baseline:            24,736,590,286 ns/op
// Add buffered reader: 13,099,883,339 ns/op
// Add buffered writer:          4,994 ns/op
// ReadFull fix                  5,137 ns/op
// Fix again                     2,106 ns/op
// PARALLEL
// Baseline:             2,683,206,126 ns/op
// Smart chunk lengths:          6,044 ns/op
func BenchmarkLargeInputFileIO(b *testing.B) {
	input, err := ioutil.TempFile("", "benchmark-input")
	if err != nil {
		b.Error("Failed to create input tmpfile", err)
	}
	defer os.Remove(input.Name())

	output, err := ioutil.TempFile("", "benchmark-output")
	if err != nil {
		b.Error("Failed to create output tmpfile", err)
	}
	defer os.Remove(output.Name())

	// Fill the input with random bytes
	data := make([]byte, 32_048_576)
	rand.Seed(time.Now().UnixNano())
	rand.Read(data)
	_, err = input.WriteAt(data, 0)
	if err != nil {
		b.Error("Failed to write data to tmp file", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = Encode(input, output)
		if err != nil {
			b.Error("Unexpected B64 error", err)
		}
	}
}
