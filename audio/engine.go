package audio

import (
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
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
)

type Engine struct {
	samplingRate int // in samples/sec default 48kHz

	input        chan Sample
	outputBuffer []int16
	outBufIdx    int
	outputIdx    int
}

func NewEngine() *Engine {
	engine := &Engine{
		samplingRate: 44100,
		input:        make(chan Sample, BUFFER_LEN),
	}

	// wiring up some sample stuff here; this needs to be done better
	outputMixer := NewMixer(engine.input)
	carrier := engine.NewOscillator(440, outputMixer.A)
	carrier.Output.Bias = 0.2
	carrier.Output.Depth = 0.2
	carrier.FM.Bias = 0.2
	carrier.FM.Depth = 0.5

	modulator := engine.NewOscillator(440, carrier.FM.CV)
	modulator.Output.Bias = 1.0
	modulator.Output.Depth = 1.0
	modulator.FM.Bias = 0.5
	modulator.FM.Depth = 0.5

	modulator2 := engine.NewOscillator(440*3, modulator.FM.CV)
	modulator2.Output.Bias = 1.0
	modulator2.Output.Depth = 1.0

	lfo := engine.NewOscillator(10, modulator2.Output.CV)
	lfo.Output.Bias = 1.0

	lfo2 := engine.NewOscillator(3, modulator.Output.CV)
	lfo2.Output.Bias = 1.0

	acarrier := engine.NewOscillator(660, outputMixer.B)
	acarrier.Output.Bias = 0.2
	acarrier.Output.Depth = 0.2
	acarrier.FM.Bias = 0.2
	acarrier.FM.Depth = 0.5

	amodulator := engine.NewOscillator(660, acarrier.FM.CV)
	amodulator.Output.Bias = 1.0
	amodulator.Output.Depth = 1.0
	amodulator.FM.Bias = 0.5
	amodulator.FM.Depth = 0.5

	amodulator2 := engine.NewOscillator(660*3, amodulator.FM.CV)
	amodulator2.Output.Bias = 1.0
	amodulator2.Output.Depth = 1.0

	alfo := engine.NewOscillator(10, amodulator2.Output.CV)
	alfo.Output.Bias = 1.0

	alfo2 := engine.NewOscillator(3, amodulator.Output.CV)
	alfo2.Output.Bias = 1.0

	go engine.Run()

	return engine
}

func (e *Engine) Run() {
	go e.runAudio()
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
		sample := <-e.input
		out[i] = sample.As16bit()
	}
}

func (e *Engine) Stop() {
	fmt.Printf("Engine.Stop() should probably do something\n")
}
