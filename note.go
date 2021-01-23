package main

import (
	"fmt"
	"math"
)

// ███╗   ██╗ ██████╗ ████████╗███████╗
// ████╗  ██║██╔═══██╗╚══██╔══╝██╔════╝
// ██╔██╗ ██║██║   ██║   ██║   █████╗
// ██║╚██╗██║██║   ██║   ██║   ██╔══╝
// ██║ ╚████║╚██████╔╝   ██║   ███████╗
// ╚═╝  ╚═══╝ ╚═════╝    ╚═╝   ╚══════╝

// Scientific scale note frequencies in Hz
var (
	C0freq      Hertz = 16.35160
	C1freq      Hertz = C0freq * 2
	C2freq      Hertz = C0freq * 4
	C3freq      Hertz = C0freq * 8
	C4freq      Hertz = C0freq * 16 // Middle C
	MiddleCfreq Hertz = C4freq
	C5freq      Hertz = C0freq * 32
	C6freq      Hertz = C0freq * 64
	C7freq      Hertz = C0freq * 128
)

// MaxNoteLen is the length of an 'infinite'/repeating note
const (
	MaxNoteLen Seconds = 86400
)

// NoteFreqs gives the fequencies in Hz for the scientic notation heptatonic Western scale
var NoteFreqs map[string]Hertz

func init() {
	NoteFreqs = make(map[string]Hertz)
	for oct, octS := range "012345678" {
		f := C0freq * Hertz(math.Pow(2, float64(oct)))
		for n, note := range "CDEFGAB" {
			NoteFreqs[string(note)+string(octS)] = f * Hertz(math.Pow(2, float64(n)*(1.0/7.0)))
		}
		for n, note := range "cdefgab" {
			NoteFreqs[string(note)+string(octS)] = f * Hertz(math.Pow(2, float64(n)*(1.0/7.0)))
		}
	}
}

// GetFreq returns the frequency of the note given in the string "A0" ... "G7"
func GetFreq(note string) Hertz {
	if len(note) != 2 {
		fmt.Printf("Wrong note %s\n", note)
		return 0
	}
	return NoteFreqs[note]
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

// Amplitude returns the signal strength at a given time
func (n Note) Amplitude(t Seconds) float64 {
	return n.Env.Amplitude(t) * n.Osc.Amplitude(t)
}
