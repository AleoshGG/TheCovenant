package game

import (
	"TheCovenant/assets"
	"context"
	"fmt"
	"image/color"
	"strconv"
	"sync/atomic" 

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type ScoreManager struct {
	defeatedChannel <-chan int // Canal de CONSUMO
	ctx             context.Context
	cancel          context.CancelFunc
	
	// Usamos 'atomic' para el score, porque una gorutina
	// escribirá en él, y el hilo principal (Draw) lo leerá.
	currentScore atomic.Uint64
}

func NewScoreManager(defeatedChan <-chan int) *ScoreManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ScoreManager{
		defeatedChannel: defeatedChan,
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (s *ScoreManager) consumer() {
	for {
		select {
		case <-s.ctx.Done():
			// Recibimos señal de parar
			return
				
		case points, ok := <-s.defeatedChannel:
			if !ok {
				// El canal se cerró
				return
			}
			// ¡Suma los puntos al score de forma segura!
			s.currentScore.Add(uint64(points))
		}
	}

}

// Start lanza la gorutina del consumidor
func (s *ScoreManager) Start() {
	go s.consumer()
}

// Stop detiene la gorutina
func (s *ScoreManager) Stop() {
	s.cancel()
}

// GetScore (lectura segura)
func (s *ScoreManager) GetScore() uint64 {
	return s.currentScore.Load()
}

// Draw dibuja el score en la pantalla
func (s *ScoreManager) Draw(screen *ebiten.Image) {
	// Convierte el score (uint64) a un string
	scoreStr := strconv.FormatUint(s.GetScore(), 10)
	
	// Prepara el texto a dibujar
	// Usamos fmt.Sprintf para formatear (ej. "SCORE: 10")
	drawText := fmt.Sprintf("SCORE: %s", scoreStr)
	
	// Usa la fuente que cargamos en assets
	// Dibuja el texto en (10, 30) con color blanco
	text.Draw(screen, drawText, assets.ScoreFont, 10, 30, color.White)
}