package assets

import (
	"bytes"
	"embed"
	"io/fs"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed all:resources
var embeddedImages embed.FS

var (
	BackgroundSprite 	 *ebiten.Image

	SpartanSprite 		 *ebiten.Image
	SpartanShootSprite 	 *ebiten.Image
	GruntSprite 		 *ebiten.Image
	GruntDejectedSprite  *ebiten.Image

	audioContext 		 *audio.Context
	ShootSniperSound 	 *audio.Player

	ScoreFont font.Face
)


func LoadAssets() {
	var err error 

	if audioContext == nil {
		audioContext = audio.NewContext(44100) 
	}

	BackgroundSprite, err = loadImage("resources/images/background.png")
	if err != nil {
		log.Fatalf("Error al cargar la imagen de fondo: %v", err)
	}

	// Carga el jugador
	SpartanSprite, err = loadImage("resources/images/chief01.png")
	if err != nil {
		log.Fatalf("Error al cargar la imagen del jugador: %v", err)
	}

	SpartanShootSprite, err = loadImage("resources/images/chief02.png")
	if err != nil {
		log.Fatalf("Error al cargar la imagen del jugador: %v", err)
	}

	// Carga el grunt
	GruntSprite, err = loadImage("resources/images/grunt01.png")
	if err != nil {
		log.Fatalf("Error al cargar la imagen del jugador: %v", err)
	}

	GruntDejectedSprite, err = loadImage("resources/images/grunt02.png")
	if err != nil {
		log.Fatalf("Error al cargar la imagen del jugador: %v", err)
	}

	ShootSniperSound, err = loadAudio("resources/audio/sniper.wav")
	if err != nil {
		log.Fatalf("Error al cargar el sonido de disparo: %v", err)
	}

	ScoreFont, err = loadFont("resources/fonts/PressStart2P-Regular.ttf")
	if err != nil {
		log.Fatalf("Error al cargar la fuente: %v", err)
	}
}

func loadImage(path string) (*ebiten.Image, error) {
	fileBytes, err := embeddedImages.ReadFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func loadAudio(path string) (*audio.Player, error) {
	// 1. Abre el archivo desde el embed
	file, err := embeddedImages.Open(path)
	if err != nil {
		log.Printf("Ruta de audio fallida: %s", path)
		return nil, err
	}

	// 2. Decodifica el .wav
	stream, err := wav.DecodeWithSampleRate(audioContext.SampleRate(), file.(fs.File))
	if err != nil {
		return nil, err
	}
	
	// 3. Crea un reproductor
	//    Usamos audio.NewPlayer para poder rebobinar y reproducir
	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		return nil, err
	}
	
	return player, nil
}

func loadFont(path string) (font.Face, error) {
	// 1. Lee los bytes del archivo
	fileBytes, err := embeddedImages.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 2. Parsea el archivo .ttf
	tt, err := opentype.Parse(fileBytes)
	if err != nil {
		return nil, err
	}

	// 3. Crea la "cara" de la fuente (Face)
	//    Usamos 48 'DPI' y un tamaño de 24 (puedes ajustar esto)
	const dpi = 72
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return fontFace, nil
}