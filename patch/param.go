package patch

import "github.com/ianmcmahon/fmsynth/fp"

type ParamId uint8

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

type Meta struct {
	patch *Patch
	label string
	// ui element enum
	// font/color overrides?
	cc byte

	// domain/range info and conversion from cc val
}

// a param wraps an fp32 or midi cc val etc
// so the algorithm has a convenient place to reference everything
// and provide upstream voices a hook to set param values
// also the param tree can be saved as a config

type Param interface {
	ID() ParamId
	Label() string
	Value() interface{}
	SetFromCC(byte)
}

type byteparam struct {
	id   ParamId
	val  byte
	meta Meta
}

func (p *byteparam) ID() ParamId {
	return p.id
}

func (p *byteparam) Label() string {
	return p.meta.label
}

func (p *byteparam) Value() interface{} {
	return p.val
}

func (p *byteparam) Set(v byte) {
	p.val = v
	p.meta.patch.update(p.id)
}

func (p *byteparam) SetFromCC(v byte) {
	p.Set(v)
}

func NewByteParam(id ParamId, defaultValue byte, meta Meta) *byteparam {
	return &byteparam{
		id:   id,
		val:  defaultValue,
		meta: meta,
	}
}

type boolparam struct {
	id   ParamId
	val  bool
	meta Meta
}

func (p *boolparam) ID() ParamId {
	return p.id
}

func (p *boolparam) Label() string {
	return p.meta.label
}

func (p *boolparam) Value() interface{} {
	return p.val
}

func (p *boolparam) Set(v bool) {
	p.val = v
	p.meta.patch.update(p.id)
}

func (p *boolparam) SetFromCC(v byte) {
	p.Set(v >= 64)
}

func NewBoolParam(id ParamId, defaultValue bool, meta Meta) *boolparam {
	return &boolparam{
		id:   id,
		val:  defaultValue,
		meta: meta,
	}
}

type uint16param struct {
	id   ParamId
	val  uint16
	meta Meta
}

func (p *uint16param) ID() ParamId {
	return p.id
}

func (p *uint16param) Label() string {
	return p.meta.label
}

func (p *uint16param) Value() interface{} {
	return p.val
}

func (p *uint16param) Set(v uint16) {
	p.val = v
	p.meta.patch.update(p.id)
}

func (p *uint16param) SetFromCC(v byte) {
	p.Set(uint16(v) << 9)
}

func NewUint16Param(id ParamId, defaultValue uint16, meta Meta) *uint16param {
	return &uint16param{
		id:   id,
		val:  defaultValue,
		meta: meta,
	}
}

type fp32param struct {
	id   ParamId
	val  fp.Fp32
	meta Meta
}

func (p *fp32param) ID() ParamId {
	return p.id
}

func (p *fp32param) Label() string {
	return p.meta.label
}

func (p *fp32param) Value() interface{} {
	return p.val
}

func (p *fp32param) Set(v fp.Fp32) {
	p.val = v
	p.meta.patch.update(p.id)
}

func (p *fp32param) SetFromCC(v byte) {
	p.Set((fp.Fp32(v) - 64) << 8)
}

func NewFp32Param(id ParamId, defaultValue float64, meta Meta) *fp32param {
	return &fp32param{
		id:   id,
		val:  fp.Float2Fp32(defaultValue),
		meta: meta,
	}
}
