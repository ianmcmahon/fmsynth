package audio

import (
	"fmt"
	"math"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/rakyll/portmidi"
)

/*
	Engine has a ticker that is the main sample clock
	all audio components instantiated from Engine
	on tick, generative components generate their sample for current t based
	on state of inputs/parameters (not sure if timing issues will bite me here)
*/

type Clocked interface {
	Tick(t time.Time)
}

const (
	NUM_VOICES = 8
	BUFFER_LEN = 512

	NoteOff         = 0x8
	NoteOn          = 0x9
	PolyAftertouch  = 0xA
	CC              = 0xB
	ProgramChange   = 0xC
	ChannelPressure = 0xD
	PitchBend       = 0xE
)

type Engine struct {
	samplingRate int // in samples/sec default 48kHz

	input      chan Sample
	midiEvents <-chan portmidi.Event

	voices   []*Voice
	voiceMap map[byte]*Voice
}

func NewEngine(midiStream <-chan portmidi.Event) *Engine {
	engine := &Engine{
		samplingRate: 44100,
		input:        make(chan Sample, BUFFER_LEN),
		midiEvents:   midiStream,
		voices:       make([]*Voice, NUM_VOICES),
		voiceMap:     make(map[byte]*Voice, 0),
	}

	mixer := NewMixer(NUM_VOICES, engine.input)

	for i := range engine.voices {
		engine.voices[i] = engine.NewSimpleVoice(i, mixer.Input(i))
	}

	go engine.Run()

	return engine
}

func (e *Engine) Run() {
	go e.runAudio()
	go e.handleMidi()
}

func (e *Engine) getVoice(note byte) *Voice {
	best := 127
	var bestV *Voice
	for _, v := range e.voices {
		curNote := v.CurNote()
		if curNote == 0 {
			return v
		}
		dist := int(math.Abs(float64(curNote) - float64(note)))
		fmt.Printf("voice %d is playing %d, we want %d, distance is %d\n", v.id, curNote, note, dist)
		if dist < best {
			best = dist
			bestV = v
		}
	}
	return bestV
}

func (e *Engine) handleMidi() {
	for event := range e.midiEvents {
		switch event.Status >> 4 {
		case NoteOn:
			note := byte(event.Data1)
			vel := byte(event.Data2)
			voice := e.getVoice(note)
			if voice != nil {
				e.voiceMap[note] = voice
				voice.NoteOn(note, vel)
			} else {
				fmt.Printf("nil voice\n")
			}
		case NoteOff:
			note := byte(event.Data1)
			if voice, ok := e.voiceMap[note]; ok {
				delete(e.voiceMap, note)
				voice.NoteOff(note)
			}
		}
	}
}

// TODO: this currently handles only mono 16bit audio
// if we implement stereo effects this will need to be changed
func (e *Engine) runAudio() {
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(e.samplingRate), BUFFER_LEN/2, e.processAudio)
	if err != nil {
		panic(err)
	}
	if err := stream.Start(); err != nil {
		panic(err)
	}
}

func (e *Engine) processAudio(_, out []int16) {
	for i := range out {
		if len(e.input) == 0 {
			fmt.Printf("input buffer underrun!\n")
			return
		}
		sample := <-e.input
		out[i] = sample.As16bit()
	}
}

func (e *Engine) Stop() {
	fmt.Printf("Engine.Stop() should probably do something\n")
}
