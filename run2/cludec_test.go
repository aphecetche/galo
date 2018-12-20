package run2_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/run2"
	"github.com/aphecetche/pigiron/mapping"
	_ "github.com/aphecetche/pigiron/mapping/impl4"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/pkg/errors"
)

func createPos(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	run2.ClusterPosStart(builder)
	run2.ClusterPosAddX(builder, float32(n)*0.1)
	run2.ClusterPosAddY(builder, float32(n)*0.2)
	run2.ClusterPosAddZ(builder, float32(n)*0.3) // to be deprecated
	return run2.ClusterPosEnd(builder)
}

func createDigit(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	run2.DigitStart(builder)
	run2.DigitAddAdc(builder, uint16(n)) // not tested
	run2.DigitAddCharge(builder, float32(n)+0.42)
	run2.DigitAddDeid(builder, 100)
	run2.DigitAddManuid(builder, 100+uint16(n))
	run2.DigitAddManuchannel(builder, byte(n))
	return run2.DigitEnd(builder)
}

func createDigits(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	var digits []flatbuffers.UOffsetT
	for i := 1; i <= n; i++ {
		digits = append(digits, createDigit(n-i+1, builder))
	}
	run2.PreClusterStartDigitsVector(builder, n)
	for _, d := range digits {
		builder.PrependUOffsetT(d)
	}
	return builder.EndVector(n)
}

func createPre(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	digits := createDigits(n, builder)
	run2.PreClusterStart(builder)
	run2.PreClusterAddDigits(builder, digits)
	return run2.PreClusterEnd(builder)
}

func createCluster(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	pos := createPos(n, builder)
	pre := createPre(n, builder)
	run2.ClusterStart(builder)
	run2.ClusterAddPre(builder, pre)
	run2.ClusterAddPos(builder, pos)
	run2.ClusterAddCharge(builder, float32(n)*100+0.42)
	return run2.ClusterEnd(builder)
}

func createClusters(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	var clusters []flatbuffers.UOffsetT
	for i := 1; i <= n; i++ {
		clusters = append(clusters, createCluster(n-i+1, builder))
	}
	run2.EventStartClustersVector(builder, n)
	for _, c := range clusters {
		builder.PrependUOffsetT(c)
	}
	return builder.EndVector(n)
}

func createEvent(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	// Event i will have i clusters of i digits.

	clusters := createClusters(n, builder)

	run2.EventStart(builder)
	run2.EventAddBc(builder, uint16(1000+n))
	run2.EventAddPeriod(builder, uint32(100+n))
	run2.EventAddNtracklets(builder, int32(n*10))
	run2.EventAddIsMB(builder, byte(n%2))
	run2.EventAddClusters(builder, clusters)

	return run2.EventEnd(builder)
}

// createFakeEvents creates enough events to fill
// at least nperdec * decoderBufferSize bytes of a buffer.
// The events are fake ones : always for DE 100
// i-th event has i clusters which has i digits
func createFakeEvents(decoderBufferSize int, nperdec int) ([]byte, int) {

	var buf bytes.Buffer

	builder := flatbuffers.NewBuilder(1024)

	maxBufSize := nperdec * decoderBufferSize

	size := make([]byte, 4)
	i := 0
	for {
		builder.Reset()
		event := createEvent(i+1, builder)
		builder.Finish(event)
		eventBuf := builder.FinishedBytes()
		binary.LittleEndian.PutUint32(size, uint32(len(eventBuf)))
		buf.Write(size)
		buf.Write(eventBuf)
		i++
		if buf.Len()+len(size)+len(eventBuf) > maxBufSize {
			break
		}
	}

	return buf.Bytes(), i
}

func compareFloats(msg string, got, want float64) error {
	if math.Abs((got-want)/want) > 1E-6 {
		return fmt.Errorf("Want %s=%7.2f Got %7.2f", msg, want, got)
	}
	return nil
}
func compareInts(msg string, got, want int) error {
	if got != want {
		return fmt.Errorf("Want %s=%v Got %v", msg, want, got)
	}
	return nil
}

func clustersAsExpected(declu *galo.DEClusters, nbase int) error {
	if len(declu.Clusters) != nbase {
		return fmt.Errorf("wrong number of clusters")
	}
	if declu.DeID != 100 {
		return fmt.Errorf("Want Deid %d - Got %d", 100, declu.DeID)
	}
	n := 1
	padloc := galo.SegCache.Segmentation(100)
	for _, clu := range declu.Clusters {
		msg := fmt.Sprintf("Cluster %d", n)

		err := compareFloats("Q", float64(clu.Q), 100.0*float64(n)+0.42)
		if err != nil {
			return errors.Wrap(err, msg)
		}
		err = compareFloats("X", float64(clu.Pos.X), 0.1*float64(n))
		if err != nil {
			return errors.Wrap(err, msg)
		}
		err = compareFloats("Y", float64(clu.Pos.Y), 0.2*float64(n))
		if err != nil {
			return errors.Wrap(err, msg)
		}

		err = compareInts("NofPads", clu.Pre.NofPads(), n)
		if err != nil {
			return errors.Wrap(err, msg)
		}

		for i := 0; i < clu.Pre.NofPads(); i++ {
			msg := fmt.Sprintf("Cluster %d Digit %d", n, i)
			d := clu.Pre.Digits[i]
			err := compareFloats("Q", float64(d.Q), float64(i+1)+0.42)
			if err != nil {
				return errors.Wrap(err, msg)
			}
			paduid := mapping.PadUID(d.ID)
			dsid := padloc.PadDualSampaID(paduid)
			err = compareInts("DsId", int(dsid), 100+i+1)
			if err != nil {
				return errors.Wrap(err, msg)
			}
			dsch := padloc.PadDualSampaChannel(paduid)
			err = compareInts("DsCh", int(dsch), i+1)
			if err != nil {
				return errors.Wrap(err, msg)
			}
		}

		n++
	}
	return nil
}

func TestCreateEvents(t *testing.T) {

	var tests = []struct {
		decsize int
		nperdec int
		want    int
	}{{1024, 3, 6},
		{1024, 6, 8},
		{2048, 6, 11}}
	for _, tp := range tests {
		// generate a flatbuffer containing enough Events
		// to get bigger than the decoder reading buffer size (decsize),
		// and then feed it to the decoder to check we're getting back
		// all our clusters
		buf, n := createFakeEvents(tp.decsize, tp.nperdec)

		if n != tp.want {
			t.Errorf("Want %v events. Got %d", tp.want, n)
		}

		br := bytes.NewReader(buf)
		dec := run2.NewClusterDecoder(br,
			func(deid mapping.DEID) mapping.PadByFEEFinder {
				return galo.SegCache.Segmentation(deid)
			}, tp.decsize)

		nread := 0

		ndec := 20
		var declusters galo.DEClusters
		for {
			err := dec.Decode(&declusters)
			if err != nil {
				break
			}
			ndec--
			if ndec == 0 {
				break
			}
			nread++
		}

		if ndec == 0 {
			t.Errorf("Reached max number of decoding")
		}
		if n != nread {
			t.Errorf("Wanted to read %d clusters. Got %d", n, nread)
		}

		if err := clustersAsExpected(&declusters, nread); err != nil {
			t.Errorf("Clusters not as expected for nread=%d err=%v", nread, err.Error())
		}
	}
}
