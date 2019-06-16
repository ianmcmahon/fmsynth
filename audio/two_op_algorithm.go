package audio

type twoOpAlgorithm struct {
	voiceId paramId
	C, M    *operator
	mEnv    *adeEnvelope

	freq fp32
}

func (a *twoOpAlgorithm) applyPatch(p *patch) {
	a.C.applyPatch(p)
	a.M.applyPatch(p)
	a.mEnv.applyPatch(p)
}

func (a *twoOpAlgorithm) Render(out []fp32) {
	for i := range out {
		// first rotate the modulator and get the sample
		mVal := a.M.rotate(a.freq, 0)

		// then rotate the carrier expected amount plus the value of M's sample
		// scaled by the envelope * index (mLevel)
		cVal := a.C.rotate(a.freq, mVal.mul(a.mEnv.ScaledIndex()))

		out[i] = cVal
	}
}

func (a *twoOpAlgorithm) Trigger(pitch fp32, velocity byte) {
	a.freq = pitch
	a.mEnv.Trigger()
}

func (a *twoOpAlgorithm) Retrigger(pitch fp32) {
	a.freq = pitch
	a.mEnv.Retrigger()
}

func (a *twoOpAlgorithm) Release() {
	a.mEnv.Release()
}

func newTwoOpAlgorithm(vId paramId) algorithm {
	return &twoOpAlgorithm{
		voiceId: vId,
		C:       Operator(GRP_C),
		M:       Operator(GRP_A),
		mEnv:    AdeEnvelope(GRP_A),
	}
}
