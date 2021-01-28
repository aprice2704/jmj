package main

import (
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
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
	SampleNo   int       // number of the last sample emitted
	Freq       Hertz     // Hz
	SR         Hertz     // Samples/Second
	Tick       Seconds   // Seconds/Sample
	DeltaPhase Angle     // Radians/Sample
	Sounds     []*Sound  // Sounds being considered for playing
	recordingL []float64
	recordingR []float64
	recordIt   bool
}

// Sound is a note played at a particular time
type Sound struct {
	*Note
	Start Seconds
	End   Seconds
}

// Amplitude is just that of the underlying note
func (snd Sound) Amplitude(t Seconds) Volts {
	return snd.Note.Amplitude(t)
}

// NewSynth makes and inits a new one
func NewSynth(t0 time.Time, f Hertz, sr Hertz) *Synth {
	syn := Synth{T0: t0, Freq: f, SR: sr}
	syn.Tick = Seconds(1 / sr)
	syn.recordingL = make([]float64, 0, 1000000)
	syn.recordingR = make([]float64, 0, 1000000)
	return &syn
}

// Now is the current 'Global Time' of the synth in Seconds since starting
// Sounds that wish to start immediately should do so at syn.Now()
func (syn Synth) Now() Seconds {
	return Seconds(float64(time.Now().Sub(syn.T0)) / float64(time.Second))
	// return syn.lastAt + syn.Tick
}

// AddSound adds a note to be played starting at time 'when'
func (syn *Synth) AddSound(n *Note, start Seconds) {
	//	fmt.Printf("Playing sound from %f to %f\n", start, start+n.Length())
	ns := &Sound{Note: n, Start: start, End: start + n.Length()}
	syn.Sounds = append(syn.Sounds, ns)
	//	sort.Slice(syn.Sounds, func(i, j int) bool { return syn.Sounds[i].End < syn.Sounds[j].End })
}

// Amplitude adds all the currently playing notes together, culls any that have completed
func (syn *Synth) Amplitude(t Seconds) Volts {
	a := Volts(0.0)
	n := 0
	for _, s := range syn.Sounds {
		if s.Start <= t && s.End >= t {
			a += s.Amplitude(t)
			n++
		}
	}
	if n > 0 {
		if math.Abs(float64(a)) > 1 { // clamp to +-1
			if math.Signbit(float64(a)) {
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

// Stream satisifies beep.Streamer, computes the instantaneous amplitude for each channel.
func (syn *Synth) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		when := Seconds(syn.SampleNo) * syn.Tick
		aR := syn.Amplitude(when)
		aL := syn.Amplitude(when)
		samples[i][0] = float64(aR)
		samples[i][1] = float64(aL)
		if syn.recordIt {
			syn.recordingR = append(syn.recordingR, float64(aR))
			syn.recordingL = append(syn.recordingL, float64(aL))
		}
		syn.SampleNo++
	}
	return len(samples), true
}

// Graphout draws an graphic of this synth
func (syn Synth) Graphout() {

	nSamples := len(syn.recordingR)
	nSecs := 1 + (nSamples / int(syn.SR))
	colH := 200
	sideH := colH / 2
	colW := 2000
	margin := 20
	yScale := float64(colH) / 2
	totalW := margin*2 + colW
	totalH := margin*2 + nSecs*(colH+margin)
	totalT := Seconds(nSamples) / Seconds(syn.SR)

	xy := func(s int) (x, y int) {
		stripe := s / int(syn.SR)
		col := margin + int(float64(colW)*float64(s%int(syn.SR))/float64(syn.SR))
		row := margin + stripe*(colH+margin) + colH/2
		return col, row
	}

	t2samp := func(t Seconds) int {
		return int(t * Seconds(syn.SR))
	}

	upLeft := image.Point{0, 0}
	lowRight := image.Point{totalW, totalH}
	all := image.Rectangle{upLeft, lowRight}
	img := image.NewRGBA(all)
	bg := image.NewUniform(imWhite)
	draw.Draw(img, all, bg, image.Pt(0, 0), draw.Over)

	// Output waveform
	for samp := 0; samp < nSamples; samp++ {
		x, y := xy(samp)
		img.Set(x, y+int(syn.recordingR[samp]*yScale), imBlue)
		img.Set(x, y, imBlack)
	}

	// 0.1s red ticks on time axis
	for t := Seconds(0); t < totalT; t += 0.1 {
		x, y := xy(t2samp(t))
		for i := -sideH / 4; i < sideH/4; i++ {
			img.Set(x, y+i, imRed)
		}
	}

	// Green lines at start of each sound
	for _, s := range syn.Sounds {
		x, y := xy(t2samp(s.Start))
		for i := -sideH; i < sideH; i++ {
			img.Set(x, y+i, imGreen)
		}
	}

	// f, _ := os.Create("line.jpg")
	// jpeg.Encode(f, img, &jpeg.Options{Quality: 95})

	f, _ := os.Create("line2.png")
	png.Encode(f, img)

}
