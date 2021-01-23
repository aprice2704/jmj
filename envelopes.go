package main

// ███████╗███╗   ██╗██╗   ██╗███████╗██╗      ██████╗ ██████╗ ███████╗███████╗
// ██╔════╝████╗  ██║██║   ██║██╔════╝██║     ██╔═══██╗██╔══██╗██╔════╝██╔════╝
// █████╗  ██╔██╗ ██║██║   ██║█████╗  ██║     ██║   ██║██████╔╝█████╗  ███████╗
// ██╔══╝  ██║╚██╗██║╚██╗ ██╔╝██╔══╝  ██║     ██║   ██║██╔═══╝ ██╔══╝  ╚════██║
// ███████╗██║ ╚████║ ╚████╔╝ ███████╗███████╗╚██████╔╝██║     ███████╗███████║
// ╚══════╝╚═╝  ╚═══╝  ╚═══╝  ╚══════╝╚══════╝ ╚═════╝ ╚═╝     ╚══════╝╚══════╝

// Envelopes are designed to modulate the amplitude of a signal.
// They should be normalized to have values between 0 and 1

import (
	"math"
)

var (
	sqrt2π float64 = math.Sqrt(τ) // simple optimization
)

// Enveloper modulates a signal
type Enveloper interface {
	Amplitude(t Seconds) float64          // Call this from outside
	OnePeriodAmplitude(t Seconds) float64 // Implement this
	Length() Seconds                      // Return overall length of the envelope
}

// Envelope is the 'base' type for envelopes
type Envelope struct {
	λ       Seconds // Period of repeat
	Repeats bool    // Does it repeat or is it single shot?
	Len     Seconds // The overall length of the envelope (might be several λ long)
}

// NewEnvelope is self-evident
func NewEnvelope(λ Seconds, reps bool, l Seconds) *Envelope {
	return &Envelope{λ: λ, Repeats: reps, Len: l}
}

// Gaussian is an envelope with height 1 at μ and RMS width of σ
// f(x) = exp(-(x-μ)^2/2σ^2) μ and σ should be specified in seconds
type Gaussian struct {
	Envelope
	μ, σ Seconds
	σσ   Seconds
}

// NewGaussian makes a new one
func NewGaussian(λ Seconds, reps bool, l Seconds, newμ, newσ Seconds) *Gaussian {
	e := NewEnvelope(λ, reps, l)
	g := &Gaussian{Envelope: *e, μ: newμ, σ: newσ, σσ: newσ * newσ}
	return g
}

// OnePeriodAmplitude fulfils Envelope interface
func (g *Gaussian) OnePeriodAmplitude(x Seconds) float64 {

	xu := float64(x - g.μ)
	return math.Exp(-xu * xu / float64(2*g.σσ))
}

// Triangle a simple /\ with period λ
type Triangle struct {
	*Envelope
}

// NewTriangle is self-evident
func NewTriangle(λ Seconds, reps bool, l Seconds) *Triangle {
	e := NewEnvelope(λ, reps, l)
	return &Triangle{Envelope: e}
}

// Amplitude is
func (tr Triangle) Amplitude(t Seconds) float64 {
	return tr.OnePeriodAmplitude(Seconds(math.Mod(float64(t), float64(tr.λ))))
}

// OnePeriodAmplitude is
func (tr Triangle) OnePeriodAmplitude(t Seconds) float64 {
	//	s := math.Mod(float64(t), float64(tr.λ))
	if Seconds(t) < (tr.λ)/2 {
		return float64(t * 2 / tr.λ)
	}
	return float64(2 - (t * 2 / tr.λ))
}

// Length is
func (tr Triangle) Length() Seconds {
	return tr.Len
}

// RepeatAmplitude is
// func (tr *Triangle) RepeatAmplitude(t Seconds) float64 {
// 	s := math.Mod(float64(t), float64(tr.λ))
// 	if Seconds(s) < (tr.λ)/2 {
// 		return s * 2 / float64(tr.λ)
// 	}
// 	return 2 - (s * 2 / float64(tr.λ))
// }
// SetPeriodandLength fulfils Enveloper interface
// func (tr Triangle) SetPeriodandLength(λ Seconds, length Seconds) {
// 	tr.Envelope.SetPeriodandLength(λ, length)
// }

// RepeatAmplitude generates a sequence of gaussian envelopes of period λ seconds
// func (g *Gaussian) OneShotAmplitude(t Seconds) float64 {

// 	s := math.Mod(float64(t), float64(g.λ))
// 	tu := s - float64(g.μ)
// 	return math.Exp(-tu * tu / float64(2*g.σσ))

// }
// OnePeriodAmplitude is the function you should implement for new types of envelope
// func (e *Envelope) OnePeriodAmplitude(t Seconds) float64 {
// 	fmt.Printf("Warning: Amplitude called on base Envelope class, probably an error\n")
// 	return 1
// }

// Amplitude is the function to call from outide the class to get either the single shot envelope or the repeated one
// If you implemented OneShotAmplitude, this should give you repeats for free
// func (e *Envelope) Amplitude(t Seconds) float64 {
// 	if e.Repeats {
// 		return e.OnePeriodAmplitude(Seconds(math.Mod(float64(t), float64(e.λ))))
// 	}
// 	return e.OnePeriodAmplitude(t)
// }

// Length is
// func (e *Envelope) Length() Seconds {
// 	return e.Len
// }

// SetPeriodandLength fulfils Enveloper interface
// func (e *Envelope) SetPeriodandLength(λ Seconds, length Seconds) {
// 	e.λ = λ
// 	e.Len = length
// }
//	SetPeriodandLength(λ Seconds, length Seconds) // Set both the period (λ) field and length
