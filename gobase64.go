package main

import (
	"bufio"
	"encoding/base64"
	"io"
	"os"
)

func main() {
	Encode(os.Stdin, os.Stdout)
}

const chanBuffers = 5
const parInBufLen = 3_000

func Encode(reader io.Reader, writer io.Writer) (int, error) {
	return EncodeSerial(reader, writer)
}

func EncodeStdlib(reader io.Reader, writer io.Writer) (int, error) {
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	input := bufio.NewReader(reader)
	inbuf := make([]byte, 3)
	written := 0
	var oerr error
	for {
		octetsIn, err := io.ReadFull(input, inbuf)
		if err == nil {
			n, werr := encoder.Write(inbuf)
			written += n
			if werr != nil {
				oerr = werr
				break
			}
		} else if err == io.ErrUnexpectedEOF {
			n, werr := encoder.Write(inbuf[:octetsIn])
			written += n
			if werr != nil {
				oerr = werr
				break
			}
		} else if err == io.EOF {
			werr := encoder.Close()
			if werr != nil {
				oerr = err
				break
			}
			oerr = nil
			break
		} else {
			oerr = err
			break
		}
	}

	n := (4 * (written / 3))
	if (written % 3) > 0 {
		n += 4
	}
	return n, oerr
}

func EncodeSerial(reader io.Reader, writer io.Writer) (int, error) {
	bufreader := bufio.NewReader(reader)
	bufwriter := bufio.NewWriter(writer)

	inbuf := make([]byte, 3)
	outbuf := make([]byte, 4)
	written := 0

	for {
		octetsIn, err := io.ReadFull(bufreader, inbuf)
		if err == nil {
			encodeTriplet(inbuf, outbuf)
		} else if err == io.ErrUnexpectedEOF {
			// This is the end of input. Encode the last 1 or 2 bytes with trailing padding.
			encodeTrailingOctets(inbuf[:octetsIn], outbuf)
		} else if err == io.EOF {
			bufwriter.Flush()
			return written, nil
		} else if err != nil {
			return written, err
		}

		n, err := bufwriter.Write(outbuf)
		written += n
		if err != nil {
			return written, err
		}
	}
}

func EncodeParallel(reader io.Reader, writer io.Writer) (int, error) {
	bufreader := bufio.NewReaderSize(reader, parInBufLen)
	bufwriter := bufio.NewWriter(writer)
	written := 0

	inchan := make(chan []byte, chanBuffers)
	outchan := make(chan []byte, chanBuffers)

	go readWorker(bufreader, inchan)
	go encodeWorker(inchan, outchan)

	for next := range outchan {
		n, err := bufwriter.Write(next)
		written += n
		if err != nil {
			return written, err
		}
	}

	bufwriter.Flush()
	return written, nil
}

func readWorker(reader io.Reader, inchan chan []byte) {
	for {
		buf := make([]byte, parInBufLen)
		octetsIn, err := io.ReadFull(reader, buf)
		if err == nil {
			// Most common case is a sucessful full-buffer read since buffer sizes are aligned.
			inchan <- buf[:octetsIn]
		} else if err == io.ErrUnexpectedEOF {
			inchan <- buf[:octetsIn]
			close(inchan)
			return
		} else if err == io.EOF {
			close(inchan)
			return
		}
	}
}

func encodeWorker(inchan chan []byte, outchan chan []byte) {
	for {
		next, more := <-inchan
		outbuf := make([]byte, 4_000)

		if !more {
			close(outchan)
			return
		}

		i := 0
		for ; i < len(next)/3; i++ {
			encodeTriplet(next[i*3:(i+1)*3], outbuf[i*4:(i+1)*4])
		}

		// Add trailing padding.
		r := len(next) % 3
		if r == 0 {
			outchan <- outbuf[:i*4]
		} else {
			encodeTrailingOctets(next[i*3:i*3+r], outbuf[i*4:(i+1)*4])
			outchan <- outbuf[:(i+1)*4]
		}
	}
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
