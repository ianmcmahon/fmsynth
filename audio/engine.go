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
	BUFFER_LEN = 16
)

type Engine struct {
	samplingRate int // in samples/sec default 48kHz
	period       time.Duration
	ticker       *time.Ticker

	clockedComponents []Clocked

	input        chan Sample
	outputBuffer []int16
	outBufIdx    int
	outputIdx    int

	samples int64
}

func NewEngine() *Engine {
	engine := &Engine{
		samplingRate:      16000,
		clockedComponents: make([]Clocked, 0),
		input:             make(chan Sample, 0),
		outputBuffer:      make([]int16, BUFFER_LEN),
		outBufIdx:         0,
		outputIdx:         0,
	}
	periodSec := 1.0 / float32(engine.samplingRate)
	engine.period = time.Duration(int(periodSec*1e9)) * time.Nanosecond
	fmt.Printf("period: %v\n", engine.period)
	engine.ticker = time.NewTicker(engine.period)

	// wiring up some sample stuff here; this needs to be done better
	engine.NewOscillator(263, engine.input)

	go engine.Run()

	return engine
}

func (e *Engine) Run() {
	go e.runAudio()
	go e.runClock()
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

	for s := range e.input {
		e.outputBuffer[e.outBufIdx] = s.As16bit()
		e.outBufIdx = (e.outBufIdx + 1) % BUFFER_LEN
	}
}

func (e *Engine) processAudio(_, out []int16) {
	for i := range out {
		out[i] = e.outputBuffer[(e.outputIdx+i)%BUFFER_LEN]
	}
}

func (e *Engine) runClock() {
	e.samples = 0
	for t := range e.ticker.C {
		for _, c := range e.clockedComponents {
			c.Tick(t)
		}
		e.samples++
	}
}

func (e *Engine) Stop() {
	fmt.Printf("ticked %d samples\n", e.samples)
	fmt.Printf("Engine.Stop() should probably do something\n")
}

func (e *Engine) addClock(component Clocked) {
	e.clockedComponents = append(e.clockedComponents, component)
}
