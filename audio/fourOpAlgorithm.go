package audio

import (
	"fmt"

	"github.com/ianmcmahon/fmsynth/fp"
	"github.com/ianmcmahon/fmsynth/patch"
)

type renderFunc func(*fourOpAlgorithm, []fp.Fp32)
type patchFunc func(*fourOpAlgorithm, *patch.Patch)

type algorithmVector struct {
	render     renderFunc
	applyPatch patchFunc
}

var algorithms []algorithmVector

type fourOpAlgorithm struct {
	voiceId      patch.ParamId
	algNum       patch.Param
	A, B1, B2, C *operator
	envA, envB   *adeEnvelope
	oprMix       patch.Param

	freq fp.Fp32
}

func (a *fourOpAlgorithm) applyPatch(p *patch.Patch) {
	a.algNum = p.ByteParam(patch.PATCH_ALGORITHM)
	algorithms[a.algNum.Value().(byte)].applyPatch(a, p)
	a.A.applyPatch(p)
	a.B1.applyPatch(p)
	a.B2.applyPatch(p)
	a.C.applyPatch(p)
	a.envA.applyPatch(p)
	a.envB.applyPatch(p)
	a.oprMix = p.Fp32Param(patch.PATCH_MIX)
}

func init() {
	fmt.Printf("four-op initializing algorithms\n")

	algorithms = make([]algorithmVector, 8)

	// gonna do my best to describe the algorithms without pictures
	// an operator with an f subscript (eg 'Af') gets the feeback param
	// A * B indicates A modulates B, A + B indicates mixing
	// x and y are inputs to the algorithm mixer

	algorithms[0] = algorithmVector{
		render: func(a *fourOpAlgorithm, out []fp.Fp32) {
			// digitone alg 1
			// y = (B2 * B1)
			// x = (Af + y) * C

			for i := range out {
				aVal := a.A.rotate(a.freq, 0).Mul(a.envA.ScaledIndex())
				b2Val := a.B2.rotate(a.freq, 0)
				b1Val := a.B1.rotate(a.freq, b2Val).Mul(a.envB.ScaledIndex())
				cMod := (aVal + b1Val) >> 1 // todo: is this atten necessary/desirable?
				cVal := a.C.rotate(a.freq, cMod)
				out[i] = crossMix(cVal, b1Val, a.oprMix.Value().(fp.Fp32))
			}
		},
		applyPatch: func(a *fourOpAlgorithm, p *patch.Patch) {
			a.A.feedback = p.Fp32Param(patch.PATCH_FEEDBACK)
		},
	}

	algorithms[1] = algorithmVector{
		render: func(a *fourOpAlgorithm, out []fp.Fp32) {
			// digitone alg 2
			// x = A * C
			// y = B2f * B1

			for i := range out {
				aVal := a.A.rotate(a.freq, 0).Mul(a.envA.ScaledIndex())
				x := a.C.rotate(a.freq, aVal)
				b2Val := a.B2.rotate(a.freq, 0)
				y := a.B1.rotate(a.freq, b2Val).Mul(a.envB.ScaledIndex())
				out[i] = crossMix(x, y, a.oprMix.Value().(fp.Fp32))
			}
		},
		applyPatch: func(a *fourOpAlgorithm, p *patch.Patch) {
			a.B2.feedback = p.Fp32Param(patch.PATCH_FEEDBACK)
		},
	}

	algorithms[2] = algorithmVector{
		render: func(a *fourOpAlgorithm, out []fp.Fp32) {
			// digitone alg 3
			// x = (Af * C) + (Af * B1)
			// y = Af * B2

			for i := range out {
				aVal := a.A.rotate(a.freq, 0).Mul(a.envA.ScaledIndex())
				y := a.B2.rotate(a.freq, aVal).Mul(a.envB.ScaledIndex())
				b1Val := a.B1.rotate(a.freq, aVal).Mul(a.envB.ScaledIndex())
				cVal := a.C.rotate(a.freq, aVal)
				x := (cVal + b1Val) >> 1
				out[i] = crossMix(x, y, a.oprMix.Value().(fp.Fp32))
			}
		},
		applyPatch: func(a *fourOpAlgorithm, p *patch.Patch) {
			a.A.feedback = p.Fp32Param(patch.PATCH_FEEDBACK)
		},
	}
}

func crossMix(a, b, mix fp.Fp32) fp.Fp32 {
	if mix > 1<<16 {
		mix = 1 << 16
	}
	if mix < 0 {
		mix = 0
	}
	invertedMix := 1<<16 - mix
	// mix level runs 0-1, 0.5 mixes them equally
	// sum and shift down a bit to attenuate
	return (a.Mul(invertedMix) + b.Mul(mix)) >> 1
}

func (a *fourOpAlgorithm) Render(out []fp.Fp32) {
	algorithms[a.algNum.Value().(byte)].render(a, out)
}

func (a *fourOpAlgorithm) Trigger(pitch fp.Fp32, velocity byte) {
	a.freq = pitch
	a.envA.Trigger()
	a.envB.Trigger()
}

func (a *fourOpAlgorithm) Retrigger(pitch fp.Fp32) {
	a.freq = pitch
	a.envA.Retrigger()
	a.envB.Retrigger()
}

func (a *fourOpAlgorithm) Release() {
	a.envA.Release()
	a.envB.Release()
}

func newFourOpAlgorithm(vId patch.ParamId) algorithm {
	return &fourOpAlgorithm{
		voiceId: vId,
		A:       Operator(patch.GRP_A),
		B1:      Operator(patch.GRP_B1),
		B2:      Operator(patch.GRP_B2),
		C:       Operator(patch.GRP_C),
		envA:    AdeEnvelope(patch.GRP_A),
		envB:    AdeEnvelope(patch.GRP_B),
	}
}
