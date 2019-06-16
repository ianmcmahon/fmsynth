package main

import (
	"fmt"

	"github.com/rakyll/portmidi"
)

var (
	automapOnline  = []byte{0xF0, 0x00, 0x20, 0x29, 0x03, 0x03, 0x11, 0x04, 0x03, 0x00, 0x01, 0x01, 0xF7}
	automapOffline = []byte{0xF0, 0x00, 0x20, 0x29, 0x03, 0x03, 0x11, 0x04, 0x03, 0x00, 0x01, 0x00, 0xF7}
	automapWrite   = []byte{0xF0, 0x00, 0x20, 0x29, 0x03, 0x03, 0x11, 0x04, 0x7F, 0x00, 0x02, 0x01}
)

type automapper struct {
	midi       *midiPort
	automap    *midiPort
	outStream  *portmidi.Stream
	eventsChan <-chan portmidi.Event
}

func Automapper(midi, automap *midiPort) *automapper {
	a := &automapper{
		midi:    midi,
		automap: automap,
	}
	go a.Start()
	return a
}

func writeScreen(row, pos byte, msg string, hostID byte) []byte {
	b := append(automapWrite, []byte{pos, row, 0x04}...)
	b[8] = hostID
	b = append(b, []byte(msg)...)
	b = append(b, 0xF7)

	fmt.Printf("%#v\n", b)
	return b
}

func (a *automapper) Start() {
	in, err := portmidi.NewInputStream(a.automap.in, 1024)
	if err != nil {
		panic(err)
	}
	a.eventsChan = in.Listen()

	a.outStream, err = portmidi.NewOutputStream(a.automap.out, 1024, 0)
	if err != nil {
		panic(err)
	}
	a.outStream.WriteSysExBytes(portmidi.Time(), automapOnline)

	a.outStream.WriteSysExBytes(portmidi.Time(), writeScreen(0, 0, "param 0", 37))
	a.outStream.WriteSysExBytes(portmidi.Time(), writeScreen(1, 0, "param 1", 37))
	a.outStream.WriteSysExBytes(portmidi.Time(), writeScreen(2, 3, "param 2", 37))

	for event := range a.eventsChan {
		fmt.Printf("automap event: %#v\n", event)
	}

}
