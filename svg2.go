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
	"io"
	"strconv"
	"strings"

	svg "github.com/lafriks/go-svg"
	"github.com/richardwilkes/toolbox/fatal"
)

var _ Drawable = &DrawableSVG2{}

// DrawableSVG2 makes an SVG conform to the Drawable interface.
type DrawableSVG2 struct {
	SVG  *SVG2
	Size Size
}

// SVG2 holds an SVG.
type SVG2 struct {
	svg   *svg.Svg
	paths []*Path2
	size  Size
}

// MustSVG2FromContentString creates a new SVG and panics if an error would be generated. The content should contain
// valid SVG file data. Note that this only reads a very small subset of an SVG currently. Specifically, the "viewBox"
// attribute and any "d" attributes from enclosed SVG "path" elements.
func MustSVG2FromContentString(content string) *SVG2 {
	s, err := NewSVG2FromContentString(content)
	fatal.IfErr(err)
	return s
}

// NewSVG2FromContentString creates a new SVG. The content should contain valid SVG file data. Note that this only reads
// a very small subset of an SVG currently. Specifically, the "viewBox" attribute and any "d" attributes from enclosed
// SVG "path" elements.
func NewSVG2FromContentString(content string) (*SVG2, error) {
	return NewSVG2FromReader(strings.NewReader(content))
}

// MustSVG2FromReader creates a new SVG and panics if an error would be generated. The reader should contain valid SVG
// file data. Note that this only reads a very small subset of an SVG currently. Specifically, the "viewBox" attribute
// and any "d" attributes from enclosed SVG "path" elements.
func MustSVG2FromReader(r io.Reader) *SVG2 {
	s, err := NewSVG2FromReader(r)
	fatal.IfErr(err)
	return s
}

// NewSVG2FromReader creates a new SVG. The reader should contain valid SVG file data. Note that this only reads a very
// small subset of an SVG currently. Specifically, the "viewBox" attribute and any "d" attributes from enclosed SVG
// "path" elements.
func NewSVG2FromReader(r io.Reader) (*SVG2, error) {
	svg, err := svg.Parse(r, svg.StrictErrorMode)
	if err != nil {
		return nil, err
	}

	var size Size
	if svg.Width != "" {
		w, err := strconv.Atoi(svg.Width)
		if err == nil {
			size.Width = float32(w)
		}
	}
	if svg.Height != "" {
		h, err := strconv.Atoi(svg.Height)
		if err == nil {
			size.Height = float32(h)
		}
	}
	if svg.ViewBox.W > 0 {
		size.Width = float32(svg.ViewBox.W)
	}
	if svg.ViewBox.W > 0 {
		size.Height = float32(svg.ViewBox.H)
	}

	s := &SVG2{
		svg:   svg,
		paths: make([]*Path2, len(svg.SvgPaths)),
		size:  size,
	}
	for i, path := range svg.SvgPaths {
		p, err := newPath2(path)
		if err != nil {
			return nil, err
		}
		s.paths[i] = p
	}

	return s, nil
}

// AspectRatio returns the SVG's width to height ratio.
func (s *SVG2) AspectRatio() float32 {
	return s.size.Width / s.size.Height
}

// LogicalSize implements the Drawable interface.
func (s *DrawableSVG2) LogicalSize() Size {
	return s.Size
}

// DrawInRect implements the Drawable interface.
func (s *DrawableSVG2) DrawInRect(canvas *Canvas, rect Rect, _ *SamplingOptions, _ *Paint) {
	canvas.Save()
	defer canvas.Restore()

	canvas.Translate(rect.X, rect.Y)
	canvas.Scale(rect.Width/s.SVG.size.Width, rect.Height/s.SVG.size.Height)

	for _, path := range s.SVG.paths {
		path.draw(canvas)
	}
}
