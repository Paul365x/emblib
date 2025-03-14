/*
** Jef adapter
** routines to read and understand Janome's Jef file format
** Creates a sequence of commands with metadata that can be run on a render engine
 */

package jef

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"os"

	"github.com/emblib/adapters/shared"
)

var fh *os.File = os.Stdout
var expand float32 = 0.5 // Expand: Size of resulting image is dependent on a specific machine
// this is a fudge factor to ensure the image fits

// Masks that do bitwise operations to help parse stitches
const (
	clr_end_mask = 0xd
	is_cmd_mask  = 0x80
)

/*
**
** Jef header parsing code
**
 */

// Jef_header stores the header portion of a jef file
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

	s.Pad1 = make([]uint32, 4) // edge for hoop
	for i := range 4 {
		s.Pad1[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}
	s.Pad2 = make([]uint32, 4) // edge for hoop
	for i := range 4 {
		s.Pad2[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}

	s.Pad3 = make([]uint32, 4) // edge for hoop
	for i := range 4 {
		s.Pad3[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}

	s.Pad4 = make([]uint32, 4) // edge for hoop
	for i := range 4 {
		s.Pad4[i] = binary.LittleEndian.Uint32(bin[count : count+4])
		count += 4
	}

	// read colours (u32) until we get to 0x0d
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
} // Parse

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

// read_cmds parses stitches to a list of render engine commands
func read_cmds(bin []byte, cols []uint32, f func() int) []shared.PCommand {

	// set the initial color
	var cmd = shared.PCommand{
		Command1: shared.ColorChg,
		Command2: 0,
		Dx:       2.0,
		Dy:       2.0,
		Color:    int(cols[f()] + 1),
	}

	var cmds []shared.PCommand // some file formats have a null first command. Add to make same
	cmds = append(cmds, cmd)

	count := uint32(0)
	var loc uint32
FORLOOP:
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
		if b0 == -128 { // all commands have -128 in b0
			switch b1 {
			case 0x10:
				// end
				break FORLOOP
			case 01:
				//color chg
				cmd.Command1 = shared.ColorChg
				cmd.Dx = float32(bin[count]) * expand
				count++
				cmd.Dy = float32(bin[count]) * expand
				count++
				cmd.Color = int(cols[f()])
			case 02:
				//jmp and trim
				cmd.Dx = float32(int8(bin[count])) * expand
				count++
				cmd.Dy = float32(int8(bin[count])*-1) * expand
				count++
				if cmd.Dx == 0 && cmd.Dy == 0 {
					cmd.Command1 = shared.Trim
				} else {
					cmd.Command1 = shared.Jump
				}
			}
		} else {
			// stitch
			cmd.Command1 = shared.Stitch
			cmd.Dx = float32(int(b0)) * expand
			cmd.Dy = float32(int(b1)*-1) * expand
		}
		if cmd.Dx == 0xff && cmd.Dy == 0xff {
			break
		}
		cmds = append(cmds, cmd)
		//	if jef_decode_cmd(cmd.Command1) != "Stitch" {
		fmt.Printf("%d %d %d\t%d\t%s\t%f %f %d\n", b0, b1, len(cmds), loc, jef_decode_cmd(cmd.Command1), cmd.Dx, cmd.Dy, cmd.Color)
		//	}
	}
	return cmds
} // read_cmds()

// decode_jef converts jef header information to useable - currently only width and height
func decode_jef(h Jef_header) shared.Payload {
	var p shared.Payload
	// two ways to get the width and height - using the extends or the hoop size
	// prefer extends
	if h.Extends[0] != 0 && h.Extends[2] != 0 && h.Extends[1] != 0 && h.Extends[3] != 0 {
		p.Width = float32(h.Extends[0] + h.Extends[2])
		p.Height = float32(h.Extends[1] + h.Extends[3])
	} else {
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

		}
	}
	return p
} // decode_jef

// inc factory to produce an incrementing closure
func inc() func() int {
	index := -1
	return func() int {
		index = index + 1
		return index
	}
}

// Read_jef reads a jef file and returns the payload ie what we are interested in
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
	pay.Palette = Janome_select()
	f := inc()
	pay.Cmds = read_cmds(bin[c:], jef.ClrChg, f)
	return &pay
} // Read_jef

/*
**
** Stitch handling
**
 */

// jef_decode_cmd converts a numeric command - render engine format - to a string
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
} // jef_decode_cmd

// defines for the colors of the standard Janome thread palette
var (
	UnknownJf      = color.RGBA{0x0, 0x0, 0x0, 255}
	BlackJf        = color.RGBA{0x0, 0x0, 0x0, 255}
	WhiteJf        = color.RGBA{0xFF, 0xFF, 0xFF, 255}
	Sunflower1Jf   = color.RGBA{0xFF, 0xFF, 0x17, 255}
	Hazel1Jf       = color.RGBA{0xFA, 0xA0, 0x60, 255}
	OliveGreenJf   = color.RGBA{0x5C, 0x76, 0x49, 255}
	GreenJf        = color.RGBA{0x40, 0xC0, 0x30, 255}
	SkyJf          = color.RGBA{0x65, 0xC2, 0xC8, 255}
	PurpleJf       = color.RGBA{0xAC, 0x80, 0xBE, 255}
	PinkJf         = color.RGBA{0xF5, 0xBC, 0xCB, 255}
	RedJf          = color.RGBA{0xFF, 0x0, 0x0, 255}
	BrownJf        = color.RGBA{0xC0, 0x80, 0x0, 255}
	BlueJf         = color.RGBA{0x0, 0x0, 0xF0, 255}
	GoldJf         = color.RGBA{0xE4, 0xC3, 0x5D, 255}
	DarkBrownJf    = color.RGBA{0xA5, 0x2A, 0x2A, 255}
	PaleVioletJf   = color.RGBA{0xD5, 0xB0, 0xD4, 255}
	PaleYellowJf   = color.RGBA{0xFC, 0xF2, 0x94, 255}
	PalePinkJf     = color.RGBA{0xF0, 0xD0, 0xC0, 255}
	PeachJf        = color.RGBA{0xFF, 0xC0, 0x0, 255}
	Beige1Jf       = color.RGBA{0xC9, 0xA4, 0x80, 255}
	WineRedJf      = color.RGBA{0x9B, 0x3D, 0x4B, 255}
	PaleSkyJf      = color.RGBA{0xA0, 0xB8, 0xCC, 255}
	YellowGreenJf  = color.RGBA{0x7F, 0xC2, 0x1C, 255}
	SilverGreyJf   = color.RGBA{0xB9, 0xB9, 0xB9, 255}
	GreyJf         = color.RGBA{0xA0, 0xA0, 0xA0, 255}
	PaleAquaJf     = color.RGBA{0x98, 0xD6, 0xBD, 255}
	BabyBlueJf     = color.RGBA{0xB8, 0xF0, 0xF0, 255}
	PowderBlueJf   = color.RGBA{0x36, 0x8B, 0xA0, 255}
	BrightBlueJf   = color.RGBA{0x4F, 0x83, 0xAB, 255}
	SlateBlueJf    = color.RGBA{0x38, 0x6A, 0x91, 255}
	NaveBlueJf     = color.RGBA{0x0, 0x20, 0x6B, 255}
	SalmonPinkJf   = color.RGBA{0xE5, 0xC5, 0xCA, 255}
	CoralJf        = color.RGBA{0xF9, 0x67, 0x6B, 255}
	BurntOrangeJf  = color.RGBA{0xE3, 0x31, 0x1F, 255}
	CinnamonJf     = color.RGBA{0xE2, 0xA1, 0x88, 255}
	UmberJf        = color.RGBA{0xB5, 0x94, 0x74, 255}
	BlondeJf       = color.RGBA{0xE4, 0xCF, 0x99, 255}
	Sunflower2Jf   = color.RGBA{0xE1, 0xCB, 0x0, 255}
	OrchidPinkJf   = color.RGBA{0xE1, 0xAD, 0xD4, 255}
	PeonyPurpleJf  = color.RGBA{0xC3, 0x0, 0x7E, 255}
	BurgundyJf     = color.RGBA{0x80, 0x0, 0x4B, 255}
	RoyalPurple1Jf = color.RGBA{0xA0, 0x60, 0xB0, 255}
	CardinalRedJf  = color.RGBA{0xC0, 0x40, 0x20, 255}
	OpalGreenJf    = color.RGBA{0xCA, 0xE0, 0xC0, 255}
	MossGreenJf    = color.RGBA{0x89, 0x98, 0x56, 255}
	MeadowGreenJf  = color.RGBA{0x0, 0xAA, 0x0, 255}
	DarkGreenJf    = color.RGBA{0x21, 0x8A, 0x21, 255}
	AquamarineJf   = color.RGBA{0x5D, 0xAE, 0x94, 255}
	EmeraldGreenJf = color.RGBA{0x4C, 0xBF, 0x8F, 255}
	PeacockGreenJf = color.RGBA{0x0, 0x77, 0x72, 255}
	DarkGreyJf     = color.RGBA{0x70, 0x70, 0x70, 255}
	IvoryWhiteJf   = color.RGBA{0xF2, 0xFF, 0xFF, 255}
	Hazel2Jf       = color.RGBA{0xB1, 0x58, 0x18, 255}
	Toast1Jf       = color.RGBA{0xCB, 0x8A, 0x7, 255}
	SalmonJf       = color.RGBA{0xF7, 0x92, 0x7B, 255}
	CocoaBrownJf   = color.RGBA{0x98, 0x69, 0x2D, 255}
	SiennaJf       = color.RGBA{0xA2, 0x71, 0x48, 255}
	Sepia1Jf       = color.RGBA{0x7B, 0x55, 0x4A, 255}
	DarkSepiaJf    = color.RGBA{0x4F, 0x39, 0x46, 255}
	VioletBlueJf   = color.RGBA{0x52, 0x3A, 0x97, 255}
	BlueInkJf      = color.RGBA{0x0, 0x0, 0xA0, 255}
	SolarBlueJf    = color.RGBA{0x0, 0x96, 0xDE, 255}
	GreenDustJf    = color.RGBA{0xB2, 0xDD, 0x53, 255}
	CrimsonJf      = color.RGBA{0xFA, 0x8F, 0xBB, 255}
	FloralPinkJf   = color.RGBA{0xDE, 0x64, 0x9E, 255}
	WineJf         = color.RGBA{0xB5, 0x50, 0x66, 255}
	OliveDrabJf    = color.RGBA{0x5E, 0x57, 0x47, 255}
	MeadowJf       = color.RGBA{0x4C, 0x88, 0x1F, 255}
	CanaryYellowJf = color.RGBA{0xE4, 0xDC, 0x79, 255}
	Toast2Jf       = color.RGBA{0xCB, 0x8A, 0x1A, 255}
	Beige2Jf       = color.RGBA{0xC6, 0xAA, 0x42, 255}
	HoneyDewJf     = color.RGBA{0xEC, 0xB0, 0x2C, 255}
	TangerineJf    = color.RGBA{0xF8, 0x80, 0x40, 255}
	OceanBlueJf    = color.RGBA{0xFF, 0xE5, 0x5, 255}
	Sepia2Jf       = color.RGBA{0xFA, 0x7A, 0x7A, 255}
	RoyalPurple2Jf = color.RGBA{0x6B, 0xE0, 0x0, 255}
	YellowOcherJf  = color.RGBA{0x38, 0x6C, 0xAE, 255}
	BeigeGreyJf    = color.RGBA{0xD0, 0xBA, 0xB0, 255}
	BambooJf       = color.RGBA{0xE3, 0xBE, 0x81, 255}
) // var

// Janome_set returns a map of the Janome thread palette indexed by name to color
func Janome_set() *map[string]color.Color {
	return &map[string]color.Color{
		"Jf_Unknown":      UnknownJf,
		"Jf_Black":        BlackJf,
		"Jf_White":        WhiteJf,
		"Jf_Sunflower1":   Sunflower1Jf,
		"Jf_Hazel1":       Hazel1Jf,
		"Jf_OliveGreen":   OliveGreenJf,
		"Jf_Green":        GreenJf,
		"Jf_Sky":          SkyJf,
		"Jf_Purple":       PurpleJf,
		"Jf_Pink":         PinkJf,
		"Jf_Red":          RedJf,
		"Jf_Brown":        BrownJf,
		"Jf_Blue":         BlueJf,
		"Jf_Gold":         GoldJf,
		"Jf_DarkBrown":    DarkBrownJf,
		"Jf_PaleViolet":   PaleVioletJf,
		"Jf_PaleYellow":   PaleYellowJf,
		"Jf_PalePink":     PalePinkJf,
		"Jf_Peach":        PeachJf,
		"Jf_Beige1":       Beige1Jf,
		"Jf_WineRed":      WineRedJf,
		"Jf_PaleSky":      PaleSkyJf,
		"Jf_YellowGreen":  YellowGreenJf,
		"Jf_SilverGrey":   SilverGreyJf,
		"Jf_Grey":         GreyJf,
		"Jf_PaleAqua":     PaleAquaJf,
		"Jf_BabyBlue":     BabyBlueJf,
		"Jf_PowderBlue":   PowderBlueJf,
		"Jf_BrightBlue":   BrightBlueJf,
		"Jf_SlateBlue":    SlateBlueJf,
		"Jf_NaveBlue":     NaveBlueJf,
		"Jf_SalmonPink":   SalmonPinkJf,
		"Jf_Coral":        CoralJf,
		"Jf_BurntOrange":  BurntOrangeJf,
		"Jf_Cinnamon":     CinnamonJf,
		"Jf_Umber":        UmberJf,
		"Jf_Blonde":       BlondeJf,
		"Jf_Sunflower2":   Sunflower2Jf,
		"Jf_OrchidPink":   OrchidPinkJf,
		"Jf_PeonyPurple":  PeonyPurpleJf,
		"Jf_Burgundy":     BurgundyJf,
		"Jf_RoyalPurple1": RoyalPurple1Jf,
		"Jf_CardinalRed":  CardinalRedJf,
		"Jf_OpalGreen":    OpalGreenJf,
		"Jf_MossGreen":    MossGreenJf,
		"Jf_MeadowGreen":  MeadowGreenJf,
		"Jf_DarkGreen":    DarkGreenJf,
		"Jf_Aquamarine":   AquamarineJf,
		"Jf_EmeraldGreen": EmeraldGreenJf,
		"Jf_PeacockGreen": PeacockGreenJf,
		"Jf_DarkGrey":     DarkGreyJf,
		"Jf_IvoryWhite":   IvoryWhiteJf,
		"Jf_Hazel2":       Hazel2Jf,
		"Jf_Toast1":       Toast1Jf,
		"Jf_Salmon":       SalmonJf,
		"Jf_CocoaBrown":   CocoaBrownJf,
		"Jf_Sienna":       SiennaJf,
		"Jf_Sepia1":       Sepia1Jf,
		"Jf_DarkSepia":    DarkSepiaJf,
		"Jf_VioletBlue":   VioletBlueJf,
		"Jf_BlueInk":      BlueInkJf,
		"Jf_SolarBlue":    SolarBlueJf,
		"Jf_GreenDust":    GreenDustJf,
		"Jf_Crimson":      CrimsonJf,
		"Jf_FloralPink":   FloralPinkJf,
		"Jf_Wine":         WineJf,
		"Jf_OliveDrab":    OliveDrabJf,
		"Jf_Meadow":       MeadowJf,
		"Jf_CanaryYellow": CanaryYellowJf,
		"Jf_Toast2":       Toast2Jf,
		"Jf_Beige2":       Beige2Jf,
		"Jf_HoneyDew":     HoneyDewJf,
		"Jf_Tangerine":    TangerineJf,
		"Jf_OceanBlue":    OceanBlueJf,
		"Jf_Sepia2":       Sepia2Jf,
		"Jf_RoyalPurple2": RoyalPurple2Jf,
		"Jf_YellowOcher":  YellowOcherJf,
		"Jf_BeigeGrey":    BeigeGreyJf,
		"Jf_Bamboo":       BambooJf,
	}
} // Janome_set

// Janome_select returns an array in index order of the Janome thread palette
func Janome_select() []color.Color {
	return []color.Color{
		UnknownJf,
		BlackJf,
		WhiteJf,
		Sunflower1Jf,
		Hazel1Jf,
		OliveGreenJf,
		GreenJf,
		SkyJf,
		PurpleJf,
		PinkJf,
		RedJf,
		BrownJf,
		BlueJf,
		GoldJf,
		DarkBrownJf,
		PaleVioletJf,
		PaleYellowJf,
		PalePinkJf,
		PeachJf,
		Beige1Jf,
		WineRedJf,
		PaleSkyJf,
		YellowGreenJf,
		SilverGreyJf,
		GreyJf,
		PaleAquaJf,
		BabyBlueJf,
		PowderBlueJf,
		BrightBlueJf,
		SlateBlueJf,
		NaveBlueJf,
		SalmonPinkJf,
		CoralJf,
		BurntOrangeJf,
		CinnamonJf,
		UmberJf,
		BlondeJf,
		Sunflower2Jf,
		OrchidPinkJf,
		PeonyPurpleJf,
		BurgundyJf,
		RoyalPurple1Jf,
		CardinalRedJf,
		OpalGreenJf,
		MossGreenJf,
		MeadowGreenJf,
		DarkGreenJf,
		AquamarineJf,
		EmeraldGreenJf,
		PeacockGreenJf,
		DarkGreyJf,
		IvoryWhiteJf,
		Hazel2Jf,
		Toast1Jf,
		SalmonJf,
		CocoaBrownJf,
		SiennaJf,
		Sepia1Jf,
		DarkSepiaJf,
		VioletBlueJf,
		BlueInkJf,
		SolarBlueJf,
		GreenDustJf,
		CrimsonJf,
		FloralPinkJf,
		WineJf,
		OliveDrabJf,
		MeadowJf,
		CanaryYellowJf,
		Toast2Jf,
		Beige2Jf,
		HoneyDewJf,
		TangerineJf,
		OceanBlueJf,
		Sepia2Jf,
		RoyalPurple2Jf,
		YellowOcherJf,
		BeigeGreyJf,
		BambooJf,
	}
} // Janome_select
