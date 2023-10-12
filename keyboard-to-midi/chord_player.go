package main

import (
	"github.com/mickmister/keyboard-to-midi/outputs"
)

type ChordCommander struct {
	currentIndex int
	MidiOut      *outputs.MidiDevice
	Chords       []outputs.QueuedMidiChord
}

func (cc ChordCommander) PlayNextChord() {
	nextIndex := (cc.currentIndex + 1) % len(cc.Chords)
	cc.currentIndex = nextIndex
	chord := cc.Chords[nextIndex]
	cc.PlayChord(chord)
}

func (cc ChordCommander) PlayChord(chord outputs.QueuedMidiChord) {
	notes := chord.GetMidiNotes()
	cc.MidiOut.SendNoteOnEvent(notes[0])
	// cc.MidiOut.SendNoteOnEvent(chord.GetMidiNotes())
}
