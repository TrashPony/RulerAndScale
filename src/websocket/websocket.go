package websocket

import (
	"github.com/TrashPony/RulerAndScale/ParseData"
	"github.com/TrashPony/RulerAndScale/TransportData"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var UsersWs = make(map[*websocket.Conn]bool)
var sendPipe = make(chan Message)
var mutex = &sync.Mutex{}

type Message struct {
	Event         string              `json:"event"`
	ScalePlatform ScalePlatform       `json:"scale_platform"`
	RulerOption   RulerOption         `json:"ruler_option"`
	Indication    Indication          `json:"indication"`
	Count         int                 `json:"count"`
	ScalePort     *TransportData.Port `json:"scale_port"`
	RulerPort     *TransportData.Port `json:"ruler_port"`
}

type ScalePlatform struct {
	Length int `json:"height"`
	Width  int `json:"width"`
}

type RulerOption struct {
	TopMax     int  `json:"top_max"`
	WidthMax   int  `json:"width_max"`
	LengthMax  int  `json:"length_max"`
	OnlyWeight bool `json:"only_weight"`
}

type Indication struct {
	Left      int `json:"left"`
	Right     int `json:"right"`
	Top       int `json:"top"`
	Back      int `json:"back"`
	WidthBox  int `json:"width_box"`
	HeightBox int `json:"height_box"`
	LengthBox int `json:"length_box"`
	Weight    int `json:"correct_weight"`
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

			time.Sleep(400 * time.Millisecond)
			rulerPort := TransportData.Ports.GetPort("ruler")
			scalePort := TransportData.Ports.GetPort("scale")

			var left, right, top, back, widthMax, heightMax, lengthMax, width, height, length, correctWeight int
			var onlyWeight bool

			if rulerPort != nil {
				rulerResponse, err := rulerPort.SendRulerCommand([]byte{0x89}, 41)
				if err != nil && err.Error() != "wrong_data" {
					TransportData.Ports.ResetPort("ruler")
				} else {
					if rulerResponse != nil {
						left, right, top, back, widthMax, heightMax, lengthMax, width, height, length, onlyWeight = ParseData.ParseRulerIndicationData(rulerResponse, []byte{0x89})
					}
				}
			}

			if scalePort != nil {
				scaleResponse, err := scalePort.SendScaleCommand()
				if err != nil && err.Error() != "wrong_data" {
					TransportData.Ports.ResetPort("scale")
				} else {
					if (err != nil && err.Error() == "wrong") || scaleResponse == nil {
						correctWeight = -1
					} else {
						correctWeight = int(ParseData.ParseScaleData(scaleResponse))
						if correctWeight == 0 {
							// иногда сериал порт посылает прошлые данные и от них надо избавится или смещает биты
							scalePort.Reconnect(0)
							scalePort.ReadBytes(5)
						}
					}
				}
			}

			sendPipe <- Message{
				Event:         "Debug",
				ScalePlatform: scalePlatform,
				RulerOption:   RulerOption{WidthMax: widthMax, TopMax: heightMax, LengthMax: lengthMax, OnlyWeight: onlyWeight},
				Indication:    Indication{Left: left, Right: right, Top: top, Back: back, WidthBox: width, HeightBox: height, LengthBox: length, Weight: correctWeight},
				RulerPort:     rulerPort,
				ScalePort:     scalePort,
			}

			continue
		}

		if msg.Event == "SetTop" {
			rulerPort := TransportData.Ports.GetPort("ruler")
			time.Sleep(500 * time.Millisecond)
			rulerPort.SendRulerCommand([]byte{0x90, byte(msg.Count)}, 0)
		}

		if msg.Event == "SetWidth" {
			rulerPort := TransportData.Ports.GetPort("ruler")
			time.Sleep(500 * time.Millisecond)
			rulerPort.SendRulerCommand([]byte{0x91, byte(msg.Count)}, 0)
		}

		if msg.Event == "SetLength" {
			rulerPort := TransportData.Ports.GetPort("ruler")
			time.Sleep(500 * time.Millisecond)
			rulerPort.SendRulerCommand([]byte{0x92, byte(msg.Count)}, 0)
		}

		if msg.Event == "ResetRuler" {
			TransportData.Ports.ResetPort("ruler")
		}

		if msg.Event == "ResetScale" {
			TransportData.Ports.ResetPort("scale")
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
