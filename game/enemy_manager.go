package game

import (
	"TheCovenant/entities"
	"image"
	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyManager gestiona todos los enemigos
type EnemyManager struct {
	activeEnemies   []*entities.Grunt
	locationChannel <-chan image.Point // Canal de CONSUMO
	defeatedChannel chan int           // Canal de PRODUCCIÓN (para el score)
	maxEnemies      int
}

// NewEnemyManager crea el gestor
func NewEnemyManager(locationChan <-chan image.Point, max int) *EnemyManager {
	return &EnemyManager{
		activeEnemies:   make([]*entities.Grunt, 0, max),
		locationChannel: locationChan,
		// Canal 'buffer' para 10 muertes (evita bloqueo)
		defeatedChannel: make(chan int, 10), 
		maxEnemies:      max,
	}
}

// DefeatedChannel expone el canal de muertes (solo lectura) para el consumidor
func (m *EnemyManager) DefeatedChannel() <-chan int {
	return m.defeatedChannel
}

// LiveEnemies retorna solo los enemigos vivos para la detección de colisión
func (m *EnemyManager) LiveEnemies() []*entities.Grunt {
	live := make([]*entities.Grunt, 0, len(m.activeEnemies))
	for _, e := range m.activeEnemies {
		if e.IsAlive() {
			live = append(live, e)
		}
	}
	return live
}

// Update se encarga de producir nuevos enemigos y limpiar los muertos
func (m *EnemyManager) Update() {
	// --- 1. Parte PRODUCTOR (de enemigos) ---
	m.spawnEnemies()

	// --- 2. Parte CONSUMIDOR (de enemigos muertos) ---
	m.updateAndCleanEnemies()
}

// spawnEnemies consume del canal de ubicaciones
func (m *EnemyManager) spawnEnemies() {
	// Revisamos si hay una nueva ubicación, pero SIN BLOQUEARNOS
	select {
	case newPos, ok := <-m.locationChannel:
		if !ok {
			return // El canal se cerró
		}
		
		// Solo producimos si no hemos llegado al límite 'n'
		if len(m.activeEnemies) < m.maxEnemies {
			newEnemy := entities.NewGrunt()
			newEnemy.SetPosition(newPos) // Usamos la nueva posición
			m.activeEnemies = append(m.activeEnemies, newEnemy)
		}

	default:
		// No hay nuevas ubicaciones, no hacemos nada
	}
}

// updateAndCleanEnemies actualiza y elimina enemigos
func (m *EnemyManager) updateAndCleanEnemies() {
	// Creamos una nueva lista para los enemigos que "sobreviven" este frame
	newList := make([]*entities.Grunt, 0, len(m.activeEnemies))

	for _, enemy := range m.activeEnemies {
		enemy.Update() // Actualiza el enemigo (movimiento, timers)
		
		if enemy.IsDefeated() {
			// 1. Está muerto Y la animación terminó
			// 2. Producimos un "1" en el canal de muertes
			m.defeatedChannel <- 1
			// 3. NO lo añadimos a la 'newList' (efectivamente borrándolo)
		} else {
			// Está vivo O está en animación de muerte.
			// Lo mantenemos en la lista.
			newList = append(newList, enemy)
		}
	}
	
	// Reemplazamos la lista vieja por la nueva
	m.activeEnemies = newList
}

// Draw dibuja todos los enemigos activos (vivos y muriendo)
func (m *EnemyManager) Draw(screen *ebiten.Image) {
	for _, enemy := range m.activeEnemies {
		enemy.Draw(screen)
	}
}