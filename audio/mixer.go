package audio

import (
	"fmt"
	"sync"
)

type mixerChannel struct {
	from  Output
	level fp32
	atten fp32
}

type Mixer struct {
	Inputs []*mixerChannel
	wg     sync.WaitGroup
}

func NewMixer(inputs int) *Mixer {
	fmt.Printf("inst'ng mixer %d inputs\n", inputs)
	mixer := &Mixer{
		Inputs: make([]*mixerChannel, inputs),
	}

	for i, _ := range mixer.Inputs {
		mixer.Inputs[i] = &mixerChannel{
			level: float2fp32(1.0),
			atten: float2fp32(1.0 / float64(inputs)),
		}
	}

	return mixer
}

func (m *Mixer) Render(out []fp32) {
	bufs := make([][]fp32, len(m.Inputs))
	for i, channel := range m.Inputs {
		bufs[i] = make([]fp32, len(out))
		channel.from.Render(bufs[i])
	}
	for i := range out {
		sum := fp32(0)
		for c, channel := range m.Inputs {
			attenuated := channel.atten.mul(bufs[c][i]).mul(channel.level)
			sum += attenuated
		}
		out[i] = sum
	}
}
