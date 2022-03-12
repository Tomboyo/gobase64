package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

type impl func(io.Reader, io.Writer) (int, error)

func Test_Characters_Serial(t *testing.T) {
	doTestCharacters(t, EncodeSerial)
}

func Test_Characters_Parallel(t *testing.T) {
	doTestCharacters(t, EncodeParallel)
}

func Test_Characters_Stdlib(t *testing.T) {
	doTestCharacters(t, EncodeStdlib)
}

func doTestCharacters(t *testing.T, f impl) {
	input := bytes.NewReader([]byte{
		0, 16, 131, 16, 81, 135, 32, 146,
		139, 48, 211, 143, 65, 20, 147, 81,
		85, 151, 97, 150, 155, 113, 215, 159,
		130, 24, 163, 146, 89, 167, 162, 154,
		171, 178, 219, 175, 195, 28, 179, 211,
		93, 183, 227, 158, 187, 243, 223, 191})
	var output bytes.Buffer

	n, err := f(input, &output)

	expected := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	o := output.String()
	if n != len(expected) || o != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: n: %v len: %v\n%q", n, output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func Test_PaddingZero_Serial(t *testing.T) {
	doTestPaddingZero(t, EncodeSerial)
}

func Test_PaddingZero_Parallel(t *testing.T) {
	doTestPaddingZero(t, EncodeParallel)
}

func Test_PaddingZero_Stdlib(t *testing.T) {
	doTestPaddingZero(t, EncodeStdlib)
}

func doTestPaddingZero(t *testing.T, f impl) {
	input := strings.NewReader("aaa")
	var output bytes.Buffer

	n, err := f(input, &output)

	expected := "YWFh"
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: len: %v\n%q", output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func Test_PaddingOne_Serial(t *testing.T) {
	doTestPaddingOne(t, EncodeSerial)
}
func Test_PaddingOne_Parallel(t *testing.T) {
	doTestPaddingOne(t, EncodeParallel)
}

func Test_PaddingOne_Stdlib(t *testing.T) {
	doTestPaddingOne(t, EncodeStdlib)
}

func doTestPaddingOne(t *testing.T, f impl) {
	input := strings.NewReader("aaaaa")
	var output bytes.Buffer

	n, err := f(input, &output)

	expected := "YWFhYWE="
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: n: %v len: %v\n%q", n, output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func Test_PaddingTwo_Serial(t *testing.T) {
	doTestPaddingTwo(t, EncodeSerial)
}

func Test_PaddingTwo_Parallel(t *testing.T) {
	doTestPaddingTwo(t, EncodeParallel)
}

func Test_PaddingTwo_Stdlib(t *testing.T) {
	doTestPaddingTwo(t, EncodeStdlib)
}

func doTestPaddingTwo(t *testing.T, f impl) {
	input := strings.NewReader("aaaa")
	var output bytes.Buffer

	n, err := f(input, &output)

	expected := "YWFhYQ=="
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n%q", len(expected), expected)
		t.Logf("\nActual: n: %v len: %v\n%q", n, output.Len(), output.String())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

func Test_LargeInput_Serial(t *testing.T) {
	doTestLargeInput(t, EncodeSerial)
}

func Test_LargeInput_Parallel(t *testing.T) {
	doTestLargeInput(t, EncodeParallel)
}

func Test_LargeInput_Stdlib(t *testing.T) {
	doTestLargeInput(t, EncodeStdlib)
}

func doTestLargeInput(t *testing.T, f impl) {
	// Input is large enough to require more than one Read from the reader.
	input := bytes.NewReader(make([]byte, 6_000))
	var output bytes.Buffer

	n, err := f(input, &output)

	expected := strings.Repeat("A", 8_000)
	if n != len(expected) || output.String() != expected || err != nil {
		t.Logf("\nExpected: len: %v\n", len(expected))
		t.Logf("\nActual: len: %v\n", output.Len())
		t.Logf("err: %v", err)
		t.Fail()
	}
}

// cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
// SERIAL
// Baseline:            24,736,590,286 ns/op
// Add buffered reader: 13,099,883,339 ns/op
// Add buffered writer:          4,994 ns/op
// ReadFull fix                  5,137 ns/op
// If-ladder order               2,106 ns/op
func Benchmark_EncodeSerial_LargeInputFileIO(b *testing.B) {
	doBenchmarkLargeInputFileIO(b, EncodeSerial)
}

// cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
// PARALLEL
// Baseline:             2,683,206,126 ns/op
// Smart chunk lengths:          6,044 ns/op
// ReadFull fix                  5,175 ns/op
// If-ladder                     4,939 ns/op
func Benchmark_EncodeParallel_LargeInputFileIO(b *testing.B) {
	doBenchmarkLargeInputFileIO(b, EncodeParallel)
}

// cpu: Intel(R) Core(TM) i7-7600U CPU @ 2.80GHz
// Buffered, baeline:    7,522,238,459 ns/op
func Benchmark_Stdlib_LargeInputFileIO(b *testing.B) {
	doBenchmarkLargeInputFileIO(b, EncodeStdlib)
}

// Benchmark runs Encode on two files. File IO is slow, so this will emphasize our buffering strategy.
func doBenchmarkLargeInputFileIO(b *testing.B, impl func(io.Reader, io.Writer) (int, error)) {
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
		// _, err = Encode(input, output)
		_, err = impl(input, output)
		if err != nil {
			b.Error("Unexpected B64 error", err)
		}
	}
}
