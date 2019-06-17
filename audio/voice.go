package audio

import (
	"fmt"
	"math"

	"github.com/ianmcmahon/fmsynth/fp"
	"github.com/ianmcmahon/fmsynth/patch"
)

type Voice struct {
	id      patch.ParamId
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

func (engine *Engine) NewSimpleVoice(id byte) *Voice {
	vId := patch.ParamId(id) << 8
	v := &Voice{
		id:      vId,
		notesOn: make([]byte, 0),
		alg:     newFourOpAlgorithm(vId),
		vca:     AdsrEnvelope(patch.GRP_VCA),
	}

	return v
}

func (v *Voice) applyPatch(p *patch.Patch) {
	v.alg.applyPatch(p)
	v.vca.applyPatch(p)
}

func (v *Voice) Render(out []fp.Fp32) {
	v.alg.Render(out)
	for i, s := range out {
		out[i] = v.vca.Scale(s)
	}
}

func (v *Voice) trigger(pitch fp.Fp32, velocity byte) {
	v.alg.Trigger(pitch, velocity)
	v.vca.Trigger()
}

func (v *Voice) retrigger(pitch fp.Fp32) {
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

func note2freq(note byte) fp.Fp32 {
	return fp.Float2Fp32(math.Pow(2, (float64(note)-69.0)/12.0) * 440.0)
}
