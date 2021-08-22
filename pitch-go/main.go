package main

import (
	"github.com/mickmister/pitch-go/server"
	"github.com/mickmister/pitch-go/view"
	"github.com/mrnikho/yingo"
	"github.com/pkg/errors"

	"fmt"

	"github.com/gordonklaus/portaudio"
)

//type Mic struct {
//chunkSize int
//bufferSize int
//bufferoptimize bool
//}

func main() {
	go func() {
		listenToGuitar()
	}()

	server.RunServer()
}

func listenToGuitar() {
	chunkSize := 2048
	bufferIncrement := 100

	pitchChannel := make(chan float32)
	MicInput(chunkSize, bufferIncrement, &pitchChannel)

	var prevPitch float32 = -1.0
	for pitch := range pitchChannel {
		// if true {continue}
		if pitch == -1.0 {
			continue
		}
		if pitch == prevPitch {
			continue
		}
		fmt.Println(pitch)
		prevPitch = pitch
		note := view.RenderFrequency(pitch)

		// go func(pitch float32) {
		// time.Sleep(500 * time.Millisecond)
		// if pitch == prevPitch {
		server.SendMessage(note)
		// }
		// }(pitch)
	}
}

func MicInput(chunkSize int, bufferincrement int, pch *chan float32) {
	//yin variables
	//m = Mic{chunkSize: chnkSze, bufferSize: 100, bufferoptimize: bfropt }
	threshold := float32(0.05)
	bufferSize := 100

	//var pitch float32

	// var bufferincrement int

	//pch := make(chan float32)

	//if bfropt {
	//	bufferincrement = 1
	//} else {
	//	bufferincrement = 100
	//}

	fmt.Println("Enter loop")
	go func() {

		portaudio.Initialize()
		defer portaudio.Terminate()

		input := make([]int16, chunkSize)
		stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(input), input)
		chkErr(err, "OpenDefaultStream")
		defer stream.Close()

		chkErr(stream.Start(), "Start")
		for {
			yin := yingo.Yin{}
			// fmt.Println("New Chunk")
			chkErr(stream.Read(), "Read")
			// fmt.Println(input)
			figurePitch(&yin, bufferincrement, input, pch, bufferSize, threshold)
		}
	}()
}

func figurePitch(yin *yingo.Yin, bufferincrement int, input []int16, pch *chan float32, bufferSize int, threshold float32) {
	var pitch float32
	// fmt.Println("Processing")
	for pitch < 10 {
		//fmt.Println(bufferSize)
		if bufferSize >= len(input) {
			// fmt.Println("Break")
			pitch = -1
			break
		}
		yin.YinInit(bufferSize, threshold)
		pitch = yin.GetPitch(&input)
		bufferSize += bufferincrement
	}
	//fmt.Println(bufferSize)
	// fmt.Println(pitch)
	*pch <- pitch
}

func chkErr(err error, message string) {
	if err != nil {
		panic(errors.Wrap(err, message))
	}
}
