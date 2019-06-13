package audio

import (
	"fmt"
	"math"
)

type AttenuvertedInput struct {
	CV      chan Sample
	DepthCV chan Sample
	Bias    Sample
	Depth   Sample
}

func Attenuverter() *AttenuvertedInput {
	a := &AttenuvertedInput{
		CV:    make(chan Sample, BUFFER_LEN),
		Bias:  0.0,
		Depth: 0.0,
	}
	return a
}

func (a *AttenuvertedInput) Value() Sample {
	cv := Sample(0.0)
	depthCv := Sample(0.0)
	if len(a.CV) > 0 {
		cv = <-a.CV
	}
	if len(a.DepthCV) > 0 {
		depthCv = <-a.DepthCV
	}
	return a.Bias + (a.Depth+depthCv)*cv
}

type Oscillator struct {
	*Engine
	table    []Sample
	freq     float64
	phaseIdx int
	log      bool

	// controls
	Output *AttenuvertedInput
	FM     *AttenuvertedInput

	output chan<- Sample
}

func (engine *Engine) NewOscillator(initialFreq float64, output chan<- Sample) *Oscillator {
	sineTable = makeSineTable(engine.samplingRate)
	osc := &Oscillator{
		Engine:   engine,
		table:    sineTable,
		freq:     initialFreq,
		phaseIdx: 0,
		Output:   Attenuverter(),
		FM:       Attenuverter(),
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
		phaseIncr := int(math.Round(o.freq + o.freq*float64(o.FM.Value())))
		o.phaseIdx = (o.phaseIdx + phaseIncr) % o.samplingRate

		// amplitude is the sum of the level "knob" setting and the level CV input times the depth
		amp := o.Output.Value()
		o.output <- amp * o.table[o.phaseIdx]
	}
}

var (
	sineTable []Sample
)

func init() {
	fmt.Printf("Oscillator init(): setting up wavetables\n")

}

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
