package websocket

import (
	"github.com/TrashPony/RulerAndScale/config"
	"github.com/TrashPony/RulerAndScale/parse_data"
	"github.com/TrashPony/RulerAndScale/transport_data"
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
	Event         string               `json:"event"`
	ScalePlatform ScalePlatform        `json:"scale_platform"`
	RulerOption   RulerOption          `json:"ruler_option"`
	Indication    Indication           `json:"indication"`
	Count         int                  `json:"count"`
	ScalePort     *transport_data.Port `json:"scale_port"`
	RulerPort     *transport_data.Port `json:"ruler_port"`
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

	//TODO есть разная платформа весов)
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
			rulerPort := transport_data.Ports.GetPort("ruler")
			scalePort := transport_data.Ports.GetPort("scale")

			var left, right, top, back, widthMax, heightMax, lengthMax, width, height, length, correctWeight int
			var onlyWeight bool

			if rulerPort != nil {
				rulerResponse, err := rulerPort.SendRulerCommand([]byte{0x89, 0x89}, 41)

				if err != nil && err.Error() != "wrong_data" {

					println("Линейка отвалилась")
					transport_data.Ports.ResetPort("ruler")

				} else {
					if rulerResponse != nil {
						left, right, top, back, widthMax, heightMax, lengthMax, width, height, length, onlyWeight = parse_data.ParseRulerIndicationData(rulerResponse)
					}
				}
			}

			if scalePort != nil {
				scaleResponse, err := scalePort.SendScaleCommand()
				if scaleResponse == nil && err.Error() != "wrong_data" {

					println("Весы отвалились")
					transport_data.Ports.ResetPort("scale")

				} else {
					if scaleResponse != nil {
						correctWeight = int(parse_data.ParseScaleData(scaleResponse))
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
			rulerPort := transport_data.Ports.GetPort("ruler")
			time.Sleep(500 * time.Millisecond)
			rulerPort.SendRulerCommand([]byte{0x90, byte(msg.Count)}, 0)

			_, width, length := config.GetConfig()
			config.WriteConfig(msg.Count, width, length)
		}

		if msg.Event == "SetWidth" {
			rulerPort := transport_data.Ports.GetPort("ruler")
			time.Sleep(500 * time.Millisecond)
			rulerPort.SendRulerCommand([]byte{0x91, byte(msg.Count)}, 0)

			top, _, length := config.GetConfig()
			config.WriteConfig(top, msg.Count, length)
		}

		if msg.Event == "SetLength" {
			rulerPort := transport_data.Ports.GetPort("ruler")
			time.Sleep(500 * time.Millisecond)
			rulerPort.SendRulerCommand([]byte{0x92, byte(msg.Count)}, 0)

			top, width, _ := config.GetConfig()
			config.WriteConfig(top, width, msg.Count)
		}

		if msg.Event == "ResetRuler" {
			transport_data.Ports.ResetPort("ruler")
		}

		if msg.Event == "ResetScale" {
			transport_data.Ports.ResetPort("scale")
		}

		if msg.Event == "Calibration" {
			rulerPort := transport_data.Ports.GetPort("ruler")
			time.Sleep(500 * time.Millisecond)
			rulerPort.SendRulerCommand([]byte{0x93, 0x93}, 0)
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
