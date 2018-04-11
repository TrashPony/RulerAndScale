package TransportData

import (
	"github.com/tarm/serial"
	"time"
)

type Port struct {
	Name string
	Config *serial.Config
	Connection *serial.Port
}

func (p *Port) Connect () (connect *serial.Port)   {
	connect, err := serial.OpenPort(p.Config)
	if err != nil {
		return nil
	} else {
		p.Connection = connect
		return connect
	}
}

func (p *Port) SendBytes(command []byte, countRead int) (data []byte)  {

	errorCount := 0

	for {

		if errorCount >= 5 {
			return nil
		}

		if p.Connection == nil {
			println("ошибка повторного подключения")
			return nil
		}

		p.Connection.Flush()

		_, err := p.Connection.Write(command)
		if err != nil {
			println("ошибка записи" + err.Error())
			errorCount++
			p.Connection.Close()
			p.Connection = p.Connect()
			continue
		}

		time.Sleep(time.Millisecond * 75) // без этой задержки байты не будут успевать приниматься

		data = make([]byte, countRead)

		n, err := p.Connection.Read(data)
		if err != nil {
			println("ошибка чтения: " + err.Error())
			p.Connection.Close()
			errorCount++
			p.Connection = p.Connect()
			continue
		}

		if n == countRead {
			return data
		}
	}
}