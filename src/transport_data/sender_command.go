package transport_data

import (
	"errors"
)

func (p *Port) SendScaleCommand() ([]byte, error) {

	countRead := 5

	err := p.SendBytes([]byte{0x4A}, countRead)
	if err != nil {
		return nil, err
	}

	data, err := p.ReadBytes(countRead)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Port) SendRulerCommand(command []byte, countRead int) ([]byte, error) {

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
