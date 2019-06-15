package audio

import (
	"fmt"
	"math"
)

type Voice struct {
	id      int
	notesOn []byte

	alg algorithm
	vca envelope
}

func (v *Voice) CurNote() byte {
	if len(v.notesOn) == 0 {
		return 0
	}
	return v.notesOn[len(v.notesOn)-1]
}

func (engine *Engine) NewSimpleVoice(id int) *Voice {
	v := &Voice{
		id:      id,
		notesOn: make([]byte, 0),
		alg:     newTwoOpAlgorithm(float2fp32(11), 1<<24),
		vca: &adsrEnvelope{
			gated:     true,
			retrigger: false,
			attack:    100,
			decay:     100,
			sustain:   float2fp32(0.3),
			release:   200,
		},
	}

	return v
}

func (v *Voice) Render(out []fp32) {
	v.alg.Render(out)
	for i, s := range out {
		out[i] = v.vca.Scale(s)
	}
}

func (v *Voice) trigger(pitch fp32, velocity byte) {
	v.alg.Trigger(pitch, velocity)
	v.vca.Trigger()
}

func (v *Voice) retrigger(pitch fp32) {
	v.alg.Retrigger(pitch)
	v.vca.Retrigger()
}

func (v *Voice) release() {
	v.alg.Release()
	v.vca.Release()
}

func (v *Voice) NoteOn(note, velocity byte) {
	if note < 0 || note > 127 {
		return
	}
	if velocity < 0 {
		velocity = 0
	}
	if velocity > 127 {
		velocity = 127
	}

	on := false
	fmt.Printf("%d: note on: %d: %v\n", v.id, note, v.notesOn)
	for _, n := range v.notesOn {
		if n == note {
			on = true
		}
	}

	if !on {
		v.notesOn = append(v.notesOn, note)
	}

	v.trigger(note2freq(note), velocity)
}

func (v *Voice) NoteOff(note byte) {
	for i, n := range v.notesOn {
		if n == note {
			copy(v.notesOn[i:], v.notesOn[i+1:])
			v.notesOn = v.notesOn[:len(v.notesOn)-1]
		}
	}

	if len(v.notesOn) > 0 {
		v.retrigger(note2freq(v.notesOn[len(v.notesOn)-1]))
	} else {
		v.release()
	}
}

func note2freq(note byte) fp32 {
	return float2fp32(math.Pow(2, (float64(note)-69.0)/12.0) * 440.0)
}
