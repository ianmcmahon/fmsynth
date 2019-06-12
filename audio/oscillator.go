package audio

import (
	"fmt"
	"math"
	"time"
)

type Oscillator struct {
	*Engine
	table    []Sample
	freq     float64
	phaseIdx int

	output chan<- Sample
}

func (engine *Engine) NewOscillator(initialFreq float64, output chan<- Sample) *Oscillator {
	osc := &Oscillator{
		Engine:   engine,
		table:    sineTable,
		freq:     initialFreq,
		phaseIdx: 0,
		output:   output,
	}

	engine.addClock(osc)

	return osc
}

func (o *Oscillator) Tick(t time.Time) {
	phaseIncr := int(math.Round(o.freq / float64(o.samplingRate) * float64(len(o.table))))
	o.phaseIdx = (o.phaseIdx + phaseIncr) % len(o.table)

	sample := o.table[o.phaseIdx]
	o.output <- sample
}

var (
	sineTable []Sample
)

func init() {
	fmt.Printf("Oscillator init(): setting up wavetables\n")

	sineTable = makeSineTable(1024)
	fmt.Printf("sinTable:\n%v\n", sineTable)
}

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
