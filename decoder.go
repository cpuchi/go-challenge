package drum

import (
	"fmt"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
// TODO: implement
func DecodeFile(path string) (*Pattern, error) {
	p := &Pattern{}
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
	Tempo float64
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