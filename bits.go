package lossless

import "io"

type bitReader struct {
	R io.Reader

	curByte byte
	curIdx  uint
}

func (b *bitReader) ReadBit() (bool, error) {
	if b.curIdx&7 == 0 {
		var buf [1]byte
		if _, err := b.R.Read(buf[:]); err != nil {
			return false, err
		}
		b.curByte = buf[0]
		b.curIdx = 0
	}
	res := (b.curByte & (1 << b.curIdx)) != 0
	b.curIdx++
	return res, nil
}

type bitWriter struct {
	W io.Writer

	curByte byte
	curIdx  uint
}

func (b *bitWriter) WriteBit(bit bool) error {
	if bit {
		b.curByte |= (1 << b.curIdx)
	}
	b.curIdx++
	if b.curIdx == 8 {
		return b.Flush()
	} else {
		return nil
	}
}

func (b *bitWriter) Flush() error {
	b.curIdx = 0
	_, err := b.W.Write([]byte{b.curByte})
	b.curByte = 0
	return err
}
