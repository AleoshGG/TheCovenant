package game

import (
	"TheCovenant/assets"
	"TheCovenant/config"
	"TheCovenant/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	spartan *entities.Spartan
}

func NewGame() (*Game, error) {
	p := entities.NewSpartan()

	return &Game{spartan: p}, nil
}


func (g *Game) Update() error {
	g.spartan.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.BackgroundSprite, nil)
	g.spartan.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
