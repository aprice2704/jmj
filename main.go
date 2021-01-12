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

// Synth is
type Synth struct {
	T0         time.Time // When this synth started playing
	Freq       float64   // Hz
	SR         float64   // Samples/Second
	Tick       float64   // Seconds/Sample
	WTick      float64   // Wavenumber/Sample, we use this to avoid disconinuities during frequency changes
	lastSample float64   // when was the last sample we made (seconds from T0)
}

// NewSynth makes and inits a new one
func NewSynth(t0 time.Time, f float64, sr float64) *Synth {
	syn := Synth{T0: t0, Freq: f, SR: sr}
	syn.Tick = 1 / sr
	syn.WTick = f * τ * syn.Tick
	syn.lastSample = 0.0
	return &syn
}

// NewFreq updates the frequency and wtick
func (syn *Synth) NewFreq(f float64) {
	syn.Freq = f
	syn.WTick = f * τ * syn.Tick
}

func main() {
	SR := float64(44100)
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
			fmt.Printf("Freq is %f\n", mySyn.Freq)
		}
		if event.Key == keyboard.KeyArrowDown {
			mySyn.NewFreq(mySyn.Freq - 1.0)
			fmt.Printf("Freq is %f\n", mySyn.Freq)
		}
	}
}

// Err satisifies beep.Streamer
func (syn Synth) Err() error {
	return nil
}

// Stream satisifies beep.Streamer
func (syn *Synth) Stream(samples [][2]float64) (n int, ok bool) {
	//	fmt.Printf("N is %d, delta T is %.9f\n", len(samples), syn.lastSample)
	wN := syn.lastSample // in wavenumber
	for i := range samples {
		samples[i][0] = 0.2 * (math.Sin(wN)) // + math.Sin(wN*2.0))
		samples[i][1] = 0.2 * (math.Cos(wN)) // *0.5) + math.Sin(wN*4.0))
		wN += syn.WTick
	}
	syn.lastSample = wN
	return len(samples), true
}
