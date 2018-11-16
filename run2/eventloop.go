package run2

import (
	"encoding/binary"
	"fmt"
	"io"
)

func buf2EventClusters(buf []byte) ([]*EventClusters, int64) {
	var off int64
	var ecs []*EventClusters
	var pstart, pend int64
	var size uint32
	for off < int64(len(buf)) {
		size = binary.LittleEndian.Uint32(buf[off : off+4])
		pstart = off + 4
		pend = pstart + int64(size)
		if pend >= int64(len(buf)) {
			break
		}
		event := GetRootAsEvent(buf[pstart:pend], 0)
		ec := GetEventClusters(event)
		ecs = append(ecs, ec)
		off = pend
	}
	if off == 0 {
		panic("reading buffer too small")
	}
	return ecs, off
}

// ForEachEvent loops over all events (until reaching maxEvents)
// in reader, converts it to EventClusters struct and
// finally executes a function for the EventClusters object.
func ForEachEvent(r io.ReaderAt, efunc func(ec *EventClusters), maxEvents int) int {
	nevents := 0
	var bufSize int64 = 8192 * 1024
	buf := make([]byte, bufSize)
	var offset int64
	var nreads int
	for nevents < maxEvents {
		nb, err := r.ReadAt(buf, offset)
		nreads++
		ecs, off := buf2EventClusters(buf[:nb])
		offset += off
		for _, ec := range ecs {
			nevents++
			efunc(ec)
			if nevents >= maxEvents {
				break
			}
		}
		if err == io.EOF {
			break
		}
	}
	fmt.Printf("nreads=%d\n", nreads)
	return nevents
}
