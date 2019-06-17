package fp

// samples, CVs, freqs etc are represented with a 32 bit fixed point number
// most significant 16 bits are the integer and least sig are mantissa

// full scale amplitude is 32767/-32768 in 16bit
// in Fp32 it's 1.0 ie 1<<16-1
// any time n samples are added they must be attenuated by 1/n
// to avoid overflow (clipping) in the final stage

type Fp32 int32

// since 1<<16 is fullscale, I'm actually holding 17 bits
// of amplitude precision
func (a Fp32) To16bit() int16 {
	// hard limiter
	if a > 1<<16 {
		a = 1 << 16
	}
	if a < -1<<16 {
		a = -1 << 16
	}
	// now that we're clamped within the limit, discard the extra bit
	// and cast down
	return int16(a >> 1)
}

func (a Fp32) Mul(b Fp32) Fp32 {
	return Fp32((int64(a) * int64(b)) >> 16)
}

func Float2Fp32(f float64) Fp32 {
	return Fp32(f * float64(1<<16))
}
