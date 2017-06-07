package types

import "testing"
import "github.com/qxnw/lib4go/ut"

func TestDecode1(t *testing.T) {
	ut.Expect(t, DecodeString(1, 2, 3), "")
	ut.Expect(t, DecodeString(2, 2, 3), "3")
	ut.Expect(t, DecodeString(1, 2, 3, 4), "4")
	ut.Expect(t, DecodeString(3, 2, 3, 4), "4")
	ut.Expect(t, DecodeString(3, 2, 3, 3, 2, 4), "2")

}
