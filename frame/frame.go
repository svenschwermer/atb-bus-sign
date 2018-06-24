package frame

import (
	"bytes"
)

// Frame is a collection of lines
// 1st index: line index
type Frame [][]byte

func New(lines, columns int) Frame {
	f := make([][]byte, lines)
	for i := range f {
		f[i] = make([]byte, (columns+7)/8)
	}
	return f
}

func (f Frame) ConcatenateLines() []byte {
	return bytes.Join(f, []byte{})
}

// Modify sets the columns [start, end) with the data provided in lines
// In "non-full" bytes, the lowermost bits are considered padding
func (f Frame) Modify(start, end int, lines [][]byte) error {
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
			destByte := f[i][destByteIdx]

			bits := min(srcBits, destBits)
			mask := byte((1 << uint(bits)) - 1)
			destShift := uint(8 - destBytePos - bits)
			srcShift := uint(8 - srcBytePos - bits)

			srcByte = (srcByte >> srcShift) & mask
			destByte &^= mask << destShift   // clear bits in dest byte
			destByte |= srcByte << destShift // apply bits from src byte
			f[i][destByteIdx] = destByte

			srcPos += bits
			destPos += bits
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
