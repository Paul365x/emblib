package pes_pec

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"os"

	"github.com/emblib/adapters/shared"
)

/*
**
** API stuff
**
 */
// parse_color_sub parses in a color structure
func parse_color_sub(bin []byte) (uint32, shared.ColorSub) {
	var col shared.ColorSub
	var count uint32
	count = 0
	col.CodeLen = bin[count]
	count++
	col.Code = bin[count : count+uint32(col.CodeLen)]
	count += uint32(col.CodeLen)
	red := bin[count]
	count++
	green := bin[count]
	count++
	blue := bin[count]
	count++
	col.Color = color.RGBA{red, green, blue, 255}
	col.U1 = bin[count]
	count++
	col.ColType = binary.LittleEndian.Uint32(bin[count : count+4])
	count += 4
	col.DescLen = bin[count]
	count++
	col.Desc = string(bin[count : count+uint32(col.DescLen)])
	count += uint32(col.DescLen)
	col.BrandLen = bin[count]
	count++
	col.Brand = string(bin[count : count+uint32(col.BrandLen)])
	count += uint32(col.BrandLen)
	col.ChartLen = bin[count]
	count++
	col.Chart = string(bin[count : count+uint32(col.ChartLen)])
	count += uint32(col.ChartLen)
	col.Count = count
	return count, col
}

/*
**
** Pes Header parsing code. While we draw using the pec format, the pes headers have some useful information
**
 */

// Preamble stores the first 12 bytes of a pes file
type Preamble struct {
	Id     string
	Ver    string
	Offset uint32
	count  uint32
}

// Preamble.Parse reads in the first 12 bytes of a pes file into the struct
func (s *Preamble) Parse(bin []byte) {
	s.Id = string(bin[:4])
	s.Ver = string(bin[4:8])
	s.Offset = binary.LittleEndian.Uint32(bin[8:12])
	s.count = 12
}

// Preamble.SizeOf returns the size in bytes - offset into the file of the byte after the preamble. Always 12 bytes
func (p Preamble) SizeOf() uint32 {
	return p.count
}

// Dump writes out this Struct
func (p Preamble) Dump() {
	fmt.Printf("Preamble:\n")
	fmt.Printf("\tId: %s\n", p.Id)
	fmt.Printf("\tVer: %s\n", p.Ver)
	fmt.Printf("\tOffset: %d 0x%X\n", p.Offset, p.Offset)
	fmt.Printf("\tcount: %d 0x%X\n", p.count, p.count)
}

// H_1 is version 1 header struct of a pes file
type H_1 struct {
	Hoop      uint16
	EDA       uint16
	Blk_count uint16
	count     uint32
}

// H_1.Parse in the version one header of a pes file
func (h1 *H_1) Parse(bin []byte) {
	h1.Hoop = binary.LittleEndian.Uint16(bin[0:2])
	h1.EDA = binary.LittleEndian.Uint16(bin[2:4])
	h1.Blk_count = binary.LittleEndian.Uint16(bin[4:6])
	h1.count = 6
}

// H_1.SizeOf returns the byte offset into the file of the next byte to read - always 6 bytes
func (h1 H_1) SizeOf() uint32 {
	return h1.count
}

// Dump writes out this Struct
func (p H_1) Dump() {
	fmt.Printf("Header1:\n")
	fmt.Printf("\tHoop: %d\n", p.Hoop)
	fmt.Printf("\tEDA: %d\n", p.EDA)
	fmt.Printf("\tBlk_count: %d 0x%X\n", p.Blk_count, p.Blk_count)
	fmt.Printf("\tcount: %d 0x%X\n", p.count, p.count)
}

// H_2 stores the version 2 header
type H_2 struct {
	HoopW uint16
	HoopH uint16
	Rot   uint16
	unk   []byte
	count uint32
}

// H_2.Parse parses a version 2 header into the struct
func (h2 *H_2) Parse(bin []byte) {
	h2.HoopW = binary.LittleEndian.Uint16(bin[0:2])
	h2.HoopH = binary.LittleEndian.Uint16(bin[2:4])
	h2.Rot = binary.LittleEndian.Uint16(bin[4:6])
	h2.unk = bin[6:24]
	h2.count = 24
}

// H_2.SizeOf returns the offset into the file of the next byte after the header - always 24 bytes
func (h2 H_2) SizeOf() uint32 {
	return h2.count
}

// Dump writes out this Struct
func (p H_2) Dump() {
	fmt.Printf("Header2:\n")
	fmt.Printf("\tHoopW: %d 0x%X\n", p.HoopW, p.HoopW)
	fmt.Printf("\tHoopH: %d 0x%X\n", p.HoopH, p.HoopH)
	fmt.Printf("\tRot: %d 0x%X\n", p.Rot, p.Rot)
	fmt.Printf("\tunk: 0x%X\n", p.unk)
	fmt.Printf("\tcount: %d 0x%X\n", p.count, p.count)
}

// H_3 version three of the header
type H_3 struct {
	u1    uint16
	SubV  uint16
	HoopW uint16
	HoopH uint16
	Rot   uint16
	u2    []byte
	count uint32
}

// H_3.Parse version three header parser
func (h3 *H_3) Parse(bin []byte) {
	h3.u1 = binary.LittleEndian.Uint16(bin[0:2])
	h3.SubV = binary.LittleEndian.Uint16(bin[2:4])
	h3.HoopW = binary.LittleEndian.Uint16(bin[4:6])
	h3.HoopH = binary.LittleEndian.Uint16(bin[6:8])
	h3.Rot = binary.LittleEndian.Uint16(bin[8:10])
	h3.u2 = bin[10:28]
	h3.count = 28
}

// H_3.SizeOf returns the offset into the file of the next byte - always 28 bytes
func (h3 H_3) SizeOf() uint32 {
	return h3.count
}

// Dump writes out this Struct
func (p H_3) Dump() {
	fmt.Printf("Header3:\n")
	fmt.Printf("\tu1: %d 0x%X\n", p.u1, p.u1)
	fmt.Printf("\tSubV: %d 0x%X\n", p.SubV, p.SubV)
	fmt.Printf("\tHoopW: %d 0x%X\n", p.HoopW, p.HoopW)
	fmt.Printf("\tHoopH: %d 0x%X\n", p.HoopH, p.HoopH)
	fmt.Printf("\tRot: %d 0x%X\n", p.Rot, p.Rot)
	fmt.Printf("\tu2: %d 0x%X\n", p.u2, p.u2)
	fmt.Printf("\tcount: %d 0x%X\n", p.count, p.count)
}

// H_4 is version 4 of the header
type H_4 struct {
	u1    uint16
	SubV  uint16
	Desc  *map[string]string
	u2    uint16
	HoopW uint16
	HoopH uint16
	Rot   uint16
	u3    []byte
	count uint32
}

// H_4.Parse is the version 4 header parser
func (h4 *H_4) Parse(bin []byte) {
	h4.u1 = binary.LittleEndian.Uint16(bin[0:2])
	h4.SubV = binary.LittleEndian.Uint16(bin[2:4])
	var count uint32
	count, h4.Desc = parse_desc(bin[4:])
	count += 4 // allow for u1 and SubV
	h4.u2 = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h4.HoopW = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h4.HoopH = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h4.Rot = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h4.u3 = bin[count : count+22]
	count += 22
	h4.count = count
}

// H_4 SizeOf returns the offset into the file of the byte after this header
func (h4 H_4) SizeOf() uint32 {
	return h4.count
}

// Dump writes out this Struct
func (p H_4) Dump() {
	fmt.Printf("Header4:\n")
	fmt.Printf("\tu1: %d 0x%X\n", p.u1, p.u1)
	fmt.Printf("\tSubV: %d 0x%X\n", p.SubV, p.SubV)
	fmt.Printf("\tDesc: ")
	fmt.Println(p.Desc)
	fmt.Printf("\tu2: %d 0x%X\n", p.u2, p.u2)
	fmt.Printf("\tHoopW: %d 0x%X\n", p.HoopW, p.HoopW)
	fmt.Printf("\tHoopH: %d 0x%X\n", p.HoopH, p.HoopH)
	fmt.Printf("\tRot: %d 0x%X\n", p.Rot, p.Rot)
	fmt.Printf("\tu3: 0x%X\n", p.u3)
	fmt.Printf("\tcount: %d 0x%X\n", p.count, p.count)
}

//
// Headers 5 & 6 - bit more complex so broken down into parts
//

// HP_1 is first section of headers 5 and 6
type HP_1 struct {
	HoopInd uint16
	SubV    uint16
	Desc    *map[string]string
	HoopChg uint16
	count   uint32
}

// HP_1.Parse parses the first section of headers 5 and 6
func (h *HP_1) Parse(bin []byte) {
	h.HoopInd = binary.LittleEndian.Uint16(bin[0:2])
	h.SubV = binary.LittleEndian.Uint16(bin[2:4])
	var count uint32
	count, h.Desc = parse_desc(bin[4:])
	count += 4
	h.HoopChg = binary.LittleEndian.Uint16(bin[count : count+2])
	h.count = count + 2
}

// HP_1.SizeOf returns the offset into the file of the next byte after this section
func (h *HP_1) SizeOf() uint32 {
	return h.count
}

// Dump writes out this Struct
func (h HP_1) Dump() {
	fmt.Printf("\tHoopInd: %d 0x%X\n", h.HoopInd, h.HoopInd)
	fmt.Printf("\tSubV: %d 0x%X\n", h.SubV, h.SubV)
	fmt.Printf("\tDesc: ")
	fmt.Println(h.Desc)
	fmt.Printf("\tHoopChg: %d 0x%X\n", h.HoopChg, h.HoopChg)
}

// HP_2 is second section of headers 5 and 6
type HP_2 struct {
	HoopH uint16
	HoopW uint16
	Rot   uint16
	count uint32
}

// HP_2.Parse parses the second section of headers 5 and 6
func (h *HP_2) Parse(bin []byte) {
	h.HoopW = binary.LittleEndian.Uint16(bin[0:2])
	h.HoopH = binary.LittleEndian.Uint16(bin[2:4])
	h.Rot = binary.LittleEndian.Uint16(bin[4:6])
	h.count = 6
}

// HP_2.SizeOf returns the offset into the file of the next byte after this section
func (h *HP_2) SizeOf() uint32 {
	return h.count
}

// Dump writes out this Struct
func (h HP_2) Dump() {
	fmt.Printf("\tHoopW: %d 0x%X\n", h.HoopW, h.HoopW)
	fmt.Printf("\tHoopH: %d 0x%X\n", h.HoopH, h.HoopH)
	fmt.Printf("\tRot: %d 0x%X\n", h.Rot, h.Rot)
}

// HP_3 is third section of headers 5 and 6
type HP_3 struct {
	BG       uint16
	FG       uint16
	Grid     uint16
	Axes     uint16
	Snap     uint16
	Interv   uint16
	u1       uint16
	OptEntEx uint16
	Imlen    uint8
	Impath   string
	Affline  []byte
	count    uint32
}

// HP_3.Parse reads the third section of headers 5 and 6
func (h *HP_3) Parse(bin []byte) {
	h.BG = binary.LittleEndian.Uint16(bin[0:2])
	h.FG = binary.LittleEndian.Uint16(bin[2:4])
	h.Grid = binary.LittleEndian.Uint16(bin[4:6])
	h.Axes = binary.LittleEndian.Uint16(bin[6:8])
	h.Snap = binary.LittleEndian.Uint16(bin[8:10])
	h.Interv = binary.LittleEndian.Uint16(bin[10:12])
	h.u1 = binary.LittleEndian.Uint16(bin[12:14])
	h.OptEntEx = binary.LittleEndian.Uint16(bin[14:16])
	h.Imlen = bin[16]
	o := h.Imlen + 17
	h.Impath = string(bin[17:o])
	h.Affline = bin[o : o+24]
	h.count = uint32(o) + 24
}

// HP_3 returns the offset into the file after the third section of headers 5 and 6
func (h *HP_3) SizeOf() uint32 {
	return h.count
}

// Dump writes out this Struct
func (h HP_3) Dump() {
	fmt.Printf("\tBackground: %d 0x%X\n", h.BG, h.BG)
	fmt.Printf("\tForeground: %d 0x%X\n", h.FG, h.FG)
	fmt.Printf("\tGrid: %d 0x%X\n", h.Grid, h.Grid)
	fmt.Printf("\tAxes: %d 0x%X\n", h.Axes, h.Axes)
	fmt.Printf("\tSnap: %d 0x%X\n", h.Snap, h.Snap)
	fmt.Printf("\tInterval: %d 0x%X\n", h.Interv, h.Interv)
	fmt.Printf("\tu1: %d 0x%X\n", h.u1, h.u1)
	fmt.Printf("\tOptEntEx: %d 0x%X\n", h.OptEntEx, h.OptEntEx)
	fmt.Printf("\tImg Path Len: %d 0x%X\n", h.Imlen, h.Imlen)
	fmt.Printf("\tImg Path: %s\n", h.Impath)
	fmt.Printf("\tAffline: %X\n", h.Affline)
}

// HP_4 is fourth section of headers 5 and 6
type HP_4 struct {
	FillCount  uint16
	Fill       []byte
	MotCount   uint16
	Motif      []byte
	FeathCount uint16
	Feather    []byte
	ColSects   uint16
	Colors     []shared.ColorSub
	Obj        uint16
	count      uint32
}

// HP_4.Parse reads the fourth section of headers 5 and 6
func (h *HP_4) Parse(bin []byte) {
	end := binary.LittleEndian.Uint16(bin[0:2])
	h.Fill = bin[2 : end+2]
	start := end + 2
	end = binary.LittleEndian.Uint16(bin[start : start+2])
	start += 2
	end += start
	h.Motif = bin[start:end]
	start = end
	end = binary.LittleEndian.Uint16(bin[start : start+2])
	start += 2
	end += start
	h.Feather = bin[start:end]
	start = end
	num := binary.LittleEndian.Uint16(bin[start : start+2])
	h.ColSects = num
	start += 2
	end += start
	for i := 0; i < int(num); i++ {
		count, col := parse_color_sub(bin[start:])
		h.Colors = append(h.Colors, col)
		start += uint16(count)
	}
	h.Obj = binary.LittleEndian.Uint16(bin[start : start+2])
	h.count = uint32(start + 2)
}

// HP_4.SizeOf returns the offset into the file of the byte after the fourth section of headers 5 and 6
func (h *HP_4) SizeOf() uint32 {
	return h.count
}

// Dump writes out this Struct
func (h HP_4) Dump() {
	fmt.Printf("\tFillCount: %d 0x%X\n", h.FillCount, h.FillCount)
	fmt.Printf("\tFill: %X", h.Fill)
	fmt.Printf("\tMotCount: %d 0x%X\n", h.MotCount, h.MotCount)
	fmt.Printf("\tMotif: %X", h.Motif)
	fmt.Printf("\tFeatherCount: %d 0x%X\n", h.FeathCount, h.FeathCount)
	fmt.Printf("\tFeather: %X\n", h.Feather)
	fmt.Printf("\tColorSects: %d 0x%X\n", h.ColSects, h.ColSects)
	for i := 0; i < int(h.ColSects); i++ {
		h.Colors[i].Dump()
	}
	fmt.Printf("\tObjects: %d 0x%X\n", h.Obj, h.Obj)
}

// H_5 version 5 of the header
type H_5 struct {
	HP_1
	HP_2
	HP_3
	HP_4
	count uint32
}

// H_5.Parse parses the version 5 header
func (h *H_5) Parse(bin []byte) {
	count := uint32(0)

	h.HP_1.Parse(bin)
	count += h.HP_1.SizeOf()

	h.HP_2.Parse(bin[count:])
	count += h.HP_2.SizeOf()

	h.HP_3.Parse(bin[count:])
	count += h.HP_3.SizeOf()

	h.HP_4.Parse(bin[count:])
	count += h.HP_4.SizeOf()

	h.count = count
}

// H_5.SizeOf returns the offset into the file of the next byte after this header
func (h H_5) SizeOf() uint32 {
	return h.count
}

// Dump displays this struct
func (h H_5) Dump() {
	fmt.Println("Header5:")
	h.HP_1.Dump()
	h.HP_2.Dump()
	h.HP_3.Dump()
	h.HP_4.Dump()
	fmt.Printf("\tcount: %d 0x%X\n", h.count, h.count)
}

// H_6 version header
type H_6 struct {
	HP_1
	Cust uint16
	HP_2
	DWidth   uint16
	DHeight  uint16
	DPWidth  uint16
	DPHeight uint16
	u1       uint16
	HP_3
	HP_4
	count uint32
}

// H_6.Parse reads the version 6 header
func (h *H_6) Parse(bin []byte) {
	count := uint32(0)

	h.HP_1.Parse(bin)
	count += h.HP_1.SizeOf()

	h.Cust = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2

	h.HP_2.Parse(bin[count:])
	count += h.HP_2.SizeOf()

	h.DWidth = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h.DHeight = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h.DPWidth = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h.DPHeight = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h.u1 = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2

	h.HP_3.Parse(bin[count:])
	count += h.HP_3.SizeOf()

	h.HP_4.Parse(bin[count:])
	count += h.HP_4.SizeOf()

	h.count = count
}

// H_6.SizeOf returns the next byte in the file after this header
func (h6 H_6) SizeOf() uint32 {
	return h6.count
}

// Dump displays this struct
func (h H_6) Dump() {
	fmt.Println("Header6:")
	h.HP_1.Dump()
	fmt.Printf("\tCust: %d 0x%X\n", h.Cust, h.Cust)
	h.HP_2.Dump()
	fmt.Printf("\tDWidth: %d 0x%X\n", h.DWidth, h.DWidth)
	fmt.Printf("\tDHeight: %d 0x%X\n", h.DHeight, h.DHeight)
	fmt.Printf("\tDPWidth: %d 0x%X\n", h.DPWidth, h.DPWidth)
	fmt.Printf("\tDPHeight: %d 0x%X\n", h.DPHeight, h.DPHeight)
	fmt.Printf("\tu1: %d 0x%X\n", h.u1, h.u1)
	h.HP_3.Dump()
	h.HP_4.Dump()
}

// Header stores the pes header in all forms
type Header struct {
	Ver   string
	P     Preamble
	H1    H_1
	H2    H_2
	H3    H_3
	H4    H_4
	H5    H_5
	H6    H_6
	count uint32
	tail  uint32
}

// Header.Parse reads in the pes header
func (Hdr *Header) Parse(bin []byte) {
	ver := string(bin[4:8])
	Hdr.Ver = ver
	switch ver {
	case "0001":
		var h_p Preamble
		var h1 H_1
		h_p.Parse(bin)
		count_p := h_p.SizeOf()
		h1.Parse(bin[count_p:])
		Hdr.P = h_p
		Hdr.H1 = h1
		Hdr.count = h1.SizeOf() + h_p.SizeOf()
	case "0020":
		var h_p Preamble
		var h H_2
		h_p.Parse(bin)
		count_p := h_p.SizeOf()
		h.Parse(bin[count_p:])
		Hdr.P = h_p
		Hdr.H2 = h
		Hdr.count = h.SizeOf() + h_p.SizeOf()
	case "0030":
		var h_p Preamble
		var h H_3
		h_p.Parse(bin)
		count_p := h_p.SizeOf()
		h.Parse(bin[count_p:])
		Hdr.P = h_p
		Hdr.H3 = h
		Hdr.count = h.SizeOf() + h_p.SizeOf()
	case "0040":
		var h_p Preamble
		var h H_4
		h_p.Parse(bin)
		count_p := h_p.SizeOf()
		h.Parse(bin[count_p:])
		Hdr.P = h_p
		Hdr.H4 = h
		Hdr.count = h.SizeOf() + h_p.SizeOf()
	case "0050":
		var h_p Preamble
		var h H_5
		h_p.Parse(bin)
		count_p := h_p.SizeOf()
		h.Parse(bin[count_p:])
		count_p += h.SizeOf()
		Hdr.P = h_p
		Hdr.H5 = h
		Hdr.count = h.SizeOf() + h_p.SizeOf()
	case "0060":
		var h_p Preamble
		var h H_6
		h_p.Parse(bin)
		count_p := h_p.SizeOf()
		h.Parse(bin[count_p:])
		count_p += h.SizeOf()
		Hdr.P = h_p
		Hdr.H6 = h
		Hdr.count = h.SizeOf() + h_p.SizeOf()
	}
	Hdr.tail = binary.LittleEndian.Uint32(bin[Hdr.count : Hdr.count+4])
	Hdr.count += 4
}

// Header.SizeOf returns the offset of the next byte after the header
func (h *Header) SizeOf() uint32 {
	return h.count
}

// Dump writes out this Struct
func (h Header) Dump() {
	h.P.Dump()
	switch h.Ver {
	case "0001":
		h.H1.Dump()
	case "0020":
		h.H2.Dump()
	case "0030":
		h.H3.Dump()
	case "0040":
		h.H4.Dump()
	case "0050":
		h.H5.Dump()
	case "0060":
		h.H6.Dump()
	}
	fmt.Printf("\ttail: %x\n", h.tail)
	fmt.Printf("\tcount: %d 0x%X\n", h.count, h.count)
}

// Helpers for header parsing
//
// parse_desc H_4 and following have a description block. this is the parser for that block
func parse_desc(bin []byte) (uint32, *map[string]string) {
	var len uint8
	var count uint32

	meta := make(map[string]string)
	count = 0
	len = uint8(bin[count])
	count++
	meta["Design"] = string(bin[count : count+uint32(len)])
	count = count + uint32(len)

	len = uint8(bin[count])
	count++
	meta["Category"] = string(bin[count : count+uint32(len)])
	count = count + uint32(len)

	len = uint8(bin[count])
	count++
	meta["Author"] = string(bin[count : count+uint32(len)])
	count = count + uint32(len)

	len = uint8(bin[count])
	count++
	meta["Keywords"] = string(bin[count : count+uint32(len)])
	count = count + uint32(len)

	len = uint8(bin[count])
	count++
	meta["Comments"] = string(bin[count : count+uint32(len)])
	count = count + uint32(len)

	return count, &meta
}

/*
**
** Pec reading code
**
 */
type H1 struct {
	Label   string
	Ret     byte
	u1      []byte
	TWidth  uint8
	THeight uint8
	u2      []byte
	NoCol   uint8
	ColIdx  []byte
	Pad     []byte
	count   uint32
}

func (h *H1) Parse(bin []byte) {
	count := uint32(0)
	h.Label = string(bin[count : count+19])
	count += 19
	h.Ret = bin[count]
	count++
	h.u1 = bin[count : count+14]
	count += 14
	h.TWidth = bin[count]
	count++
	h.THeight = bin[count]
	count++
	h.u2 = bin[count : count+12]
	count += 12
	h.NoCol = bin[count]
	count++
	h.ColIdx = bin[count : count+uint32(h.NoCol)+1]
	count += uint32(h.NoCol) + 1
	sz := 462 - uint32(h.NoCol)
	h.Pad = bin[count : count+sz]
	count += sz
	h.count = count
}

func (h H1) SizeOf() uint32 {
	return h.count
}

func (h *H1) Dump() {
	fmt.Printf("Header1:\n")
	fmt.Printf("\tLabel: %s\n", h.Label)
	fmt.Printf("\tRet: %d 0x%X\n", h.Ret, h.Ret)
	fmt.Printf("\tu1: %X\n", h.u1)
	fmt.Printf("\tTWidth: %d 0x%X\n", h.TWidth, h.TWidth)
	fmt.Printf("\tTHeight: %d 0x%X\n", h.THeight, h.THeight)
	fmt.Printf("\tu2: %X\n", h.u2)
	fmt.Printf("\tNumCols: %d 0x%X\n", h.NoCol, h.NoCol)
	fmt.Printf("\tColIdx: %X\n", h.ColIdx)
	fmt.Printf("\tPad: %X\n", h.Pad)
	fmt.Printf("\tcount: %d 0x%X\n", h.count, h.count)
}

type H2 struct {
	u1     uint16
	TOffs  uint16
	u2     uint32
	Width  int16
	Height int16
	u3     []byte
	count  uint32
}

func (h *H2) Parse(bin []byte) {
	count := uint32(0)
	h.u1 = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h.TOffs = binary.LittleEndian.Uint16(bin[count : count+2])
	count += 2
	h.u2 = binary.LittleEndian.Uint32(bin[count : count+4])
	count += 4
	h.Width = int16(binary.LittleEndian.Uint16(bin[count : count+2]))
	count += 2
	h.Height = int16(binary.LittleEndian.Uint16(bin[count : count+2]))
	count += 2
	h.u3 = bin[count : count+8]
	count += 8
	h.count = count
}

func (h H2) SizeOf() uint32 {
	return h.count
}

func (h *H2) Dump() {
	fmt.Printf("Header2:\n")
	fmt.Printf("\tu1: %d 0x%x\n", h.u1, h.u1)
	fmt.Printf("\tTOffs: %d 0x%x\n", h.TOffs, h.TOffs)
	fmt.Printf("\tu2: %d 0x%X\n", h.u2, h.u2)
	fmt.Printf("\tWidth: %d 0x%X\n", h.Width, h.Width)
	fmt.Printf("\tHeight: %d 0x%X\n", h.Height, h.Height)
	fmt.Printf("\tu3: %X\n", h.u3)
	fmt.Printf("\tcount: %d 0x%X\n", h.count, h.count)
}

/*
**
** Stitch handling - ie pec body
**
 */

const (
	is_cmd_mask    = 128
	cmd_mask       = 112
	long_mask      = 4094
	long_c_mask    = 15
	color_flag     = 254
	end_flag       = 0xff
	long_test_neg  = 0b1000000000000000
	short_test_neg = 0b1000000
)

type PCommand struct {
	Command1 int
	Command2 int
	Dx       float32
	Dy       float32
	Color    int
}

func pec_decode_cmd(c int) string {
	switch c {
	case shared.Stitch:
		return "Stitch"
	case shared.Jump:
		return "Jump"
	case shared.ColorChg:
		return "ColorChg"
	case shared.End:
		return "End"
	}
	return "unk"
}

func (p PCommand) Dump() {
	fmt.Printf("\tCommand1: %s\n", pec_decode_cmd(p.Command1))
	fmt.Printf("\tCommand2: %s\n", pec_decode_cmd(p.Command2))
	fmt.Printf("%08b\n", p.Dx)
	fmt.Printf("\tDx : %f 0x%X\n", p.Dx, p.Dx)
	fmt.Printf("%08b\n", p.Dy)
	fmt.Printf("\tDy : %f 0x%X\n", p.Dy, p.Dy)
	fmt.Printf("\tColor : %d 0x%X\n", p.Color, p.Color)
}

func next_chunk(bin []byte) []byte {

	count := 0
	var pl []byte
	if bin[count] == end_flag {
		pl = append(pl, bin[count:count+1]...)
		count += 1
	} else if bin[count] == color_flag {
		pl = append(pl, bin[count:count+3]...)
		count += 3
	} else {
		for range 2 {
			if bin[count]&is_cmd_mask > 0 {
				pl = append(pl, bin[count:count+2]...)
				count += 2
			} else {
				pl = append(pl, bin[count:count+1]...)
				count += 1
			}
		}
	}
	return pl

}

func decode_color(c []byte) int {
	return int(c[2])
}

func decode_short(c []byte) (float32, float32) {
	val1 := int16(c[0])
	val2 := int16(c[1])
	if val1 >= 0x40 {
		val1 -= 0x80
	}
	if val2 >= 0x40 {
		val2 -= 0x80
	}
	f1 := float32(val1) * 0.1
	f2 := float32(val2) * 0.1
	return f1, f2
}

func decode_long(c []byte) (int, float32) {
	var cmd int
	flag := c[0] & cmd_mask
	flag = flag >> 4
	switch flag {
	case 1:
		cmd = shared.Jump
	case 2:
		cmd = shared.Trim
	}
	val := int16(c[0] & 0x0F)
	val = val << 8
	val += int16(c[1])
	if val&0x800 > 0 {
		val -= 0x1000
	}
	f := float32(val) * 0.1
	return cmd, f

}

func next_command(bin []byte) (int, *PCommand) {

	var p PCommand

	c := next_chunk(bin)
	count := len(c)

	if c[0] == end_flag {
		p.Command1 = shared.End
	} else {
		switch len(c) {
		case 2:
			// two short coords
			p.Command1 = shared.Stitch
			p.Dx, p.Dy = decode_short(c)
		case 3:
			// short and long or color
			if c[0] == color_flag {
				p.Command1 = shared.ColorChg
				p.Color = decode_color(c)
			} else if c[0]&is_cmd_mask > 0 {
				p.Command1, p.Dx = decode_long(c[0:2])
				p.Dy, _ = decode_short(c[2:])
			} else {
				p.Dx, _ = decode_short(c)
				p.Command2, p.Dy = decode_long(c[1:3])
			}

		case 4:
			p.Command1, p.Dx = decode_long(c[0:2])
			p.Command2, p.Dy = decode_long(c[2:4])
		}
	}
	return count, &p
}

type Payload struct {
	Width   float32
	Height  float32
	Rot     uint16
	Desc    map[string]string
	BG      color.Color
	Path    string
	ColList []shared.ColorSub
	Palette []color.Color
	Head    string
	Cmds    []PCommand
}

func (p *Payload) decode_pes(h Header) {
	// hoop scaling. Apparently coords and sizes are 10 times the real size
	switch h.Ver {
	case "0001":
		if h.H1.Hoop == 0 {
			p.Width = float32(100)
			p.Height = float32(100)
		} else if h.H1.Hoop == 1 {
			p.Width = float32(130)
			p.Height = float32(180)
		}
	case "0020":
		p.Height = float32(h.H2.HoopH)
		p.Width = float32(h.H2.HoopW)
		p.Rot = h.H2.Rot
	case "0030":
		p.Height = float32(h.H3.HoopH)
		p.Width = float32(h.H3.HoopW)
		p.Rot = h.H3.Rot
	case "0040":
		p.Height = float32(h.H4.HoopH)
		p.Width = float32(h.H4.HoopW)
		p.Rot = h.H4.Rot
		p.Desc = *h.H4.Desc
	case "0050":
		p.Height = float32(h.H5.HoopH)
		p.Width = float32(h.H5.HoopW)
		p.Rot = h.H5.Rot
		p.Desc = *h.H5.Desc
		p.Path = h.H5.Impath
		p.ColList = h.H5.Colors
	case "0060":
		p.Height = float32(h.H6.HoopH)
		p.Width = float32(h.H6.HoopW)
		p.Rot = h.H6.Rot
		p.Desc = *h.H6.Desc
		p.Path = h.H6.Impath
		p.ColList = h.H6.Colors
	}
}

func Read_pes(file string) *Payload {
	var pay Payload

	// get the actual file contents
	reader, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	bin, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	reader.Close()

	// get all the information we can from the pes header
	var pes_hdr Header
	pes_hdr.Parse(bin)

	// get what we want from pes header into our payload
	pay.decode_pes(pes_hdr)

	// parse the pec section for a little metadata and the stitches
	PecBin := bin[pes_hdr.P.Offset:]
	var H1 H1
	H1.Parse(PecBin)

	var H2 H2
	count := H1.SizeOf()
	H2.Parse(PecBin[count:])

	pay.Head = H1.Label[2:]
	pay.Palette = convert_colors(pay.ColList, H1.ColIdx)

	var cmds []PCommand

	niggle := 0
	count = 0
	l := H1.SizeOf() + H2.SizeOf()
	StBin := PecBin[l:]
	for {
		niggle++
		b, p := next_command(StBin[count:])
		cmds = append(cmds, *p)
		count += uint32(b)
		if p.Command1 == shared.End {
			break
		}
	}
	pay.Cmds = cmds
	return &pay
}

func convert_colors(c []shared.ColorSub, p []byte) []color.Color {
	var cols []color.Color
	if c != nil {
		for h := range c {
			cols = append(cols, c[h].Color)
		}
	} else {
		palette := Brother_select()
		for i := 0; i < len(p); i++ {
			cols = append(cols, palette[p[i]-1]) // 1 based index not 0
		}
	}
	return cols
}

var (
	PrussianBlueBr = color.RGBA{0x1a, 0x0a, 0x94, 255}
	BlueBr         = color.RGBA{0x0f, 0x75, 0xff, 255}
	TealGreenBr    = color.RGBA{0x00, 0x93, 0x4c, 255}
	CFBlueBr       = color.RGBA{0xba, 0xbd, 0xfe, 255}
	RedBr          = color.RGBA{0xec, 0x00, 0x00, 255}
	RedBrownBr     = color.RGBA{0xe4, 0x99, 0x5a, 255}
	MagentaBr      = color.RGBA{0xcc, 0x48, 0xab, 255}
	LightLilacBr   = color.RGBA{0xfd, 0xc4, 0xfa, 255}
	LilacBr        = color.RGBA{0xdd, 0x84, 0xab, 255}
	MintGreenBr    = color.RGBA{0x6b, 0xd3, 0x8a, 255}
	DeepGoldBr     = color.RGBA{0xe4, 0xa9, 0x45, 255}
	OrangeBr       = color.RGBA{0xff, 0xbd, 0x42, 255}
	YellowBr       = color.RGBA{0xff, 0xe6, 0x00, 255}
	LimeGreenBr    = color.RGBA{0x6c, 0xd9, 0x00, 255}
	BrassBr        = color.RGBA{0xc1, 0xa9, 0x41, 255}
	SilverBr       = color.RGBA{0xb5, 0xad, 0x97, 255}
	RussetBrownBr  = color.RGBA{0xba, 0x9c, 0x5f, 255}
	CreamBrownBr   = color.RGBA{0xfa, 0xf5, 0x9e, 255}
	PewterBr       = color.RGBA{0x80, 0x80, 0x80, 255}
	BlackBr        = color.RGBA{0x0, 0x0, 0x0, 255}
	UltraMarineBr  = color.RGBA{0x00, 0x1c, 0xdf, 255}
	RoyaPurpleBr   = color.RGBA{0xdf, 0x00, 0xb8, 255}
	DarkGrayBr     = color.RGBA{0x62, 0x62, 0x62, 255}
	DarkBrownBr    = color.RGBA{0x69, 0x26, 0x0d, 255}
	DeepRoseBr     = color.RGBA{0xff, 0x00, 0x60, 255}
	LightBrownBr   = color.RGBA{0xbf, 0x82, 0x00, 255}
	SalmonPinkBr   = color.RGBA{0xf3, 0x91, 0x78, 255}
	VermilionBr    = color.RGBA{0xff, 0x68, 0x05, 255}
	WhiteBr        = color.RGBA{0xf0, 0xf0, 0xf0, 255}
	VioletBr       = color.RGBA{0xc8, 0x32, 0xcd, 255}
	SeaCrestBr     = color.RGBA{0xb0, 0xbf, 0x9b, 255}
	SkyBlueBr      = color.RGBA{0x65, 0xbf, 0xeb, 255}

	PumpkinBr     = color.RGBA{0xff, 0xba, 0x04, 255}
	CreamYellowBr = color.RGBA{0xff, 0xf0, 0x6c, 255}
	KhakiBr       = color.RGBA{0xfe, 0xca, 0x15, 255}
	ClayBrownBr   = color.RGBA{0xf3, 0x81, 0x01, 255}
	LeafGreenBr   = color.RGBA{0x37, 0xa9, 0x23, 255}
	PeacockBlueBr = color.RGBA{0x23, 0x46, 0x5f, 255}
)

func Brother_set() *map[string]color.Color {
	return &map[string]color.Color{
		"Br_PrussianBlue":   PrussianBlueBr,
		"Br_Blue":           BlueBr,
		"Br_TealGreen":      TealGreenBr,
		"Br_CornFlowerBlue": CFBlueBr,
		"Br_Red":            RedBr,
		"Br_RedBrown":       RedBrownBr,
		"Br_Magenta":        MagentaBr,
		"Br_LightLilac":     LightLilacBr,
		"Br_Lilac":          LilacBr,
		"Br_MintGreen":      MintGreenBr,
		"Br_DeepGold":       DeepGoldBr,
		"Br_Orange":         OrangeBr,
		"Br_Yellow":         YellowBr,
		"Br_LimeGreen":      LimeGreenBr,
		"Br_Brass":          BrassBr,
		"Br_Silver":         SilverBr,
		"Br_RussetBrown":    RussetBrownBr,
		"Br_CreamBrown":     CreamBrownBr,
		"Br_Pewter":         PewterBr,
		"Br_Black":          BlackBr,
		"Br_UltraMarine":    UltraMarineBr,
		"Br_RoyalPurple":    RoyaPurpleBr,
		"Br_DarkGray":       DarkGrayBr,
		"Br_DarkBrown":      DarkBrownBr,
		"Br_DeepRose":       DeepRoseBr,
		"Br_LightBrown":     LightBrownBr,
		"Br_SalmonPink":     SalmonPinkBr,
		"Br_Vermilion":      VermilionBr,
		"Br_White":          WhiteBr,
		"Br_Violet":         VioletBr,
		"Br_SeaCrest":       SeaCrestBr,
		"Br_SkyBlue":        SkyBlueBr,
		"Br_Pumpkin":        PumpkinBr,
		"Br_CreamYellow":    CreamYellowBr,
		"Br_Khaki":          KhakiBr,
		"Br_ClayBrown":      ClayBrownBr,
		"Br_LeafGreen":      LeafGreenBr,
		"Br_PeacockBlue":    PeacockBlueBr,
	}
} //brother_set()

// order is important
func Brother_select() []color.Color {
	return []color.Color{
		PrussianBlueBr,
		BlueBr,
		TealGreenBr,
		CFBlueBr,
		RedBr,
		RedBrownBr,
		MagentaBr,
		LightLilacBr,
		LilacBr,
		MintGreenBr,
		DeepGoldBr,
		OrangeBr,
		YellowBr,
		LimeGreenBr,
		BrassBr,
		SilverBr,
		RussetBrownBr,
		CreamBrownBr,
		PewterBr,
		BlackBr,
		UltraMarineBr,
		RoyaPurpleBr,
		DarkGrayBr,
		DarkBrownBr,
		DeepRoseBr,
		LightBrownBr,
		SalmonPinkBr,
		VermilionBr,
		WhiteBr,
		VioletBr,
		SeaCrestBr,
		SkyBlueBr,
		PumpkinBr,
		CreamYellowBr,
		KhakiBr,
		ClayBrownBr,
		LeafGreenBr,
		PeacockBlueBr,
	}
} //Brother_Select()
