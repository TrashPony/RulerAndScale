package websocket

import (
	"github.com/TrashPony/RulerAndScale/ParseData"
	"github.com/TrashPony/RulerAndScale/TransportData"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var UsersWs = make(map[*websocket.Conn]bool)
var sendPipe = make(chan Message)
var mutex = &sync.Mutex{}

type Message struct {
	Event         string        `json:"event"`
	ScalePlatform ScalePlatform `json:"scale_platform"`
	RulerOption   RulerOption   `json:"ruler_option"`
	Indication    Indication    `json:"indication"`
	Count         int           `json:"count"`
}

type ScalePlatform struct {
	Length int `json:"height"`
	Width  int `json:"width"`
}

type RulerOption struct {
	TopMax    int `json:"top_max"`
	WidthMax  int `json:"width_max"`
	LengthMax int `json:"length_max"`
}

type Indication struct {
	Left      int `json:"left"`
	Right     int `json:"right"`
	Top       int `json:"top"`
	Back      int `json:"back"`
	WidthBox  int `json:"width_box"`
	HeightBox int `json:"height_box"`
	LengthBox int `json:"length_box"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}

	ws, err := upgrader.Upgrade(w, r, nil) // запрос GET для перехода на протокол websocket
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	UsersWs[ws] = true

	go Reader(ws)
}

func Reader(ws *websocket.Conn) {

	//TODO если разная платформа весов)
	scalePlatform := ScalePlatform{Length: 50, Width: 40}

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil { // Если есть ошибка при чтение из сокета вероятно клиент отключился, удаляем его сессию
			delete(UsersWs, ws)
			ws.Close()
			return
		}

		if msg.Event == "Debug" {
			rulerPort := TransportData.Ports.GetPort("ruler")

			if rulerPort != nil {
				rulerResponse, _ := rulerPort.SendRulerCommand([]byte{0x89}, 41)
				if rulerResponse == nil && err.Error() != "wrong_data" {
					println("Линейка отвалилась")
					TransportData.Ports.ResetPort("ruler")
					continue
				} else {
					if err != nil {
						continue
					}
				}

				left, right, top, back, widthMax, heightMax, lengthMax, width, height, length := ParseData.ParseRulerIndicationData(rulerResponse, []byte{0x89})

				sendPipe <- Message{
					Event:         "Debug",
					ScalePlatform: scalePlatform,
					RulerOption:   RulerOption{WidthMax: widthMax, TopMax: heightMax, LengthMax: lengthMax},
					Indication:    Indication{Left: left, Right: right, Top: top, Back: back, WidthBox: width, HeightBox: height, LengthBox: length},
				}
			} else {
				// TODO лийнека не подключена
			}
			continue
		}

		if msg.Event == "SetTop" {
			rulerPort := TransportData.Ports.GetPort("ruler")
			rulerPort.SendBytes([]byte{0x90, byte(msg.Count)}, 0, 0)
		}

		if msg.Event == "SetWidth" {
			rulerPort := TransportData.Ports.GetPort("ruler")
			rulerPort.SendBytes([]byte{0x91, byte(msg.Count)}, 0, 0)
		}

		if msg.Event == "SetLength" {
			rulerPort := TransportData.Ports.GetPort("ruler")
			rulerPort.SendBytes([]byte{0x92, byte(msg.Count)}, 0, 0)
		}
	}
}

func Sender() {
	for {
		msg := <-sendPipe

		mutex.Lock()
		for ws, _ := range UsersWs {

			err := ws.WriteJSON(msg)

			if err != nil {
				delete(UsersWs, ws)
				ws.Close()
			}
		}
		mutex.Unlock()
	}
}
