package patch

import (
	"fmt"

	"github.com/ianmcmahon/fmsynth/fp"
)

type ParamId uint8

func (p ParamId) AsString() string {
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

func (p ParamId) opr() string {
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

func (p ParamId) oprParam() string {
	switch p & 0xFC {
	case OPR_RATIO:
		return "RATIO"
	case OPR_FEEDBACK:
		return "FEEDBACK"
	}
	return "undef"
}

func (p ParamId) env() string {
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

func (p ParamId) envParam() string {
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

func (p ParamId) patchParam() string {
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
	GRP_A   ParamId = 0x0
	GRP_B   ParamId = 0x1
	GRP_C   ParamId = 0x2
	GRP_D   ParamId = 0x3
	GRP_B1  ParamId = 0x1 // B1 is an alias for B
	GRP_B2  ParamId = 0x3 // B2 is an alias for D (digitone names)
	GRP_VCA ParamId = 0x3 // VCA is an alias for D (alg has three envelopes, A, B, and VCA)

	PATCH_TYPE ParamId = 0x0 << 2
	OPR_TYPE   ParamId = 0x1 << 2
	ENV_TYPE   ParamId = 0x2 << 2

	PATCH_ALGORITHM ParamId = 0x0<<4 | PATCH_TYPE
	PATCH_FEEDBACK  ParamId = 0x1<<4 | PATCH_TYPE
	PATCH_MIX       ParamId = 0x2<<4 | PATCH_TYPE

	OPR_RATIO    ParamId = 0x0<<4 | OPR_TYPE
	OPR_FEEDBACK ParamId = 0x1<<4 | OPR_TYPE

	ENV_ATTACK    ParamId = 0x0<<4 | ENV_TYPE
	ENV_DECAY     ParamId = 0x1<<4 | ENV_TYPE
	ENV_ENDLEVEL  ParamId = 0x2<<4 | ENV_TYPE
	ENV_INDEX     ParamId = 0x3<<4 | ENV_TYPE
	ENV_GATED     ParamId = 0x4<<4 | ENV_TYPE
	ENV_RETRIGGER ParamId = 0x5<<4 | ENV_TYPE
	ENV_SUSTAIN   ParamId = 0x6<<4 | ENV_TYPE
	ENV_RELEASE   ParamId = 0x7<<4 | ENV_TYPE
)

type ccMeta struct {
	ccNum   byte
	convert func(byte) interface{}
}

// a param wraps an fp32 or midi cc val etc
// so the algorithm has a convenient place to reference everything
// and provide upstream voices a hook to set param values
// also the param tree can be saved as a config

type Param interface {
	ID() ParamId
	CC() ccMeta
	Value() interface{}
}

type byteparam struct {
	id  ParamId
	val byte
	cc  ccMeta
}

func (p *byteparam) ID() ParamId {
	return p.id
}

func (p *byteparam) Value() interface{} {
	return p.val
}

func (p *byteparam) Set(v byte) {
	p.val = v
}

func (p *byteparam) CC() ccMeta {
	return p.cc
}

func NewByteParam(id ParamId, defaultValue byte, cc ccMeta) *byteparam {
	return &byteparam{
		id:  id,
		val: defaultValue,
		cc:  cc,
	}
}

type boolparam struct {
	id  ParamId
	val bool
	cc  ccMeta
}

func (p *boolparam) ID() ParamId {
	return p.id
}

func (p *boolparam) Value() interface{} {
	return p.val
}

func (p *boolparam) Set(v bool) {
	p.val = v
}

func (p *boolparam) CC() ccMeta {
	return p.cc
}

func NewBoolParam(id ParamId, defaultValue bool) *boolparam {
	return &boolparam{
		id:  id,
		val: defaultValue,
	}
}

type uint16param struct {
	id  ParamId
	val uint16
	cc  ccMeta
}

func (p *uint16param) ID() ParamId {
	return p.id
}

func (p *uint16param) Value() interface{} {
	return p.val
}

func (p *uint16param) Set(v uint16) {
	p.val = v
}

func (p *uint16param) CC() ccMeta {
	return p.cc
}

func NewUint16Param(id ParamId, defaultValue uint16) *uint16param {
	return &uint16param{
		id:  id,
		val: defaultValue,
	}
}

type fp32param struct {
	id  ParamId
	val fp.Fp32
	cc  ccMeta
}

func (p *fp32param) ID() ParamId {
	return p.id
}

func (p *fp32param) Value() interface{} {
	return p.val
}

func (p *fp32param) Set(v fp.Fp32) {
	p.val = v
}

func (p *fp32param) CC() ccMeta {
	return p.cc
}

func NewFp32Param(id ParamId, defaultValue float64) *fp32param {
	return &fp32param{
		id:  id,
		val: fp.Float2Fp32(defaultValue),
	}
}
