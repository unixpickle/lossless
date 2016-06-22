package lossless

import "io"

func Compress(p Predictor, input io.Reader, output io.Writer) error {
	w := &bitWriter{W: output}
	p.Reset()
	for {
		dist := p.Predictions()
		var buf [1]byte
		if _, err := input.Read(buf[:]); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if err := encodeByte(w, dist, buf[0]); err != nil {
			return err
		}
		p.SawByte(buf[0])
	}
}

func Decompress(p Predictor, input io.Reader, output io.Writer) error {
	r := &bitReader{R: input}
	p.Reset()
	for {
		dist := p.Predictions()
		b, err := decodeByte(r, dist)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if _, err := output.Write([]byte{b}); err != nil {
			return err
		}
		p.SawByte(b)
	}
}
