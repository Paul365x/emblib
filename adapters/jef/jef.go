package jef

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/emblib/adapters/shared"
)

var fh *os.File = os.Stdout

/*
**
** Jef header parsing code
**
 */
type Jef_header struct {
	Offset  uint32 // file offset to stitches
	unk1    uint32
	Date    string // date created
	Ver     string // software version
	unk2    uint8
	ClrCnt  uint32   // number of colors
	PtsLen  uint32   // length of points - halved
	Hoop    uint32   // hoop used
	Extends []uint32 // offsets from centre of hoop (4 values)
	Pad1    []uint32 // padding inset from hoop 110x110 edge (4 values)
	Pad2    []uint32 // padding inset from hoop 50x50 edge (4 values)
	Pad3    []uint32 // padding inset from hoop 140x200 edge (4 values)
	Pad4    []uint32 // padding inset from custom hoop edge (4 values)
	ClrChg  []uint32 // color changes
	count   uint32   // bytes in struct
}

// Jef_header.Parse reads in the header of a jef file into the struct
func (s *Jef_header) Parse(bin []byte) {
	s.Offset = binary.LittleEndian.Uint32(bin[0:4])
	s.unk1 = binary.LittleEndian.Uint32(bin[4:8])
	s.Date = string(bin[8:22])
	s.Ver = string(bin[22:23])
	s.unk2 = uint8(bin[23])
	s.ClrCnt = binary.LittleEndian.Uint32(bin[24:28])
	s.PtsLen = binary.LittleEndian.Uint32(bin[28:32])
	s.Hoop = binary.LittleEndian.Uint32(bin[32:36])
	s.Extends = make([]uint32, 4)
	var count uint32 = 36
	for i := range 4 {
		s.Extends[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}

	s.Pad1 = make([]uint32, 4)
	for i := range 4 {
		s.Pad1[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}
	s.Pad2 = make([]uint32, 4)
	for i := range 4 {
		s.Pad2[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}

	s.Pad3 = make([]uint32, 4)
	for i := range 4 {
		s.Pad3[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}

	s.Pad4 = make([]uint32, 4)
	for i := range 4 {
		s.Pad4[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}

	var clr uint32
	for {
		clr = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
		if clr == clr_end_mask {
			break
		}
		s.ClrChg = append(s.ClrChg, clr)
	}
	for {
		clr = binary.LittleEndian.Uint32(bin[count : count+4])
		if clr != clr_end_mask {
			break
		}
		count += 4
	}
	s.count = count
}

// Preamble.SizeOf returns the size in bytes - offset into the file of the byte after the preamble. Always 12 bytes
func (p Jef_header) SizeOf() uint32 {
	return p.count
}

// Dump writes out this Struct
func (p Jef_header) Dump() {
	fmt.Fprintf(fh, "Header:\n")
	fmt.Fprintf(fh, "\tOffset: %d 0x%X\n", p.Offset, p.Offset)
	fmt.Fprintf(fh, "\tunk1: %d 0x%X\n", p.unk1, p.unk1)
	fmt.Fprintf(fh, "\tDate: %s\n", p.Date)
	fmt.Fprintf(fh, "\tVer: %s\n", p.Ver)
	fmt.Fprintf(fh, "\tunk2: %d 0x%X\n", p.unk2, p.unk2)
	fmt.Fprintf(fh, "\tColorCnt: %d 0x%X\n", p.ClrCnt, p.ClrCnt)
	fmt.Fprintf(fh, "\tPtsLen: %d 0x%X\n", p.PtsLen, p.PtsLen)
	fmt.Fprintf(fh, "\tHoop: %d 0x%X\n", p.Hoop, p.Hoop)
	fmt.Fprintf(fh, "\tExtends: %X\n", p.Extends)
	fmt.Fprintf(fh, "\tPad1: %X\n", p.Pad1)
	fmt.Fprintf(fh, "\tPad2: %X\n", p.Pad2)
	fmt.Fprintf(fh, "\tPad3: %X\n", p.Pad3)
	fmt.Fprintf(fh, "\tPad4: %X\n", p.Pad4)
	fmt.Fprintf(fh, "\tColorChg: %X\n", p.ClrChg)

	fmt.Fprintf(fh, "\tcount: %d 0x%X\n\n", p.count, p.count)
}

func read_cmds(bin []byte) []shared.PCommand {
	var cmd shared.PCommand
	var cmds []shared.PCommand
	count := uint32(0)
	var loc uint32
	for {
		loc = count
		b0 := int8(bin[count])
		count++
		b1 := int8(bin[count])
		// end of file is marked by a series of 0xff or -1 at this point
		if b0 == -1 && b1 == -1 {
			break
		}
		count++
		if b0 == -127 {
			switch b1 {
			case 10:
				// end
				break
			case 01:
				//color chg
				cmd.Command1 = shared.ColorChg
				cmd.Dx = float32(bin[count])
				count++
				cmd.Dy = float32(bin[count])
				count++
			case 02:
				//jmp
				cmd.Dx = float32(bin[count])
				count++
				cmd.Dy = float32(bin[count])
				count++
				if cmd.Dx == 0 && cmd.Dy == 0 {
					cmd.Command1 = shared.Trim
				} else {
					cmd.Command1 = shared.Jump
				}
			}
		} else {
			cmd.Command1 = shared.Stitch
			cmd.Dx = float32(int(b0))
			cmd.Dy = float32(int(b1) * -1)
		}
		if cmd.Dx == 0xff && cmd.Dy == 0xff {
			break
		}
		cmds = append(cmds, cmd)
		fmt.Printf("%d\t%d\t%s\t%f %f\n", len(cmds), loc, jef_decode_cmd(cmd.Command1), cmd.Dx, cmd.Dy)
	}
	return cmds
}

func decode_jef(h Jef_header) shared.Payload {
	var p shared.Payload
	// two ways to get the width and height - using the extends or the hoop size
	switch h.Hoop {
	case 0:
		// 110 x 110 mm
		p.Width = 110
		p.Height = 110
	case 1:
		// 50 x 50 mm
		p.Width = 50
		p.Height = 50
	case 2:
		// 140 x 200 mm
		p.Width = 140
		p.Height = 200
	case 3:
		// 126 x 110 mm
		p.Width = 126
		p.Height = 110
	case 4:
		// 200 x 200 mm
		p.Width = 200
		p.Height = 200
	default:
		p.Width = float32(h.Extends[0] + h.Extends[2])
		p.Height = float32(h.Extends[1] + h.Extends[3])
	}
	return p
}

func Read_jef(file string) *shared.Payload {
	var pay shared.Payload

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

	var jef Jef_header
	jef.Parse(bin)
	jef.Dump()
	c := jef.SizeOf()
	pay = decode_jef(jef)

	pay.Cmds = read_cmds(bin[c:])
	return &pay
}

const (
	clr_end_mask = 0xd
	is_cmd_mask  = 0x80
)

/*
func main() {
	pay := read_jef("2024SewNow.JEF")
	cmds := pay.Cmds

	myApp := app.New()
	w := myApp.NewWindow("Lines")
	c_prev := PCommand{
		Command1: shared.Stitch,
		Command2: shared.Stitch,
		Dx:       0,
		Dy:       0,
		Color:    0,
	}

	prev := fyne.NewPos(0, 0)
	img := container.NewWithoutLayout()

	// all stitches are offset from centre
	origin_x := pay.Width / 2.0
	origin_y := pay.Height / 2.0

	fmt.Println(origin_x, origin_y, pay.Width, pay.Height)

	//cols := pay.Palette
	//col_idx := 0
	for p := range cmds {
		if p == 0 && cmds[p].Dx == 0 && cmds[p].Dy == 0 {
			continue
		}
		if p == 1 {
			c_prev = cmds[p]
			prev = fyne.NewPos(cmds[p].Dx+origin_x, cmds[p].Dy+origin_y)
			continue
		}
		x := cmds[p].Dx + c_prev.Dx
		y := cmds[p].Dy + c_prev.Dy
		switch cmds[p].Command1 {
		case shared.ColorChg:
		//	col_idx = cmds[p].Color - 1 // 1 indexed in file
		case shared.Trim: // jump without line/thread
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y
		default:
			if p > 5 {
				//		break
			}
			c_prev = cmds[p]
			c_prev.Dx = x
			c_prev.Dy = y
			fmt.Println("coord: ", x, y)
			l := canvas.NewLine(color.Black)
			l.Position1 = prev
			prev = fyne.NewPos(x+origin_x, y+origin_y)
			l.Position2 = prev
			l.StrokeWidth = 2
			img.Add(l)
		}
	}

	w.SetContent(img)

	w.Resize(fyne.NewSize(1000, 1000))
	w.ShowAndRun()

}
*/
/*
**
** Stitch handling
**
 */

func jef_decode_cmd(c int) string {
	switch c {
	case shared.Trim:
		return "Trim"
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

/*


func (p PCommand) Dump() {
	fmt.Fprintf(fh, "\tCommand1: %s\n", decode_cmd(p.Command1))
	fmt.Fprintf(fh, "\tCommand2: %s\n", decode_cmd(p.Command2))
	fmt.Printf("%08b\n", p.Dx)
	fmt.Fprintf(fh, "\tDx : %f 0x%X\n", p.Dx, p.Dx)
	fmt.Printf("%08b\n", p.Dy)
	fmt.Fprintf(fh, "\tDy : %f 0x%X\n", p.Dy, p.Dy)
	fmt.Fprintf(fh, "\tColor : %d 0x%X\n", p.Color, p.Color)
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
		cmd = Jump
	case 2:
		cmd = Trim
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
		p.Command1 = End
	} else {
		switch len(c) {
		case 2:
			// two short coords
			p.Command1 = Stitch
			p.Dx, p.Dy = decode_short(c)
		case 3:
			// short and long or color
			if c[0] == color_flag {
				p.Command1 = ColorChg
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
	ColList []ColorSub
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

func read_pes(file string) *Payload {
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
		if p.Command1 == End {
			break
		}
	}
	pay.Cmds = cmds
	return &pay
}

func convert_colors(c []ColorSub, p []byte) []color.Color {
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
*/
