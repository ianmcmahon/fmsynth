package audio

import "time"

type Mixer struct {
	A      chan Sample
	B      chan Sample
	Output chan<- Sample
}

func NewMixer(output chan<- Sample) *Mixer {
	mixer := &Mixer{
		A:      make(chan Sample, BUFFER_LEN),
		B:      make(chan Sample, BUFFER_LEN),
		Output: output,
	}

	go func() {
		for {
			var a Sample
			var b Sample
			var do bool
			if len(mixer.A) > 0 {
				a = <-mixer.A
				do = true
			}
			if len(mixer.B) > 0 {
				b = <-mixer.B
				do = true
			}
			if do {
				mixer.Output <- (a * 0.9) + (b * 0.1)
			}
			time.Sleep(22 * time.Nanosecond)
		}
	}()

	return mixer
}
