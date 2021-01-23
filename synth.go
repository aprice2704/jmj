package main

import (
	"math"
	"time"
)

// ███████╗██╗   ██╗███╗   ██╗████████╗██╗  ██╗
// ██╔════╝╚██╗ ██╔╝████╗  ██║╚══██╔══╝██║  ██║
// ███████╗ ╚████╔╝ ██╔██╗ ██║   ██║   ███████║
// ╚════██║  ╚██╔╝  ██║╚██╗██║   ██║   ██╔══██║
// ███████║   ██║   ██║ ╚████║   ██║   ██║  ██║
// ╚══════╝   ╚═╝   ╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝

// Synth is
type Synth struct {
	T0         time.Time // When this synth started playing
	Freq       Hertz     // Hz
	SR         Hertz     // Samples/Second
	Tick       Seconds   // Seconds/Sample
	DeltaPhase Angle     // Radians/Sample
	lastSample Angle     // phase of the last sample we made. Used to avoid disconinuities during frequency changes
	lastAt     Seconds   // When we made the last sample
	lastTime   time.Time //
	Sounds     []*Sound  // Sounds being considered for playing
}

// Sound is a note played at a particular time
type Sound struct {
	*Note
	Start Seconds
	End   Seconds
}

// Amplitude is just that of the underlying note
func (snd Sound) Amplitude(t Seconds) float64 {
	return snd.Note.Amplitude(t)
}

// NewSynth makes and inits a new one
func NewSynth(t0 time.Time, f Hertz, sr Hertz) *Synth {
	syn := Synth{T0: t0, Freq: f, SR: sr}
	syn.Tick = Seconds(1 / sr)
	//	syn.DeltaPhase = Angle(Seconds(f) * τ * syn.Tick)
	//	syn.lastSample = 0.0
	syn.lastAt = 0.0
	return &syn
}

// Now returns the time the synth considers itself to be at, which is in fact the
// next time (in seconds from starting) at which a sample will be generated.
// Sounds that wish to start immediately should do so at syn.Now()
func (syn Synth) Now() Seconds {
	return syn.lastAt + syn.Tick
}

// AddSound adds a note to be played starting at time 'when'
func (syn *Synth) AddSound(n *Note, when Seconds) {
	ns := &Sound{Note: n, Start: when, End: when + n.Length()}
	syn.Sounds = append(syn.Sounds, ns)
	//	sort.Slice(syn.Sounds, func(i, j int) bool { return syn.Sounds[i].End < syn.Sounds[j].End })
}

// Amplitude adds all the currently playing notes together, culls any that have completed
func (syn *Synth) Amplitude(t Seconds) float64 {
	a := 0.0
	n := 0
	for _, s := range syn.Sounds {
		if s.Start <= t && s.End >= t {
			a += float64(s.Amplitude(t))
			n++
		}
	}
	if n > 0 {
		if math.Abs(a) > 1 { // clamp to +-1
			if math.Signbit(a) {
				return -1
			}
			return 1
		}
		return a // its ok, in range -1...+1
	}
	return 0.0 // dead air
}

// PruneSounds removes any from the list that have finished playing
func (syn *Synth) PruneSounds(t Seconds) {
	newSounds := []*Sound{}
	for _, n := range syn.Sounds {
		if n.End <= t {
			newSounds = append(newSounds, n)
		}
	}
	syn.Sounds = newSounds
}
