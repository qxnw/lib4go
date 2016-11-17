package draw

import (
	"testing"
)

func TestNewDraw(t *testing.T) {
	img := NewDraw(100, 100)
	if img == nil {
		t.Error("test fail")
	}

	img = NewDraw(0, 0)
	if img == nil {
		t.Error("test fail")
	}

	img = NewDraw(-100, 100)
	if img == nil {
		t.Error("test fail")
	}

	img = NewDraw(-100, -100)
	if img == nil {
		t.Error("test fail")
	}

	img = NewDraw(100, -100)
	if img == nil {
		t.Error("test fail")
	}
}

func TestNewDrawFromFile(t *testing.T) {
	path := "/home/champly/picture.png"
	img, err := NewDrawFromFile(100, 100, path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if img == nil {
		t.Error("test fail")
	}

	path = "/home/champly/picture.jpg"
	img, err = NewDrawFromFile(-100, 100, path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if img == nil {
		t.Error("test fail")
	}

	path = "/home/champly/picture.jpg"
	img, err = NewDrawFromFile(-100, -100, path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if img == nil {
		t.Error("test fail")
	}

	path = "/home/champly/picture.jpg"
	img, err = NewDrawFromFile(100, -100, path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if img == nil {
		t.Error("test fail")
	}

	path = "/home/champly/err_picture.jpg"
	img, err = NewDrawFromFile(100, 100, path)
	if err == nil {
		t.Error("test fail")
	}

	path = "/home/champly/picture.gif"
	img, err = NewDrawFromFile(100, 100, path)
	if err == nil {
		t.Error("test fail")
	}
}

func TestDrawFont(t *testing.T) {
	path := "/home/champly/picture.jpg"
	img, err := NewDrawFromFile(1920, 1080, path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if img == nil {
		t.Error("test fail")
	}

	fontPath := "/usr/share/fonts/msyhbd.ttf"
	text := "Hello World"
	col := "155"
	fontSize := 16.0
	img.DrawFont(fontPath, text, col, fontSize, 100, 300)
	err = img.Save("/home/champly/picture_test.png")
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	fontPath = "/usr/share/fonts/msyhbdsdfadf.ttf"
	text = "Hello World"
	col = "155"
	fontSize = 16.0
	img.DrawFont(fontPath, text, col, fontSize, 100, 300)
	err = img.Save("/home/champly/picture_test.png")
	if err != nil {
		t.Errorf("test fail %v", err)
	}
}

func TestDrawImage(t *testing.T) {
	path := "/home/champly/picture_test.png"
	img, err := NewDrawFromFile(1920, 1080, path)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	if img == nil {
		t.Error("test fail")
	}

	err = img.DrawImageWithScale("/home/champly/baidu.png", 100, 100, 200, 200, 100, 100)
	if err != nil {
		t.Errorf("test fail %v", err)
	}
	err = img.Save("/home/champly/picture_test.png")
	if err != nil {
		t.Errorf("test fail %v", err)
	}

	// 路径不正确保存文件
	err = img.Save("/home/champly/picture/picture_test.png")
	if err == nil {
		t.Error("test fail")
	}
}
