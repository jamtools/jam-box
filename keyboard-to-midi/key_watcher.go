package main

import (
	"time"

	"github.com/MarinX/keylogger"
)

type keyWatcher struct {
	keys []string
	ch   chan uint16
}

var _ KeyWatcher = &keyWatcher{}

type KeyWatcher interface {
	ProcessKeyEvent(ke keylogger.InputEvent)
	GetCurrentlyHeldKeys() []string
	ReadKeyPresses() chan uint16
}

func NewKeyWatcher() KeyWatcher {
	return &keyWatcher{
		keys: []string{},
		ch:   make(chan uint16),
	}
}

func (kw *keyWatcher) ReadKeyPresses() chan uint16 {
	return kw.ch
}

func (kw *keyWatcher) GetCurrentlyHeldKeys() []string {
	return kw.keys
}

func (kw *keyWatcher) ProcessKeyEvent(ke keylogger.InputEvent) {
	if !ke.KeyPress() && !ke.KeyRelease() {
		return
	}

	if ke.KeyString() == "" {
		return
	}

	if ke.KeyPress() {
		found := false
		for _, storedKey := range kw.keys {
			if storedKey == ke.KeyString() {
				found = true
				break
			}
		}

		if !found {
			kw.keys = append(kw.keys, ke.KeyString())
			kw.ch <- ke.Code
		}
	} else {
		found := false
		toSave := []string{}
		for _, storedKey := range kw.keys {
			if storedKey == ke.KeyString() {
				found = true
			} else {
				toSave = append(toSave, storedKey)
			}
		}

		if found {
			kw.keys = toSave
		}
	}

	// if ke.KeyPress() {
	// 	fmt.Println("Pressed", ke.KeyString())
	// }

	// if ke.KeyRelease() {
	// 	fmt.Println("Released", ke.KeyString())
	// }

	// fmt.Println(strings.Join(kw.GetCurrentlyHeldKeys(), " "))
	time.Sleep(time.Millisecond)
}
