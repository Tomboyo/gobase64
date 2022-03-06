package main

import (
	"container/list"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math/bits"
	"os"
)

// our worker count
const WORKERS = 4

type B64Chars struct {
	value []byte
}

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

func btoa_parallel(reader io.Reader, writer io.Writer) {
	c := make(chan *B64Chars, 16)

	// load 3^10 bytes
	var inbuf = make([]byte, 3)
	var bl = 0
	outbuffer := list.New()
	outbuffer.Init()
	for {
		r, rerr := reader.Read(inbuf)
		bl += r
		if rerr != nil && rerr != io.EOF {
			log.Fatal("Unexpected IO error", rerr)
		}
		// offload the getChars to a byte array
		go func() {
			s, err := getChars(&inbuf, 3-r)
			if err != nil {
				log.Fatal("Unexpected Base64 encoding error", err)
			}
			chars := new(B64Chars)
			chars.value = s
			c <- chars
		}()
		outbuffer.PushBack(<-c)
		if rerr == io.EOF {
			break
		}
	}
	for e := outbuffer.Front(); e != nil; e = e.Next() {
		writer.Write(e.Value.(*B64Chars).value)
	}
	var p = (bl % 3)
	if p == 1 {
		writer.Write([]byte{'=', '='})
	} else if p == 2 {
		writer.Write([]byte{'='})
	}
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
	rawArgs := os.Args
	if len(rawArgs) > 1 {
		args := rawArgs[1:]
		if args[0] == "p" {
			btoa_parallel(reader, writer)
		} else {
			log.Fatal("Unknown flag", args[0])
		}
	} else {
		btoa(reader, writer)
	}
}
