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
	p.Tempo = tempo
	return p, nil
}

type Track struct{
	Id int
	Instrument string
	Steps [16]bool
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
		s += fmt.Sprintf("(%v) %v\t", track.Id, track.Instrument)
		stepsString := "|"
		for i, step := range track.Steps{
			if step {
				stepsString += "x"
			} else {
				stepsString += "-"
			}
			if i%4 == 3{
				stepsString += "|"
			}
		}
		s += stepsString + "\n"
	}
	return s
}