package main

import "math"

//  ██████╗ ███████╗ ██████╗██╗██╗     ██╗      █████╗ ████████╗ ██████╗ ██████╗ ███████╗
// ██╔═══██╗██╔════╝██╔════╝██║██║     ██║     ██╔══██╗╚══██╔══╝██╔═══██╗██╔══██╗██╔════╝
// ██║   ██║███████╗██║     ██║██║     ██║     ███████║   ██║   ██║   ██║██████╔╝███████╗
// ██║   ██║╚════██║██║     ██║██║     ██║     ██╔══██║   ██║   ██║   ██║██╔══██╗╚════██║
// ╚██████╔╝███████║╚██████╗██║███████╗███████╗██║  ██║   ██║   ╚██████╔╝██║  ██║███████║
//  ╚═════╝ ╚══════╝ ╚═════╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚══════╝

// Oscillator is a periodically repeating waveform
type Oscillator struct {
	ν       Hertz                 // Fundamental frequency
	Phase   Angle                 // Last known phase
	PhaseAt Seconds               // When that phase occurred
	Wave    func(a Angle) float64 // Function that describes wave shape
}

// Waveform is a function that encodes the shape of the cycles of a waveform in *angle*
type Waveform func(a Angle) float64

// Amplitude returns the strength of the waveform at a given *time* for a particular *oscillator*, which may change frequency
func (osc *Oscillator) Amplitude(t Seconds) float64 {
	dT := t - osc.PhaseAt
	dA := Angle(dT) * τ * Angle(osc.ν) // convert time to phase angle at new freq
	osc.Phase += dA
	osc.PhaseAt = t
	return osc.Wave(osc.Phase)
}

// NewSine returns a new sine wave oscillator
func NewSine(newν Hertz) *Oscillator {
	return &Oscillator{
		ν:       newν,
		Phase:   0,
		PhaseAt: 0,
		Wave: func(a Angle) float64 {
			return math.Sin(float64(a))
		},
	}
}

// NewFreq updates the frequency and Phase
func (osc *Oscillator) NewFreq(ν Hertz) {
	osc.ν = ν
	//	syn.DeltaPhase = Angle(Seconds(f) * τ * syn.Tick)
}
