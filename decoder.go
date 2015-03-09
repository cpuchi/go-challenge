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
// TODO: implement
func DecodeFile(path string) (*Pattern, error) {
	b, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		fmt.Println("ioutil.ReadFile failed:", readErr)
	}
	p := &Pattern{}
	p.Version = fmt.Sprintf("%s", b[14:25])
	var tempo float32
	buf := bytes.NewReader(b[46:50])
	err := binary.Read(buf, binary.LittleEndian, &tempo)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	//tracks := make([]Track, 2)
	//for the rest->next byte is id
	//chill until find byte > 31, grab bytes until no longer 31
	//immediate next 16 bytes are steps
	//repeat
	//how to handle weird value in pattern_5? --> read each track, throw error when unexpected
	p.Tempo = tempo
	track, trackError := makeTrack(b[50:])
	fmt.Println(track)
	fmt.Println(trackError)
	return p, nil
}

func makeTrack(b []byte) (*Track, error){
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
		steps[i] = int(value) != 0
	}
	return &Track{id, instrument, steps}, nil
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
	s += "\n"
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
		s += fmt.Sprint("%s", track)
	}
	return s
}