package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math/bits"
	"os"
)

// Not immutable but runtime array slicing to create a pseudo-immutable array is slower
var charset = [64]byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
	'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1',
	'2', '3', '4', '5', '6', '7', '8', '9', '+', '/'}

// sextuple to chars
// assumes output is initially empty
func getChars(buffer *[]byte, padding int) (string, error) {
	var out string
	var in = *buffer
	var l = len(in)
	if l != 3 {
		return "", errors.New("expected 3 bytes")
	}
	if l > 0 {
		var unsigned = binary.BigEndian.Uint32(append([]byte{0x00}, in...))
		var part1, part2, part3, part4 string
		if padding > 0 {
			part1 = "="
		} else {
			part1 = string(charset[(unsigned & 0x003F)])
		}
		if padding > 1 {
			part2 = "="
		} else {
			part2 = string(charset[bits.RotateLeft32((unsigned&0x0FC0), -6)])
		}
		part3 = string(charset[bits.RotateLeft32((unsigned&0x3F000), -12)])
		part4 = string(charset[bits.RotateLeft32((unsigned&0xFC0000), -18)])
		out = part4 + part3 + part2 + part1
	}
	return out, nil
}
func btoa(reader io.Reader, writer io.Writer) {
	var inbuf = make([]byte, 3)

	for {
		r, err := reader.Read(inbuf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Unexpected IO error", err)
		}
		out, err := getChars(&inbuf, 3-r)
		if err != nil {
			log.Fatal("Unexpected Base64 encoding error", err)
		}
		writer.Write([]byte(out))
	}
}

func main() {
	var reader = os.Stdin
	var writer = os.Stdout

	btoa(reader, writer)
}
