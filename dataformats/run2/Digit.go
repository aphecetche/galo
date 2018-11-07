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