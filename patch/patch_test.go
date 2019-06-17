package patch

import (
	"fmt"
	"math"
	"testing"
)

func TestCCByteRange(t *testing.T) {
	fullRange := byteRange(0, 127)
	halfRange := byteRange(0, 63)
	doubleRange := byteRange(0, 255)
	algorithmRange := byteRange(0, 7)
	oddRange := byteRange(5, 15)

	var cc byte
	for cc = 0; cc <= 127; cc++ {
		assertEqual(t, fullRange(cc), cc, "")
		assertEqual(t, halfRange(cc), cc>>1, "")
		assertEqual(t, doubleRange(cc), cc<<1, "")
		assertEqual(t, algorithmRange(cc), cc>>4, "")
	}
	assertEqual(t, oddRange(0), byte(5), "")
	assertEqual(t, oddRange(60), byte(10), "")
	assertEqual(t, oddRange(90), byte(12), "")
	assertEqual(t, oddRange(127), byte(15), "")
}

func TestUint16ByteRange(t *testing.T) {
	fullRange := uint16Range(0, math.MaxUint16)
	oddRange := uint16Range(1000, 2000)

	var cc byte
	for cc = 0; cc <= 127; cc++ {
		assertEqual(t, fullRange(cc), uint16(cc)<<9, "")
	}

	assertEqual(t, fullRange(0), uint16(0), "")
	assertEqual(t, fullRange(64), uint16(32768), "")

	assertEqual(t, oddRange(0), uint16(1000), "")
	assertEqual(t, oddRange(60), uint16(1469), "")
	assertEqual(t, oddRange(90), uint16(1703), "")
	assertEqual(t, oddRange(127), uint16(1993), "")

}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}
