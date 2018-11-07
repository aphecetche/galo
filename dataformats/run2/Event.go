// automatically generated by the FlatBuffers compiler, do not modify

package run2

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Event struct {
	_tab flatbuffers.Table
}

func GetRootAsEvent(buf []byte, offset flatbuffers.UOffsetT) *Event {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Event{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Event) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Event) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Event) Bc() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutateBc(n uint16) bool {
	return rcv._tab.MutateUint16Slot(4, n)
}

func (rcv *Event) Period() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutatePeriod(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *Event) Ntracklets() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutateNtracklets(n int32) bool {
	return rcv._tab.MutateInt32Slot(8, n)
}

func (rcv *Event) IsMB() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutateIsMB(n byte) bool {
	return rcv._tab.MutateByteSlot(10, n)
}

func (rcv *Event) Clusters(obj *Cluster, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Event) ClustersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func EventStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func EventAddBc(builder *flatbuffers.Builder, bc uint16) {
	builder.PrependUint16Slot(0, bc, 0)
}
func EventAddPeriod(builder *flatbuffers.Builder, period uint32) {
	builder.PrependUint32Slot(1, period, 0)
}
func EventAddNtracklets(builder *flatbuffers.Builder, ntracklets int32) {
	builder.PrependInt32Slot(2, ntracklets, 0)
}
func EventAddIsMB(builder *flatbuffers.Builder, isMB byte) {
	builder.PrependByteSlot(3, isMB, 0)
}
func EventAddClusters(builder *flatbuffers.Builder, clusters flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(clusters), 0)
}
func EventStartClustersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func EventEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
