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

	triggered   bool
	inDecay     bool
	sampleCount uint32
	current     fp32
}

func (e *adeEnvelope) Trigger() {
	e.triggered = true
	e.current = 0
	e.sampleCount = 0
}

func (e *adeEnvelope) Retrigger() {
	if !e.retrigger && e.triggered {
		return
	}
	e.Trigger()
}

func (e *adeEnvelope) Release() {
	e.triggered = false
}

func (e *adeEnvelope) Scale(s fp32) fp32 {
	// attack and decay are times in units of 1024 samples (about 22.8us for 44.1kHz)
	// this way I can shift down the sample count 10 bits and divide
	e.sampleCount++

	if e.triggered && !e.inDecay {
		// attack phase
		// if we're saturated, short circuit
		if e.current >= 1<<16 {
			e.current = 1 << 16
			if !e.gated {
				e.inDecay = true
			}
			return e.current
		}
		// attack phase needs to complete in a.attack * 1024 samples
		// units of time so far is sampleCount/1024
		// proportion of attack phase completed is timeSoFar/a.attack
		// I'm shifting up six bits which is equivalent to converting to fp32
		// and then shifting down 10 bits
		// this gives us an amplitude value scaled 0.0-1.0 in fp32
		e.current = fp32((e.sampleCount << 6) / uint32(e.attack))
		return e.current
	} else {
		if e.current > 1 {
			e.inDecay = true
		}
	}

	if e.inDecay {
		// release phase
		// if we're saturated, that's end of phase
		if e.current <= e.endLevel {
			e.current = e.endLevel
			e.inDecay = false
			return e.current
		}
		// for decay, we do the same as attack, but current is 1 - (samples/1024)/decay
		e.current = fp32(1<<16 - (e.sampleCount<<6)/uint32(e.decay))
		return e.current
	}

	return e.current
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
	fmt.Printf("trigger, going to ATTACK\n")
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
	fmt.Printf("release, going to RELEASE\n")
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
				fmt.Printf("going to DECAY\n")
				e.state = DECAY
				e.sampleCount = 0
			} else {
				fmt.Printf("going to RELEASE\n")
				e.state = RELEASE
				e.sampleCount = 0
			}
		} else {
			e.current = fp32((e.sampleCount << 6) / uint32(e.attack))
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
			fmt.Printf("going to SUSTAIN\n")
		} else {
			e.current = fp32(1<<16) - fp32((e.sampleCount<<6)/uint32(e.decay)).mul(1<<16-e.sustain)
		}
	case SUSTAIN:
		e.current = e.sustain
	case RELEASE:
		if e.current <= 0 {
			e.current = 0
			e.state = COMPLETE
			e.sampleCount = 0
			fmt.Printf("going to COMPLETE\n")
		} else {
			e.current = fp32(1<<16 - (e.sampleCount<<6)/uint32(e.release)).mul(e.ref)
		}
	case COMPLETE:
	}

	if e.state != COMPLETE && e.sampleCount%1024 == 0 {
		fmt.Printf("current: %.2f\n", float64(e.current)/65536.0)
	}

	return s.mul(e.current)
}
