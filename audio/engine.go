package audio

import (
	"fmt"
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

	voice *Voice
}

func NewEngine(midiStream <-chan portmidi.Event) *Engine {
	engine := &Engine{
		samplingRate: 44100,
		input:        make(chan Sample, BUFFER_LEN),
		midiEvents:   midiStream,
	}

	engine.voice = engine.NewSimpleVoice(engine.input)

	go engine.Run()

	return engine
}

func (e *Engine) Run() {
	go e.runAudio()
	go e.handleMidi()
}

func (e *Engine) handleMidi() {
	for event := range e.midiEvents {
		switch event.Status >> 4 {
		case NoteOn:
			e.voice.NoteOn(byte(event.Data1), byte(event.Data2))
		case NoteOff:
			e.voice.NoteOff(byte(event.Data1))
		}
		time.Sleep(20 * time.Nanosecond)
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
