package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/go-vgo/robotgo"
	"github.com/otiai10/gosseract"
	hook "github.com/robotn/gohook"
)

var (
	startX  int
	startY  int
	finishX int
	finishY int
)

func bitmap() {
	println(startX, startY, finishX, finishY)

	if finishY-startY < 0 {
		finishY, startY = startY, finishY
	}
	if finishX-startX < 0 {
		finishX, startX = startX, finishX
	}

	robotgo.SaveCapture("1.png", startX, startY, finishX-startX, finishY-startY)
}

func addMouse() {

	robotgo.EventHook(hook.KeyUp, []string{}, func(e hook.Event) {
		fmt.Println("ctrl up")
		if e.Button == 0 {
			finishX, finishY = robotgo.GetMousePos()
			//finishX, finishY = int(e.X), int(e.Y)
			bitmap()

			client := gosseract.NewClient()
			defer client.Close()
			_ = client.SetImage("1.png")
			text, _ := client.Text()
			_ = clipboard.WriteAll(text)
			fmt.Println(text)

			err := beeep.Notify("OCR result copied to the clipboard", text, "")
			if err != nil {
				panic(err)
			}

			robotgo.EventEnd()

		}
	})
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("A tray application for snipping screen and extracting text value from it to the clipboard, App")
	systray.SetTooltip("Snip OCR")

	mUrl := systray.AddMenuItem("Take Snip", "take snip")
	mQuit := systray.AddMenuItem("Quit", "quit")

	for {
		select {
		case <-mUrl.ClickedCh:

			println("select screen part for snip")
			ctrl := robotgo.AddEvent("ctrl")
			if ctrl == true {
				fmt.Println("ctrl down")
				startX, startY = robotgo.GetMousePos()
				addMouse()
			}

		case <-mQuit.ClickedCh:
			systray.Quit()
			fmt.Println("Quit2 now...")
			return
		}
	}
}

func onExit() {
	println("On Exit")
	// clean up here

}

func main() {
	systray.Run(onReady, onExit)
}
