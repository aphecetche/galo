// automatically generated by the FlatBuffers compiler, do not modify

package run2

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Digit struct {
	_tab flatbuffers.Table
}

func GetRootAsDigit(buf []byte, offset flatbuffers.UOffsetT) *Digit {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Digit{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Digit) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Digit) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Digit) Adc() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Digit) MutateAdc(n uint16) bool {
	return rcv._tab.MutateUint16Slot(4, n)
}

func (rcv *Digit) Deid() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Digit) MutateDeid(n uint16) bool {
	return rcv._tab.MutateUint16Slot(6, n)
}

func (rcv *Digit) Manuid() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Digit) MutateManuid(n uint16) bool {
	return rcv._tab.MutateUint16Slot(8, n)
}

func (rcv *Digit) Manuchannel() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Digit) MutateManuchannel(n byte) bool {
	return rcv._tab.MutateByteSlot(10, n)
}

func DigitStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func DigitAddAdc(builder *flatbuffers.Builder, adc uint16) {
	builder.PrependUint16Slot(0, adc, 0)
}
func DigitAddDeid(builder *flatbuffers.Builder, deid uint16) {
	builder.PrependUint16Slot(1, deid, 0)
}
func DigitAddManuid(builder *flatbuffers.Builder, manuid uint16) {
	builder.PrependUint16Slot(2, manuid, 0)
}
func DigitAddManuchannel(builder *flatbuffers.Builder, manuchannel byte) {
	builder.PrependByteSlot(3, manuchannel, 0)
}
func DigitEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type ClusterPos struct {
	_tab flatbuffers.Table
}

func GetRootAsClusterPos(buf []byte, offset flatbuffers.UOffsetT) *ClusterPos {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ClusterPos{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ClusterPos) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ClusterPos) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ClusterPos) X() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *ClusterPos) MutateX(n float32) bool {
	return rcv._tab.MutateFloat32Slot(4, n)
}

func (rcv *ClusterPos) Y() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *ClusterPos) MutateY(n float32) bool {
	return rcv._tab.MutateFloat32Slot(6, n)
}

func (rcv *ClusterPos) Z() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *ClusterPos) MutateZ(n float32) bool {
	return rcv._tab.MutateFloat32Slot(8, n)
}

func ClusterPosStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ClusterPosAddX(builder *flatbuffers.Builder, x float32) {
	builder.PrependFloat32Slot(0, x, 0.0)
}
func ClusterPosAddY(builder *flatbuffers.Builder, y float32) {
	builder.PrependFloat32Slot(1, y, 0.0)
}
func ClusterPosAddZ(builder *flatbuffers.Builder, z float32) {
	builder.PrependFloat32Slot(2, z, 0.0)
}
func ClusterPosEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type PreCluster struct {
	_tab flatbuffers.Table
}

func GetRootAsPreCluster(buf []byte, offset flatbuffers.UOffsetT) *PreCluster {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &PreCluster{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *PreCluster) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PreCluster) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *PreCluster) Digits(obj *Digit, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *PreCluster) DigitsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func PreClusterStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func PreClusterAddDigits(builder *flatbuffers.Builder, digits flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(digits), 0)
}
func PreClusterStartDigitsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func PreClusterEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Cluster struct {
	_tab flatbuffers.Table
}

func GetRootAsCluster(buf []byte, offset flatbuffers.UOffsetT) *Cluster {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Cluster{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Cluster) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Cluster) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Cluster) Pre(obj *PreCluster) *PreCluster {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(PreCluster)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Cluster) Pos(obj *ClusterPos) *ClusterPos {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(ClusterPos)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func ClusterStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func ClusterAddPre(builder *flatbuffers.Builder, pre flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(pre), 0)
}
func ClusterAddPos(builder *flatbuffers.Builder, pos flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(pos), 0)
}
func ClusterEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
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
