package frame

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

// Frame is a collection of lines
// 1st index: line index
type Frame struct {
	data  [][]byte
	width int
}

func New(lines, columns int) *Frame {
	data := make([][]byte, lines)
	for i := range data {
		data[i] = make([]byte, (columns+7)/8)
	}
	return &Frame{data, columns}
}

func FromText(str string, minimumWidth int) *Frame {
	width := utf8.RuneCountInString(str) * font5x8.width
	width = max(width, minimumWidth)
	f := New(font5x8.height, width)
	f.Text(0, str)
	return f
}

func (f *Frame) Width() int {
	return f.width
}

func (f *Frame) ConcatenateLines() []byte {
	return bytes.Join(f.data, []byte{})
}

// Modify sets the columns [start, end) with the data provided in lines
// In "non-full" bytes, the lowermost bits are considered padding
func (f *Frame) Modify(start, end int, lines [][]byte) error {
	for i, srcLine := range lines {
		srcPos := 0
		for destPos := start; destPos < end; {
			// every iteration covers the range of bits between two byte-borders
			// in both source and destination line arrays
			srcBytePos := srcPos % 8 // counting from MSB to LSB
			srcByteIdx := srcPos / 8
			srcBits := min(end-destPos, 8-srcBytePos)
			srcByte := srcLine[srcByteIdx]

			destBytePos := destPos % 8
			destByteIdx := destPos / 8
			destBits := 8 - destBytePos
			destByte := f.data[i][destByteIdx]

			bits := min(srcBits, destBits)
			mask := byte((1 << uint(bits)) - 1)
			destShift := uint(8 - destBytePos - bits)
			srcShift := uint(8 - srcBytePos - bits)

			srcByte = (srcByte >> srcShift) & mask
			destByte &^= mask << destShift   // clear bits in dest byte
			destByte |= srcByte << destShift // apply bits from src byte
			f.data[i][destByteIdx] = destByte

			srcPos += bits
			destPos += bits
		}
	}
	return nil
}

func (f *Frame) Text(pos int, str string) error {
	for _, c := range str {
		f.Modify(pos, pos+font5x8.width, font5x8.Get(c))
		pos += font5x8.width
	}
	return nil
}

func (f *Frame) SubFrame(start, end int) *Frame {
	width := end - start
	out := f.Duplicate()
	for i := range f.data {
		shift := uint(start)
		out.data[i] = out.data[i][shift/8 : len(out.data[i])] // drop whole leading bytes
		shift %= 8
		if shift > 0 {
			mask := byte((1 << shift) - 1)
			for j, b := range out.data[i] {
				out.data[i][j] = b << shift
				if j > 0 {
					overflow := (b >> (8 - shift)) & mask
					out.data[i][j-1] |= overflow
				}
			}
		}
		out.data[i] = out.data[i][0 : (width+7)/8] // drop whole trailing bytes
	}
	out.width = width
	return out
}

func (f *Frame) Duplicate() *Frame {
	out := New(len(f.data), 2*f.width)
	for i := range f.data {
		l := len(f.data[i])
		copy(out.data[i][0:l], f.data[i][0:l])
		out.Modify(f.width, 2*f.width, f.data)
	}
	return out
}

func (f *Frame) String() (s string) {
	for i := 0; i < f.width; i++ {
		s += fmt.Sprintf("%2d ", i)
	}
	s += "\n"
	for _, line := range f.data {
		for j, pixels := range line {
			for i := 0; i < 8 && j*8+i < f.width; i++ {
				if pixels&(1<<uint(7-i)) != 0x00 {
					s += " █ "
				} else {
					s += " . "
				}
			}
		}
		s += "\n"
	}
	return
}

func (f *Frame) CompactString() (s string) {
	for _, line := range f.data {
		for j, pixels := range line {
			for i := 0; i < 8 && j*8+i < f.width; i++ {
				if pixels&(1<<uint(7-i)) != 0x00 {
					s += "█"
				} else {
					s += "."
				}
			}
		}
		s += "\n"
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
