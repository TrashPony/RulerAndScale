package InputData

import (
	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
	"log"
	"runtime"
	"time"
)

func ToClipBoard(data string) {

	println(data)

	err := clipboard.WriteAll(data)

	if err != nil {
		log.Fatal(err)
	}

	pressCtrlV()
	pressEnter()
	time.Sleep(time.Millisecond * 400)
}

func pressCtrlV() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(200 * time.Millisecond)
	}

	//set keys
	kb.SetKeys(keybd_event.VK_V)
	kb.HasCTRL(true)

	//launch
	err = kb.Launching()
	if err != nil {
		panic(err)
	}
}

func pressEnter() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(200 * time.Millisecond)
	}

	//set keys
	kb.SetKeys(keybd_event.VK_ENTER)

	//launch
	err = kb.Launching()
	if err != nil {
		panic(err)
	}
}
