package InputData

import (
	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
	"log"
	"runtime"
	"sync"
	"time"
)

func PrintResult(data string) {
	keyMX.Lock()
	defer keyMX.Unlock()

	println(data)

	toClipBoard(data)
	toClipBoard("_ESC_Save")
}

func toClipBoard(data string) {
	err := clipboard.WriteAll(data)

	if err != nil {
		log.Fatal(err)
	}

	pressCtrlV()
	time.Sleep(time.Millisecond * 1000)
}

var kb = initKeyBD()
var keyMX sync.Mutex

func initKeyBD() keybd_event.KeyBonding {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2000 * time.Millisecond)
	}

	return kb
}

func pressCtrlV() {
	kb.Clear()

	// CTRL-V
	kb.SetKeys(keybd_event.VK_V)
	kb.HasCTRL(true)
	err := kb.Launching()
	if err != nil {
		panic(err)
	}

	kb.Clear()

	//Enter
	kb.SetKeys(keybd_event.VK_ENTER)
	err = kb.Launching()
	if err != nil {
		panic(err)
	}

	kb.Clear()
}
