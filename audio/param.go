package audio

import "fmt"

type paramId uint8

func (p paramId) AsString() string {
	typ := p & 0xC

	switch typ {
	case PATCH_TYPE:
		return fmt.Sprintf("PATCH_%s", p.patchParam())
	case OPR_TYPE:
		return fmt.Sprintf("OPR_%s-%s", p.opr(), p.oprParam())
	case ENV_TYPE:
		return fmt.Sprintf("ENV_%s-%s", p.env(), p.envParam())
	default:
		fmt.Printf("got %x as a param id, with a type of %x\n", p, typ)
	}
	return "undef"
}

func (p paramId) opr() string {
	switch p & 0x03 {
	case GRP_A:
		return "A"
	case GRP_B1:
		return "B1"
	case GRP_B2:
		return "B2"
	case GRP_C:
		return "C"
	}
	return "undef"
}

func (p paramId) oprParam() string {
	switch p & 0xFC {
	case OPR_RATIO:
		return "RATIO"
	case OPR_FEEDBACK:
		return "FEEDBACK"
	}
	return "undef"
}

func (p paramId) env() string {
	switch p & 0x03 {
	case GRP_A:
		return "A"
	case GRP_B:
		return "B"
	case GRP_VCA:
		return "VCA"
	}
	return "undef"
}

func (p paramId) envParam() string {
	switch p & 0xFC {
	case ENV_ATTACK:
		return "ENV_ATTACK"
	case ENV_DECAY:
		return "ENV_DECAY"
	case ENV_SUSTAIN:
		return "ENV_SUSTAIN"
	case ENV_RELEASE:
		return "ENV_RELEASE"
	case ENV_ENDLEVEL:
		return "ENV_ENDLEVEL"
	case ENV_INDEX:
		return "ENV_INDEX"
	case ENV_GATED:
		return "ENV_GATED"
	case ENV_RETRIGGER:
		return "ENV_RETRIGGER"
	default:
		fmt.Printf("got %x as a param id, with a top nybble of %x\n", p, p&0xF0)
	}
	return "undef"
}

func (p paramId) patchParam() string {
	switch p & 0xFC {
	case PATCH_ALGORITHM:
		return "ALGORITHM"
	case PATCH_MIX:
		return "MIX"
	case PATCH_FEEDBACK:
		return "FEEDBACK"
	}
	return "undef"
}

// these constants are combined together to get the unique param id for a particular param,
// for instance envelope B's decay is ENV_DECAY|GRP_B, and operator B2's feedback would be OPR_FEEDBACK|GRP_B2
// note that each operator has a feedback param, but in the digitone scheme only one operator
// gets feedback.  The "patch level" param that the user can tweak is OPR_FEEDBACK|PATCH_TYPE
// and the algorithm select wires it to the appropriate operator
const (
	GRP_A   paramId = 0x0
	GRP_B   paramId = 0x1
	GRP_C   paramId = 0x2
	GRP_D   paramId = 0x3
	GRP_B1  paramId = 0x1 // B1 is an alias for B
	GRP_B2  paramId = 0x3 // B2 is an alias for D (digitone names)
	GRP_VCA paramId = 0x3 // VCA is an alias for D (alg has three envelopes, A, B, and VCA)

	PATCH_TYPE paramId = 0x0 << 2
	OPR_TYPE   paramId = 0x1 << 2
	ENV_TYPE   paramId = 0x2 << 2

	PATCH_ALGORITHM paramId = 0x0<<4 | PATCH_TYPE
	PATCH_FEEDBACK  paramId = 0x1<<4 | PATCH_TYPE
	PATCH_MIX       paramId = 0x2<<4 | PATCH_TYPE

	OPR_RATIO    paramId = 0x0<<4 | OPR_TYPE
	OPR_FEEDBACK paramId = 0x1<<4 | OPR_TYPE

	ENV_ATTACK    paramId = 0x0<<4 | ENV_TYPE
	ENV_DECAY     paramId = 0x1<<4 | ENV_TYPE
	ENV_ENDLEVEL  paramId = 0x2<<4 | ENV_TYPE
	ENV_INDEX     paramId = 0x3<<4 | ENV_TYPE
	ENV_GATED     paramId = 0x4<<4 | ENV_TYPE
	ENV_RETRIGGER paramId = 0x5<<4 | ENV_TYPE
	ENV_SUSTAIN   paramId = 0x6<<4 | ENV_TYPE
	ENV_RELEASE   paramId = 0x7<<4 | ENV_TYPE
)

// a param wraps an fp32 or midi cc val etc
// so the algorithm has a convenient place to reference everything
// and provide upstream voices a hook to set param values
// also the param tree can be saved as a config

type param interface {
	ID() paramId
}

type byteparam struct {
	id  paramId
	val byte
}

func (p *byteparam) ID() paramId {
	return p.id
}

func (p *byteparam) Value() byte {
	return p.val
}

func (p *byteparam) Set(v byte) {
	p.val = v
}

func newByteParam(id paramId, defaultValue byte) *byteparam {
	return &byteparam{
		id:  id,
		val: defaultValue,
	}
}

type boolparam struct {
	id  paramId
	val bool
}

func (p *boolparam) ID() paramId {
	return p.id
}

func (p *boolparam) Value() bool {
	return p.val
}

func (p *boolparam) Set(v bool) {
	p.val = v
}

func newBoolParam(id paramId, defaultValue bool) *boolparam {
	return &boolparam{
		id:  id,
		val: defaultValue,
	}
}

type uint16param struct {
	id  paramId
	val uint16
}

func (p *uint16param) ID() paramId {
	return p.id
}

func (p *uint16param) Value() uint16 {
	return p.val
}

func newUint16Param(id paramId, defaultValue uint16) *uint16param {
	return &uint16param{
		id:  id,
		val: defaultValue,
	}
}

func (p *uint16param) Set(v uint16) {
	p.val = v
}

type fp32param struct {
	id  paramId
	val fp32
}

func (p *fp32param) ID() paramId {
	return p.id
}

func (p *fp32param) Value() fp32 {
	return p.val
}

func (p *fp32param) Set(v fp32) {
	p.val = v
}

func newFp32Param(id paramId, defaultValue float64) *fp32param {
	return &fp32param{
		id:  id,
		val: float2fp32(defaultValue),
	}
}
