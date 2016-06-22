package lossless

import (
	"sort"
	"sync"
)

// ByteProbs is a list where the i-th element is the
// probability that a given byte is byte(i).
// In other words, ByteProbs is a probability
// distribution over byte values.
//
// For example, if every byte is equally likely, then
// the corresponding ByteProbs would be full of 1/256.
type ByteProbs [256]float64

// A Predictor is any algorithm with can predict the
// next byte in a stream of data.
type Predictor interface {
	// Predictions returns the current predictions for
	// the next byte in the sequence.
	Predictions() ByteProbs

	// SawByte tells the model which byte was actually
	// seen, allowing the model to update its predictions
	// for the byte after this one.
	SawByte(b byte)

	// Reset resets the model's memory, indicating that
	// a new stream of data is being processed.
	Reset()
}

var predictorTable = map[string]Predictor{}
var tableLock sync.Mutex

// RegisterPredictor registers a Predictor for the
// given unique identifier.
//
// It is safe to call this concurrently with the
// other Predictor management functions.
func RegisterPredictor(id string, p Predictor) {
	tableLock.Lock()
	defer tableLock.Unlock()
	predictorTable[id] = p
}

// GetPredictor returns the Predictor for the unique
// identifier, or nil if the ID is not registered.
//
// It is safe to call this concurrently with the
// other Predictor management functions.
func GetPredictor(id string) Predictor {
	tableLock.Lock()
	defer tableLock.Unlock()
	return predictorTable[id]
}

// PredictorIDs returns an alphabetically sorted list
// of Predictor unique identifiers.
//
// It is safe to call this concurrently with the
// other Predictor management functions.
func PredictorIDs() []string {
	tableLock.Lock()
	defer tableLock.Unlock()
	res := make([]string, 0, len(predictorTable))
	for id := range predictorTable {
		res = append(res, id)
	}
	sort.Strings(res)
	return res
}
