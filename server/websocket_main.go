package server

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"github.com/poslegm/blockchain-chat/network"
	"github.com/poslegm/blockchain-chat/db"
	"github.com/poslegm/blockchain-chat/message"
)

var WebSocketQueue = make(chan WebSocketMessage)

func receive(ws *websocket.Conn, handle func (m WebSocketMessage)) {
	// чтение не должно прекращаться
	ws.SetReadDeadline(time.Time{})

	for {
		msg := WebSocketMessage{}

		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSockets.receive: " + err.Error())
			break
		}

		handle(msg)
	}
	ws.Close()
}

// выбирает ответ на сообщение в зависимости от типа и кладёт его в очередь сообщений
func switchTypes(msg WebSocketMessage) {
	fmt.Println("Websockets.swithTypes: ", msg)

	switch msg.Type {
	case "GetMessages":
		networkMessages, err := db.GetAllMessages()
		if err != nil {
			fmt.Println("Websockets.switchTypes: ", err.Error())
			return
		}
		fmt.Println(networkMessages)
		chatMessages := make([]ChatMessage, 0)
		for _, networkMsg := range networkMessages {
			textMsg, err := networkMsg.AsTextMessage()
			if err != nil {
				if err.Error() != "unsuitable-pair" {
					fmt.Println("Websockets.switchTypes: ", err.Error())
				}
				continue
			}

			chatMessages = append(chatMessages, ChatMessage{
				textMsg.Receiver, textMsg.Sender, textMsg.Text, false,
			})
		}

		WebSocketQueue <- WebSocketMessage{Type:"AllMessages", Messages:chatMessages}
	case "SendMessage":
		if len(msg.Messages) != 1 {
			fmt.Printf("WebSocket.switchTypes: incorrect message - %#v\n", msg)
			return
		}
		chatMsg := msg.Messages[0]

		if chatMsg.NewPublicKey {
			err := chatMsg.addNewPublicKeyToDb()
			if err != nil {
				fmt.Println("Websocket.swithTypes: can't add new public key ", err.Error())
			}
		}

		kp, err := db.GetKeyByAddress(chatMsg.Receiver)
		if chatMsg.NewPublicKey {
			kp = &message.KeyPair{[]byte(chatMsg.Receiver), []byte{}, []byte{}}
		}
		if err != nil {
			fmt.Println("WebSockets.swithTypes: can't get kp from db ", err.Error())
			return
		} else if kp == nil {
			fmt.Println("WebSockets.swithTypes: there is no kp in db")
			return
		}

		networkMsg, err := network.CreateTextNetworkMessage(
			kp.GetBase58Address(),
			chatMsg.Sender,
			chatMsg.Text,
			kp.PublicKey,
		)

		if err != nil {
			fmt.Println("Websockets.switchTypes: can't send message ", err.Error())
		} else {
			go network.CurrentNetworkUser.SendMessage(networkMsg)
		}

		if chatMsg.NewPublicKey {
			WebSocketQueue <- WebSocketMessage{
				Type: "NewKeyHash",
				Key: kp.GetBase58Address(),
				Messages: []ChatMessage{chatMsg },
			}
		}
	case "GetMyKey":
		publicKey, err := db.GetPublicKey()
		if err != nil {
			fmt.Println("Websockets.switchTypes: can't send public key ", err.Error())
		} else {
			WebSocketQueue <- WebSocketMessage{Type:"Key", Key:string(publicKey)}
		}
	}
}

func handleMessagesQueue(ws *websocket.Conn) {
	for {
		select {
		case msg := <-WebSocketQueue:
			ws.SetWriteDeadline(time.Now().Add(30 * time.Second))
			err := ws.WriteJSON(&msg)
			if err != nil {
				fmt.Println("WebSocket.handleMessagesQueue: " + err.Error())
				// если сообщение не удалось отправить, оно добавляется обратно в конец очереди
				WebSocketQueue <- msg
			} else {
				fmt.Println("WebSocket.handleMessagesQueue: sended ", msg)
			}
		case msg := <-network.CurrentNetworkUser.IncomingMessages:
			textMsg, err := msg.AsTextMessage()
			if err != nil {
				if err.Error() != "unsuitable-pair" {
					fmt.Println("Websockts.switchTypes: ", err.Error())
				}
				continue
			} else {
				WebSocketQueue <- WebSocketMessage{
					Type:"NewMessage",
					Messages:[]ChatMessage{{
						Receiver: textMsg.Receiver,
						Sender: textMsg.Sender,
						Text: textMsg.Text,
						NewPublicKey: false,
					}},
				}
			}
		}
	}
}

func createConnection(ws *websocket.Conn, handle func (m WebSocketMessage)) {
	go receive(ws, handle)
	go handleMessagesQueue(ws)
}

func createWSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}

		createConnection(ws, switchTypes)
	}
}