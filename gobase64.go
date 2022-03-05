package main

import (
	"errors"
	"io"
	"log"
	"os"
)

func main() {
	reader := os.Stdin
	writer := os.Stdout

	inbuf := make([]byte, 1024)
	outbuf := make([]byte, 1368)

	for {
		octetsIn, err := reader.Read(inbuf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Unexpected IO error", err)
		}

		n, err := Encode(&inbuf, octetsIn, &outbuf)

		if err != nil {
			log.Fatal("Unexpected Base64 encoding error", err)
		}

		writer.Write(outbuf[:n])
	}
}

func Encode(input *[]byte, octets int, output *[]byte) (int, error) {
	triplets := octets / 3
	remainder := octets % 3
	var outlen int
	if remainder == 0 {
		outlen = triplets * 4
	} else {
		outlen = triplets*4 + 4
	}

	if outlen > len(*output) {
		return 0, errors.New("output buffer is too small")
	}

	in := *input
	out := *output

	// 1. read all triplets
	for triplet := 0; triplet < triplets; triplet++ {
		i := triplet * 3
		j := triplet * 4
		out[j] = toB64((in[i] & u6m) >> 2)
		out[j+1] = toB64(((in[i] & l2m) << 4) | ((in[i+1] & u4m) >> 4))
		out[j+2] = toB64(((in[i+1] & l4m) << 2) | ((in[i+2] & u2m) >> 6))
		out[j+3] = toB64(in[i+2] & l6m)
	}

	// 2. read remaining octets and add padding
	if remainder != 0 {
		i := triplets
		j := i * 4

		out[j] = toB64((in[i] & u6m) >> 2)

		if remainder == 1 {
			out[j+1] = toB64((in[i] & l2m) << 4) // zero-pad lower 4 bits
			out[j+2] = pad
			out[j+3] = pad
		} else {
			out[j+1] = toB64(((in[i] & l2m) << 4) | ((in[i+1] & u4m) >> 4))
			out[j+2] = toB64((in[i+1] & l4m) << 2) // zero-pad lower 2 bits
			out[j+3] = pad
		}
	}

	return outlen, nil
}

var u6m = byte(253) // 1111 1100
var l2m = byte(3)   // 0000 0011
var u4m = byte(241) // 1111 0000
var l4m = byte(15)  // 0000 1111
var u2m = byte(192) // 1100 0000
var l6m = byte(63)  // 0011 1111

func toB64(in byte) byte {
	return chars[int(in)]
}

var pad = []byte("=")[0]

var chars = []byte{
	[]byte("A")[0],
	[]byte("B")[0],
	[]byte("C")[0],
	[]byte("D")[0],
	[]byte("E")[0],
	[]byte("F")[0],
	[]byte("G")[0],
	[]byte("H")[0],
	[]byte("I")[0],
	[]byte("J")[0],
	[]byte("K")[0],
	[]byte("L")[0],
	[]byte("M")[0],
	[]byte("N")[0],
	[]byte("O")[0],
	[]byte("P")[0],
	[]byte("Q")[0],
	[]byte("R")[0],
	[]byte("S")[0],
	[]byte("T")[0],
	[]byte("U")[0],
	[]byte("V")[0],
	[]byte("W")[0],
	[]byte("X")[0],
	[]byte("Y")[0],
	[]byte("Z")[0],
	[]byte("a")[0],
	[]byte("b")[0],
	[]byte("c")[0],
	[]byte("d")[0],
	[]byte("e")[0],
	[]byte("f")[0],
	[]byte("g")[0],
	[]byte("h")[0],
	[]byte("i")[0],
	[]byte("j")[0],
	[]byte("k")[0],
	[]byte("l")[0],
	[]byte("m")[0],
	[]byte("n")[0],
	[]byte("o")[0],
	[]byte("p")[0],
	[]byte("q")[0],
	[]byte("r")[0],
	[]byte("s")[0],
	[]byte("t")[0],
	[]byte("u")[0],
	[]byte("v")[0],
	[]byte("w")[0],
	[]byte("x")[0],
	[]byte("y")[0],
	[]byte("z")[0],
	[]byte("0")[0],
	[]byte("1")[0],
	[]byte("2")[0],
	[]byte("3")[0],
	[]byte("4")[0],
	[]byte("5")[0],
	[]byte("6")[0],
	[]byte("7")[0],
	[]byte("8")[0],
	[]byte("9")[0],
	[]byte("+")[0],
	[]byte("/")[0],
}
