package state

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"secure-chat-client/env"
	"sync"
)

type WsSendNewMessage struct {
	ToUserID           uuid.UUID `json:"toUserId"`
	CipherText         string    `json:"cipherText"`
	Nonce              string    `json:"nonce"`
	SenderIdentityId   string    `json:"senderIdentityId"`
	ReceiverIdentityId string    `json:"ReceiverIdentityId"`
}

type WsNewMessageRecieved struct {
	FromUserID         uuid.UUID `json:"fromUserId"`
	CipherText         string    `json:"cipherText"`
	SenderIdentityId   string    `json:"senderIdentityId"`
	ReceiverIdentityId string    `json:"ReceiverIdentityId"`
}

var (
	wsConn    *websocket.Conn
	sendChan  chan []byte
	recvChan  chan []byte
	once      sync.Once
	closeChan chan struct{}
)

func ConnectWebSocket(token string) error {
	var err error
	//we should not use once.Do since, we can logout and log back in!!
	once.Do(func() {
		header := http.Header{}
		header.Set("Authorization", "Bearer "+token)

		wsConn, _, err = websocket.DefaultDialer.Dial(env.Opts.WebSocketUrl, header)
		if err != nil {
			return
		}

		sendChan = make(chan []byte, 4096)
		recvChan = make(chan []byte, 4096)
		closeChan = make(chan struct{})

		go readLoop()
		go writeLoop()
	})
	return err
}

func SendMessage(msg []byte) {
	select {
	case sendChan <- msg:
	default:
		slog.Info("Send channel full, dropping message")
	}
}

func Subscribe() <-chan []byte {
	return recvChan
}

func readLoop() {
	defer wsConn.Close()
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			slog.Error("WebSocket read error:", err)
			close(closeChan)
			return
		}

		var m WsNewMessageRecieved
		unmarshalErr := json.Unmarshal(msg, &m)
		if unmarshalErr != nil {
			slog.Error("WebSocket unmarshal error:", unmarshalErr)
		} else {
			str, strErr := json.Marshal(m)
			if strErr != nil {
				slog.Error("WebSocket marshal error:", strErr)
			} else {
				slog.Info("WebSocket received:", string(str))
			}
		}

		select {
		case recvChan <- msg:
		default:
			slog.Warn("Receive channel full, dropping message")
		}
	}
}

func writeLoop() {
	defer wsConn.Close()
	for {
		select {
		case msg := <-sendChan:
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				slog.Error("WebSocket write error:", err)
				return
			}
		case <-closeChan:
			return
		}
	}
}
