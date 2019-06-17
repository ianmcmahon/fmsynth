package audio

import (
	"fmt"
	"math"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/rakyll/portmidi"
)

const (
	NUM_VOICES    = 8
	BUFFER_LEN    = 128
	SAMPLING_RATE = 44100

	NoteOff         = 0x8
	NoteOn          = 0x9
	PolyAftertouch  = 0xA
	CC              = 0xB
	ProgramChange   = 0xC
	ChannelPressure = 0xD
	PitchBend       = 0xE
)

// anything can be an output if it makes noise
// components get their input data a bufferful at a time
// by calling Render() on the upstream outputs
type Output interface {
	Render(out []fp32)
}

type Engine struct {
	samplingRate int // in samples/sec default 48kHz

	input      Output
	midiEvents <-chan portmidi.Event

	voices   []*Voice
	voiceMap map[byte]*Voice

	// one global patch right now; this will need to be somewhere else once we're multitimbral
	patch *patch

	audioChan chan fp32
}

func NewEngine(midiStream <-chan portmidi.Event) *Engine {
	engine := &Engine{
		samplingRate: SAMPLING_RATE,
		midiEvents:   midiStream,
		voices:       make([]*Voice, NUM_VOICES),
		voiceMap:     make(map[byte]*Voice, 0),
		patch:        initialPatch(),
		audioChan:    make(chan fp32, BUFFER_LEN*2),
	}

	mixer := LevelMixer(NUM_VOICES)
	engine.input = mixer

	for i := range engine.voices {
		engine.voices[i] = engine.NewSimpleVoice(byte(i))
		engine.voices[i].applyPatch(engine.patch)
		mixer.Inputs[i].from = engine.voices[i]
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
		//		fmt.Printf("voice %d is playing %d, we want %d, distance is %d\n", v.id, curNote, note, dist)
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
		case CC:
			num := byte(event.Data1)
			val := byte(event.Data2)
			e.patch.HandleCC(num, val)
		default:
			fmt.Printf("unknown message: %x %x %x\n", event.Status, event.Data1, event.Data2)
		}
	}
}

func (e *Engine) HandleCC(num, val byte) {
	fmt.Printf("CC %x -> %x\n", num, val)

	// arbitrarily mapping the 4 knobs on my kbd to VCA ADSR
	switch num {
	case 0x14:
		e.patch.Uint16Param(ENV_ATTACK | GRP_VCA).Set(uint16(val << 2))
	case 0x15:
		e.patch.Uint16Param(ENV_DECAY | GRP_VCA).Set(uint16(val << 2))
	case 0x16:
		s := fp32(val) << 9
		fmt.Printf("setting sustain to %x -> %.2f\n", val, float64(s)/float64(1<<16))
		e.patch.Fp32Param(ENV_SUSTAIN | GRP_VCA).Set(s)
	case 0x17:
		e.patch.Uint16Param(ENV_RELEASE | GRP_VCA).Set(uint16(val << 2))
	}
}

// TODO: this currently handles only mono 16bit audio
// if we implement stereo effects this will need to be changed
func (e *Engine) runAudio() {
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(e.samplingRate), BUFFER_LEN, e.processAudio)
	if err != nil {
		panic(err)
	}
	if err := stream.Start(); err != nil {
		panic(err)
	}

	renderTime := make([]time.Duration, 100)
	go func() {
		for {
			var sum time.Duration
			for _, d := range renderTime {
				sum = sum + d
			}
			avg := sum / time.Duration(len(renderTime))
			fmt.Printf("average render time: %s\n", avg)
			time.Sleep(5 * time.Second)
		}
	}()

	// audioChan will block when buffer is full
	// when portaudio requests a chunk, processAudio consumes from the channel
	// and this will unblock
	buf := make([]fp32, BUFFER_LEN)
	for {
		start := time.Now()
		e.input.Render(buf)
		elapsed := time.Now().Sub(start)
		renderTime = append(renderTime[:len(renderTime)-1], elapsed)
		for _, s := range buf {
			e.audioChan <- s
		}
	}
}

var underrun bool

func (e *Engine) processAudio(_, out []int16) {
	for i := range out {
		sample := <-e.audioChan
		out[i] = sample.to16bit()
	}
}

func (e *Engine) Stop() {
	fmt.Printf("Engine.Stop() should probably do something\n")
}
