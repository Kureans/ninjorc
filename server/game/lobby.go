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
	m.lobbies[0] = Lobby{
		isOpen: true,
	}
	m.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// start lobby 1 for now
	go m.lobbies[0].run()
}

func (m *LobbyManager) HandleNewClient(w http.ResponseWriter, r *http.Request) {
	m.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	lobbyCh := make(chan LobbyInput)
	gameCh := make(chan GameInputBatch)
	//add to lobby 1 for now
	conn := Connection{
		socket:  socket,
		gameCh:  gameCh,
		lobbyCh: lobbyCh,
	}

	player := &Player{
		id:      len(m.lobbies[0].players) + 1,
		isReady: true,
		conn:    conn,
		//init orc later when game starts
		controller: PlayerController{
			inputCh: gameCh,
		},
	}

	m.lobbies[0].players = append(m.lobbies[0].players, player)
	go player.routeInputs()
	go player.handleLobbyInputs(lobbyCh)
	print("added a client\n")
}

type Lobby struct {
	isOpen  bool
	game    Game
	players []*Player
}

func (l *Lobby) run() {
	for l.isOpen {
		for !(l.hasEnoughPlayers() && l.arePlayersReady()) {
		}
		print("enough players and all ready, starting game\n")
		l.notifyPlayers()
		l.startGame()
		print("game ended, back to lobby")
		l.resetReadyStatus()
	}
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

func (l *Lobby) resetReadyStatus() {
	for _, player := range l.players {
		player.isReady = false
	}
}

func (l *Lobby) notifyPlayers() {
	payloadArr := make([]interface{}, 1)
	payloadArr[0] = LobbyInput{CanStartGame: true}
	notificationPacket := Packet{
		Id: 1, Type: "L", Size: 1, Data: payloadArr,
	}
	for idx, player := range l.players {
		notificationPacket.Id = idx
		player.conn.sendPacket(notificationPacket)
	}

}

func (l *Lobby) startGame() {
	l.game.initResponseChannels(&l.players)
	l.game.gameState.init(&l.players)
	l.game.run()
}
