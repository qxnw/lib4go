package html

import (
	"strings"
	"testing"
)

func Test(t *testing.T) {
	htmlEnInput := "<div>"
	htmlEnExcept := "&lt;div&gt;"
	htmlEnActual := HTMLEncode(htmlEnInput)
	if !strings.EqualFold(htmlEnExcept, htmlEnActual) {
		t.Errorf("HTMLEncode fail %s to %s", htmlEnInput, htmlEnActual)
	}

	htmlDeInput := "&lt;div&gt;"
	htmlDeExcept := "<div>"
	htmlDeActual := HTMLDecode(htmlDeInput)
	if !strings.EqualFold(htmlDeExcept, htmlDeActual) {
		t.Errorf("HTMLDecode fail %s to %s", htmlDeExcept, htmlDeActual)
	}
}
