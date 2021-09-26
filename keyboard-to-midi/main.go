package main

import (
	"encoding/json"
	"fmt"
	"syscall"
	"time"

	// "github.com/MarinX/keylogger"
	"github.com/MarinX/keylogger"
	"github.com/mickmister/keyboard-to-midi/outputs"
	"gopkg.in/music-theory.v0/chord"
	"gopkg.in/music-theory.v0/note"
)

func makeInputEvent(key string, press bool) keylogger.InputEvent {
	eventType := keylogger.KeyPress
	if !press {
		eventType = keylogger.KeyRelease
	}

	return keylogger.InputEvent{
		Time:  syscall.NsecToTimeval(int64(time.Now().Nanosecond())), // syscall.Timeval
		Type:  keylogger.EvKey,
		Code:  inversedKeyCodeMap[key],
		Value: int32(eventType),
	}
}

func main() {
	c := outputs.QueuedMidiChord{
		Notes: []outputs.Note{
			{
				MidiNumber: 20,
			},
		},
		Chord: chord.Chord{
			Root:      note.A,
			AdjSymbol: note.Sharp,
			Tones: map[chord.Interval]note.Class{
				chord.I1: note.A,
			},
		},
	}

	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println(string(b))
	fmt.Println(c.GetMidiNotes())
}

func main3() {
	// InputEvent is the keyboard event structure itself

	kw := NewKeyWatcher()

	ch := kw.ReadKeyPresses()
	go func() {
		for {
			code := <-ch
			fmt.Println(code)
		}
	}()

	var p = func(key string, pressed bool) {
		kw.ProcessKeyEvent(makeInputEvent(key, pressed))

		// processKeyEvent("fake device", makeInputEvent(key, pressed))
	}

	p("A", true)
	p("A", true)
	p("C", true)
	p("A", false)
	p("B", true)
	p("B", false)
}

func main2() {
	// keyboardNames := keylogger.FindAllKeyboardDevices()

	// fmt.Println(keyboardNames)

	devices := []string{
		"/dev/input/event0",
		"/dev/input/event1",
		"/dev/input/event2",
		"/dev/input/event3",
	}

	devices = []string{
		"/dev/input/by-id/usb-04f3_0103-event-kbd",
	}

	for _, deviceName := range devices {
		keylog, err := keylogger.New(deviceName)
		if err != nil {
			panic(err)
		}

		// fmt.Println(keylog)
		ch := keylog.Read()

		fmt.Println("reading")
		go func(deviceName string) {
			for {
				k := <-ch
				processKeyEvent(deviceName, k)
			}
		}(deviceName)
	}

	end := make(chan bool)
	<-end

	// event3
}

func processKeyEvent(deviceName string, k keylogger.InputEvent) {
	out := [][]interface{}{
		{
			"Code",
			k.Code,
		},
		{
			"Time",
			k.Time,
		},
		{
			"Type",
			k.Type,
		},
		{
			"Value",
			k.Value,
		},
		{
			"KeyPress()",
			k.KeyPress(),
		},
		{
			"KeyRelease()",
			k.KeyRelease(),
		},
		{
			"KeyString()",
			k.KeyString(),
		},
	}

	b, _ := json.MarshalIndent(out, "", "  ")

	if false {
		fmt.Println(deviceName, string(b))
	}

	if k.KeyString() == "" {
		return
	}

	if k.KeyPress() {
		fmt.Println("Pressed", k.KeyString())
	}

	if k.KeyRelease() {
		fmt.Println("Released", k.KeyString())
	}
}
