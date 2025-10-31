package entities

import (
	"TheCovenant/assets"
	"TheCovenant/config"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const defeathFrameDuration = 10

type Grunt struct {
	idleImg *ebiten.Image
	defeathImg *ebiten.Image
	X 	float64
	Y float64
	Speed float64
	Opts *ebiten.DrawImageOptions

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
		Opts: &ebiten.DrawImageOptions{},
	}
}

func (g *Grunt) Update() {
	if g.defeathTimer > 0 {
		g.defeathTimer--
	}
}

func (g *Grunt) Hit() {
	// Solo activa la animación si no está ya derrotado
	if g.defeathTimer == 0 {
		g.defeathTimer = defeathFrameDuration
		// (En el futuro, aquí podrías reproducir un sonido de "grunt herido")
	}
}

func (g *Grunt) Draw(screen *ebiten.Image) {
	g.Opts.GeoM.Reset()
	g.Opts.GeoM.Translate(g.X, g.Y)

	if g.defeathTimer > 0 {
		screen.DrawImage(g.defeathImg, g.Opts)
	} else {
		screen.DrawImage(g.idleImg, g.Opts)
	}
}

// BoundingBox retorna el "hitbox" actual del Grunt
func (g *Grunt) BoundingBox() image.Rectangle {
	// Obtiene el tamaño de la imagen
	width, height := g.idleImg.Size()

	// Crea un rectángulo usando la posición (X,Y) y el tamaño
	// image.Rect crea un Rect(x0, y0, x1, y1)
	return image.Rect(
		int(g.X),                 // x0 (esquina superior izquierda)
		int(g.Y),                 // y0 (esquina superior izquierda)
		int(g.X) + width,         // x1 (esquina inferior derecha)
		int(g.Y) + height,        // y1 (esquina inferior derecha)
	)
}