package url

import (
	"strings"
	"testing"
)

func TestURLEncode(t *testing.T) {
	urlEnInput := "www.baidu.com?name=tom"
	urlEnExcept := "www.baidu.com%3Fname%3Dtom"
	urlEnActual := URLEncode(urlEnInput)
	if !strings.EqualFold(urlEnExcept, urlEnActual) {
		t.Errorf("URLEncode fail %s to %s", urlEnInput, urlEnActual)
	}

	urlDeInput := "www.baidu.com%3Fname%3Dtom"
	urlDeExcept := "www.baidu.com?name=tom"
	urlDeActual, err := URLDecode(urlDeInput)
	if err != nil {
		t.Error("URLDecode fail")
		return
	}
	if !strings.EqualFold(urlDeExcept, urlDeActual) {
		t.Errorf("URLEncode fail %s to %s", urlDeInput, urlDeActual)
	}

	errDeInput := "!@#!!@"
	errDeActual, err := URLDecode(errDeInput)
	if err != nil {
		t.Error("URLDecode fail")
	}
	if !strings.EqualFold(errDeActual, errDeInput) {
		t.Error("URLDecode fail")
	}

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
