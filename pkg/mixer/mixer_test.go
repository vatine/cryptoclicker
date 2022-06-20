package mixer

import (
	"testing"

	"crypto/sha256"
)

func TestMixer1(t *testing.T) {
	var blank [64]byte

	blank[0] = 0x80

	own := Mix(blank)
	real := sha256.Sum256([]byte{})

	for ix, self := range own {
		if self != real[ix] {
			t.Errorf("ix %d: mismatched bits: %08b  %08b", ix, self, real[ix])
		}
	}
}

func TestMixer2(t *testing.T) {
	real := [32]byte{0xB9, 0x4D, 0x27, 0xB9, 0x93, 0x4D, 0x3E, 0x08, 0xA5, 0x2E, 0x52, 0xD7, 0xDA, 0x7D, 0xAB, 0xFA, 0xC4, 0x84, 0xEF, 0xE3, 0x7A, 0x53, 0x80, 0xEE, 0x90, 0x88, 0xF7, 0xAC, 0xE2, 0xEF, 0xCD, 0xE9}

	var in [64]byte

	in[0] = 'h'
	in[1] = 'e'
	in[2] = 'l'
	in[3] = 'l'
	in[4] = 'o'
	in[5] = ' '
	in[6] = 'w'
	in[7] = 'o'
	in[8] = 'r'
	in[9] = 'l'
	in[10] = 'd'
	in[11] = 0x80

	in[63] = 88

	own := Mix(in)

	for ix, self := range own {
		if self != real[ix] {
			t.Errorf("ix %d: mismatched bits: saw %08b  want %08b", ix, self, real[ix])
		}
	}

}
