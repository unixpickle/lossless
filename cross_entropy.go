package lossless

import (
	"io"
	"math"

	"github.com/unixpickle/num-analysis/kahan"
)

const maxCrossEntropy = 256.0

// CrossEntropy computes the Shannon cross-entropy for
// the predictor on the given input stream.
//
// If the stream ends early with an error, this will
// return said error along with the cross entropy up
// to the error.
func CrossEntropy(p Predictor, input io.Reader) (float64, error) {
	result := kahan.NewSummer64()
	buffer := make([]byte, bufferSize)

	p.Reset()
	for {
		n, readErr := input.Read(buffer)
		if n == 0 {
			if readErr == io.EOF {
				break
			} else if readErr != nil {
				return result.Sum(), readErr
			}
			continue
		}

		for _, x := range buffer[:n] {
			pred := p.Predictions()[x]
			if pred == 0 {
				result.Add(maxCrossEntropy)
			} else {
				result.Add(-math.Log2(pred))
			}
			p.SawByte(x)
		}

		if readErr == io.EOF {
			break
		} else if readErr != nil {
			return result.Sum(), readErr
		}
	}

	return result.Sum(), nil
}
