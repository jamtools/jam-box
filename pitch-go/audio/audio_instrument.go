package audio

type AudioInstrument interface {
	Start()
	Stop()
	GetCurrentFrequency() (freq float32)
	RegisterFrequencyChange(func(freq float32))
}
