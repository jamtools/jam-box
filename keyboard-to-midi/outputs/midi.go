package outputs

import (
	"strings"

	"gopkg.in/music-theory.v0/chord"
)

type MidiEvent struct {
	Type string `json:"type"`
	Note int    `json:"note"`
}

type MidiDevice struct {
	Name         string
	heldDownKeys []int
}

func (md *MidiDevice) SendNoteOnEvent(note int) {
}

func (md *MidiDevice) SendNoteOffEvent(note int) {
}

func (md *MidiDevice) GetHeldDownKeys() []int {
	return md.heldDownKeys
}

type Note struct {
	MidiNumber int
}

type ChordProgression struct {
	Name   string            `json:"name"`
	Chords []QueuedMidiChord `json:"chords"`
}

func (prog ChordProgression) String() string {
	if prog.Name == "" {
		return prog.Name
	}

	names := []string{}
	for _, chord := range prog.Chords {
		names = append(names, chord.Name)
	}

	return strings.Join(names, ", ")
}

type QueuedMidiChord struct {
	chord.Chord
	Name  string
	Notes []Note
}

func (c QueuedMidiChord) GetMidiNotes() []int {
	notes := []int{}
	for _, n := range c.Notes {
		notes = append(notes, n.MidiNumber)
	}
	return notes
}

func ConvertMidiNumberToChord(midiNumber int) QueuedMidiChord {

}
