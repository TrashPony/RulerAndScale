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

	p.SendBytes([]byte{0x4A}, countRead)
	data, err := p.ReadBytes(countRead)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Port) SendRulerCommand(command []byte, countRead int) ([]byte, error) {
	//
	//if p.commandID > 200 {
	//	p.commandID = 1
	//}
	//p.commandID++
	//
	//err := p.Reconnect(0)
	//if err != nil {
	//	return nil, err
	//}

	// запрос габаритов коробки
	err := p.SendBytes(command, countRead)
	if err != nil {
		return nil, err
	}

	data, err := p.ReadBytes(countRead)
	if err != nil {
		return nil, err
	}
	if data != nil && len(data) > 0 {
		return data, nil
	}

	return nil, errors.New("wrong_data")
}
