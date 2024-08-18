// Copyright (c) 2021-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package unison

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/lafriks/go-svg"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/fatal"
	"github.com/richardwilkes/unison/enums/strokecap"
	"github.com/richardwilkes/unison/enums/strokejoin"
	"github.com/richardwilkes/unison/internal/skia"
)

var _ Drawable = &DrawableSVG{}

// Pre-defined SVG images used by Unison.
var (
	//go:embed resources/images/broken_image.svg
	brokenImageSVG string
	BrokenImageSVG = MustSVGFromContentString(brokenImageSVG)

	//go:embed resources/images/circled_chevron_right.svg
	circledChevronRightSVG string
	CircledChevronRightSVG = MustSVGFromContentString(circledChevronRightSVG)

	//go:embed resources/images/circled_exclamation.svg
	circledExclamationSVG string
	CircledExclamationSVG = MustSVGFromContentString(circledExclamationSVG)

	//go:embed resources/images/circled_question.svg
	circledQuestionSVG string
	CircledQuestionSVG = MustSVGFromContentString(circledQuestionSVG)

	//go:embed resources/images/checkmark.svg
	checkmarkSVG string
	CheckmarkSVG = MustSVGFromContentString(checkmarkSVG)

	//go:embed resources/images/chevron_right.svg
	chevronRightSVG string
	ChevronRightSVG = MustSVGFromContentString(chevronRightSVG)

	//go:embed resources/images/circled_x.svg
	circledXSVG string
	CircledXSVG = MustSVGFromContentString(circledXSVG)

	//go:embed resources/images/dash.svg
	dashSVG string
	DashSVG = MustSVGFromContentString(dashSVG)

	//go:embed resources/images/document.svg
	documentSVG string
	DocumentSVG = MustSVGFromContentString(documentSVG)

	//go:embed resources/images/sort_ascending.svg
	sortAscendingSVG string
	SortAscendingSVG = MustSVGFromContentString(sortAscendingSVG)

	//go:embed resources/images/sort_descending.svg
	sortDescendingSVG string
	SortDescendingSVG = MustSVGFromContentString(sortDescendingSVG)

	//go:embed resources/images/triangle_exclamation.svg
	triangleExclamationSVG string
	TriangleExclamationSVG = MustSVGFromContentString(triangleExclamationSVG)

	//go:embed resources/images/window_maximize.svg
	windowMaximizeSVG string
	WindowMaximizeSVG = MustSVGFromContentString(windowMaximizeSVG)

	//go:embed resources/images/window_restore.svg
	windowRestoreSVG string
	WindowRestoreSVG = MustSVGFromContentString(windowRestoreSVG)
)

var strokeCaps = map[svg.CapMode]strokecap.Enum{
	svg.NilCap:    strokecap.Butt,
	svg.ButtCap:   strokecap.Butt,
	svg.SquareCap: strokecap.Square,
	svg.RoundCap:  strokecap.Round,
}

var strokeJoins = map[svg.JoinMode]strokejoin.Enum{
	svg.Round: strokejoin.Round,
	svg.Bevel: strokejoin.Bevel,
	svg.Miter: strokejoin.Miter,
}

// DrawableSVG makes an SVG conform to the Drawable interface.
type DrawableSVG struct {
	SVG  *SVG
	Size Size
}

// SVG holds an SVG.
type SVG struct {
	dom  skia.SVGDOM
	size Size
}

// SVGOption is an option that may be passed to SVG construction functions.
type SVGOption func(s *svgOptions) error

type svgOptions struct {
	parseErrorMode    svg.ErrorMode
	ignoreUnsupported bool
}

// SVGOptionIgnoreParseErrors is an option that will ignore some errors when parsing an SVG.
// If the XML is not well formed an error will still be generated.
func SVGOptionIgnoreParseErrors() SVGOption {
	return func(opts *svgOptions) error {
		opts.parseErrorMode = svg.IgnoreErrorMode
		return nil
	}
}

// SVGOptionWarnParseErrors is an option that will issue warnings to the log for some errors when parsing an SVG.
// If the XML is not well formed an error will still be generated.
func SVGOptionWarnParseErrors() SVGOption {
	return func(opts *svgOptions) error {
		opts.parseErrorMode = svg.WarnErrorMode
		return nil
	}
}

// SVGOptionIgnoreUnsupported is an option that will ignore unsupported SVG features that might be encountered
// in an SVG.
//
// If this option is not present, then unsupported features will result in an error when constructing the SVG.
func SVGOptionIgnoreUnsupported() SVGOption {
	return func(opts *svgOptions) error {
		opts.ignoreUnsupported = true
		return nil
	}
}

// MustSVG creates a new SVG the given svg path string (the contents of a single "d" attribute from an SVG "path"
// element) and panics if an error would be generated. The 'size' should be gotten from the original SVG's 'viewBox'
// parameter.
//
// Note: It is probably better to use one of the other Must... methods that take the full SVG content.
func MustSVG(size Size, svg string) *SVG {
	s, err := NewSVG(size, svg)
	fatal.IfErr(err)
	return s
}

// NewSVG creates a new SVG the given svg path string (the contents of a single "d" attribute from an SVG "path"
// element). The 'size' should be gotten from the original SVG's 'viewBox' parameter.
//
// Note: It is probably better to use one of the other New... methods that take the full SVG content.
func NewSVG(size Size, path string) (*SVG, error) {
	return NewSVGFromContentString(fmt.Sprintf(`
		<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %f %f" version="1.1">
			<path d="%s"/>
		</svg>
	`,
		size.Width, size.Height, path))
}

// MustSVGFromContentString creates a new SVG and panics if an error would be generated. The content should contain
// valid SVG file data. Note that this only reads a very small subset of an SVG currently. Specifically, the "viewBox"
// attribute and any "d" attributes from enclosed SVG "path" elements.
func MustSVGFromContentString(content string, options ...SVGOption) *SVG {
	s, err := NewSVGFromContentString(content, options...)
	fatal.IfErr(err)
	return s
}

// NewSVGFromContentString creates a new SVG. The content should contain valid SVG file data. Note that this only reads
// a very small subset of an SVG currently. Specifically, the "viewBox" attribute and any "d" attributes from enclosed
// SVG "path" elements.
func NewSVGFromContentString(content string, options ...SVGOption) (*SVG, error) {
	return NewSVGFromReader(strings.NewReader(content), options...)
}

// MustSVGFromReader creates a new SVG and panics if an error would be generated. The reader should contain valid SVG
// file data. Note that this only reads a very small subset of an SVG currently. Specifically, the "viewBox" attribute
// and any "d" attributes from enclosed SVG "path" elements.
func MustSVGFromReader(r io.Reader, options ...SVGOption) *SVG {
	s, err := NewSVGFromReader(r, options...)
	fatal.IfErr(err)
	return s
}

// NewSVGFromReader creates a new SVG. The reader should contain valid SVG file data. Note that this only reads a very
// small subset of an SVG currently. Specifically, the "viewBox" attribute and any "d" attributes from enclosed SVG
// "path" elements.
func NewSVGFromReader(r io.Reader, options ...SVGOption) (*SVG, error) {
	s := &SVG{}
	opts := svgOptions{parseErrorMode: svg.StrictErrorMode}
	for _, option := range options {
		if err := option(&opts); err != nil {
			return nil, err
		}
	}

	var r2 bytes.Buffer
	r1 := io.TeeReader(r, &r2)
	sData, err := svg.Parse(r1, opts.parseErrorMode)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	s.size = NewSize(float32(sData.ViewBox.W), float32(sData.ViewBox.H))
	data, err := io.ReadAll(&r2)
	if err != nil {
		return nil, errs.NewWithCause("failed to read SVG", err)
	}
	stream := skia.MemoryStreamMakeCopy(data)
	if stream == nil {
		return nil, errs.New("failed to create skia stream for SVG source")
	}
	s.dom = skia.SVGDOMMakeFromStream(stream)
	if s.dom == nil {
		return nil, errs.New("failed to create skia svg dom")
	}
	skia.MemoryStreamDelete(stream)

	runtime.SetFinalizer(s, func(obj *SVG) {
		ReleaseOnUIThread(func() {
			skia.SVGDOMDelete(obj.dom)
		})
	})
	return s, nil
}

// Size returns the original size.
func (s *SVG) Size() Size {
	return s.size
}

// OffsetToCenterWithinScaledSize returns the scaled offset values to use to keep the image centered within the given
// size.
func (s *SVG) OffsetToCenterWithinScaledSize(size Size) Point {
	scale := min(size.Width/s.size.Width, size.Height/s.size.Height)
	return Point{X: (size.Width - s.size.Width*scale) / 2, Y: (size.Height - s.size.Height*scale) / 2}
}

// AspectRatio returns the SVG's width to height ratio.
func (s *SVG) AspectRatio() float32 {
	return s.size.Width / s.size.Height
}

// DrawInRect draws this SVG resized to fit in the given rectangle. If paint is not nil, the SVG paths will be drawn
// with the provided paint, ignoring any fill or stroke attributes within the source SVG. Be sure to set the Paint's
// style (fill or stroke) as desired.
func (s *SVG) DrawInRect(canvas *Canvas, rect Rect, _ *SamplingOptions, _ *Paint) {
	canvas.Save()
	defer canvas.Restore()
	offset := s.OffsetToCenterWithinScaledSize(rect.Size)
	canvas.Translate(rect.X+offset.X, rect.Y+offset.Y)
	canvas.Scale(rect.Width/s.size.Width, rect.Height/s.size.Height)
	skia.SVGDOMSetContainerSize(s.dom, s.size.Width, s.size.Height)
	skia.SVGDOMRender(s.dom, canvas.canvas)
}

// DrawInRectPreservingAspectRatio draws this SVG resized to fit in the given rectangle, preserving the aspect ratio.
func (s *SVG) DrawInRectPreservingAspectRatio(canvas *Canvas, rect Rect, opts *SamplingOptions, paint *Paint) {
	ratio := s.AspectRatio()
	w := rect.Width
	h := w / ratio
	if h > rect.Height {
		h = rect.Height
		w = h * ratio
	}
	s.DrawInRect(canvas, NewRect(rect.X+(rect.Width-w)/2, rect.Y+(rect.Height-h)/2, w, h), opts, paint)
}

// LogicalSize implements the Drawable interface.
func (s *DrawableSVG) LogicalSize() Size {
	return s.Size
}

// DrawInRect implements the Drawable interface.
//
// If paint is not nil, the SVG paths will be drawn with the provided paint, ignoring any fill or stroke attributes
// within the source SVG. Be sure to set the Paint's style (fill or stroke) as desired.
func (s *DrawableSVG) DrawInRect(canvas *Canvas, rect Rect, opts *SamplingOptions, paint *Paint) {
	s.SVG.DrawInRectPreservingAspectRatio(canvas, rect, opts, paint)
}
