package audio

import "math"

var sineTable = makeSineTable(SAMPLING_RATE)

// an operator is a single oscillator that can be phase modulated
type operator struct {
	group paramId
	ratio *fp32param
	phase fp32
}

func Operator(group paramId) *operator {
	return &operator{group: group}
}

func (o *operator) applyPatch(p *patch) {
	o.ratio = p.Fp32Param(OPR_RATIO | o.group)
}

// increments the phase based on frequency and returns the next sample
func (o *operator) rotate(freq, mod fp32) fp32 {
	f := freq.mul(o.ratio.Value()) + mod

	// phase (pitch / table_freq) * (table_len / sampling_rate) I believe
	// but since our table is 1Hz at sampling_rate, len/rate = 1 and pitch/table = pitch
	o.phase += f >> 16
	if o.phase >= SAMPLING_RATE {
		o.phase -= SAMPLING_RATE
	}
	return sineTable[o.phase]
}

// an algorithm is a particular configuration of operators and envelopes
// that exposes parameter inputs
type algorithm interface {
	Trigger(pitch fp32, velocity byte)
	Retrigger(pitch fp32)
	Release()
	Render(out []fp32)
	applyPatch(p *patch)
}

type digitoneFourOpAlgorithm struct {
	A, B1, B2, C *operator
	envA, envB   *envelope
}

func makeSineTable(tableLength int) []fp32 {
	table := make([]fp32, tableLength)
	for i := range table {
		phase := float64(i) / float64(tableLength)
		table[i] = float2fp32(math.Sin(2 * math.Pi * phase))
		// this is a 1Hz wave across the table,
		// so a complete rotation of phase is one cycle thru table
	}
	return table
}
