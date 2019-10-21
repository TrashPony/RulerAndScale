package TransportData

import (
	"errors"
	"github.com/jacobsa/go-serial/serial"
)

type ScaleResponse struct {
	ReadyAndDiscreteness []byte
	Weight               []byte
}

func SendScaleCommand(port *Port) *ScaleResponse {

	var response ScaleResponse
	countRead := 2

	// не/готовность 0/128 и дискретность 0х00-1г,0х01-0.1г,0х04-0.01кг,0х05-0.1кг
	response.ReadyAndDiscreteness = port.SendBytes([]byte{0x48}, countRead, 0)

	//вес в виде 2х байтов n х n
	response.Weight = port.SendBytes([]byte{0x45}, countRead, 0)

	if response.ReadyAndDiscreteness != nil && response.Weight != nil {
		return &response
	} else {
		return nil
	}
}

func (p *Port) SendRulerCommand(command []byte, countRead int) ([]byte, error) {

	p.mx.Lock()
	defer p.mx.Unlock()

	if p.commandID > 200 {
		p.commandID = 1
	}
	p.commandID++

	p.Connection.Close()
	p.Config.MinimumReadSize = uint(countRead)

	var err error
	p.Connection, err = serial.Open(*p.Config)
	if err != nil {
		println("serial.Open: %v", err.Error())
		return nil, err
	}

	// запрос габаритов коробки
	data := p.SendBytes(command, countRead, p.commandID)

	if data != nil && data[0] == p.commandID {
		return data, nil
	} else {
		// иногда сериал порт посылает прошлые команды
		countTryRead := 5
		for countTryRead == 0 {

			data = make([]byte, countRead)
			_, err = p.Connection.Read(data)
			if err != nil {
				println("ошибка чтения: " + err.Error())
			}

			if data[0] == p.commandID {
				return data, nil
			}
			countTryRead--
		}
	}

	return nil, errors.New("wrong_data")
}
