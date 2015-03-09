package drum

import (
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/binary"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	b, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		fmt.Println("ioutil.ReadFile failed:", readErr)
	}
	p := &Pattern{}
	p.Version = func(b []byte) string{
		var s string
		for _, value := range b {
			if(value <= 31){
				return s
			}
			s += string(value)
		}
		return s
	}(b[14:25])
	var tempo float32
	buf := bytes.NewReader(b[46:50])
	err := binary.Read(buf, binary.LittleEndian, &tempo)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	p.Tempo = tempo
	tracks := make([]Track, 0)
	for track, nextSlice := makeTrack(b[50:]); track != nil; track, nextSlice = makeTrack(nextSlice){
		tracks = append(tracks, *track)
	}
	p.Tracks = tracks
	return p, nil
}

func makeTrack(b []byte) (*Track, []byte){
	if len(b) == 0 {
		return nil, nil
	}
	id := int(b[0])
	b = b[1:]
	var instrument string
	stepsStart := func(bytes []byte) int {
		inString := false
		for i, value := range bytes{
			if (inString && value <= 31){
				return i
			}
			if value >= 31{
				inString = true
				instrument += string(value)
			}
		}
		return -1
	}(b)
	var steps [16]bool
	for i, value := range b[stepsStart:stepsStart+16]{
		integer := int(value)
		if integer != 0 && integer != 1 {
			return nil, nil
		}
		steps[i] = int(value) != 0
	}
	return &Track{id, instrument, steps}, b[stepsStart+16:]
}

type Track struct{
	Id int
	Instrument string
	Steps [16]bool
}

func (t *Track) String() string {
	var s string
	s += fmt.Sprintf("(%v) %v\t|", t.Id, t.Instrument)
	for i, step := range t.Steps{
		if step {
			s+= "x"
		} else {
			s += "-"
		}
		if i%4 == 3{
			s += "|"
		}
	}
	return s
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct{
	Version string
	Tempo float32
	Tracks []Track
}

func (p *Pattern) String() string {
	s := fmt.Sprintf("Saved with HW Version: %v\n", p.Version)
	s += fmt.Sprintf("Tempo: %v\n", p.Tempo)
	for _, track := range p.Tracks{
		s += track.String() + "\n"
	}
	return s
}