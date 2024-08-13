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
		bg: unison.MustColorDecode("#ff0090"),
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

	paint := unison.NewPaint()
	svg.DrawInRect(canvas, unison.NewRect(0, 0, 100, 100), nil, paint)
}
