package audio

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
}

type adeEnvelope struct {
	gated     bool
	retrigger bool
	attack    uint16
	decay     uint16
	endLevel  fp32

	state       State
	sampleCount uint32
	current     fp32
}

func (e *adeEnvelope) Trigger() {
	e.state = ATTACK
	e.current = 0
	e.sampleCount = 0
}

func (e *adeEnvelope) Retrigger() {
	if !e.retrigger {
		return
	}
	e.Trigger()
}

func (e *adeEnvelope) Release() {
	if e.gated {
		e.state = DECAY
		e.sampleCount = 0
	}
}

func (e *adeEnvelope) Scale(s fp32) fp32 {
	// attack and decay are times in units of 1024 samples (about 22.8us for 44.1kHz)
	// this way I can shift down the sample count 10 bits and divide
	e.sampleCount++

	switch e.state {
	case ATTACK:
		if e.current >= 1<<16 {
			e.current = 1 << 16
			if e.gated {
				e.state = SUSTAIN
				e.sampleCount = 0
			} else {
				e.state = DECAY
				e.sampleCount = 0
			}
		} else {
			e.current = fp32((e.sampleCount << 11) / uint32(e.attack))
		}
	case DECAY:
		if e.current <= e.endLevel {
			e.current = e.endLevel
			e.state = COMPLETE
			e.sampleCount = 0
		} else {
			e.current = fp32(1<<16) - fp32((e.sampleCount<<11)/uint32(e.decay)).mul(1<<16-e.endLevel)
		}
	case SUSTAIN:
		e.current = 1 << 16
	case COMPLETE:
		e.current = e.endLevel
	}

	return s.mul(e.current)
}

type adsrEnvelope struct {
	gated     bool
	retrigger bool
	attack    uint16
	decay     uint16
	sustain   fp32
	release   uint16

	state       State
	sampleCount uint32
	current     fp32
	ref         fp32
}

func (e *adsrEnvelope) Trigger() {
	e.state = ATTACK
	e.ref = e.current
	e.sampleCount = 0
}

func (e *adsrEnvelope) Retrigger() {
	if !e.retrigger {
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
		if e.current >= 1<<16 {
			e.current = 1 << 16
			if e.gated {
				e.state = DECAY
				e.sampleCount = 0
			} else {
				e.state = RELEASE
				e.sampleCount = 0
			}
		} else {
			e.current = fp32((e.sampleCount << 11) / uint32(e.attack))
			// this nice little hack ensures that if we trigger during the release of a previous cycle,
			// the level stays continuous at where it was until the rise catches up
			// to avoid a click at the discontinuity when it drops to 0
			if e.ref > e.current {
				e.current = e.ref
			}
		}
	case DECAY:
		if e.current <= e.sustain {
			e.current = e.sustain
			e.state = SUSTAIN
			e.ref = e.current
			e.sampleCount = 0
		} else {
			e.current = fp32(1<<16) - fp32((e.sampleCount<<11)/uint32(e.decay)).mul(1<<16-e.sustain)
		}
	case SUSTAIN:
		e.current = e.sustain
	case RELEASE:
		if e.current <= 0 {
			e.current = 0
			e.state = COMPLETE
			e.sampleCount = 0
		} else {
			e.current = fp32(1<<16 - (e.sampleCount<<11)/uint32(e.release)).mul(e.ref)
		}
	case COMPLETE:
	}

	return s.mul(e.current)
}
