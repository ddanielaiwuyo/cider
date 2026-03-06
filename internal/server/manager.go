package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"

	"github.com/jackc/pgx/v5"
)

type GameSession struct {
	id      int
	players [2]client
}

type manager struct {
	register      chan client
	remove        chan int
	deliver       chan Message
	sessions      map[int]*GameSession
	serverDeliver chan Message
	dbconn        *pgx.Conn
	query         chan Query
}

func NewManager(dbConn *pgx.Conn) *manager {
	return &manager{
		register: make(chan client),
		remove:   make(chan int),
		deliver:  make(chan Message, 10),
		sessions: make(map[int]*GameSession),
		query:    make(chan Query),
		dbconn:   dbConn,
	}
}

func (m *manager) Listen(ctx context.Context) {
	parentCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	play := make(chan Message, 2)
	for {
		select {
		case client := <-m.register:
			slog.Info("adding new client", "", client.id)
			connectedUsers[client.id] = client.conn
		case id := <-m.remove:
			slog.Info("removing client with id:", "", id)
			delete(connectedUsers, id)
		case msg := <-m.deliver:

			switch msg.MessageType {
			case NewGameMessage:
				log.Println("create new game")
				ssid := m.createGameSession(msg.From, msg.Dest)
				if ssid < 0 {
					log.Println("sssid got violated??")
					continue
				}
				go m.startGame(parentCtx, ssid, play)
			case GameMessage:
				log.Println("game message")
				play <- msg
				sendMessage(m, msg)
				// so if the game message is here
				// can we pipe into the the startGameFunc?
			case ChatMessage:
				log.Println("normal message")
				sendMessage(m, msg)
			default:
				log.Println("unsupported message perhaps?", msg)
			}

		case q := <-m.query:
			// returns result to caller via channel
			m.executeQuery(ctx, q)

		case <-ctx.Done():
			slog.Info("exiting manager:", "", ctx.Err().Error())
			return
		}
	}
}

// TODO we cant determine the kind of query to
// parse the struct into, is it better we return
// the rows to scan for the caller?
func (mgr *manager) executeQuery(ctx context.Context, q Query) {
	// yea GO had to fight me w this? why doesnt string
	// satisfy any??

	// args := make([]any, len(q.params))
	// for i, v := range q.params {
	// 	args[i] = v
	// }
	// pgxConn.QueryRow(ctx, query, args...)
	rows := mgr.dbconn.QueryRow(ctx, q.query, q.params...)
	q.result <- rows
}

func (mgr *manager) startGame(ctx context.Context, ssid int, play chan Message) {
	log.Println(" [debug] starting game")
	session, ok := mgr.sessions[ssid]
	if !ok {
		slog.Info(" [start_game] could not find ssid for game", "id", ssid)
		return
	}

	log.Println(" [debug] sending welcome msg")
	_ = session

	// welcomeMsg := `
	// The game is about to be start...
	// Buckle up brochachos
	// `

	// conn1, conn2 := session.players[0], session.players[1]
	// msg := Message{
	// 	From:        serverId,
	// 	MessageType: GameMessage,
	// 	Content:     welcomeMsg,
	// 	Dest:        conn1.id,
	// }

	// mgr.deliver <- msg
	// msg.Dest = conn2.id
	// mgr.deliver <- msg

	log.Println(" [debug] inside go routine")
	for {
		select {
		case newplay, ok := <-play:
			if !ok {
				log.Println(" [debug] play channel closed!")
				return
			}

			log.Printf(" [debug] [NEW-PLAY]: %+v\n", newplay)

		case <-ctx.Done():
			log.Println(" [game-sess] calling home, parent exiting")
			return
		}
	}

}

func (mgr *manager) createGameSession(from int, dest int) int {
	homeConn, found := connectedUsers[from]
	if !found {
		slog.Info("home-conn is possibly unregistered, and not found", "", "")
	}

	if homeConn == nil {
		slog.Info("home-conn is nil", "", homeConn)
		return -1
	}
	awayConn, destfound := connectedUsers[dest]
	if !destfound {
		slog.Info("away-conn is possibly unregistered, and not found", "", "")
	}

	if awayConn == nil {
		slog.Info("away-conn is nil", "", homeConn)
		return -1
	}

	sessionId := rand.Intn(8200)
	mgr.sessions[sessionId] = &GameSession{
		id: sessionId,
		players: [2]client{
			{id: from, conn: homeConn},
			{id: dest, conn: awayConn},
		},
	}

	log.Println("[debug] created game session succesfully")
	welcomeMsg := `
	The game is about to be start...
	Buckle up brochachos
	`

	// conn1, conn2 := session.players[0], session.players[1]
	msg := Message{
		From:        serverId,
		MessageType: ChatMessage,
		Content:     fmt.Sprintf("%s;ssid%d", welcomeMsg, sessionId),
		Dest:        from,
	}

	mgr.deliver <- msg
	msg.Dest = dest
	mgr.deliver <- msg
	return sessionId
}
