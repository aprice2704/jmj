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
	"fmt"
	"math"
)

// PlanckTime is the shortest possible interval of time
const (
	PlanckTime Seconds = 1e-20
)

var (
	sqrt2π float64 = math.Sqrt(τ) // simple optimization
)

// Enveloper modulates a signal
type Enveloper interface {
	Amplitude(t Seconds) Volts // Call this from outside
	Length() Seconds           // Return overall length of the envelope
}

// Envelope is the 'base' type for envelopes
type Envelope struct {
	T0 Seconds // *Global* time when the envelope starts
	λ  Seconds // Period of repeat
	//	Repeats bool    // Does it repeat or is it single shot?
	Len Seconds // The overall length of the envelope (might be several λ long)
}

// NewEnvelope is
func NewEnvelope(t0 Seconds, λ Seconds, reps bool, l Seconds) *Envelope {
	return &Envelope{T0: t0, λ: λ, Len: l}
}

// ADSR is a classic ADSR envelope
type ADSR struct {
	Envelope
	Ta          Seconds      // Attack time (0->1)
	Td          Seconds      // Decay (1->Ls)
	Ls          Volts        // Sustain level (Ls)
	TsMax       Seconds      // Maximum sustain time
	TsMin       Seconds      // Minimum sustain time
	Tr          Seconds      // Release time (Ls->0)
	sStart      LocalSeconds // when did sustain begin?
	knowRelease bool         // Is release time known?
	releaseAt   LocalSeconds // when released
	tsActual    Seconds      // Actual sustain time (derived from keyup etc.)
}

// NewADSR makes a new one, pass ts as zero if not known at creation
func NewADSR(t0 Seconds, reps bool, ta Seconds, td Seconds, ls Volts, tsmax Seconds, tsmin Seconds, ts Seconds) *ADSR {
	if tsmax < PlanckTime {
		tsmax = MaxNoteLen
	}
	adsr := ADSR{Ta: ta, Td: td, Ls: ls, TsMax: tsmax, TsMin: tsmin, tsActual: ts, sStart: LocalSeconds(ta + td)}
	adsr.Envelope = Envelope{T0: t0}
	if ts > PlanckTime { // we know when release happens
		adsr.releaseAt = LocalSeconds(ts + ta + td)
		adsr.knowRelease = true
	}
	return &adsr
}

// Release triggers the release at the given global time. Strangeness will result if called after release should have started
func (adsr ADSR) Release(t Seconds) {
	tLocal := LocalSeconds(t - adsr.T0)
	ts := tLocal - adsr.sStart                             // length of the sustain
	tsact := max(min(Seconds(ts), adsr.TsMin), adsr.TsMax) // clip into valid range
	adsr.releaseAt = LocalSeconds(tsact)
	adsr.knowRelease = true
}

// Amplitude is
func (adsr ADSR) Amplitude(t Seconds) Volts {
	localT := LocalSeconds(t - adsr.T0)
	var a Volts
	switch {
	case localT < LocalSeconds(adsr.Ta):
		a = Volts(localT / LocalSeconds(adsr.Ta))
	case localT < adsr.sStart:
		a = 1 - Volts((localT-LocalSeconds(adsr.Ta))*(1-LocalSeconds(adsr.Ls))/LocalSeconds(adsr.Td))
	case localT < adsr.releaseAt:
		a = adsr.Ls
	case localT > adsr.releaseAt:
		if !adsr.knowRelease {
			fmt.Printf("ADSR envelope error: in release stage of unreleased envelope")
		}
	default:
		fmt.Printf("ADSR envelope error: in unknown portion of envelope")
	}
	return a
}

func max(a, b Seconds) Seconds {
	if a > b {
		return a
	}
	return b
}

func min(a, b Seconds) Seconds {
	if a < b {
		return a
	}
	return b
}

// Triangle a simple /\ with period λ
type Triangle struct {
	Envelope
}

// NewTriangle makes one
func NewTriangle(t Seconds, λ Seconds, reps bool, l Seconds) *Triangle {
	tr := Triangle{}
	tr.Envelope = Envelope{T0: t, λ: λ, Len: l}
	//	fmt.Printf("New triangle at %f\n", t)
	return &tr
}

// Amplitude is
func (tr Triangle) Amplitude(t Seconds) Volts {
	localT := t - tr.T0
	return tr.onePeriodAmplitude(LocalSeconds(math.Mod(float64(localT), float64(tr.λ))))
}

// onePeriodAmplitude is
func (tr Triangle) onePeriodAmplitude(t LocalSeconds) Volts {
	if Seconds(t) < (tr.λ)/2 {
		return Volts(Seconds(t) * 2 / tr.λ)
	}
	return Volts(2 - (Seconds(t) * 2 / tr.λ))
}

// Length is
func (tr Triangle) Length() Seconds {
	return tr.Len
}

// Gaussian is an envelope with height 1 at μ and RMS width of σ
// f(x) = exp(-(x-μ)^2/2σ^2) μ and σ should be specified in seconds
type Gaussian struct {
	Envelope
	μ, σ Seconds
	σσ   Seconds
}

// NewGaussian makes a new one
func NewGaussian(globalT Seconds, λ Seconds, reps bool, l Seconds, newμ, newσ Seconds) *Gaussian {
	e := NewEnvelope(globalT, λ, reps, l)
	g := &Gaussian{Envelope: *e, μ: newμ, σ: newσ, σσ: newσ * newσ}
	return g
}

// OnePeriodAmplitude fulfils Envelope interface
func (g *Gaussian) onePeriodAmplitude(localT Seconds) Volts {
	xu := float64(localT - g.μ)
	return Volts(math.Exp(-xu * xu / float64(2*g.σσ)))
}
