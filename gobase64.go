package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	Encode(os.Stdin, os.Stdout)
}

func Encode(reader io.Reader, writer io.Writer) (int, error) {
	bufreader := bufio.NewReader(reader)
	bufwriter := bufio.NewWriter(writer)

	inbuf := make([]byte, 3)
	outbuf := make([]byte, 4)
	written := 0

	for {
		octetsIn, err := bufreader.Read(inbuf)
		if err == io.EOF {
			break
		} else if err != nil {
			return written, err
		}

		if octetsIn == 3 {
			encodeTriplet(inbuf, outbuf)
		} else {
			// This is the end of input. Encode the last 1 or 2 bytes with trailing padding.
			encodeTrailingOctets(inbuf[:octetsIn], outbuf)
		}

		n, err := bufwriter.Write(outbuf[:])
		written += n
		if err != nil {
			return written, err
		}
	}

	bufwriter.Flush()

	return written, nil
}

// in is a 3-byte slice
// out is a 4-byte slice
func encodeTriplet(in []byte, out []byte) {
	// Take the appropriate masks across byte boundaries and join the results with bitwise or
	out[0] = toB64((in[0] & u6m) >> 2)
	out[1] = toB64(((in[0] & l2m) << 4) | ((in[1] & u4m) >> 4))
	out[2] = toB64(((in[1] & l4m) << 2) | ((in[2] & u2m) >> 6))
	out[3] = toB64(in[2] & l6m)
}

// in is a 1- or 2-byte slice
// out is a 4-byte slice
func encodeTrailingOctets(in []byte, out []byte) {
	out[0] = toB64((in[0] & u6m) >> 2)
	if len(in) == 1 {
		out[1] = toB64((in[0] & l2m) << 4) // zero-pad lower 4 bits
		out[2] = pad
		out[3] = pad
	} else {
		out[1] = toB64(((in[0] & l2m) << 4) | ((in[1] & u4m) >> 4))
		out[2] = toB64((in[1] & l4m) << 2) // zero-pad lower 2 bits
		out[3] = pad
	}
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

var pad = byte('=')
var chars = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd',
	'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
	'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', '+', '/',
}
