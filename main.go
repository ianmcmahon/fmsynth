package main

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
	"github.com/ianmcmahon/fmsynth/audio"
	"github.com/ianmcmahon/fmsynth/ui"
	"github.com/rakyll/portmidi"
	wde "github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/cocoa"
)

type midiPort struct {
	name string
	in   portmidi.DeviceID
	out  portmidi.DeviceID
}

func midiPorts() map[string]*midiPort {
	m := make(map[string]*midiPort, 0)
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))

		if _, ok := m[info.Name]; !ok {
			m[info.Name] = &midiPort{name: info.Name}
		}
		if info.IsInputAvailable {
			m[info.Name].in = portmidi.DeviceID(i)
		}
		if info.IsOutputAvailable {
			m[info.Name].out = portmidi.DeviceID(i)
		}
	}

	return m
}

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	portmidi.Initialize()
	defer portmidi.Terminate()

	/*
		ports := midiPorts()
		Automapper(ports["ReMOTE ZeRO SL Port 1"], ports["ReMOTE ZeRO SL Port 3"])

		for {
		}

		os.Exit(1)
	*/

	var ch <-chan portmidi.Event

	device := portmidi.DefaultInputDeviceID()
	in, err := portmidi.NewInputStream(device, 1024)
	if err != nil {
		fmt.Printf("error opening midi: %v\n", err)
	} else {
		defer in.Close()

		in.SetChannelMask(1)

		ch = in.Listen()
	}

	engine := audio.NewEngine(ch)

	go ui.Start()

	wde.Run()

	engine.Stop()
	fmt.Printf("exiting\n")
}
