package main

import "math"

// ███╗   ██╗ ██████╗ ████████╗███████╗
// ████╗  ██║██╔═══██╗╚══██╔══╝██╔════╝
// ██╔██╗ ██║██║   ██║   ██║   █████╗
// ██║╚██╗██║██║   ██║   ██║   ██╔══╝
// ██║ ╚████║╚██████╔╝   ██║   ███████╗
// ╚═╝  ╚═══╝ ╚═════╝    ╚═╝   ╚══════╝

// Scientific scale note frequencies in Hz
const (
	C0freq      = 16.35160
	C1freq      = C0freq * 2
	C2freq      = C0freq * 4
	C3freq      = C0freq * 8
	C4freq      = C0freq * 16 // Middle C
	MiddleCfreq = C4freq
	C5freq      = C0freq * 32
	C6freq      = C0freq * 64
	C7freq      = C0freq * 128
)

// MaxNoteLen is the length of an 'infinite'/repeating note
const (
	MaxNoteLen Seconds = 86400
)

// NoteFreqs gives the fequencies in Hz for the scientic notation heptatonic Western scale
var NoteFreqs map[string]float64

func init() {
	NoteFreqs = make(map[string]float64)
	for oct, octS := range "012345678" {
		f := C0freq * math.Pow(2, float64(oct))
		for n, note := range "CDEFGAB" {
			NoteFreqs[string(note)+string(octS)] = f * math.Pow(2, float64(n)*(1.0/7.0))
		}
	}
}

// Note is an instance of a voice, played with an envelope
type Note struct {
	BaseFreq Hertz
	Env      Enveloper
	Osc      *Oscillator // just for now
	//	Voice *Voicer // TODO
}

// Length returns that of the underlying envelope
func (n Note) Length() Seconds {
	return n.Env.Length()
}

// Amplitude returns the singal strength at a given time
func (n Note) Amplitude(t Seconds) float64 {
	return n.Env.Amplitude(t) * n.Osc.Amplitude(t)
}
