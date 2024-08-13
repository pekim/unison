package main

import (
	"github.com/richardwilkes/unison"
)

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
		bg: unison.MustColorDecode("#ff80e0"),
	}
	ss.DrawCallback = ss.defaultDraw
	ss.SetSizer(ss.defaultSizes)
	return ss
}

func (ss *svgs) defaultSizes(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
	size := unison.Size{
		Width:  800,
		Height: 800,
	}
	return size, size, size
}

func (ss *svgs) defaultDraw(canvas *unison.Canvas, _ unison.Rect) {
	r := ss.ContentRect(false)
	unison.DrawRectBase(canvas, r, ss.bg, unison.Transparent)

	svg := unison.DrawableSVG{
		SVG: unison.CircledChevronRightSVG,
		// Size: unison.NewSize(100, 100),
	}

	svg2 := unison.DrawableSVG2{
		SVG: unison.MustSVG2FromContentString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
    <path d="M256 8c137 0 248 111 248 248S393 504 256 504 8 393 8 256 119 8 256 8zm113.9 231L234.4 103.5c-9.4-9.4-24.6-9.4-33.9 0l-17 17c-9.4 9.4-9.4 24.6 0 33.9L285.1 256 183.5 357.6c-9.4 9.4-9.4 24.6 0 33.9l17 17c9.4 9.4 24.6 9.4 33.9 0L369.9 273c9.4-9.4 9.4-24.6 0-34z"/>
</svg>`),
		// Size: unison.NewSize(100, 100),
	}

	paint := unison.NewPaint()
	svg.DrawInRect(canvas, unison.NewRect(0, 0, 100, 100), nil, paint)
	svg2.DrawInRect(canvas, unison.NewRect(150, 0, 100, 100), nil, paint)
}
