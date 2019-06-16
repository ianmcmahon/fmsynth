package audio

import "math"

var sineTable = makeSineTable(SAMPLING_RATE)

// an operator is a single oscillator that can be phase modulated
type operator struct {
	freq  *fp32param
	ratio *fp32param
	phase fp32
}

// increments the phase based on frequency and returns the next sample
// uses the param value of freq (which will incorporate mods eventually)
// times the param val of ratio, plus phase increment (frequency modulation)
func (o *operator) rotate(phaseIncr fp32) fp32 {
	freq := o.freq.Value().mul(o.ratio.Value()) + phaseIncr

	// phase (pitch / table_freq) * (table_len / sampling_rate) I believe
	// but since our table is 1Hz at sampling_rate, len/rate = 1 and pitch/table = pitch
	o.phase += freq >> 16
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
