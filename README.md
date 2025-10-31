# The Covenant

Este es un juego 2D desarrollado en Go con la librería Ebiten. El objetivo del juego es simple: dispara a los enemigos (Grunts) que aparecen en la pantalla para sumar puntos.

## Arquitectura del Proyecto

El proyecto está estructurado en varios paquetes, cada uno con una responsabilidad clara:

- `main.go`: El punto de entrada de la aplicación. Se encarga de inicializar y correr el juego.
- `assets/`: Maneja la carga de todos los recursos del juego, como imágenes, fuentes y audio.
- `config/`: Contiene la configuración básica del juego, como las dimensiones de la pantalla y el título.
- `entities/`: Define las entidades del juego, como el jugador (`Spartan`) y los enemigos (`Grunt`).
- `game/`: Contiene la lógica principal del juego.

### El Paquete `game`

El paquete `game` es el núcleo del juego y está dividido en los siguientes componentes:

- `game.go`: Orquesta el flujo principal del juego, actualizando y dibujando todos los elementos en la pantalla.
- `spawner.go`: Un productor que genera ubicaciones aleatorias para los nuevos enemigos.
- `enemy_manager.go`: Un consumidor de las ubicaciones generadas por `spawner.go`. También actúa como productor, notificando al `score_manager.go` cuando un enemigo es derrotado.
- `score_manager.go`: Un consumidor que actualiza la puntuación cada vez que un enemigo es derrotado.

## Flujo del Juego

1. El juego comienza y el `main.go` inicializa todos los componentes.
2. El `LocationSpawner` (`spawner.go`) comienza a generar coordenadas `(x, y)` aleatorias y las envía a través de un canal.
3. El `EnemyManager` (`enemy_manager.go`) está escuchando en ese canal. Cada vez que recibe una nueva ubicación, crea un nuevo `Grunt` en esa posición.
4. El jugador (`Spartan`) puede moverse y disparar.
5. Cuando un disparo alcanza a un `Grunt`, el `EnemyManager` lo marca como "derrotado".
6. Una vez que la animación de derrota del `Grunt` termina, el `EnemyManager` envía una señal (un entero) a través de otro canal para notificar que un enemigo ha sido eliminado.
7. El `ScoreManager` (`score_manager.go`) recibe esta señal y actualiza la puntuación del jugador.
8. El ciclo se repite, con nuevos enemigos apareciendo constantemente.

## Patrón de Concurrencia: Productor-Consumidor

Este proyecto utiliza el patrón de concurrencia **productor-consumidor** para gestionar la aparición de enemigos y la actualización de la puntuación.

### ¿Qué es el Patrón Productor-Consumidor?

El patrón productor-consumidor es un patrón de diseño de concurrencia que desacopla a los "productores" de datos de los "consumidores" de datos a través de una cola o canal compartido.

- **Productor**: Es el responsable de generar datos y ponerlos en la cola.
- **Consumidor**: Es el responsable de tomar datos de la cola y procesarlos.

Este patrón es útil porque permite que el productor y el consumidor trabajen a diferentes velocidades. Por ejemplo, el productor puede generar datos rápidamente y el consumidor puede procesarlos más lentamente, sin que uno bloquee al otro.

### Implementación en "The Covenant"

En este juego, el patrón productor-consumidor se utiliza en dos lugares:

1. **Aparición de Enemigos**:

    - **Productor**: `LocationSpawner` genera ubicaciones `(image.Point)` y las envía a un canal.
    - **Consumidor**: `EnemyManager` recibe las ubicaciones del canal y crea nuevos `Grunt` en esas posiciones.

    Esto permite que la lógica de generación de ubicaciones esté completamente separada de la lógica de gestión de enemigos.

2. **Actualización de la Puntuación**:

    - **Productor**: `EnemyManager` envía una señal (un entero) a un canal cada vez que un enemigo es derrotado.
    - **Consumidor**: `ScoreManager` recibe la señal del canal y actualiza la puntuación.

    Esto desacopla la lógica de puntuación de la lógica de los enemigos, lo que hace que el código sea más modular y fácil de mantener.
