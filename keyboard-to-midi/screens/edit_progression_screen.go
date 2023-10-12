package screens

import (
	"github.com/mickmister/keyboard-to-midi/music"
	"github.com/mickmister/keyboard-to-midi/outputs"
)

type editProgressionScreen struct {
	functions    []ScreenFunction
	heldDownKeys []int
}

func newEditProgressionScreen() Screen {
	s := &editProgressionScreen{heldDownKeys: []int{}}
	s.functions = []ScreenFunction{
		{
			Name: "Add Chord",
			Run:  s.addChord,
		},
		{
			Name: "Remove Chord",
			Run:  s.removeChord,
		},
		{
			Name: "Save Progression",
			Run:  s.saveProgression,
		},
	}

	return s
}

func (s *editProgressionScreen) addChord(state State) (State, error) {
	notes := state.MidiDevices[0].GetHeldDownKeys()
	chord := music.ConvertMidiNumbersToChord(notes)
	state.CurrentProgression.Chords = append(state.CurrentProgression.Chords, chord)

	return state, nil
}

func (s *editProgressionScreen) saveProgression(state State) (State, error) {
	found := false
	current := state.CurrentProgression

	result := []outputs.ChordProgression{}

	if current.Name == "" {
		current.Name = current.String()
	}

	for _, prog := range state.Progressions {
		if prog.Name == current.Name {
			result = append(result, current)
			found = true
		} else {
			result = append(result, prog)
		}
	}

	if !found {
		result = append(result, current)
	}

	state.Progressions = result

	return state, nil
}

func (s *editProgressionScreen) removeChord(state State) (State, error) {
	state.CurrentProgression.Chords = state.CurrentProgression.Chords[0 : len(state.CurrentProgression.Chords)-1]
	return state, nil
}

func (s *editProgressionScreen) HandleMidiEvent(midiEvent outputs.MidiEvent, midiDevice *outputs.MidiDevice, stateIn State) State {
	_ = midiDevice.GetHeldDownKeys()
	return stateIn
}

func (s *editProgressionScreen) GetFunctions() []ScreenFunction {
	return s.functions
}
