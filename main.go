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

// Hertz is a fequency (1/Seconds)
type Hertz Seconds

// Angle is a portion of a wave, typically a phase, also radians
type Angle float64

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
		// if event.Key == keyboard.KeyArrowUp {
		// 	mySyn.NewFreq(mySyn.Freq + 1.0)
		// }
		// if event.Key == keyboard.KeyArrowDown {
		// 	mySyn.NewFreq(mySyn.Freq - 1.0)
		// }
		// if event.Key == keyboard.KeyPgdn {
		// 	mySyn.NewFreq(mySyn.Freq - 10.0)
		// }
		// if event.Key == keyboard.KeyPgup {
		// 	mySyn.NewFreq(mySyn.Freq + 10.0)
		// }
		//		fmt.Printf("Freq is %f\n", mySyn.Freq)
		if event.Rune == 'c' {
			myOsc := NewSine(MiddleCfreq)
			myEnv := Triangle{Envelope{λ: 3, Repeats: false, Len: 3}}
			myNote := &Note{BaseFreq: MiddleCfreq, Env: myEnv, Osc: myOsc}
			//	mySound := Sound{Note: myNote, Start: mySyn.Now(), End: mySyn.Now() + myEnv.Len}
			mySyn.AddSound(myNote, mySyn.Now())
		}
	}

}

// Err satisifies beep.Streamer
func (syn Synth) Err() error {
	return nil
}

// Stream satisifies beep.Streamer, computes the instantaneous amplitude for each channel.
func (syn *Synth) Stream(samples [][2]float64) (n int, ok bool) {
	//	fmt.Printf("N is %d, delta T is %.9f\n", len(samples), syn.lastSample)
	//	phase := syn.lastSample
	when := syn.lastAt
	for i := range samples {
		when += syn.Tick
		samples[i][0] = syn.Amplitude(when)
		samples[i][1] = syn.Amplitude(when)
		// samples[i][0] = GaussianRepeat(1, 1, 2, when) * math.Sin(float64(phase))
		// samples[i][1] = GaussianRepeat(1, 1, 2, when) * math.Sin(float64(phase))
		// samples[i][0] = GaussianRepeat(0.05, 0.02, 0.1, when) * math.Sin(float64(phase))
		// samples[i][1] = GaussianRepeat(0.05, 0.02, 0.11, when) * math.Sin(float64(phase))
		// 		samples[i][0] = TriangleRepeat(0.2, when) * math.Sin(float64(phase))
		// 		samples[i][1] = TriangleRepeat(0.2, when) * math.Sin(float64(phase+1))
		//		phase += syn.DeltaPhase
	}
	//	syn.lastSample = phase
	syn.lastAt = when
	return len(samples), true
}
