package audio

import (
	"math"
	"time"
)

type ADE struct {
	Gated       bool
	Retrigger   bool
	Attack      time.Duration
	Decay       time.Duration
	EndLevel    float64
	outputParam *param
}

func (a ADE) trig() {
	go func() {
		attackSteps := a.Attack.Nanoseconds() / 1000 / 1000 // step envelopes in 1ms steps
		for i := int64(0); i <= attackSteps; i++ {
			ampVal := float64(i) / float64(attackSteps)
			a.outputParam.value = ampVal
			//fmt.Printf("rising step %d of %d amp: %.2f\n", i, attackSteps, ampVal)
			time.Sleep(1 * time.Millisecond)
		}
		if !a.Gated {
			decaySteps := a.Decay.Nanoseconds() / 1000 / 1000
			for i := int64(0); i <= decaySteps; i++ {
				ampVal := 1.0 - (float64(i) / float64(decaySteps))
				a.outputParam.value = ampVal
				//	fmt.Printf("trig decay step %d of %d amp: %.2f\n", i, decaySteps, ampVal)
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}

func (a ADE) retrig() {
	if a.Retrigger {
		a.trig()
	}
}

func (a ADE) release() {
	if a.Gated {
		go func() {
			decaySteps := a.Decay.Nanoseconds() / 1000 / 1000
			for i := int64(0); i <= decaySteps; i++ {
				ampVal := 1.0 - (float64(i) / float64(decaySteps))
				a.outputParam.value = ampVal
				//	fmt.Printf("gated decay step %d of %d amp: %.2f\n", i, decaySteps, ampVal)
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}
}

type Oscillator struct {
	*Engine
	table    []Sample
	phaseIdx int

	pitch    *param
	amp      *param
	envelope ADE

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
		envelope: ADE{true, true, 100 * time.Millisecond, 200 * time.Millisecond, 0.25, nil},
		output:   output,
	}

	osc.envelope.outputParam = osc.amp

	go osc.oscillate()

	return osc
}

func (o *Oscillator) setPitch(freq float64) {
	o.pitch.value = freq // handle portamento here probably
	o.phaseIdx = 0.0
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
		time.Sleep(1 * time.Nanosecond)
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
