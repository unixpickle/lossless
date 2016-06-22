package lossless

type MarkovPredictor struct {
	Grams int

	lastBlock []byte
	history   map[string][]int
}

func (m *MarkovPredictor) New() Predictor {
	return &MarkovPredictor{Grams: m.Grams}
}

func (m *MarkovPredictor) Predictions() ByteProbs {
	if len(m.lastBlock) < m.Grams {
		return EqualByteProbs
	}

	if m.history == nil {
		m.history = map[string][]int{}
	}

	hist, ok := m.history[string(m.lastBlock)]
	if !ok {
		return EqualByteProbs
	}

	var res ByteProbs
	var histTotal int
	for _, x := range hist {
		histTotal += x
	}
	scale := 1.0 / float64(histTotal)
	for i, x := range hist {
		res[i] = float64(x) * scale
	}
	return res
}

func (m *MarkovPredictor) SawByte(b byte) {
	if len(m.lastBlock) < m.Grams {
		m.lastBlock = append(m.lastBlock, b)
		return
	}

	if m.history == nil {
		m.history = map[string][]int{}
	}

	val, ok := m.history[string(m.lastBlock)]
	if !ok {
		val = make([]int, 0x100)
		m.history[string(m.lastBlock)] = val
	}
	val[int(b)]++

	copy(m.lastBlock, m.lastBlock[:len(m.lastBlock)-1])
	m.lastBlock[len(m.lastBlock)-1] = b
}

func (m *MarkovPredictor) Reset() {
	m.history = nil
	m.lastBlock = nil
}
