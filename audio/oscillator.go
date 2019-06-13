package audio

import (
	"math"
	"time"
)

type param struct {
	value   float64
	cv      chan float64
	depth   float64
	depthcv chan float64
}

func Param() *param {
	return &param{
		cv:      make(chan float64, BUFFER_LEN),
		depthcv: make(chan float64, BUFFER_LEN),
	}
}

func (p *param) Value() float64 {
	cv := 0.0
	depthcv := 0.0
	if len(p.cv) > 0 {
		cv = <-p.cv
	}
	if len(p.depthcv) > 0 {
		depthcv = <-p.depthcv
	}
	return p.value + (p.depth+depthcv)*cv
}

type Oscillator struct {
	*Engine
	table    []Sample
	phaseIdx int

	pitch *param
	amp   *param

	output chan<- Sample
}

func (engine *Engine) NewOscillator(output chan<- Sample) *Oscillator {
	sineTable = makeSineTable(engine.samplingRate)
	osc := &Oscillator{
		Engine:   engine,
		table:    sineTable,
		phaseIdx: 0,
		pitch:    Param(),
		amp:      Param(),
		output:   output,
	}

	go osc.oscillate()

	return osc
}

func (o *Oscillator) oscillate() {
	// just start generating samples based on the sampling rate
	// and trust in the channel buffering to maintain our rate

	// the state is the phase angle
	// phase angle is evenly rotated across the table for one cycle

	// to get a 440Hz sine wave, need to rotate through this table
	// 440 times in SAMPLING_RATE cycles
	for {
		// this should be freq * (table_len / sampling_rate) I believe
		// but since our table is 1Hz at sampling_rate, len/rate = 1
		phaseIncr := int(math.Round(o.pitch.Value()))
		o.phaseIdx = (o.phaseIdx + phaseIncr) % o.samplingRate
		o.output <- Sample(o.amp.Value()) * o.table[o.phaseIdx]
		time.Sleep(20 * time.Nanosecond)
	}
}

var (
	sineTable []Sample
)

// TODO: make this cache generated tables
func makeSineTable(tableLength int) []Sample {
	table := make([]Sample, tableLength)
	for i := range table {
		phase := float64(i) / float64(tableLength)
		table[i] = Sample(math.Sin(2 * math.Pi * phase))
		// this is a 1Hz wave across the table,
		// so a complete rotation of phase is one cycle thru table
	}
	return table
}
