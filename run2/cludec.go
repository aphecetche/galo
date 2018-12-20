package run2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/pkg/errors"
)

type run2ClusterDecoder struct {
	r             io.ReaderAt
	buf           []byte
	offset        int64
	curdeclu      int
	declu         []galo.DEClusters
	maxsize       int
	eof           bool
	padfinderfunc mapping.PadByFEEFinderFunc
}

var _ galo.DEClustersDecoder = (*run2ClusterDecoder)(nil)

func NewClusterDecoder(src io.ReaderAt, padfinderfunc mapping.PadByFEEFinderFunc, bufSize int) *run2ClusterDecoder {
	if bufSize == 0 {
		bufSize = 8192
	}
	return &run2ClusterDecoder{
		r:             src,
		buf:           make([]byte, bufSize),
		offset:        0,
		maxsize:       2 << 23, // around 8MB
		declu:         nil,
		eof:           false,
		padfinderfunc: padfinderfunc}
}

func getDEClusters(e *Event, padfinderfunc mapping.PadByFEEFinderFunc) *galo.DEClusters {
	var clusters []galo.Cluster
	var ci Cluster
	var deid mapping.DEID
	for i := 0; i < e.ClustersLength(); i++ {
		e.Clusters(&ci, i)
		prei := ci.Pre(nil)
		var pre galo.PreCluster
		var d Digit
		for j := 0; j < prei.DigitsLength(); j++ {
			prei.Digits(&d, j)
			if deid == 0 {
				deid = mapping.DEID(d.Deid())
			} else {
				if deid != mapping.DEID(d.Deid()) {
					panic("mixing detection element ids")
				}
			}
			padfinder := padfinderfunc(deid)
			id, err := padfinder.FindPadByFEE(mapping.DualSampaID(d.Manuid()),
				mapping.DualSampaChannelID(d.Manuchannel()))
			if err != nil {
				panic(err)
			}
			pre.Digits = append(pre.Digits, galo.Digit{ID: id, Q: float64(d.Charge())})
		}
		clu := galo.Cluster{
			Q:   galo.ClusterCharge(ci.Charge()),
			Pos: galo.ClusterPos{X: float64(ci.Pos(nil).X()), Y: float64(ci.Pos(nil).Y())},
			Pre: pre,
		}
		clusters = append(clusters, clu)
	}
	return &galo.DEClusters{
		DeID:     deid,
		Clusters: clusters,
	}
}
func buf2DEClusters(buf []byte, padfinderfunc mapping.PadByFEEFinderFunc) ([]galo.DEClusters, int64) {
	var off int64
	var declu []galo.DEClusters
	var pstart, pend int64
	var size uint32
	for off < int64(len(buf)) {
		size = binary.LittleEndian.Uint32(buf[off : off+4])
		pstart = off + 4
		pend = pstart + int64(size)
		if pend > int64(len(buf)) {
			break
		}
		event := GetRootAsEvent(buf[pstart:pend], 0)
		dc := getDEClusters(event, padfinderfunc)
		declu = append(declu, *dc)
		off = pend
	}
	return declu, off
}

func (dec *run2ClusterDecoder) Decode(declu *galo.DEClusters) error {
	if dec.eof {
		return io.EOF
	}

	if dec.curdeclu >= len(dec.declu) {
		n, err := dec.r.ReadAt(dec.buf, dec.offset)
		if n == 0 {
			return errors.Wrap(err, "No data read in")
		}
		if err != nil && err != io.EOF {
			return err
		}
		clusters, offset := buf2DEClusters(dec.buf[:n], dec.padfinderfunc)
		if offset == 0 {
			// reading buffer too small, try to
			// increase it
			// as long as we don't go over maxsize
			newsize := len(dec.buf) * 2
			if newsize < dec.maxsize {
				dec.buf = make([]byte, newsize)
				return dec.Decode(declu)
			}
			return fmt.Errorf("Event too big for our buffer")
		}
		dec.declu = clusters
		dec.offset += offset
		dec.curdeclu = 0
		if err == io.EOF {
			dec.eof = true
		}
	}
	*declu = dec.declu[dec.curdeclu]
	dec.curdeclu++
	return nil
}

func (dec *run2ClusterDecoder) Close() {
}
