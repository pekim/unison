package unison

import (
	"runtime"

	"github.com/richardwilkes/unison/internal/skia"
)

// TextBlob represents runs of text for a font, that may be drawn on a Canvas.
type TextBlob struct {
	skia.TextBlob
}

func newTextBlob(textBlob skia.TextBlob) *TextBlob {
	if textBlob == nil {
		return nil
	}
	tb := &TextBlob{TextBlob: textBlob}
	runtime.SetFinalizer(tb, func(obj *TextBlob) {
		tb.Unref()
	})

	return tb
}

func (tb *TextBlob) Unref() {
	runtime.SetFinalizer(tb, nil)
	ReleaseOnUIThread(func() {
		if tb.TextBlob != nil {
			skia.TextBlobUnref(tb.TextBlob)
			tb.TextBlob = nil
		}
	})
}
