package audio

import "fmt"

type State byte

const (
	COMPLETE = iota
	ATTACK
	DECAY
	SUSTAIN
	RELEASE
)

// an envelope returns a CV for a given time based params
// attack and decay are times in units of 32 samples (about 0.7ms)
type envelope interface {
	Trigger()
	Retrigger()
	Release()
	Scale(fp32) fp32
	applyPatch(p *patch)
}

type adeEnvelope struct {
	group     paramId
	gated     *boolparam
	retrigger *boolparam
	attack    *uint16param
	decay     *uint16param
	endLevel  *fp32param
	index     *fp32param // this is stored with the envelope in the digitone style algorithm, it scales the envelope output

	state       State
	sampleCount uint32
	current     fp32
}

func AdeEnvelope(group paramId) *adeEnvelope {
	return &adeEnvelope{group: group}
}

func (e *adeEnvelope) applyPatch(p *patch) {
	e.gated = p.BoolParam(ENV_GATED | e.group)
	e.retrigger = p.BoolParam(ENV_RETRIGGER | e.group)
	e.attack = p.Uint16Param(ENV_ATTACK | e.group)
	e.decay = p.Uint16Param(ENV_DECAY | e.group)
	e.endLevel = p.Fp32Param(ENV_ENDLEVEL | e.group)
	e.index = p.Fp32Param(ENV_INDEX | e.group)
}

func (e *adeEnvelope) Trigger() {
	e.state = ATTACK
	e.current = 0
	e.sampleCount = 0
}

func (e *adeEnvelope) Retrigger() {
	if !e.retrigger.Value() {
		return
	}
	e.Trigger()
}

func (e *adeEnvelope) Release() {
	if e.gated.Value() {
		e.state = DECAY
		e.sampleCount = 0
	}
}

// the modulation index is stored on the envelope in this scheme
// this function scales the index parameter by the current envelope amplitude
func (e *adeEnvelope) ScaledIndex() fp32 {
	return e.Scale(e.index.Value())
}

func (e *adeEnvelope) Scale(s fp32) fp32 {
	// attack and decay are times in units of 1024 samples (about 22.8us for 44.1kHz)
	// this way I can shift down the sample count 10 bits and divide
	e.sampleCount++

	switch e.state {
	case ATTACK:
		// skip the phase if time is 0
		if e.attack.Value() == 0 || e.current >= 1<<16 {
			e.current = 1 << 16
			if e.gated.Value() {
				e.state = SUSTAIN
				e.sampleCount = 0
			} else {
				e.state = DECAY
				e.sampleCount = 0
			}
		} else {
			e.current = fp32((e.sampleCount << 11) / uint32(e.attack.Value()))
		}
	case DECAY:
		if e.decay.Value() == 0 || e.current <= e.endLevel.Value() {
			e.current = e.endLevel.Value()
			e.state = COMPLETE
			e.sampleCount = 0
		} else {
			e.current = fp32(1<<16) - fp32((e.sampleCount<<11)/uint32(e.decay.Value())).mul(1<<16-e.endLevel.Value())
		}
	case SUSTAIN:
		e.current = 1 << 16
	case COMPLETE:
		e.current = e.endLevel.Value()
	}

	return s.mul(e.current)
}

type adsrEnvelope struct {
	group     paramId
	gated     *boolparam
	retrigger *boolparam
	attack    *uint16param
	decay     *uint16param
	sustain   *fp32param
	release   *uint16param

	state       State
	sampleCount uint32
	current     fp32
	ref         fp32
}

func AdsrEnvelope(group paramId) *adsrEnvelope {
	return &adsrEnvelope{group: group}
}

func (e *adsrEnvelope) applyPatch(p *patch) {
	e.gated = p.BoolParam(ENV_GATED | e.group)
	e.retrigger = p.BoolParam(ENV_RETRIGGER | e.group)
	e.attack = p.Uint16Param(ENV_ATTACK | e.group)
	e.decay = p.Uint16Param(ENV_DECAY | e.group)
	e.release = p.Uint16Param(ENV_RELEASE | e.group)
	e.sustain = p.Fp32Param(ENV_SUSTAIN | e.group)
}

func (e *adsrEnvelope) Trigger() {
	e.state = ATTACK
	e.ref = e.current
	e.sampleCount = 0
	fmt.Printf("triggering adsr: %d, %d, %d, %d\n", e.attack.Value(), e.decay.Value(), e.sustain.Value(), e.release.Value())
}

func (e *adsrEnvelope) Retrigger() {
	if !e.retrigger.Value() {
		return
	}
	e.Trigger()
}

func (e *adsrEnvelope) Release() {
	e.state = RELEASE
	e.ref = e.current
	e.sampleCount = 0
}

func (e *adsrEnvelope) Scale(s fp32) fp32 {
	e.sampleCount++

	switch e.state {
	case ATTACK:
		if e.attack.Value() == 0 || e.current >= 1<<16 {
			e.current = 1 << 16
			if e.gated.Value() {
				e.state = DECAY
				e.sampleCount = 0
			} else {
				e.state = RELEASE
				e.sampleCount = 0
			}
		} else {
			e.current = fp32((e.sampleCount << 11) / uint32(e.attack.Value()))
			// this nice little hack ensures that if we trigger during the release of a previous cycle,
			// the level stays continuous at where it was until the rise catches up
			// to avoid a click at the discontinuity when it drops to 0
			if e.ref > e.current {
				e.current = e.ref
			}
		}
	case DECAY:
		if e.decay.Value() == 0 || e.current <= e.sustain.Value() {
			e.current = e.sustain.Value()
			e.state = SUSTAIN
			e.ref = e.current
			e.sampleCount = 0
		} else {
			e.current = fp32(1<<16) - fp32((e.sampleCount<<11)/uint32(e.decay.Value())).mul(1<<16-e.sustain.Value())
		}
	case SUSTAIN:
		e.current = e.sustain.Value()
	case RELEASE:
		if e.release.Value() == 0 || e.current <= 0 {
			e.current = 0
			e.state = COMPLETE
			e.sampleCount = 0
		} else {
			e.current = fp32(1<<16 - (e.sampleCount<<11)/uint32(e.release.Value())).mul(e.ref)
		}
	case COMPLETE:
	}

	return s.mul(e.current)
}
