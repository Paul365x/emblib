/*
** Shared contains code required by all adapters - essentially the guts of the api
**
 */

package shared

import (
	"image/color"
	"os"
)

var fh *os.File = os.Stdout

// PCommand is a struct to hold a command - jump, trim, stitch etc
type PCommand struct {
	Command1 int
	Command2 int
	Dx       float32
	Dy       float32
	Color    int
}

// Payload captures metadata from file headers and also the stitch commands
type Payload struct {
	Width        float32
	Height       float32
	Rot          uint16
	Desc         map[string]string
	BG           color.Color
	Path         string
	Palette      []color.Color
	Palette_type bool
	Head         string
	Cmds         []PCommand
}

// ColorSub stores the color structure
type ColorSub struct {
	CodeLen  uint8
	Code     []byte
	Color    color.Color
	U1       uint8
	ColType  uint32
	DescLen  uint8
	Desc     string
	BrandLen uint8
	Brand    string
	ChartLen uint8
	Chart    string
	Count    uint32
}

// command constants
const (
	Stitch = iota + 1
	Jump
	Trim
	ColorChg
	End
)
