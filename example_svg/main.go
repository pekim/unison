package main

import (
	_ "embed"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

//go:embed gopher.svg
var gopherSvg string

//go:embed tiger.svg
var tigerSvg string

//go:embed black-king.svg
var blackKingSvg string

//go:embed linear-gradient.svg
var linearGradientSvg string

//go:embed linear-gradient-stroke.svg
var linearGradientStrokeSvg string

//go:embed radial-gradient.svg
var radialGradientSvg string

//go:embed two-point-conical-gradient.svg
var twoPointConicalGradientSvg string

var circledChevronRight = unison.DrawableSVG{
	SVG: unison.CircledChevronRightSVG,
}

var blackKing = unison.DrawableSVG{
	SVG: unison.MustSVGFromContentString(blackKingSvg),
}

var gopher = unison.DrawableSVG{
	SVG: unison.MustSVGFromContentString(gopherSvg),
}

var tiger = unison.DrawableSVG{
	SVG: unison.MustSVGFromContentString(tigerSvg),
}

var linearGradient = unison.DrawableSVG{
	SVG: unison.MustSVGFromContentString(linearGradientSvg),
}

var linearGradientStroke = unison.DrawableSVG{
	SVG: unison.MustSVGFromContentString(linearGradientStrokeSvg),
}

var radialGradient = unison.DrawableSVG{
	SVG: unison.MustSVGFromContentString(radialGradientSvg),
}

var twoPointConicalGradient = unison.DrawableSVG{
	SVG: unison.MustSVGFromContentString(twoPointConicalGradientSvg),
}

func main() {
	unison.AttachConsole()
	unison.Start(unison.StartupFinishedCallback(func() {
		wnd, err := unison.NewWindow("example")
		if err != nil {
			panic(err)
		}

		unison.DefaultMenuFactory().BarForWindow(wnd, func(m unison.Menu) {
			unison.InsertStdMenus(m, func(_ unison.MenuItem) {}, nil, nil)
		})

		content := wnd.Content()
		content.SetBorder(unison.NewEmptyBorder(unison.NewUniformInsets(20)))
		content.SetLayout(&unison.FlexLayout{
			Columns:  1,
			HSpacing: unison.StdHSpacing,
			VSpacing: 10,
			// HAlign:   align.Middle,
			// VAlign:   align.Middle,
		})

		content.AddChild(newSVGs())

		wnd.Pack()
		wnd.ToFront()
	}))

}

type svgs struct {
	Drawable unison.Drawable
	unison.Panel
	bg unison.Color
}

func newSVGs() *svgs {
	ss := &svgs{
		// bg: unison.MustColorDecode("#ffc0f0"),
		bg: unison.MustColorDecode("#ffffff"),
	}
	ss.DrawCallback = ss.defaultDraw
	ss.SetSizer(ss.defaultSizes)
	return ss
}

func (ss *svgs) defaultSizes(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
	size := unison.Size{
		Width:  1000,
		Height: 1000,
	}
	return size, size, size
}

func (ss *svgs) defaultDraw(canvas *unison.Canvas, _ unison.Rect) {
	r := ss.ContentRect(false)
	unison.DrawRectBase(canvas, r, ss.bg, unison.Transparent)

	{
		circledChevronRight.DrawInRect(canvas, unison.NewRect(0, 0, 100, 100), nil, nil)
	}

	{
		paint := unison.NewPaint()
		paint.SetColor(unison.RGB(0xff, 0x00, 0x00))
		paint.SetStyle(paintstyle.Fill)
		circledChevronRight.DrawInRect(canvas, unison.NewRect(150, 0, 100, 100), nil, paint)
	}

	{
		width := float32(120)
		height := width / blackKing.SVG.AspectRatio()

		bg := unison.NewPaint()
		bg.SetColor(unison.RGB(0x00, 0xff, 0x00))
		rect := unison.NewRect(300, 0, width, height)
		canvas.DrawRect(rect, bg)
		blackKing.DrawInRect(canvas, rect, nil, nil)
	}

	{
		width := float32(200)
		height := width / gopher.SVG.AspectRatio()
		gopher.DrawInRect(canvas, unison.NewRect(50, 150, width, height), nil, nil)
	}

	{
		width := float32(400)
		height := width / tiger.SVG.AspectRatio()
		tiger.DrawInRect(canvas, unison.NewRect(300, 150, width, height), nil, nil)
	}

	{
		width := float32(150)
		height := width / linearGradient.SVG.AspectRatio()
		linearGradient.DrawInRect(canvas, unison.NewRect(50, 450, width, height), nil, nil)
	}

	{
		width := float32(400)
		height := width / linearGradientStroke.SVG.AspectRatio()
		linearGradientStroke.DrawInRect(canvas, unison.NewRect(250, 600, width, height), nil, nil)
	}

	{
		width := float32(180)
		height := width / radialGradient.SVG.AspectRatio()
		radialGradient.DrawInRect(canvas, unison.NewRect(50, 800, width, height), nil, nil)
	}

	{
		width := float32(180)
		height := width / twoPointConicalGradient.SVG.AspectRatio()
		twoPointConicalGradient.DrawInRect(canvas, unison.NewRect(350, 800, width, height), nil, nil)
	}
}
