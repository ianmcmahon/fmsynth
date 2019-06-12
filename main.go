package main

import (
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/ianmcmahon/fmsynth/audio"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	engine := audio.NewEngine()

	time.Sleep(10 * time.Second)
	engine.Stop()
	fmt.Printf("exiting\n")
}
