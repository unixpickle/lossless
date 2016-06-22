package lossless

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var headerByteOrder = binary.LittleEndian

const bufferSize = 4096

func Compress(p Predictor, input io.Reader, output io.Writer) error {
	p.Reset()
	buffer := make([]byte, bufferSize)
	for {
		n, readErr := input.Read(buffer)
		if n == 0 {
			if readErr == io.EOF {
				return nil
			} else if readErr != nil {
				return readErr
			} else {
				continue
			}
		}

		sizeField := uint32(n)
		if err := binary.Write(output, headerByteOrder, &sizeField); err != nil {
			return err
		}
		encoded := compressBlock(p, buffer[:n])
		if _, err := output.Write(encoded); err != nil {
			return err
		}

		if readErr == io.EOF {
			return nil
		} else if readErr != nil {
			return readErr
		}
	}
}

func Decompress(p Predictor, input io.Reader, output io.Writer) error {
	p.Reset()
	for {
		var blockSize uint32
		if err := binary.Read(input, headerByteOrder, &blockSize); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		block, err := decompressBlock(p, input, int(blockSize))
		if err != nil {
			return errors.New("buffer underflow: " + err.Error())
		}
		if _, err := output.Write(block); err != nil {
			return err
		}
	}
}

func compressBlock(p Predictor, block []byte) []byte {
	var outBlock bytes.Buffer
	w := &bitWriter{W: &outBlock}

	for _, b := range block {
		dist := p.Predictions()
		encodeByte(w, dist, b)
		p.SawByte(b)
	}

	w.Flush()
	return outBlock.Bytes()
}

func decompressBlock(p Predictor, input io.Reader, blockLen int) ([]byte, error) {
	outBlock := make([]byte, blockLen)

	r := &bitReader{R: input}
	for i := 0; i < blockLen; i++ {
		dist := p.Predictions()
		b, err := decodeByte(r, dist)
		if err != nil {
			return nil, err
		}
		outBlock[i] = b
		p.SawByte(b)
	}

	return outBlock, nil
}
