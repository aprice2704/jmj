package main

import (
	"fmt"
	"math"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

const (
	π = math.Pi
	τ = 2 * math.Pi
)

// Seconds are quantities of time
type Seconds float64

// Hertz is a fequency
type Hertz Seconds

// Angle is a portion of a wave, typically a phase
type Angle float64

// Synth is
type Synth struct {
	T0         time.Time // When this synth started playing
	Freq       Hertz     // Hz
	SR         Hertz     // Samples/Second
	Tick       Seconds   // Seconds/Sample
	DeltaPhase Angle     // We use this to avoid disconinuities during frequency changes
	lastSample Angle     // phase of the last sample we made
	lastAt     Seconds   // When we made the last sample
}

// NewSynth makes and inits a new one
func NewSynth(t0 time.Time, f Hertz, sr Hertz) *Synth {
	syn := Synth{T0: t0, Freq: f, SR: sr}
	syn.Tick = Seconds(1 / sr)
	syn.DeltaPhase = Angle(Seconds(f) * τ * syn.Tick)
	syn.lastSample = 0.0
	syn.lastAt = 0.0
	return &syn
}

// NewFreq updates the frequency and Phase
func (syn *Synth) NewFreq(f Hertz) {
	syn.Freq = f
	syn.DeltaPhase = Angle(Seconds(f) * τ * syn.Tick)
}

func main() {

	SR := Hertz(44100)
	mySyn := NewSynth(time.Now(), 330, SR)
	sr := beep.SampleRate(SR)
	speaker.Init(sr, sr.N(time.Second/100))
	speaker.Play(mySyn)
	keysEvents, err := keyboard.GetKeys(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEsc {
			break
		}
		if event.Key == keyboard.KeyArrowUp {
			mySyn.NewFreq(mySyn.Freq + 1.0)
		}
		if event.Key == keyboard.KeyArrowDown {
			mySyn.NewFreq(mySyn.Freq - 1.0)
		}
		if event.Key == keyboard.KeyPgdn {
			mySyn.NewFreq(mySyn.Freq - 10.0)
		}
		if event.Key == keyboard.KeyPgup {
			mySyn.NewFreq(mySyn.Freq + 10.0)
		}
		fmt.Printf("Freq is %f\n", mySyn.Freq)
	}
}

// Err satisifies beep.Streamer
func (syn Synth) Err() error {
	return nil
}

// Stream satisifies beep.Streamer
func (syn *Synth) Stream(samples [][2]float64) (n int, ok bool) {
	//	fmt.Printf("N is %d, delta T is %.9f\n", len(samples), syn.lastSample)
	phase := syn.lastSample
	when := syn.lastAt
	for i := range samples {
		// samples[i][0] = GaussianRepeat(1, 1, 2, when) * math.Sin(float64(phase))
		// samples[i][1] = GaussianRepeat(1, 1, 2, when) * math.Sin(float64(phase))
		samples[i][0] = GaussianRepeat(0.05, 0.02, 0.1, when) * math.Sin(float64(phase))
		samples[i][1] = GaussianRepeat(0.05, 0.02, 0.1, when) * math.Sin(float64(phase))
		phase += syn.DeltaPhase
		when += syn.Tick
	}
	syn.lastSample = phase
	syn.lastAt = when
	return len(samples), true
}
