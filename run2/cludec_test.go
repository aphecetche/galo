package run2_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/aphecetche/galo"
	"github.com/aphecetche/galo/run2"
	flatbuffers "github.com/google/flatbuffers/go"
)

func createPos(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	run2.ClusterPosStart(builder)
	run2.ClusterPosAddX(builder, float32(n)*0.1)
	run2.ClusterPosAddY(builder, float32(n)*0.2)
	run2.ClusterPosAddZ(builder, float32(n)*0.3)
	return run2.ClusterPosEnd(builder)
}

func createDigit(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	run2.DigitStart(builder)
	run2.DigitAddAdc(builder, uint16(n))
	run2.DigitAddDeid(builder, 100)
	run2.DigitAddManuid(builder, 100+uint16(n))
	run2.DigitAddManuchannel(builder, byte(n))
	return run2.DigitEnd(builder)
}

func createDigits(n int, builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	var digits []flatbuffers.UOffsetT
	for i := 0; i < n; i++ {
		digits = append(digits, createDigit(i, builder))
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
	run2.PreClusterAddDigits(builder, 0)
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
	for i := 0; i < n; i++ {
		clusters = append(clusters, createCluster(n, builder))
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
		event := createEvent(i, builder)
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

func TestCreateEvents(t *testing.T) {

	var tests = []struct {
		decsize int
		nperdec int
		want    int
	}{{1024, 3, 7},
		{1024, 6, 9},
		{2048, 6, 12}}

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
		dec := run2.NewClusterDecoder(br, tp.decsize)

		nread := 0

		ndec := 20
		for {
			var declusters galo.DEClusters
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
	}
}
