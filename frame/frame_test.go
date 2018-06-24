package frame

import (
	"testing"
	"time"
)

func TestFrameModify(t *testing.T) {
	f := New(8, 32)
	t.Logf("Before:\n%v", f)
	modifier := [][]byte{
		{0x81, 0x8f},
		{0x7e, 0x8f},
	}
	f.Modify(5, 18, modifier)
	t.Logf("After:\n%v", f)
}

func TestMovingBar(t *testing.T) {
	f := New(8, 32)
	mod := [][]byte{
		{0x40}, {0x20}, {0x10}, {0x08}, {0x08}, {0x10}, {0x20}, {0x40},
	}

	f.data[0][0] = 0x80
	f.data[1][0] = 0x40
	f.data[2][0] = 0x20
	f.data[3][0] = 0x10
	f.data[4][0] = 0x10
	f.data[5][0] = 0x20
	f.data[6][0] = 0x40
	f.data[7][0] = 0x80

	t.Logf("Init:\n%v", f)
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 28; i++ {
		f.Modify(i, i+5, mod)
		t.Logf("Step #%d:\n%v", i+1, f)
	}
}

func TestFrameText(t *testing.T) {
	f := New(8, 32)
	f.Text(0, "11 5m")
	t.Logf("11 5m:\n%v", f)
}

func TestSubFrame(t *testing.T) {
	f := FromText("test string", 32)
	t.Logf("Original:\n%v", f)

	for i := 0; i < f.Width(); i++ {
		sub := f.SubFrame(i, i+32)
		t.Logf("Pos %d:\n%v", i, sub)
	}
}
