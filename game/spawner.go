package game

import (
	"TheCovenant/config"
	"context"
	"image"
	"math/rand"
	"time"
)

// Este es el encargado de estar generando nuevas ubicaciones aleatorias tomando en cuenta
// la pantalla de juego, un tiempo de espera, y la cantidad de ubicaciones que tiene que
// generar, esto a corde a la cantidad de enemigos

type LocationSpawner struct {
	producerChannel chan image.Point
	context 		context.Context
	cancel 			context.CancelFunc
}

func NewLocationSpawner() *LocationSpawner {
	ctx, cancel := context.WithCancel(context.Background())

	producerChannel := make(chan image.Point, 10)
	
	return &LocationSpawner{
		producerChannel: producerChannel,
		context: ctx,
		cancel: cancel,
	}
}

func (ls *LocationSpawner) Spawner(ticker time.Ticker, source rand.Source, rng rand.Rand) {
	defer func () {
		ticker.Stop()
		close(ls.producerChannel)
	} ()

	for {
		select {
			case <-ls.context.Done():
				// Recibimos la señal de parar
				return 
			
			case <-ticker.C:
				// ¡Tiempo de generar una nueva ubicación!
				
				// Algoritmo: Spawnear en el borde derecho,
				// en una altura (Y) aleatoria.
				
				x := config.ScreenWidth - 500 // 50 píxeles fuera de la pantalla
				
				// Rango de Y (ej. entre 10% y 90% de la altura)
				minY := int(float64(config.ScreenHeight) * 0.1)
				maxY := int(float64(config.ScreenHeight) * 0.9)
				y := rng.Intn(maxY-minY) + minY
				
				newPosition := image.Point{X: x, Y: y}
				
				// Intenta enviar la posición al canal,
				// pero no te bloquees si está lleno.
				select {
				case ls.producerChannel <- newPosition:
					// Posición enviada
				default:
					// El canal está lleno, no pasa nada,
					// descartamos esta posición e intentamos en el próximo tick.
				}
			}
	}

}

func (ls *LocationSpawner) Start() {
	// Usamos un generador de números aleatorios propio
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	ticker := time.NewTicker(2 * time.Second)
	go ls.Spawner(*ticker, source, *rng)
}

// Stop detiene la gorutina de forma segura
func (ls *LocationSpawner) Stop() {
	ls.cancel()
}

// LocationChannel expone el canal como solo-lectura (<-chan)
// Así, el "consumidor" solo puede leer de él.
func (ls *LocationSpawner) LocationChannel() <-chan image.Point {
	return ls.producerChannel
}