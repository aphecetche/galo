package run2

import (
	"encoding/binary"
	"io"
)

func ForEachEvent(r io.Reader, efunc func(e *Event), maxEvents int) int {
	nevents := 0
	for nevents < maxEvents {
		sb := make([]byte, 4)
		nb, err := r.Read(sb)
		if nb != 4 {
			break
		}
		size := binary.LittleEndian.Uint32(sb)
		buf := make([]byte, size)
		nb, err = r.Read(buf)
		if uint32(nb) != size {
			panic(err)
		}
		event := GetRootAsEvent(buf, 0)
		efunc(event)
		nevents++
	}
	return nevents
}
