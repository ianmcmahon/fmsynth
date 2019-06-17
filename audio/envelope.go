package audio

import (
	"github.com/ianmcmahon/fmsynth/fp"
	"github.com/ianmcmahon/fmsynth/patch"
)

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
	Scale(fp.Fp32) fp.Fp32
	applyPatch(p *patch.Patch)
}

type adeEnvelope struct {
	group     patch.ParamId
	gated     patch.Param
	retrigger patch.Param
	attack    patch.Param
	decay     patch.Param
	endLevel  patch.Param
	index     patch.Param // this is stored with the envelope in the digitone style algorithm, it scales the envelope output

	state       State
	sampleCount uint32
	current     fp.Fp32
}

func AdeEnvelope(group patch.ParamId) *adeEnvelope {
	return &adeEnvelope{group: group}
}

func (e *adeEnvelope) applyPatch(p *patch.Patch) {
	e.gated = p.BoolParam(patch.ENV_GATED | e.group)
	e.retrigger = p.BoolParam(patch.ENV_RETRIGGER | e.group)
	e.attack = p.Uint16Param(patch.ENV_ATTACK | e.group)
	e.decay = p.Uint16Param(patch.ENV_DECAY | e.group)
	e.endLevel = p.Fp32Param(patch.ENV_ENDLEVEL | e.group)
	e.index = p.Fp32Param(patch.ENV_INDEX | e.group)
}

func (e *adeEnvelope) Trigger() {
	e.state = ATTACK
	e.current = 0
	e.sampleCount = 0
}

func (e *adeEnvelope) Retrigger() {
	if !e.retrigger.Value().(bool) {
		return
	}
	e.Trigger()
}

func (e *adeEnvelope) Release() {
	if e.gated.Value().(bool) {
		e.state = DECAY
		e.sampleCount = 0
	}
}

// the modulation index is stored on the envelope in this scheme
// this function scales the index parameter by the current envelope amplitude
func (e *adeEnvelope) ScaledIndex() fp.Fp32 {
	return e.Scale(e.index.Value().(fp.Fp32))
}

func (e *adeEnvelope) Scale(s fp.Fp32) fp.Fp32 {
	// attack and decay are times in units of 1024 samples (about 22.8us for 44.1kHz)
	// this way I can shift down the sample count 10 bits and divide
	e.sampleCount++

	attack := e.attack.Value().(uint16)
	decay := e.decay.Value().(uint16)
	endlevel := e.endLevel.Value().(fp.Fp32)
	switch e.state {
	case ATTACK:
		// can't have a zero attack or we get phase discontinuities (clicks)
		if attack == 0 {
			attack++
		}
		if e.current >= 1<<16 {
			e.current = 1 << 16
			if e.gated.Value().(bool) {
				e.state = SUSTAIN
				e.sampleCount = 0
			} else {
				e.state = DECAY
				e.sampleCount = 0
			}
		} else {
			e.current = fp.Fp32((e.sampleCount << 11) / uint32(attack))
		}
	case DECAY:
		if decay == 0 {
			decay++
		}
		if decay == 0 || e.current <= endlevel {
			e.current = endlevel
			e.state = COMPLETE
			e.sampleCount = 0
		} else {
			e.current = fp.Fp32(1<<16) - fp.Fp32((e.sampleCount<<11)/uint32(decay)).Mul(1<<16-endlevel)
		}
	case SUSTAIN:
		e.current = 1 << 16
	case COMPLETE:
		e.current = endlevel
	}

	return s.Mul(e.current)
}

type adsrEnvelope struct {
	group     patch.ParamId
	gated     patch.Param
	retrigger patch.Param
	attack    patch.Param
	decay     patch.Param
	sustain   patch.Param
	release   patch.Param

	state       State
	sampleCount uint32
	current     fp.Fp32
	ref         fp.Fp32
}

func AdsrEnvelope(group patch.ParamId) *adsrEnvelope {
	return &adsrEnvelope{group: group}
}

func (e *adsrEnvelope) applyPatch(p *patch.Patch) {
	e.gated = p.BoolParam(patch.ENV_GATED | e.group)
	e.retrigger = p.BoolParam(patch.ENV_RETRIGGER | e.group)
	e.attack = p.Uint16Param(patch.ENV_ATTACK | e.group)
	e.decay = p.Uint16Param(patch.ENV_DECAY | e.group)
	e.release = p.Uint16Param(patch.ENV_RELEASE | e.group)
	e.sustain = p.Fp32Param(patch.ENV_SUSTAIN | e.group)
}

func (e *adsrEnvelope) Trigger() {
	e.state = ATTACK
	e.ref = e.current
	e.sampleCount = 0
}

func (e *adsrEnvelope) Retrigger() {
	if !e.retrigger.Value().(bool) {
		return
	}
	e.Trigger()
}

func (e *adsrEnvelope) Release() {
	e.state = RELEASE
	e.ref = e.current
	e.sampleCount = 0
}

func (e *adsrEnvelope) Scale(s fp.Fp32) fp.Fp32 {
	e.sampleCount++

	attack := e.attack.Value().(uint16)
	decay := e.decay.Value().(uint16)
	sustain := e.sustain.Value().(fp.Fp32)
	release := e.release.Value().(uint16)
	switch e.state {
	case ATTACK:
		if attack == 0 {
			attack++
		}
		if attack == 0 || e.current >= 1<<16 {
			e.current = 1 << 16
			if e.gated.Value().(bool) {
				e.state = DECAY
				e.sampleCount = 0
			} else {
				e.state = RELEASE
				e.sampleCount = 0
			}
		} else {
			e.current = fp.Fp32((e.sampleCount << 11) / uint32(attack))
			// this nice little hack ensures that if we trigger during the release of a previous cycle,
			// the level stays continuous at where it was until the rise catches up
			// to avoid a click at the discontinuity when it drops to 0
			if e.ref > e.current {
				e.current = e.ref
			}
		}
	case DECAY:
		if decay == 0 {
			decay++
		}
		if e.current <= sustain {
			e.current = sustain
			e.state = SUSTAIN
			e.ref = e.current
			e.sampleCount = 0
		} else {
			e.current = fp.Fp32(1<<16) - fp.Fp32((e.sampleCount<<11)/uint32(decay)).Mul(1<<16-sustain)
		}
	case SUSTAIN:
		e.current = sustain
	case RELEASE:
		if release == 0 {
			release++
		}
		if e.current <= 0 {
			e.current = 0
			e.state = COMPLETE
			e.sampleCount = 0
		} else {
			e.current = fp.Fp32(1<<16 - (e.sampleCount<<11)/uint32(release)).Mul(e.ref)
		}
	case COMPLETE:
	}

	return s.Mul(e.current)
}
