package audio

import (
	"fmt"
	"math"
)

type Voice struct {
	id      int
	notesOn []byte

	A  *Oscillator
	B1 *Oscillator
	B2 *Oscillator
	C  *Oscillator
}

func (v *Voice) CurNote() byte {
	if len(v.notesOn) == 0 {
		return 0
	}
	return v.notesOn[len(v.notesOn)-1]
}

func (engine *Engine) NewSimpleVoice(id int, output chan<- Sample) *Voice {
	mixer := NewMixer(2, output)
	v := &Voice{
		id:      id,
		notesOn: make([]byte, 0),
		A:       engine.NewOscillator(mixer.Input(0)),
		C:       engine.NewOscillator(mixer.Input(1)),
	}

	return v
}

func (v *Voice) Trigger(pitch, velocity float64) {
	fmt.Printf("triggering %.2fhz\n", pitch)
	v.A.pitch.value = pitch * 2.0
	v.A.phaseIdx = 0
	v.A.amp.value = velocity * 0.7
	v.C.pitch.value = pitch
	v.C.phaseIdx = 0
	v.C.amp.value = velocity
}

func (v *Voice) Retrigger(pitch float64) {
	v.A.pitch.value = pitch * 2.0
	v.A.phaseIdx = 0
	v.C.pitch.value = pitch
	v.C.phaseIdx = 0
}

func (v *Voice) Release() {
	v.C.amp.value = 0.0
	v.A.amp.value = 0.0
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

	v.Trigger(note2freq(note), float64(velocity)/127.0)
}

func (v *Voice) NoteOff(note byte) {
	for i, n := range v.notesOn {
		if n == note {
			copy(v.notesOn[i:], v.notesOn[i+1:])
			v.notesOn = v.notesOn[:len(v.notesOn)-1]
		}
	}

	if len(v.notesOn) > 0 {
		v.Retrigger(note2freq(v.notesOn[len(v.notesOn)-1]))
	} else {
		v.Release()
	}
}

func note2freq(note byte) float64 {
	return math.Pow(2, (float64(note)-69.0)/12.0) * 440.0
}
