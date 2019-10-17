package websocket

import (
	"github.com/TrashPony/RulerAndScale/ParseData"
	"github.com/TrashPony/RulerAndScale/TransportData"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var usersWs = make(map[*websocket.Conn]bool)
var sendPipe = make(chan Message)
var mutex = &sync.Mutex{}

type Message struct {
	Event         string        `json:"event"`
	ScalePlatform ScalePlatform `json:"scale_platform"`
	RulerOption   RulerOption   `json:"ruler_option"`
	Indication    Indication    `json:"indication"`
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

	usersWs[ws] = true

	go Reader(ws)
}

func Reader(ws *websocket.Conn) {

	//TODO если разная платформа весов)
	scalePlatform := ScalePlatform{Length: 50, Width: 40}

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil { // Если есть ошибка при чтение из сокета вероятно клиент отключился, удаляем его сессию
			delete(usersWs, ws)
			ws.Close()
			return
		}

		if msg.Event == "Debug" {
			rulerPort := TransportData.Ports.GetPort("ruler")

			data := rulerPort.SendRulerCommand([]byte{0x89}, 13)
			widthMax, heightMax, lengthMax := ParseData.ParseRulerData(data, []byte{0x89})

			data = rulerPort.SendRulerCommand([]byte{0x80}, 17)
			left, right, top, back := ParseData.ParseRulerIndicationData(data, []byte{0x80})

			data = rulerPort.SendRulerCommand([]byte{0x88}, 13)
			width, height, length := ParseData.ParseRulerData(data, []byte{0x89})

			sendPipe <- Message{
				Event:         "Debug",
				ScalePlatform: scalePlatform,
				RulerOption:   RulerOption{WidthMax: widthMax, TopMax: heightMax, LengthMax: lengthMax},
				Indication:    Indication{Left: left, Right: right, Top: top, Back: back, WidthBox: width, HeightBox: height, LengthBox: length},
			}
		}
	}
}

func Sender() {
	for {
		msg := <-sendPipe

		mutex.Lock()
		for ws, _ := range usersWs {

			err := ws.WriteJSON(msg)

			if err != nil {
				delete(usersWs, ws)
				ws.Close()
			}
		}
		mutex.Unlock()
	}
}
