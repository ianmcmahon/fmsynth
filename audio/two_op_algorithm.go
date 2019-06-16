package audio

type twoOpAlgorithm struct {
	voiceId paramId
	C, M    *operator
	ratio   *fp32param
	mLevel  *fp32param
	mEnv    envelope
}

func (a *twoOpAlgorithm) Render(out []fp32) {
	for i := range out {
		// first rotate the modulator and get the sample
		mVal := a.M.rotate(0)

		// then rotate the carrier expected amount plus the value of M's sample
		// scaled by the envelope * index (mLevel)
		cVal := a.C.rotate(mVal.mul(a.mEnv.Scale(a.mLevel.Value())))

		out[i] = cVal
	}
}

func (a *twoOpAlgorithm) Trigger(pitch fp32, velocity byte) {
	a.C.freq.Set(pitch)
	a.M.freq.Set(pitch)
	a.mEnv.Trigger()
}

func (a *twoOpAlgorithm) Retrigger(pitch fp32) {
	a.C.freq.Set(pitch)
	a.M.freq.Set(pitch)
	a.mEnv.Retrigger()
}

func (a *twoOpAlgorithm) Release() {
	a.mEnv.Release()
}

func newTwoOpAlgorithm(vId paramId) algorithm {
	return &twoOpAlgorithm{
		voiceId: vId,
		C: &operator{
			freq:  newFp32Param(vId|C_FREQ, 440.0),
			ratio: newFp32Param(vId|C_RATIO, 1.0),
		},
		M: &operator{
			freq:  newFp32Param(vId|A_FREQ, 440.0),
			ratio: newFp32Param(vId|A_RATIO, 1.0),
		},
		ratio:  newFp32Param(vId|A_RATIO, 0.0),
		mLevel: newFp32Param(vId|AENV_LEVEL, 0.0),
		mEnv: &adeEnvelope{
			gated:     newBoolParam(vId|AENV_GATED, false),
			retrigger: newBoolParam(vId|AENV_RETRIGGER, true),
			attack:    newUint16Param(vId|AENV_ATTACK, 200),
			decay:     newUint16Param(vId|AENV_DECAY, 800),
			endLevel:  newFp32Param(vId|AENV_ENDLEVEL, 0.1),
		},
	}
}
