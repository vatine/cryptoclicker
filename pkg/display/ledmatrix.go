package display

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"

	"github.com/vatine/cryptoclicker/pkg/mixer"
)

const ledSize = 8

var led *ebiten.Image

func DrawMatrix(f *ButtonField, i *ebiten.Image) {
	startBytes := f.AsArray()
	bytes := mixer.Mix(startBytes)

	for ix, val := range bytes {
		y := (ix / 2) * ledSize
		xBase := (ix%2)*8*ledSize + 290

		for b := 0; b < 8; b++ {
			mask := byte(1) << (7 - b)
			x := xBase + b*ledSize
			log.WithFields(log.Fields{
				"x":    x,
				"y":    y,
				"b":    b,
				"mask": mask,
				"val":  val,
				"ix":   ix,
			}).Debug("DrawMatrix")
			if (val & mask) == mask {
				var g ebiten.GeoM
				g.Translate(float64(x), float64(y))
				opts := &ebiten.DrawImageOptions{
					GeoM: g,
				}
				i.DrawImage(led, opts)
			}
		}
	}

}

func init() {
	led = ebiten.NewImage(ledSize, ledSize)
	ledCol := color.RGBA{R: 127, G: 255, B: 10, A: 255}
	// blank := color.RGBA{R: 0, G: 0, B: 0, A: 0}

	half := ledSize / 2
	for x := 0; x < ledSize; x++ {
		for y := 0; y < ledSize; y++ {
			dx := x - half
			dy := y - half

			if (dx*dx)+(dy*dy) < (half * half) {
				led.Set(x, y, ledCol)
			}
		}
	}
}
