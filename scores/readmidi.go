package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/reader"
)

type printer struct{}

func (pr printer) noteOn(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
}

func (pr printer) noteOff(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOff (ch %v: key %v)\n", p.Track, p.AbsoluteTicks, channel, key)
}

func main() {

	var p printer

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(reader.NoLogger(),
		// set the functions for the messages you are interested in
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
	)

	f := "Rendez-vous_III_Laser_Harpe.mid"
	err := reader.ReadSMFFile(rd, f)

	if err != nil {
		fmt.Printf("could not read SMF file %v\n", f)
	}

}
