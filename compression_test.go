package lossless

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"testing"
)

func BenchmarkEnglishCompression(b *testing.B) {
	var buffer bytes.Buffer
	for i := 0; i < b.N; i++ {
		buffer.WriteByte(byte(rand.Intn(0x100)))
	}
	b.ResetTimer()
	Decompress(&EnglishPredictor{}, &buffer, ioutil.Discard)
}
