package audio

type twoOpAlgorithm struct {
	C, M   *operator
	ratio  fp32
	mLevel fp32
	mEnv   envelope
}

func (a *twoOpAlgorithm) Render(out []fp32) {
	for i := range out {
		// first increment the phase of the modulator and look up the sample
		mVal := a.M.rotate(a.M.freq)

		// then increment the phase of the carrier both the frequency
		// and the output sample from M, scaled by the level and envelope

		cVal := a.C.rotate(a.C.freq + mVal.mul(a.mEnv.Scale(a.mLevel)))

		out[i] = cVal
	}
}

func (a *twoOpAlgorithm) Trigger(pitch fp32, velocity byte) {
	a.C.freq = pitch
	a.M.freq = pitch.mul(a.ratio)
	a.mEnv.Trigger()
}

func (a *twoOpAlgorithm) Retrigger(pitch fp32) {
	a.C.freq = pitch
	a.M.freq = pitch.mul(a.ratio)
	a.mEnv.Retrigger()
}

func (a *twoOpAlgorithm) Release() {
	a.mEnv.Release()
}

func newTwoOpAlgorithm() algorithm {
	return &twoOpAlgorithm{
		C:      &operator{},
		M:      &operator{},
		ratio:  fp32(1 << 16),
		mLevel: fp32(1 << 24),
		mEnv: &adeEnvelope{
			gated:     true,
			retrigger: true,
			attack:    5,
			decay:     20,
			endLevel:  float2fp32(0.25),
		},
	}
}
