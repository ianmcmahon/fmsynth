package audio

import "math"

// 1.0 corresponds to full scale output, ie 32767, -1.0 is -32768
type Sample float32

func (s Sample) As16bit() int16 {
	return int16(math.Round(float64(s * Sample(math.MaxInt16))))
}
