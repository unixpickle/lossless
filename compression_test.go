package lossless

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"testing"
)

func TestEnglishCompression(t *testing.T) {
	rand.Seed(123)
	predictor := &EnglishPredictor{}
	for i := 0; i < 5; i++ {
		bufLen := 0x3000 + i
		var buffer bytes.Buffer
		var origCheck bytes.Buffer
		for i := 0; i < bufLen; i++ {
			b := byte(rand.Intn(0x100))
			buffer.WriteByte(b)
			origCheck.WriteByte(b)
		}
		var output bytes.Buffer
		if err := Compress(predictor, &buffer, &output); err != nil {
			t.Fatal("failed to compress:", err)
		} else {
			var orig bytes.Buffer
			if err := Decompress(predictor, &output, &orig); err != nil {
				t.Fatal("failed to decompress:", err)
			}
			for i, x := range origCheck.Bytes() {
				a := orig.Bytes()[i]
				if a != x {
					t.Fatal("decompressed data differs at byte", i)
				}
			}
		}
	}
}

func BenchmarkEnglishCompression(b *testing.B) {
	var buffer bytes.Buffer
	for i := 0; i < b.N; i++ {
		buffer.WriteByte(byte(rand.Intn(0x100)))
	}
	b.ResetTimer()
	Compress(&EnglishPredictor{}, &buffer, ioutil.Discard)
}
