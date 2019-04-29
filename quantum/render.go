package quantum

import (
	"image"

	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	RenderLineSpace = 30
	RenderGateSize  = 26
	RenderDotSize   = 3
)

func RenderGate(numBits int, g Gate) *image.RGBA {
	switch g := g.(type) {
	case *HGate:
		return renderTextGate(numBits, g.Bit, "H")
	case *TGate:
		text := "T"
		if g.Conjugate {
			text = "T*"
		}
		return renderTextGate(numBits, g.Bit, text)
	case *XGate:
		return renderTextGate(numBits, g.Bit, "X")
	case *YGate:
		return renderTextGate(numBits, g.Bit, "Y")
	case *ZGate:
		return renderTextGate(numBits, g.Bit, "Z")
	case Circuit:
		return renderCircuit(numBits, g)
	}
	panic("do not know how to render gate")
}

func gateWidth(numBits int, g Gate) int {
	return RenderGateSize + 10
}

func renderTextGate(numBits, bit int, text string) *image.RGBA {
	// TODO: measure text here for width
	dest := image.NewRGBA(image.Rect(0, 0, RenderGateSize+10, RenderLineSpace*(numBits+1)))
	// ctx := draw2dimg.NewGraphicContext(dest)

	// TODO: this.

	return dest
}

func renderCircuit(numBits int, c Circuit) *image.RGBA {
	var images []*image.RGBA
	var totalWidth int
	for _, subGate := range c {
		img := RenderGate(numBits, subGate)
		images = append(images, img)
		totalWidth += img.Bounds().Dx()
	}
	dest := image.NewRGBA(image.Rect(0, 0, totalWidth, RenderLineSpace*(numBits+1)))
	ctx := draw2dimg.NewGraphicContext(dest)
	x := 0
	for _, img := range images {
		ctx.Save()
		ctx.Translate(float64(x), 0)
		ctx.DrawImage(img)
		ctx.Restore()
		x += img.Bounds().Dx()
	}
	return dest
}
