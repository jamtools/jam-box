package screens

import (
	"github.com/mickmister/keyboard-to-midi/outputs"
)

type State struct {
	MidiDevices        []*outputs.MidiDevice
	CurrentProgression outputs.ChordProgression
	Progressions       []outputs.ChordProgression
	Screen             Screen
}

type ScreenFunction struct {
	Name string
	Run  func(state State) (State, error)
}

type Screen interface {
	GetFunctions() []ScreenFunction
}
