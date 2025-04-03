package pdf

import (
	"fmt"
	"image/color"

	"github.com/jung-kurt/gofpdf"
)

// NamedColor stores a thread color and its name
type NamedColor struct {
	Name  string
	Color color.Color
}

// ImgType stores the image file path and its displayed dimensions
type ImgType struct {
	DispWidth  float64
	DispHeight float64
	File       string
}

// Palette is a test slice of colors
var palette []NamedColor

// fillP is a test helper function to populate some colors
func fillP() {
	palette = append(palette, NamedColor{Name: "VermillionBr", Color: color.RGBA{0xFF, 0x68, 0x05, 255}})
	palette = append(palette, NamedColor{Name: "RedBrownBr", Color: color.RGBA{0xEC, 0, 0, 255}})
	palette = append(palette, NamedColor{Name: "KhakiBr", Color: color.RGBA{0xFE, 0xCA, 0x15, 255}})
	palette = append(palette, NamedColor{Name: "SalmonPinkBr", Color: color.RGBA{0xF3, 0x91, 0x78, 255}})
	palette = append(palette, NamedColor{Name: "EmeraldGreenBr", Color: color.RGBA{0x22, 0x89, 0x27, 255}})
	palette = append(palette, NamedColor{Name: "DarkGrayBr", Color: color.RGBA{0x62, 0x62, 0x62, 255}})
	palette = append(palette, NamedColor{Name: "SkyBlueBr", Color: color.RGBA{0x65, 0xBF, 0xEB, 255}})
	palette = append(palette, NamedColor{Name: "GrayBr", Color: color.RGBA{0xA6, 0xA6, 0x95, 255}})
	palette = append(palette, NamedColor{Name: "DarkBrownBr", Color: color.RGBA{0x69, 0x26, 0x0d, 255}})

	palette = append(palette, NamedColor{Name: "VermillionBr", Color: color.RGBA{0xFF, 0x68, 0x05, 255}})
	palette = append(palette, NamedColor{Name: "RedBrownBr", Color: color.RGBA{0xEC, 0, 0, 255}})
	palette = append(palette, NamedColor{Name: "KhakiBr", Color: color.RGBA{0xFE, 0xCA, 0x15, 255}})
	palette = append(palette, NamedColor{Name: "SalmonPinkBr", Color: color.RGBA{0xF3, 0x91, 0x78, 255}})
	palette = append(palette, NamedColor{Name: "EmeraldGreenBr", Color: color.RGBA{0x22, 0x89, 0x27, 255}})
	palette = append(palette, NamedColor{Name: "DarkGrayBr", Color: color.RGBA{0x62, 0x62, 0x62, 255}})
	palette = append(palette, NamedColor{Name: "SkyBlueBr", Color: color.RGBA{0x65, 0xBF, 0xEB, 255}})
	palette = append(palette, NamedColor{Name: "GrayBr", Color: color.RGBA{0xA6, 0xA6, 0x95, 255}})
	palette = append(palette, NamedColor{Name: "DarkBrownBr", Color: color.RGBA{0x69, 0x26, 0x0d, 255}})

}

// EmbPdf collects everything required to make a pdf from an embroidery image
type EmbPdf struct {
	Img     ImgType
	Palette []NamedColor
	Title   string
	Desc    map[string]string
	Meta    string
	pdf     *gofpdf.Fpdf
	used    float64
}

// NewEmbPdf is a constructor for a mainly empty EmbPdf
func NewEmbPdf(file string) *EmbPdf {
	p := gofpdf.New("P", "mm", PageSize, "")
	p.AddPage()
	p.SetMargins(Margin, Margin, Margin)
	i := ImgType{DispWidth: 0, DispHeight: 0, File: file}
	return &EmbPdf{
		Img:     i,
		Palette: nil,
		Title:   file,
		Desc:    nil,
		Meta:    "",
		pdf:     p,
		used:    0.0,
	}
}

// Layout constants that control how the pdf is layed out
const (
	PageSize      = "A4"    // We are Australian so metric
	A4Width       = 190     // 200mm - 2 margins
	A4Height      = 297     // A4 in Portrait is 297mm high
	Font          = "Arial" // Font family
	TitleSize     = 26      // size in points
	TitleWeight   = "B"     // B is bold, I italic, U underscore, S strikethrough or "" for normal
	ImgConstraint = 90      // max dimension that the image is displayed
	Margin        = 10      // margin size in mm

	ImgBlock = 30 // y value for the image and side blocks (mm)

	SideBlock    = 110 // x value for the side block (mm)
	SideHeight   = 5   // height of the colour rectangle (mm)
	SideWidth    = 10  // width of the colour rectangle (mm)
	SidePad      = 5   // horizontal padding between elements of side block (mm)
	SideFontSize = 10  // Size of the text in the side block (points)
	Tween        = 1   // vertical padding between elements of side block (mm)
)

// SetTitle draws the title into the PDF
func (p *EmbPdf) SetTitle() {
	p.pdf.SetFont(Font, TitleWeight, TitleSize)
	p.pdf.CellFormat(A4Width, 0, p.Title, "", 2, "C", false, 0, "")
	p.used = 3 * p.pdf.PointConvert(TitleSize)
}

// SetImage draws the embroidery image into the pdf
func (p *EmbPdf) SetImage() {
	options := gofpdf.ImageOptions{
		ImageType:             "",
		ReadDpi:               true,
		AllowNegativePosition: false,
	}
	info := p.pdf.RegisterImageOptions(p.Img.File, options)
	wd, h := info.Extent()
	if wd > h { // this is wrong
		p.Img.DispWidth = ImgConstraint
		p.Img.DispHeight = ImgConstraint * h / wd
	} else {
		p.Img.DispWidth = ImgConstraint * wd / h
		p.Img.DispHeight = ImgConstraint
	}

	p.pdf.ImageOptions(
		p.Img.File,
		Margin,
		p.used,
		p.Img.DispWidth,
		p.Img.DispHeight,
		false,
		options,
		0,
		"",
	)
}

// SetColor writes out a color splatch and string at the requested location
func (p *EmbPdf) SetColor(idx int, offset float64, y float64) {

	// create the color rect
	c := p.Palette[idx].Color
	p.pdf.SetFillColor(int(c.(color.RGBA).R), int(c.(color.RGBA).G), int(c.(color.RGBA).B))
	p.pdf.Rect(offset, float64(y), SideWidth, SideHeight, "DF")

	// convert points to mm and then position in relation to colour rect
	texty := SideHeight - ((SideHeight - p.pdf.PointConvert(SideFontSize)) / 2)

	// create hex color code
	str := convertHex(c)
	x := float64(offset + SideWidth + SidePad)
	y = y + texty
	p.pdf.SetFont("Courier", "", SideFontSize)
	p.pdf.Text(float64(x), float64(y), str)

	// create the color name
	x = x + p.pdf.GetStringWidth(str) + SidePad
	p.pdf.SetFont(Font, "", SideFontSize)
	p.pdf.Text(x, y, p.Palette[idx].Name)
}

// SetColors draws the palette of colors into the pdf in the side box and then under the image
func (p *EmbPdf) SetColors() {

	// work out how many colors can be put in the sideblock
	maxCols := int(p.Img.DispHeight / (SideHeight + Tween))
	here := 0.0
	// write them to pdf
	for i := 0; i < maxCols; i++ {
		y := float64(p.used + (float64(i) * (Tween + SideHeight)))
		p.SetColor(i, SideBlock, y)
		here = y
	}

	yoffset := p.used + Tween + p.Img.DispHeight
	here = yoffset
	for i := maxCols; i < len(p.Palette); i++ {
		y := yoffset + (float64(i-maxCols) * (Tween + SideHeight))
		p.SetColor(i, Margin, y)
		here = y
	}
	p.used = here + Tween + SideHeight + (p.pdf.PointConvert(SideFontSize) / 2)

}

// Layout creates draws the pdf contents and saves it
func (p *EmbPdf) Layout() {

	p.SetTitle()
	p.SetImage()
	p.SetColors()
	p.pdf.Text(Margin, p.used, "we are here")
	err := p.pdf.OutputFileAndClose("hello.pdf")
	if err == nil {
		fmt.Println("PDF generated successfully")
	}

}

// convertHex is a helper function that translates an image/color into a web style hex string
func convertHex(c color.Color) string {
	R, G, B, _ := c.RGBA()
	r := R >> 8
	g := G >> 8
	b := B >> 8
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func main() {
	p := NewEmbPdf("D1124.jpg")
	fillP()
	//var err error
	p.Palette = palette
	p.Layout()
}
