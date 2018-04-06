package TransportData

import (
	"github.com/tarm/serial"
	"time"
)

type ScaleResponse struct {
	ReadyAndDiscreteness []byte
	Weight               []byte
}

type RulerResponse struct {
	Width  []byte
	Height []byte
	Length []byte
}

func SendScaleCommand(port *serial.Port) (*ScaleResponse) {

	var response ScaleResponse
	countRead := 2

	// не/готовность 0/128 и дискретность 0х00-1г,0х01-0.1г,0х04-0.01кг,0х05-0.1кг
	response.ReadyAndDiscreteness = SendCommand([]byte{0x48}, port, countRead)

	//вес в виде 2х байтов n х n
	response.Weight = SendCommand([]byte{0x45}, port, countRead)

	if response.ReadyAndDiscreteness != nil && response.Weight != nil {
		return &response
	} else {
		return nil
	}
}

func SendRulerCommand(port *serial.Port) (*RulerResponse) {

	var response RulerResponse
	countRead := 5

	// запрос ширины коробки
	response.Width = SendCommand([]byte{0x88}, port, countRead)

	// запрос высоты коробки
	response.Height = SendCommand([]byte{0x99}, port, countRead)

	// запрос длинны коробки
	response.Length = SendCommand([]byte{0x77}, port, countRead)

	if response.Width != nil && response.Height != nil && response.Length != nil {
		return &response
	} else {
		return nil
	}
}

func SendCommand(command []byte, port *serial.Port, countRead int) ([]byte) {
	for {
		_, err := port.Write(command)
		if err != nil {
			port.Close()
			return nil
		}

		time.Sleep(time.Millisecond * 100)

		data := make([]byte, countRead)
		n, err := port.Read(data)
		if err != nil {
			port.Close()
			return nil
		}

		if n == countRead {
			return data
		}
	}
}
