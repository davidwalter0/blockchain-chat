package server

import (
	"github.com/davidwalter0/blockchain-chat/db"
	"github.com/davidwalter0/blockchain-chat/message"
)

type WebSocketMessage struct {
	Type     string
	Messages []ChatMessage
	Key      string
	Contacts []string

	// для добавление новой пары ключей
	PublicKey  string
	PrivateKey string
	Passphrase string
}

type ChatMessage struct {
	Receiver     string // хеш публичного ключа получателя или сам публичный ключ, если NewPublicKey == true
	Sender       string // хеш публичного ключа отправителья
	Text         string
	Time         int64
	NewPublicKey bool // true, если в получателя нет в контактах
}

func (msg ChatMessage) addNewPublicKeyToDb() error {
	return db.AddContacts([]*message.KeyPair{{
		PublicKey:  []byte(msg.Receiver),
		PrivateKey: []byte{},
		Passphrase: []byte{},
	}})
}
