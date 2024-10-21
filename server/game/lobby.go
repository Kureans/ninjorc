package game

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type LobbyManager struct {
	lobbies  []Lobby
	upgrader websocket.Upgrader
}

func (m *LobbyManager) Init() {
	m.lobbies = make([]Lobby, 10)
	m.lobbies[0] = Lobby{}
	m.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// start lobby 1 for now
	go m.lobbies[0].run()
}

func (m *LobbyManager) InitClient(w http.ResponseWriter, r *http.Request) {
	m.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientCh := make(chan ClientInput)
	playerCh := make(chan PlayerInput)
	//add to lobby 1 for now
	conn := Connection{
		socket:   socket,
		clientCh: clientCh,
		playerCh: playerCh,
	}

	player := Player{
		isReady: true,
		conn:    conn,
		controller: PlayerController{
			inputCh: playerCh,
		},
	}

	m.lobbies[0].players = append(m.lobbies[0].players, player)
	print("added a client")
}

type Lobby struct {
	game    Game
	players []Player
}

func (l *Lobby) run() {
	for !(l.hasEnoughPlayers() && l.arePlayersReady()) {
	}
	print("enough players and all ready, starting game")

	l.startGame()
}

func (l *Lobby) hasEnoughPlayers() bool {
	return len(l.players) >= 2
}

func (l *Lobby) arePlayersReady() bool {
	status := true
	for _, player := range l.players {
		status = status && player.isReady
	}
	return status
}

func (l *Lobby) startGame() {
	l.game.run()
}
