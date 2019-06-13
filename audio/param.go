package audio

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
