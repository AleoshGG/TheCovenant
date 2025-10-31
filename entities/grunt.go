package entities

import (
	"TheCovenant/assets"
	"TheCovenant/config"

	"github.com/hajimehoshi/ebiten/v2"
)

const defeathFrameDuration = 10

type Grunt struct {
	idleImg *ebiten.Image
	defeathImg *ebiten.Image
	X 	float64
	Y float64
	Speed float64
	Opst *ebiten.DrawImageOptions

	defeathTimer int
}

func NewGrunt() *Grunt {
	idle := assets.GruntSprite
	defeath := assets.GruntDejectedSprite

	width, height := idle.Size()

	return &Grunt{
		idleImg: idle,
		defeathImg: defeath,
		X: float64(config.ScreenWidth - width) / 2,
		Y: float64(config.ScreenHeight - height) / 2,
		Speed: 4,
		Opst: &ebiten.DrawImageOptions{},
	}
}

func (g *Grunt) Update() {
}

func (g *Grunt) Draw(screen *ebiten.Image) {
	g.Opst.GeoM.Reset()
	g.Opst.GeoM.Translate(g.X, g.Y)

	if g.defeathTimer > 0 {
		screen.DrawImage(g.defeathImg, g.Opst)
	} else {
		screen.DrawImage(g.idleImg, g.Opst)
	}
}