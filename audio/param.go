package audio

type paramId uint16

const (
	C_FREQ paramId = iota
	C_RATIO
	B_FREQ
	B_RATIO
	A_FREQ
	A_RATIO
	AENV_ATTACK
	AENV_DECAY
	AENV_ENDLEVEL
	AENV_LEVEL // A index
	AENV_GATED
	AENV_RETRIGGER
	BENV_ATTACK
	BENV_DECAY
	BENV_ENDLEVEL
	BENV_LEVEL // B1/B2 index
	BENV_GATED
	BENV_RETRIGGER
	VCA_ATTACK
	VCA_DECAY
	VCA_SUSTAIN
	VCA_RELEASE
	VCA_GATED
	VCA_RETRIGGER
)

// a param wraps an fp32 or midi cc val etc
// so the algorithm has a convenient place to reference everything
// and provide upstream voices a hook to set param values
// also the param tree can be saved as a config

type boolparam struct {
	id  paramId
	val bool
}

func (p *boolparam) Value() bool {
	return p.val
}

func newBoolParam(id paramId, defaultValue bool) *boolparam {
	return &boolparam{
		id:  id,
		val: defaultValue,
	}
}

type uint16param struct {
	id   paramId
	val  uint16
	mods []*uint16
}

func (p *uint16param) Value() uint16 {
	return p.val
}

func newUint16Param(id paramId, defaultValue uint16) *uint16param {
	return &uint16param{
		id:   id,
		val:  defaultValue,
		mods: make([]*uint16, 0),
	}
}

type fp32param struct {
	id   paramId
	val  fp32
	mods []*fp32
}

func (p *fp32param) Value() fp32 {
	return p.val
}

func (p *fp32param) Set(v fp32) {
	p.val = v
}

func newFp32Param(id paramId, defaultValue float64) *fp32param {
	return &fp32param{
		id:   id,
		val:  float2fp32(defaultValue),
		mods: make([]*fp32, 0),
	}
}
