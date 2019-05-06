package quantum

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

type RenderParams struct {
	// The number of qubit lines to draw.
	NumBits int

	// The spacing between each line, and bordering the
	// top and bottom of the image.
	LineSpace int

	// The radius of dots when wires are connected.
	DotSize int

	// The size for text gates.
	FontSize int
}

func (r *RenderParams) QubitY(bit int) int {
	return bit*r.LineSpace + r.LineSpace/2
}

func (r *RenderParams) Height() int {
	return r.LineSpace * r.NumBits
}

func DefaultRenderParams(numBits int) *RenderParams {
	return &RenderParams{
		NumBits:   numBits,
		LineSpace: 50,
		DotSize:   3,
		FontSize:  18,
	}
}

// Renderer is any gate that can be rendered to an image.
// All gates should render horizontal lines for the
// qubits.
type Renderer interface {
	Render(params *RenderParams) (*image.RGBA, error)
}

// RenderText renders a gate as a text box around a
// certain qubit.
func RenderText(params *RenderParams, bit int, text string) (*image.RGBA, error) {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(font, &truetype.Options{Size: float64(params.FontSize)})

	ctx := gg.NewContext(1, 1)
	ctx.SetFontFace(face)
	width, height := ctx.MeasureString(text)
	if width < height {
		width = height
	}

	dest := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(width))+params.FontSize*2,
		params.Height()))

	RenderQubits(params, dest)

	// Erase the spot behind the rectangle.
	x1 := params.FontSize / 2
	y1 := params.QubitY(bit) - int(height/2) - params.FontSize/2
	w := int(width) + params.FontSize
	h := int(height) + params.FontSize
	for x := x1; x < x1+w; x++ {
		for y := y1; y < y1+h; y++ {
			dest.SetRGBA(x, y, color.RGBA{})
		}
	}

	ctx = gg.NewContextForRGBA(dest)
	ctx.SetFontFace(face)
	ctx.SetRGB(0, 0, 0)
	ctx.DrawStringAnchored(text, float64(dest.Bounds().Dx())/2, float64(params.QubitY(bit))-2,
		0.5, 0.5)

	ctx.MoveTo(float64(x1), float64(y1))
	ctx.LineTo(float64(x1+w), float64(y1))
	ctx.LineTo(float64(x1+w), float64(y1+h))
	ctx.LineTo(float64(x1), float64(y1+h))
	ctx.ClosePath()
	ctx.Stroke()

	return dest, nil
}

// RenderQubits draws the qubits in an image.
func RenderQubits(params *RenderParams, img *image.RGBA) {
	ctx := gg.NewContextForRGBA(img)
	for i := 0; i < params.NumBits; i++ {
		y := float64(params.QubitY(i))
		ctx.MoveTo(0, y)
		ctx.LineTo(float64(img.Bounds().Dx()), y)
	}
	ctx.Stroke()
}

func (h *HGate) Render(params *RenderParams) (*image.RGBA, error) {
	return RenderText(params, h.Bit, "H")
}

func (t *TGate) Render(params *RenderParams) (*image.RGBA, error) {
	if t.Conjugate {
		return RenderText(params, t.Bit, "T*")
	} else {
		return RenderText(params, t.Bit, "T")
	}
}

func (x *XGate) Render(params *RenderParams) (*image.RGBA, error) {
	return RenderText(params, x.Bit, "X")
}

func (y *YGate) Render(params *RenderParams) (*image.RGBA, error) {
	return RenderText(params, y.Bit, "Y")
}

func (z *ZGate) Render(params *RenderParams) (*image.RGBA, error) {
	return RenderText(params, z.Bit, "Z")
}

func (c Circuit) Render(params *RenderParams) (*image.RGBA, error) {
	var images []*image.RGBA
	var totalWidth int
	for _, subGate := range c {
		renderer, ok := subGate.(Renderer)
		if !ok {
			return nil, fmt.Errorf("cannot render %T", subGate)
		}
		img, err := renderer.Render(params)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
		totalWidth += img.Bounds().Dx()
	}
	dest := image.NewRGBA(image.Rect(0, 0, totalWidth, params.Height()))
	x := 0
	for _, img := range images {
		draw.Draw(dest, image.Rect(x, 0, x+img.Bounds().Dx(), img.Bounds().Dy()),
			img, image.Point{0, 0}, draw.Over)
		x += img.Bounds().Dx()
	}
	return dest, nil
}
