package TransportData

import (
	"errors"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"sync"
)

type Port struct {
	Name       string
	Config     *serial.OpenOptions
	Connection io.ReadWriteCloser
	commandID  byte
	mx         sync.Mutex
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

func (p *Port) SendBytes(command []byte, countRead int, id byte) {

	if id > 0 {
		command = append(command, id)
	}

	_, err := p.Connection.Write(command)
	if err != nil {
		println("ошибка записи" + err.Error())
	}

	return
}

func (p *Port) ReadBytes(countRead int) ([]byte, error) {

	data := make([]byte, countRead)

	n, err := p.Connection.Read(data)
	if err != nil {
		println("ошибка чтения: " + err.Error())
	}

	if n == countRead {
		return data, nil
	} else {
		return nil, errors.New("wrong")
	}
}
