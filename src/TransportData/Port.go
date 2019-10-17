package TransportData

import (
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

func (p *Port) SendBytes(command []byte, countRead int, id byte) (data []byte) {

	if id > 0 {
		command = append(command, id)
	}

	_, err := p.Connection.Write(command)
	if err != nil {
		println("ошибка записи" + err.Error())
		return data
	}

	//time.Sleep(150 * time.Millisecond)
	data = make([]byte, countRead)

	_, err = p.Connection.Read(data)
	if err != nil {
		println("ошибка чтения: " + err.Error())
	}

	return
}
