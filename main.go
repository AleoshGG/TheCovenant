package main

import (
	"TheCovenant/assets"
	"TheCovenant/config"
	"TheCovenant/game"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	assets.LoadAssets()
	
	// Usa los valores de config para la ventana
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle(config.Title)

	// Crea la instancia del juego desde el paquete 'game'
	myGame, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	// Ejecuta el juego
	if err := ebiten.RunGame(myGame); err != nil {
		log.Fatal(err)
	}
}