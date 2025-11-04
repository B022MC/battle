package game

import (
	"errors"
	"unicode/utf16"
)

const (
	SizeTcpBuffer = 1024 * 1024
)

type Packer struct {
	Head     CmdHead
	buffer   [SizeTcpBuffer]byte
	BufSize  int
	dataSize int
}

type CmdInfo struct {
	version    byte
	checkCode  byte
	packetSize uint16
}

type CmdCommand struct {
	MainCmdID uint16
	SubCmdID  uint16
}

type CmdHead struct {
	info CmdInfo
	Cmd  CmdCommand
}

func (p *Packer) SetCmd(main int, sub int) {
	p.PushDWord(0)
	p.PushWord(uint16(main))
	p.PushWord(uint16(sub))
	p.dataSize = 0
	p.Head.Cmd.MainCmdID = uint16(main)
	p.Head.Cmd.SubCmdID = uint16(sub)
}

func (p *Packer) Load(data []byte) error {
	if len(data) < 8 {
		return errors.New("invalid data")
	}

	p.BufSize = len(data)
	copy(p.buffer[:], data)
	// for i := range data {
	// 	p.buffer[i] = data[i]
	// }

	p.Head.info.version = p.ReadByte(0)
	p.Head.info.checkCode = p.ReadByte(1)
	p.Head.info.packetSize = p.ReadWord(2)
	p.dataSize = int(p.Head.info.packetSize - 8)
	p.BufSize = int(p.Head.info.packetSize)
	return nil
}

func (p *Packer) BufferSize() int {
	return p.BufSize
}

func (p *Packer) SetBufferSize(size int) {
	p.BufSize = size
}

func (p *Packer) DataSize() int {
	return p.dataSize
}

func (p *Packer) Bytes() []byte {
	var res []byte
	for i := 0; i < p.BufSize; i++ {
		res = append(res, p.buffer[i])
	}

	return res
}

func (p *Packer) Data() []byte {
	return p.buffer[8:p.BufSize]
}

func (p *Packer) ReadByte(offset int) byte {
	return p.buffer[offset]
}

func (p *Packer) ReadWord(offset int) uint16 {
	lo := uint16(p.buffer[offset])
	offset++
	hi := uint16(p.buffer[offset])
	return hi<<8 + lo
}

func (p *Packer) ReadDWord(offset int) uint32 {
	b1 := uint32(p.buffer[offset])
	offset++
	b2 := uint32(p.buffer[offset])
	offset++
	b3 := uint32(p.buffer[offset])
	offset++
	b4 := uint32(p.buffer[offset])
	return b4<<24 + b3<<16 + b2<<8 + b1
}

func (p *Packer) WriteByte(offset int, value byte) {
	p.buffer[offset] = value
}

func (p *Packer) WriteDWord(offset int, data uint32) {
	b := data & 0xff
	p.buffer[offset] = byte(b)
	offset++

	b = (data & 0xff00) >> 8
	p.buffer[offset] = byte(b)
	offset++

	b = (data & 0xff0000) >> 16
	p.buffer[offset] = byte(b)
	offset++

	b = (data & 0xff000000) >> 24
	p.buffer[offset] = byte(b)
}

func (p *Packer) WriteWord(offset int, data uint16) {
	hi := (data & 0xff00) >> 8
	low := data & 0xff
	p.buffer[offset] = byte(low)
	p.buffer[offset+1] = byte(hi)
}

func (p *Packer) PushByte(data byte) {
	p._pushByte(data)
}

func (p *Packer) PushBytes(data []byte) {
	for i := len(data) - 1; i >= 0; i-- {
		p._pushByte(data[i])
	}
}

func (p *Packer) PushString(data string, sz int) {
	var bytes [2]byte
	runes := utf16.Encode([]rune(data))
	for _, r := range runes {
		bytes[1] = byte(r >> 8)
		bytes[0] = byte(r & 255)
		p._pushByte(bytes[0])
		p._pushByte(bytes[1])
	}
	diff := sz - len([]byte(data))
	for i := 0; i < diff; i++ {
		p.PushWord(0)
	}
}

func (p *Packer) ReadString(offset int, sz int) string {
	var codes []uint16
	for i := 0; i < sz/2; i++ {
		b0 := p.ReadByte(offset)
		offset++
		b1 := p.ReadByte(offset)
		offset++

		r := (uint16(b1)&255)<<8 + uint16(b0)
		codes = append(codes, r)
	}

	runes := utf16.Decode(codes)
	return string(runes)
}

func (p *Packer) _pushByte(data byte) {
	p.buffer[p.BufSize] = data
	p.BufSize++
	p.dataSize++
}

func (p *Packer) PushWord(data uint16) {
	low := data & 0xff
	p._pushByte(byte(low))
	hi := (data & 0xff00) >> 8
	p._pushByte(byte(hi))
}

func (p *Packer) PushDWord(data uint32) {
	b := data & 0xff
	p._pushByte(byte(b))

	b = (data & 0xff00) >> 8
	p._pushByte(byte(b))

	b = (data & 0xff0000) >> 16
	p._pushByte(byte(b))

	b = (data & 0xff000000) >> 24
	p._pushByte(byte(b))
}

func (p *Packer) InsertDWord(offset int, data uint32) {
	p.BufSize += 4
	for i := int(p.BufSize) - 1; i >= offset; i-- {
		p.buffer[i] = p.buffer[i-4]
	}

	b := data & 0xff
	p.buffer[offset] = byte(b)
	offset++

	b = (data & 0xff00) >> 8
	p.buffer[offset] = byte(b)
	offset++

	b = (data & 0xff0000) >> 16
	p.buffer[offset] = byte(b)
	offset++

	b = (data & 0xff000000) >> 24
	p.buffer[offset] = byte(b)
	offset++
}

func (p *Packer) PaddingTo(sz int) {
	diff := sz - p.dataSize
	for i := 0; i < diff; i++ {
		p._pushByte(0)
	}
}
