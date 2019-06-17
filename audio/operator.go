package audio

import (
	"math"

	"github.com/ianmcmahon/fmsynth/fp"
	"github.com/ianmcmahon/fmsynth/patch"
)

var sineTable = makeSineTable(SAMPLING_RATE)

// an operator is a single oscillator that can be phase modulated
type operator struct {
	group    patch.ParamId
	ratio    patch.Param
	feedback patch.Param
	phase    fp.Fp32
}

func Operator(group patch.ParamId) *operator {
	return &operator{group: group}
}

func (o *operator) applyPatch(p *patch.Patch) {
	o.ratio = p.Fp32Param(patch.OPR_RATIO | o.group)
}

// increments the phase based on frequency and returns the next sample
func (o *operator) rotate(freq, mod fp.Fp32) fp.Fp32 {
	f := freq.Mul(o.ratio.Value().(fp.Fp32)) + mod

	// phase (pitch / table_freq) * (table_len / sampling_rate) I believe
	// but since our table is 1Hz at sampling_rate, len/rate = 1 and pitch/table = pitch
	o.phase += f >> 16
	if o.phase >= SAMPLING_RATE {
		o.phase -= SAMPLING_RATE
	}
	if o.phase < 0 {
		o.phase = 0
	}
	sample := sineTable[o.phase]

	if o.feedback != nil && o.feedback.Value().(fp.Fp32) != 0 {
		// now apply feedback
		o.phase += sample.Mul(o.feedback.Value().(fp.Fp32)) >> 16
		if o.phase >= SAMPLING_RATE {
			o.phase -= SAMPLING_RATE
		}
		if o.phase < 0 {
			o.phase = 0
		}
		sample = sineTable[o.phase]
	}
	return sample
}

// an algorithm is a particular configuration of operators and envelopes
// that exposes parameter inputs
type algorithm interface {
	Trigger(pitch fp.Fp32, velocity byte)
	Retrigger(pitch fp.Fp32)
	Release()
	Render(out []fp.Fp32)
	applyPatch(p *patch.Patch)
}

type digitoneFourOpAlgorithm struct {
	A, B1, B2, C *operator
	envA, envB   *envelope
}

func makeSineTable(tableLength int) []fp.Fp32 {
	table := make([]fp.Fp32, tableLength)
	for i := range table {
		phase := float64(i) / float64(tableLength)
		table[i] = fp.Float2Fp32(math.Sin(2 * math.Pi * phase))
		// this is a 1Hz wave across the table,
		// so a complete rotation of phase is one cycle thru table
	}
	return table
}
