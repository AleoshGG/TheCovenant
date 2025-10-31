package entities

import (
	"TheCovenant/assets"
	"TheCovenant/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const shootFrameDuration = 15

type Spartan struct {
	idleImg  *ebiten.Image // Imagen normal
	shootImg *ebiten.Image // Imagen disparando
	X        float64
	Y        float64
	Speed    float64
	Opts     *ebiten.DrawImageOptions

	// Temporizador para el frame de disparo
	shootTimer int
}

func NewSpartan() *Spartan {
// Obtenemos AMBAS imágenes de los assets
	idle := assets.SpartanSprite
	shoot := assets.SpartanShootSprite

	// Usamos la imagen 'idle' para calcular el tamaño
	width, height := idle.Size()

	return &Spartan{
		idleImg: idle,
		shootImg: shoot,
		X: float64(config.ScreenWidth - width) / 12,
		Y: float64(config.ScreenHeight - height) / 12,
		Speed: 10,
		Opts: &ebiten.DrawImageOptions{},
	}
}

func (p *Spartan) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Y -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Y += p.Speed
	}

	if p.shootTimer > 0 {
		p.shootTimer--
	}

	if p.shootTimer == 0 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonLeft)) {
		
		// Acción 1: Poner sonido
		assets.ShootSniperSound.Rewind() // Rebobina el sonido al inicio
		assets.ShootSniperSound.Play()   // ¡Suena!

		// Acción 2: Cambiar el frame (iniciando el temporizador)
		p.shootTimer = shootFrameDuration
	}

}

func (p *Spartan) Draw(screen *ebiten.Image) {
	p.Opts.GeoM.Reset()
	p.Opts.GeoM.Translate(p.X, p.Y)

	// Decidimos qué imagen dibujar
	if p.shootTimer > 0 {
		// Si el temporizador está activo, dibuja el frame de disparo
		screen.DrawImage(p.shootImg, p.Opts)
	} else {
		// Si no, dibuja el frame normal (idle)
		screen.DrawImage(p.idleImg, p.Opts)
	}
}