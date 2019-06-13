package audio

import (
	"fmt"
	"time"
)

type mixerChannel struct {
	Input chan Sample
	Level *param
	atten float64
}

type Mixer struct {
	Inputs []*mixerChannel
	Output chan<- Sample
}

func (m *Mixer) Input(n int) chan<- Sample {
	return m.Inputs[n].Input
}

func NewMixer(inputs int, output chan<- Sample) *Mixer {
	fmt.Printf("inst'ng mixer %d inputs, output chan: %v\n", inputs, output)
	mixer := &Mixer{
		Inputs: make([]*mixerChannel, inputs),
		Output: output,
	}

	for i, _ := range mixer.Inputs {
		mixer.Inputs[i] = &mixerChannel{
			Input: make(chan Sample, BUFFER_LEN),
			Level: Param(),
			atten: 1.0 / float64(inputs),
		}
		mixer.Inputs[i].Level.value = 1.0
	}

	go func() {
		for {
			var sum float64
			for _, channel := range mixer.Inputs {
				if len(channel.Input) > 0 {
					sample := <-channel.Input
					attenuated := (float64(sample) * channel.atten) * channel.Level.Value()
					sum += attenuated
				}
			}
			mixer.Output <- Sample(sum)
			time.Sleep(10 * time.Nanosecond)
		}
	}()

	return mixer
}
