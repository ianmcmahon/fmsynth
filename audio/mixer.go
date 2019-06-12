package audio

type Mixer struct {
	input  <-chan Sample
	output chan<- Sample
}

func NewMixer(output chan<- Sample) *Mixer {
	mixer := &Mixer{
		input:  make(chan Sample, 0),
		output: output,
	}

	return mixer
}
