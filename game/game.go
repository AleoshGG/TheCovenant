package game

import (
	"TheCovenant/assets"
	"TheCovenant/config"
	"TheCovenant/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxEnemies = 10
)

type Game struct {
	spartan *entities.Spartan
	
	locationSpawner *LocationSpawner
	enemyManager    *EnemyManager
	scoreManager    *ScoreManager
}

func NewGame() (*Game, error) {
	spartan := entities.NewSpartan()

	locationSpawner := NewLocationSpawner()
	enemyManager := NewEnemyManager(locationSpawner.LocationChannel(), maxEnemies)

	scoreManager := NewScoreManager(enemyManager.DefeatedChannel())

	locationSpawner.Start()
	scoreManager.Start()

	return &Game{
		spartan: spartan, 
		locationSpawner: locationSpawner,
		enemyManager: enemyManager,
		scoreManager: scoreManager,
	}, nil
}


// Update delega la lógica a los sistemas
func (g *Game) Update() error {
	// 1. Actualiza al jugador
	shotFired := g.spartan.Update()

	// 2. Actualiza el gestor de enemigos
	//    (Esto spawnea nuevos y limpia los muertos)
	g.enemyManager.Update()
	
	// 3. Lógica de Colisión (la conexión)
	if shotFired {
		playerBox := g.spartan.BoundingBox()

		// ¡Iteramos sobre los enemigos VIVOS!
		for _, enemy := range g.enemyManager.LiveEnemies() {
			targetBox := enemy.BoundingBox()

			// Comprueba alineación Y
			isAlignedY := (playerBox.Min.Y < targetBox.Max.Y) &&
				(playerBox.Max.Y > targetBox.Min.Y)

			// Comprueba que está a la izquierda
			isToTheLeft := playerBox.Max.X < targetBox.Min.X

			if isAlignedY && isToTheLeft {
				enemy.Hit()
				// (Opcional: 'break' si la bala solo golpea a uno)
			}
		}
	}

	return nil
}

// Draw delega el dibujado
func (g *Game) Draw(screen *ebiten.Image) {
	// 1. Dibuja el fondo
	screen.DrawImage(assets.BackgroundSprite, nil)

	// 2. Dibuja al jugador
	g.spartan.Draw(screen)

	// 3. Dibuja todos los enemigos (vivos y muriendo)
	g.enemyManager.Draw(screen)

	// 4. Dibuja el score
	g.scoreManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
