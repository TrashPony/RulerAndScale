package transport_data

import (
	"errors"
	"github.com/jacobsa/go-serial/serial"
	"io"
)

type Port struct {
	Name       string
	Config     *serial.OpenOptions
	Connection io.ReadWriteCloser
	Init       bool
}

func (p *Port) Reconnect(countRead int) error {

	if p.Connection == nil {
		return errors.New("no port")
	}

	p.Connection.Close()
	p.Config.MinimumReadSize = uint(countRead)
	var err error
	p.Connection, err = serial.Open(*p.Config)
	if err != nil {
		println("serial.Open: %v", err.Error())
		return err
	}

	return nil
}

func (p *Port) SendBytes(command []byte, countRead int) error {

	_, err := p.Connection.Write(command)
	if err != nil {
		println("ошибка записи" + err.Error())
		return err
	}

	return nil
}

func (p *Port) ReadBytes(countRead int) ([]byte, error) {

	data := make([]byte, countRead)

	n, err := p.Connection.Read(data)
	if err != nil {
		println("ошибка чтения: " + err.Error())
		return nil, err
	}

	// countRead говорит сколько должно быть считано байт, если не удалось то это считается ошибкой
	// нельзя использовать MinimumReadSize из конфига порта, т.к. если устройство не ответит это будет дедлок
	if n == countRead {
		return data, nil
	} else {
		return nil, errors.New("wrong_data")
	}
}
