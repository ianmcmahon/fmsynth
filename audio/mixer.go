package audio

import (
	"fmt"
	"sync"

	"github.com/ianmcmahon/fmsynth/fp"
)

type mixerChannel struct {
	from  Output
	level fp.Fp32
	atten fp.Fp32
}

type Mixer interface {
	Render(out []fp.Fp32)
}

type levelMixer struct {
	Inputs []*mixerChannel
	wg     sync.WaitGroup
}

func LevelMixer(inputs int) *levelMixer {
	fmt.Printf("inst'ng mixer %d inputs\n", inputs)
	mixer := &levelMixer{
		Inputs: make([]*mixerChannel, inputs),
	}

	for i, _ := range mixer.Inputs {
		mixer.Inputs[i] = &mixerChannel{
			level: fp.Float2Fp32(1.0),
			atten: fp.Float2Fp32(1.0 / float64(inputs)),
		}
	}

	return mixer
}

func (m *levelMixer) Render(out []fp.Fp32) {
	bufs := make([][]fp.Fp32, len(m.Inputs))
	for i, channel := range m.Inputs {
		bufs[i] = make([]fp.Fp32, len(out))
		channel.from.Render(bufs[i])
	}
	for i := range out {
		sum := fp.Fp32(0)
		for c, channel := range m.Inputs {
			attenuated := channel.atten.Mul(bufs[c][i]).Mul(channel.level)
			sum += attenuated
		}
		out[i] = sum
	}
}
