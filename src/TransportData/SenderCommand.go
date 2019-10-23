package TransportData

import (
	"errors"
)

func (p *Port) SendScaleCommand() ([]byte, error) {

	countRead := 5

	err := p.Reconnect(0)
	if err != nil {
		return nil, err
	}

	p.SendBytes([]byte{0x4A}, countRead, 0)
	data, err := p.ReadBytes(countRead)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Port) SendRulerCommand(command []byte, countRead int) ([]byte, error) {

	if p.commandID > 200 {
		p.commandID = 1
	}
	p.commandID++

	err := p.Reconnect(countRead)
	if err != nil {
		return nil, err
	}

	// запрос габаритов коробки
	p.SendBytes(command, countRead, p.commandID)
	data, _ := p.ReadBytes(countRead)
	if data != nil && len(data) > 0 && data[0] == p.commandID {
		return data, nil
	} else {

		// иногда сериал порт посылает прошлые данные и от них надо избавится
		err := p.Reconnect(0)
		if err != nil {
			return nil, err
		}

		data, _ = p.ReadBytes(countRead)
	}

	return nil, errors.New("wrong_data")
}
