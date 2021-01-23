package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

const (
	fontPath = "Go-Mono.ttf"
	fontSize = 24
)

var recordingL = make([]float64, 0, 1000000)
var recordingR = make([]float64, 0, 1000000)
var recordIt bool

var (
	imCyan  = color.RGBA{100, 200, 200, 0xff}
	imRed   = color.RGBA{255, 0, 0, 0xff}
	imBlue  = color.RGBA{0, 0, 255, 0xff}
	imGreen = color.RGBA{0, 255, 0, 0xff}
	imBlack = color.RGBA{0, 0, 0, 0xff}
	imWhite = color.RGBA{255, 255, 255, 0xff}
)

var (
	red   = sdl.Color{R: 255, G: 0, B: 0, A: 255}
	green = sdl.Color{R: 0, G: 255, B: 0, A: 255}
	blue  = sdl.Color{R: 0, G: 0, B: 255, A: 255}
	black = sdl.Color{R: 0, G: 0, B: 0, A: 255}
)

func main() {

	SR := Hertz(44100)
	mySyn := NewSynth(time.Now(), 330, SR)
	sr := beep.SampleRate(SR)
	speaker.Init(sr, sr.N(time.Second/5))
	speaker.Play(mySyn)

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("JMJ", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	if err = ttf.Init(); err != nil {
		fmt.Printf("Error: SDL: Could not init ttf lib: %s\n", err)
		return
	}
	defer ttf.Quit()

	// Load the font for our text
	var font *ttf.Font
	if font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		fmt.Printf("Error: SDL: Could not open font: %s\n", err)
		return
	}
	defer font.Close()

	mainSurf, err := window.GetSurface()
	if err != nil {
		fmt.Printf("Error: SDL: Could not get surface: %s\n", err)
		return
	}

	mainSurf.FillRect(nil, 0)
	textAt(font, red, black, mainSurf, 2, 2, "JMJ")
	textAt(font, green, black, mainSurf, 2, 32, "JMJ too")
	window.UpdateSurface()

	lowRowIn := "zxcvbnm"
	lowRowOut := "ABCDEFG"

	running := true
	recordIt = true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				// fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
				// 	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case *sdl.MouseButtonEvent:
				// fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
				// 	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			case *sdl.MouseWheelEvent:
				// fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
				// 	t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyboardEvent:
				//				typeName := "?"
				switch t.Type {
				case 768:
					//					typeName = "KeyDown"
					freq := MiddleCfreq
					c := fmt.Sprintf("%c", t.Keysym.Sym)
					if c == "q" {
						running = false
					}
					p := strings.Index(lowRowIn, c)
					ns := "?"
					if p != -1 {
						ns = lowRowOut[p:p+1] + "4"
						freq = GetFreq(ns)
					}
					fmt.Printf("Adding sound %s at %f from key %s (index %d)\n", ns, freq, c, p)
					myOsc := NewSine(freq)
					myEnv := NewTriangle(1, false, 1)
					myNote := &Note{BaseFreq: freq, Env: myEnv, Osc: myOsc}
					mySyn.AddSound(myNote, mySyn.Now())
				case 769:
					//					typeName = "KeyUp"
				}
				// fmt.Printf("[%d ms] Keyboard\ttype: %s (%d)\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				// t.Timestamp, typeName, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			default:
				// fmt.Printf("Unknown event type: %d\n", event)
			}
			textAt(font, blue, black, mainSurf, 2, 62, fmt.Sprintf("Sounds: %d", len(mySyn.Sounds)))
			window.UpdateSurface()
			time.Sleep(time.Millisecond)
		}
	}

	if recordIt {
		width := 800
		w2 := float64(width) / 2
		step := 5
		height := len(recordingR) / step
		upLeft := image.Point{0, 0}
		lowRight := image.Point{width, height}
		all := image.Rectangle{upLeft, lowRight}
		img := image.NewRGBA(all)
		bg := image.NewUniform(imWhite)
		draw.Draw(img, all, bg, image.Pt(0, 0), draw.Over)
		row := 0
		for samp := 0; samp < len(recordingR); samp += step {
			lt := int(w2)
			rt := int((1 + recordingR[samp]) * w2)
			if lt > rt {
				lt, rt = rt, lt
			}
			for s := lt; s < rt; s++ {
				img.Set(s, row, imBlue)
			}
			img.Set(width/2, row, imBlack)
			row++
		}
		f, _ := os.Create("line.jpg")
		jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	}

}

// Err satisifies beep.Streamer
func (syn Synth) Err() error {
	return nil
}

func textAt(f *ttf.Font, fgColor sdl.Color, bgColor sdl.Color, s *sdl.Surface, x int32, y int32, txt string) {

	var textSur *sdl.Surface
	var err error
	textSur, err = f.RenderUTF8Shaded(txt, fgColor, bgColor)
	if err != nil {
		fmt.Printf("Error: SDL: Could not render text: %s\n", err)
		return
	}
	defer textSur.Free()
	err = textSur.Blit(nil, s, &sdl.Rect{X: x, Y: y, W: 0, H: 0})
	if err != nil {
		fmt.Printf("Error: SDL: Blitting text failed: %s", err)
		return
	}

	return
}

var lastprint time.Time

// Stream satisifies beep.Streamer, computes the instantaneous amplitude for each channel.
func (syn *Synth) Stream(samples [][2]float64) (n int, ok bool) {
	when := syn.lastAt
	for i := range samples {
		when += syn.Tick
		aR := syn.Amplitude(when)
		aL := syn.Amplitude(when)
		samples[i][0] = aR
		samples[i][1] = aL
		if recordIt {
			recordingR = append(recordingR, aR)
			recordingL = append(recordingL, aL)
		}
	}
	syn.lastAt = when
	syn.lastTime = time.Now()
	return len(samples), true
}

// line := charts.NewLine()
// line.SetGlobalOptions(
// 	charts.WithTitleOpts(opts.Title{Title: "Sound Output", Subtitle: "Final result, right channel only"}),
// 	charts.WithDataZoomOpts(opts.DataZoom{Type: "inside", Start: 0, End: 100}),
// 	charts.WithToolboxOpts(opts.Toolbox{Feature: &opts.ToolBoxFeature{DataZoom: &opts.ToolBoxFeatureDataZoom{Show: true}}}),
// )
// line.SetXAxis(opts.XAxis{Name: "Sample #", Type: "time"})
// d := make([]opts.LineData, len(recordingR), len(recordingR))
// for i, v := range recordingR {
// 	d[i].Value = v
// }
// line.AddSeries("Right Amplitude", d)
