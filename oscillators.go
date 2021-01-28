package main

import (
	"math"
)

//  ██████╗ ███████╗ ██████╗██╗██╗     ██╗      █████╗ ████████╗ ██████╗ ██████╗ ███████╗
// ██╔═══██╗██╔════╝██╔════╝██║██║     ██║     ██╔══██╗╚══██╔══╝██╔═══██╗██╔══██╗██╔════╝
// ██║   ██║███████╗██║     ██║██║     ██║     ███████║   ██║   ██║   ██║██████╔╝███████╗
// ██║   ██║╚════██║██║     ██║██║     ██║     ██╔══██║   ██║   ██║   ██║██╔══██╗╚════██║
// ╚██████╔╝███████║╚██████╗██║███████╗███████╗██║  ██║   ██║   ╚██████╔╝██║  ██║███████║
//  ╚═════╝ ╚══════╝ ╚═════╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚══════╝

// Osciller is an interface for oscillators
type Osciller interface {
	Amplitude(t Seconds) Volts
}

// Oscillator is a periodically repeating waveform
type Oscillator struct {
	T0      Seconds             // Global time when this osc started
	ν       Hertz               // Fundamental frequency
	Phase   Angle               // Last known phase
	PhaseAt Seconds             // When that phase occurred
	Wave    func(a Angle) Volts // Function that describes wave shape
}

// Waveform is a function that encodes the shape of the cycles of a waveform in *angle*
type Waveform func(a Angle) Volts

// Amplitude returns the strength of the waveform (which may change frequency) at a given global time
func (osc *Oscillator) Amplitude(t Seconds) Volts {
	ot := t - osc.T0 // local time
	dT := ot - osc.PhaseAt
	dA := Angle(dT) * τ * Angle(osc.ν) // convert time to phase angle at new freq
	osc.Phase += dA
	osc.PhaseAt = ot
	return osc.Wave(osc.Phase)
}

// NewSine returns a new sine wave oscillator starting at global time t
func NewSine(t Seconds, newν Hertz) *Oscillator {
	//	fmt.Printf("New sine osc at %f\n", t)
	return &Oscillator{
		T0:      t, // Global start time
		ν:       newν,
		Phase:   0,
		PhaseAt: 0,
		Wave: func(a Angle) Volts {
			return Volts(math.Sin(float64(a)))
		},
	}
}

// NewFreq updates the frequency and Phase
func (osc *Oscillator) NewFreq(ν Hertz) {
	osc.ν = ν
}
