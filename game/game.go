package game

import (
	"TheCovenant/assets"
	"TheCovenant/config"
	"TheCovenant/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	spartan *entities.Spartan
	grunt   *entities.Grunt
}

func NewGame() (*Game, error) {
	p := entities.NewSpartan()
	g := entities.NewGrunt()

	return &Game{spartan: p, grunt: g}, nil
}


func (g *Game) Update() error {
	// 1. Llama al update del jugador y captura el evento de disparo
	shotFired := g.spartan.Update()

	// 2. Llama al update del enemigo (para su contador)
	g.grunt.Update()

	if shotFired {
		// A. Obtén los rectángulos de ambos
		spartanBox := g.spartan.BoundingBox()
		targetBox := g.grunt.BoundingBox()

		// B. Comprueba si se solapan en el eje Y
		//    (Si el jugador está alineado horizontalmente con el enemigo)
		
		// La lógica es:
		// - El 'top' del jugador debe estar por ENCIMA del 'bottom' del enemigo
		// - El 'bottom' del jugador debe estar por DEBAJO del 'top' del enemigo
		isAlignedY := (spartanBox.Min.Y < targetBox.Max.Y) && 
		              (spartanBox.Max.Y > targetBox.Min.Y)

		// C. (Opcional pero recomendado) Comprueba que el jugador está
		//    a la izquierda del enemigo, para que el disparo sea "hacia adelante".
		isToTheLeft := spartanBox.Max.X < targetBox.Min.X

		// D. ¡Golpea si ambas condiciones son verdaderas!
		if isAlignedY && isToTheLeft {
			g.grunt.Hit()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.BackgroundSprite, nil)
	g.spartan.Draw(screen)
	g.grunt.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
