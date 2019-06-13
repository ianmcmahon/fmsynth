package main

import (
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/ianmcmahon/fmsynth/audio"
	"github.com/rakyll/portmidi"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	portmidi.Initialize()
	defer portmidi.Terminate()

	device := portmidi.DefaultInputDeviceID()
	in, err := portmidi.NewInputStream(device, 1024)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	in.SetChannelMask(1)

	ch := in.Listen()

	engine := audio.NewEngine(ch)

	for {
		time.Sleep(1 * time.Second)
	}
	engine.Stop()
	fmt.Printf("exiting\n")
}
