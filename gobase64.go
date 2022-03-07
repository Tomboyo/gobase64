package main

import (
	"bytes"
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
func getChars(buffer *[]byte, padding int) ([]byte, error) {
	if padding > 2 {
		return []byte{}, nil
	}
	var out []byte
	var in = *buffer
	var l = len(in)
	if l != 3 {
		return nil, errors.New("expected 3 bytes")
	}
	if padding > 0 {
		in[2] = 0x00
	}
	if padding > 1 {
		in[1] = 0x00
	}
	if l > 0 {
		var unsigned = binary.BigEndian.Uint32(append([]byte{0x00}, in...))
		out = append(out, charset[bits.RotateLeft32((unsigned&0xFC0000), -18)])
		out = append(out, charset[bits.RotateLeft32((unsigned&0x3F000), -12)])
		if padding < 2 {
			out = append(out, charset[bits.RotateLeft32((unsigned&0x0FC0), -6)])
		}
		if padding < 1 {
			out = append(out, charset[(unsigned&0x003F)])
		}
	}
	return out, nil
}

func btoa(reader io.Reader, writer io.Writer) {
	var inbuf = make([]byte, 3)
	var bl = 0
	for {
		r, rerr := reader.Read(inbuf)
		bl += r
		if rerr != nil && rerr != io.EOF {
			log.Fatal("Unexpected IO error", rerr)
		}
		out, err := getChars(&inbuf, 3-r)
		if err != nil {
			log.Fatal("Unexpected Base64 encoding error", err)
		}
		writer.Write(out)
		if rerr == io.EOF {
			break
		}
	}
	var p = (bl % 3)
	if p == 1 {
		writer.Write([]byte{'=', '='})
	} else if p == 2 {
		writer.Write([]byte{'='})
	}
}

func main() {
	var reader = os.Stdin
	var writer = os.Stdout

	// greedy read from stdin
	inSize := 30 * 1024
	pbuffer := make([]byte, inSize)
	for {
		inpRead, err := reader.Read(pbuffer)
		if err != nil && err != io.EOF {
			log.Fatal("Unexpected IO error", err)
		} else if err == io.EOF {
			break
		}
		workBuffer := pbuffer[0:inpRead]
		btoa(bytes.NewBuffer(workBuffer), writer)
	}
}
