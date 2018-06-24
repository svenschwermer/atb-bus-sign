package frame_test

import (
	"fmt"
	"testing"
	"time"

	"svenschwermer.de/bus-sign/frame"
)

func TestFrameModify(t *testing.T) {
	f := frame.New(8, 32)
	t.Logf("Before:\n%s", drawFrame(f))
	modifier := [][]byte{
		{0x81, 0x8f},
		{0x7e, 0x8f},
	}
	f.Modify(5, 18, modifier)
	t.Logf("After:\n%s", drawFrame(f))
}

func TestMovingBar(t *testing.T) {
	f := frame.New(8, 32)
	mod := [][]byte{
		{0x40}, {0x20}, {0x10}, {0x08}, {0x08}, {0x10}, {0x20}, {0x40},
	}

	f[0][0] = 0x80
	f[1][0] = 0x40
	f[2][0] = 0x20
	f[3][0] = 0x10
	f[4][0] = 0x10
	f[5][0] = 0x20
	f[6][0] = 0x40
	f[7][0] = 0x80

	t.Logf("Init:\n%s", drawFrame(f))
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 28; i++ {
		f.Modify(i, i+5, mod)
		t.Logf("Step #%d:\n%s", i+1, drawFrame(f))
	}
}

func TestFrameText(t *testing.T) {
	f := frame.New(8, 32)
	f.Text(0, "11 5m")
	t.Logf("11 5m:\n%s", drawFrame(f))
}

func drawFrame(f frame.Frame) (s string) {
	for i := 0; i < len(f[0])*8; i++ {
		s += fmt.Sprintf("%2d ", i)
	}
	s += "\n"
	for _, line := range f {
		for _, pixels := range line {
			for i := 7; i >= 0; i-- {
				if pixels&(1<<uint(i)) != 0x00 {
					s += " # "
				} else {
					s += " . "
				}
			}
		}
		s += "\n"
	}
	return
}
