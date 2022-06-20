package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	log "github.com/sirupsen/logrus"

	"github.com/vatine/cryptoclicker/pkg/display"
)

type Game struct {
	field *display.ButtonField
}

var black = color.RGBA{R: 0, G: 0, B: 0, A: 255}

func NewGame() *Game {
	log.Debug("Creating new game")
	g := Game{
		field: display.NewField(),
	}

	return &g
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(black)

	g.field.Draw(screen)
	display.DrawMatrix(g.field, screen)
}

func (g *Game) Update() error {
	pressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if pressed {
		x, y := ebiten.CursorPosition()
		log.WithFields(log.Fields{
			"pressed": pressed,
			"x":       x,
			"y":       y,
		}).Debug("Update()")
		g.field.Click(x, y)
	}

	return nil
}

func (g *Game) Layout(x, y int) (int, int) {
	return 640, 480
}
