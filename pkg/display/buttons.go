package display

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
)

const width = 8

var buttons map[bool]*ebiten.Image
var blank = color.RGBA{R: 0, G: 0, B: 0, A: 0}

type ButtonGroup struct {
	vals uint8
}

type ButtonRow struct {
	rowNo   int
	buttons [4]*ButtonGroup
}

type ButtonField struct {
	rows [16]*ButtonRow
}

func NewField() *ButtonField {
	var rv ButtonField

	for r := 0; r < 16; r++ {
		rv.rows[r] = NewButtonRow(r)
	}

	return &rv
}

func NewButtonRow(row int) *ButtonRow {
	var rv ButtonRow

	rv.rowNo = row
	for ix := 0; ix < 4; ix++ {
		rv.buttons[ix] = new(ButtonGroup)
	}

	return &rv
}

func (f *ButtonField) Click(x, y int) {
	xPos := x / (width + 1)
	yPos := y / (width + width + 1)
	log.WithFields(log.Fields{
		"x":    x,
		"y":    y,
		"xPos": xPos,
		"yPos": yPos,
	}).Debug("click")

	if yPos <= len(f.rows) {
		f.rows[yPos].Click(xPos)
	}
}

func (f *ButtonField) AsArray() [64]byte {
	var rv [64]byte

	for ix := 0; ix < 16; ix++ {
		base := ix * 4
		for offset, val := range f.rows[ix].AsArray() {
			rv[base+offset] = val
		}
	}

	return rv
}

// Activate a specific button within a ButtonGroup,
// The buttons are numbered 7 to 0 from left to right.
func (b *ButtonGroup) Activate(p int) {
	b.vals = b.vals | (1 << (7 - p))
}

// Dectivate a specific button within a ButtonGroup,
// The buttons are numbered 0 to 7 from left to right.
func (b *ButtonGroup) Dectivate(p int) {
	b.vals = b.vals & ^(1 << (7 - p))
}

// Toggle a specific button within a ButtonGroup,
// The buttons are numbered 0 to 7 from left to right.
func (b *ButtonGroup) Toggle(p int) {
	log.WithFields(log.Fields{
		"b.vals": b.vals,
	}).Debug("ButtonGroup.Toggle before")
	b.vals = b.vals ^ (1 << (7 - p))
	log.WithFields(log.Fields{
		"b.vals": b.vals,
	}).Debug("ButtonGroup.Toggle after")
}

// Return the states of a ButtonGroup as an array of bools, each one
// reflecting the activation status of the button.
func (b ButtonGroup) AsArray() [8]bool {
	var rv [8]bool

	for n := 0; n < 8; n++ {
		mask := uint8(1 << (7 - n))
		rv[n] = (b.vals & mask) == mask
	}

	return rv
}

func (b *ButtonRow) AsArray() [4]byte {
	var rv [4]byte

	for ix := 0; ix < 4; ix++ {
		rv[ix] = b.buttons[ix].vals
	}

	return rv
}

func (b *ButtonRow) Click(x int) {
	button := x % 8
	group := x / 8

	log.WithFields(log.Fields{
		"row":    b.rowNo,
		"button": button,
		"group":  group,
	}).Debug("ButtonRow click")
	if group < len(b.buttons) {
		b.buttons[group].Toggle(button)
	}
}

// Initialize the unset and set button images
// For the moment, hard-coded for a width of 8 and a height of 16
//
//  01234567
// 0  ****
// 1 *    *
// 2*      *
// 3*      *
func init() {
	red := color.RGBA{R: 64, G: 0, B: 0, A: 255}
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	gray := color.RGBA{R: 192, G: 192, B: 192, A: 255}

	off := ebiten.NewImage(width, 2*width)
	off.Fill(blank)

	on := ebiten.NewImage(width, 2*width)
	on.Fill(blank)

	frame := func(i *ebiten.Image) {
		i.Set(2, 0, black)
		i.Set(3, 0, black)
		i.Set(4, 0, black)
		i.Set(5, 0, black)
		i.Set(1, 1, black)
		i.Set(6, 1, black)
		for y := 2; y <= 13; y++ {
			i.Set(0, y, black)
			i.Set(7, y, black)
		}
		i.Set(1, 14, black)
		i.Set(6, 14, black)
		i.Set(2, 15, black)
		i.Set(3, 15, black)
		i.Set(4, 15, black)
		i.Set(5, 15, black)
	}
	fill := func(c color.Color, i *ebiten.Image) {
		for y := 1; y < 15; y++ {
			left := 1
			right := 6
			if (y == 1) || (y == 14) {
				left = 2
				right = 5
			}

			for x := left; x <= right; x++ {
				i.Set(x, y, c)
			}
		}
	}
	frame(off)
	fill(red, off)
	frame(on)
	fill(green, on)

	on.Set(3, 2, gray)
	on.Set(4, 2, gray)
	on.Set(3, 3, gray)
	on.Set(4, 3, gray)

	off.Set(3, 12, gray)
	off.Set(4, 12, gray)
	off.Set(3, 13, gray)
	off.Set(4, 13, gray)

	buttons = map[bool]*ebiten.Image{false: off, true: on}
}

func (b ButtonRow) Draw(i *ebiten.Image) {
	yBase := float64(b.rowNo * (width + width + 1))
	xBase := 0.0
	xDelta := float64(width + 1)

	for _, bg := range b.buttons {
		settings := bg.AsArray()
		for _, bv := range settings {
			var g ebiten.GeoM
			g.Translate(xBase, yBase)
			opts := &ebiten.DrawImageOptions{
				GeoM: g,
			}
			src := buttons[bv]
			i.DrawImage(src, opts)

			xBase += xDelta
		}
	}
}

func (bf *ButtonField) Draw(i *ebiten.Image) {
	for _, br := range bf.rows {
		br.Draw(i)
	}
}
