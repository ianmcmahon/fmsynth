package patch

import (
	"fmt"
	"math"
)

/*
	A patch holds all the parameters and mod matrix info for a track
	A track can use as many voices as are available, and when a voice
	is assigned to a track, the patch is applied to the voice
	which 'wires' all the parameter pointers to the patch's concrete vals

	In the digitone, 128 patches are stored in the sound pool
	they can be loaded from flash storage, but the sound pool is part of
	the project.
	patches from the sound pool can be applied to the track per-trig
*/

type Patch struct {
	params map[ParamId]Param
	byCC   map[byte]Param
}

func (p *Patch) HandleCC(num, val byte) {

}

func InitialPatch() *Patch {
	p := &Patch{
		params: make(map[ParamId]Param, 0),
		byCC:   make(map[byte]Param, 0),
	}

	p.addByte(PATCH_ALGORITHM, 0, ccMeta{3, byteRange(0, 7)})
	p.addFp32(PATCH_FEEDBACK, 0.0)
	p.addFp32(PATCH_MIX, 0.5)

	p.addFp32(OPR_RATIO|GRP_A, 1.0)
	p.addFp32(OPR_RATIO|GRP_B1, 1.0)
	p.addFp32(OPR_RATIO|GRP_B2, 1.0)
	p.addFp32(OPR_RATIO|GRP_C, 1.0)

	p.addBool(ENV_GATED|GRP_A, true)
	p.addBool(ENV_RETRIGGER|GRP_A, true)
	p.addUint16(ENV_ATTACK|GRP_A, 0)
	p.addUint16(ENV_DECAY|GRP_A, 0)
	p.addFp32(ENV_ENDLEVEL|GRP_A, 0.0)
	p.addFp32(ENV_INDEX|GRP_A, 1.0)

	p.addBool(ENV_GATED|GRP_B, true)
	p.addBool(ENV_RETRIGGER|GRP_B, true)
	p.addUint16(ENV_ATTACK|GRP_B, 0)
	p.addUint16(ENV_DECAY|GRP_B, 0)
	p.addFp32(ENV_ENDLEVEL|GRP_B, 0.0)
	p.addFp32(ENV_INDEX|GRP_B, 1.0)

	p.addBool(ENV_GATED|GRP_VCA, true)
	p.addBool(ENV_RETRIGGER|GRP_VCA, false)
	p.addUint16(ENV_ATTACK|GRP_VCA, 0)
	p.addUint16(ENV_DECAY|GRP_VCA, 0)
	p.addUint16(ENV_RELEASE|GRP_VCA, 0)
	p.addFp32(ENV_SUSTAIN|GRP_VCA, 1.0)

	return p
}

func (p *Patch) ByteParam(id ParamId) *byteparam {
	if v, ok := p.params[id].(*byteparam); ok {
		return v
	}
	fmt.Printf("%s is a %T, expected byte\n", id.AsString(), p.params[id])
	return nil
}

func (p *Patch) BoolParam(id ParamId) *boolparam {
	if v, ok := p.params[id].(*boolparam); ok {
		return v
	}
	fmt.Printf("%s is a %T, expected bool\n", id.AsString(), p.params[id])
	return nil
}

func (p *Patch) Uint16Param(id ParamId) *uint16param {
	if v, ok := p.params[id].(*uint16param); ok {
		return v
	}
	fmt.Printf("%s is a %T, expected uint16\n", id.AsString(), p.params[id])
	return nil
}

func (p *Patch) Fp32Param(id ParamId) *fp32param {
	if v, ok := p.params[id].(*fp32param); ok {
		return v
	}
	panic(fmt.Errorf("%s is a %T, expected fp32param\n", id.AsString(), p.params[id]))
}

func (p *Patch) addByte(id ParamId, v byte, cc ccMeta) {
	p.params[id] = NewByteParam(id, v, cc)
	p.byCC[cc.ccNum] = p.params[id]
}

func (p *Patch) addBool(id ParamId, v bool) {
	p.params[id] = NewBoolParam(id, v)
	//p.byCC[cc.ccNum] = p.params[id]
}

func (p *Patch) addUint16(id ParamId, v uint16) {
	p.params[id] = NewUint16Param(id, v)
	//p.byCC[cc.ccNum] = p.params[id]
}

func (p *Patch) addFp32(id ParamId, v float64) {
	p.params[id] = NewFp32Param(id, v)
	//p.byCC[cc.ccNum] = &(p.params[id])
}

func byteRange(min, max byte) func(byte) interface{} {
	// this converter returns bytes, so range is limited to 0-255 (2x upscaled)
	if min < 0 {
		min = 0
	}
	if max > math.MaxUint8 {
		max = math.MaxUint8
	}
	// short circuit obvious stuff for performance
	if min == 0 && max == 127 {
		return func(b byte) interface{} {
			return b
		}
	}
	if min == 0 && max == 255 {
		return func(b byte) interface{} {
			return b << 1
		}
	}

	return func(b byte) interface{} {
		return byte(((uint16(b) * uint16(max+1-min)) >> 7) + uint16(min))
	}
}

func uint16Range(min, max uint16) func(byte) interface{} {
	if min == 0 && max == math.MaxUint16 {
		return func(b byte) interface{} {
			return uint16(b) << 9
		}
	}
	return func(b byte) interface{} {
		return uint16(((uint32(b) * uint32(max+1-min)) >> 7) + uint32(min))
	}
}
