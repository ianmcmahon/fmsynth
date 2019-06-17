package audio

import "testing"

func testPatch() *patch {
	patch := &patch{
		params: make(map[paramId]param, 0),
	}

	patch.addFp32(OPR_RATIO|GRP_A, 1.0)

	return patch
}

func TestRotate(t *testing.T) {
	// how the hell do I test this lol
	//op := Operator(GRP_A)
}
