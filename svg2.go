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
	_ "embed"
	"io"
	"strconv"
	"strings"

	svg "github.com/lafriks/go-svg"
	"github.com/richardwilkes/toolbox/fatal"
)

var _ Drawable = &DrawableSVG2{}

// DrawableSVG makes an SVG conform to the Drawable interface.
type DrawableSVG2 struct {
	SVG  *SVG2
	Size Size
}

// SVG holds an SVG.
type SVG2 struct {
	svg   *svg.Svg
	paths []*Path2
	size  Size
}

// // MustSVG2 creates a new SVG the given svg path string (the contents of a single "d" attribute from an SVG "path"
// // element) and panics if an error would be generated. The 'size' should be gotten from the original SVG's 'viewBox'
// // parameter.
// // func MustSVG2(size Size, svg string) *SVG2 {
// // 	s, err := NewSVG2(size, svg)
// // 	fatal.IfErr(err)
// // 	return s
// // }

// // // NewSVG2 creates a new SVG the given svg path string (the contents of a single "d" attribute from an SVG "path"
// // // element). The 'size' should be gotten from the original SVG's 'viewBox' parameter.
// // func NewSVG2(size Size, svg string) (*SVG2, error) {
// // 	unscaledPath, err := NewPathFromSVGString(svg)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return &SVG2{
// // 		size:          size,
// // 		unscaledPath:  unscaledPath,
// // 		scaledPathMap: make(map[Size]*Path),
// // 	}, nil
// // }

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

// // MustSVG2FromReader creates a new SVG and panics if an error would be generated. The reader should contain valid SVG
// // file data. Note that this only reads a very small subset of an SVG currently. Specifically, the "viewBox" attribute
// // and any "d" attributes from enclosed SVG "path" elements.
// func MustSVG2FromReader(r io.Reader) *SVG2 {
// 	s, err := NewSVG2FromReader(r)
// 	fatal.IfErr(err)
// 	return s
// }

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

	// // var svgXML struct {
	// // 	ViewBox string `xml:"viewBox,attr"`
	// // 	Paths   []struct {
	// // 		Path string `xml:"d,attr"`
	// // 	} `xml:"path"`
	// // }
	// // if err := xml.NewDecoder(r).Decode(&svgXML); err != nil {
	// // 	return nil, errs.NewWithCause("unable to decode SVG", err)
	// // }
	// // svg := &SVG2{scaledPathMap: make(map[Size]*Path)}
	// // var width, height string
	// // if parts := strings.Split(svgXML.ViewBox, " "); len(parts) == 4 {
	// // 	width = parts[2]
	// // 	height = parts[3]
	// // }
	// // v, err := strconv.ParseFloat(width, 64)
	// // if err != nil || v < 1 || v > 4096 {
	// // 	return nil, errs.NewWithCause("unable to determine SVG width", err)
	// // }
	// // svg.size.Width = float32(v)
	// // v, err = strconv.ParseFloat(height, 64)
	// // if err != nil || v < 1 || v > 4096 {
	// // 	return nil, errs.NewWithCause("unable to determine SVG height", err)
	// // }
	// // svg.size.Height = float32(v)
	// // for i, svgPath := range svgXML.Paths {
	// // 	var p *Path
	// // 	if p, err = NewPathFromSVGString(svgPath.Path); err != nil {
	// // 		return nil, errs.NewWithCausef(err, "unable to decode SVG: path element #%d", i)
	// // 	}
	// // 	if svg.unscaledPath == nil {
	// // 		svg.unscaledPath = p
	// // 	} else {
	// // 		svg.unscaledPath.Path(p, false)
	// // 	}
	// // }
	// return svg, nil
}

// func (s *SVG2) addPaths() {
// 	s.unscaledPaths = make([]*Path, len(s.svg.SvgPaths))
// 	for i, path := range s.svg.SvgPaths {
// 		s.unscaledPaths[i] = newPath(
// 			s.createPath(path),
// 			s.createFillPaint(path),
// 			s.createStrokePaint(path),
// 		)
// 	}
// }

// func (s *SVG2) createPath(path svg.SvgPath) skia.Path {
// 	skPath := skia.PathNew()

// 	for _, op := range path.Path {
// 		switch op := op.(type) {
// 		case svg.OpMoveTo:
// 			skia.PathMoveTo(skPath, float32(op.X/64.0), float32(op.Y/64.0))
// 		case svg.OpLineTo:
// 			C.sk_path_line_to(skPath, C.float(op.X)/64.0, C.float(op.Y)/64.0)
// 		case svg.OpQuadTo:
// 			C.sk_path_quad_to(skPath,
// 				C.float(op[0].X)/64.0, C.float(op[0].Y)/64.0,
// 				C.float(op[1].X)/64.0, C.float(op[1].Y)/64.0,
// 			)
// 		case svg.OpCubicTo:
// 			C.sk_path_cubic_to(skPath,
// 				C.float(op[0].X)/64.0, C.float(op[0].Y)/64.0,
// 				C.float(op[1].X)/64.0, C.float(op[1].Y)/64.0,
// 				C.float(op[2].X)/64.0, C.float(op[2].Y)/64.0,
// 			)
// 		case svg.OpClose:
// 			C.sk_path_close(skPath)
// 		}
// 	}

// 	if path.Style.Transform != svg.Identity {
// 		C.sk_path_transform(skPath, &C.sk_matrix_t{
// 			scaleX: C.float(path.Style.Transform.A),
// 			skewX:  C.float(path.Style.Transform.C),
// 			transX: C.float(path.Style.Transform.E),
// 			skewY:  C.float(path.Style.Transform.B),
// 			scaleY: C.float(path.Style.Transform.D),
// 			transY: C.float(path.Style.Transform.F),
// 			persp0: 0,
// 			persp1: 0,
// 			persp2: 1,
// 		})
// 	}
// 	return skPath
// }

// func (s *SVG2) createFillPaint(path svg.SvgPath) *Paint {
// 	if path.Style.FillerColor == nil ||
// 		path.Style.FillOpacity == 0 {
// 		return nil
// 	}

// 	paint := NewPaint()
// 	paint.SetAntialias(true)
// 	paint.SetStyle(PAINTSTYLE_FILL)

// 	if color, ok := path.Style.FillerColor.(svg.PlainColor); ok {
// 		paint.SetColor(ColorARGB(color.A, color.R, color.G, color.B))
// 	} else {
// 		fmt.Fprintf(os.Stderr, "unsupported path fill style %T\n", path.Style.FillerColor)
// 	}

// 	return &paint
// }

// func (s *SVG2) createStrokePaint(path svg.SvgPath) *Paint {
// 	if path.Style.LinerColor == nil ||
// 		path.Style.LineOpacity == 0 ||
// 		path.Style.LineWidth == 0 {
// 		return nil
// 	}

// 	paint := NewPaint()
// 	paint.SetAntialias(true)
// 	paint.SetStrokeCap(path.Style.Join.TrailLineCap)
// 	paint.SetStrokeJoin(path.Style.Join.LineJoin)
// 	paint.SetStrokeWidth(float32(path.Style.LineWidth))
// 	paint.SetStyle(PAINTSTYLE_STROKE)

// 	if color, ok := path.Style.LinerColor.(svg.PlainColor); ok {
// 		paint.SetColor(ColorARGB(color.A, color.R, color.G, color.B))
// 	} else {
// 		fmt.Fprintf(os.Stderr, "unsupported path stroke style %T\n", path.Style.LinerColor)
// 		fmt.Println(path.Style.FillerColor)
// 	}

// 	return &paint
// }

// // Size returns the original size.
// func (s *SVG2) Size() Size {
// 	return s.size
// }

// // OffsetToCenterWithinScaledSize returns the scaled offset values to use to keep the image centered within the given
// // size.
// func (s *SVG2) OffsetToCenterWithinScaledSize(size Size) Point {
// 	scale := min(size.Width/s.size.Width, size.Height/s.size.Height)
// 	return Point{X: (size.Width - s.size.Width*scale) / 2, Y: (size.Height - s.size.Height*scale) / 2}
// }

// // PathScaledTo returns the path with the specified scaling. You should not modify this path, as it is cached.
// func (s *SVG2) PathScaledTo(scale float32) *Path {
// 	if scale == 1 {
// 		return s.unscaledPath
// 	}
// 	scaledSize := Size{Width: scale, Height: scale}
// 	p, ok := s.scaledPathMap[scaledSize]
// 	if !ok {
// 		p = s.unscaledPath.NewScaled(scale, scale)
// 		s.scaledPathMap[scaledSize] = p
// 	}
// 	return p
// }

// // PathForSize returns the path scaled to fit in the specified size. You should not modify this path, as it is cached.
// func (s *SVG2) PathForSize(size Size) *Path {
// 	return s.PathScaledTo(min(size.Width/s.size.Width, size.Height/s.size.Height))
// }

// LogicalSize implements the Drawable interface.
func (s *DrawableSVG2) LogicalSize() Size {
	return s.Size
}

// DrawInRect implements the Drawable interface.
func (s *DrawableSVG2) DrawInRect(canvas *Canvas, rect Rect, _ *SamplingOptions, paint *Paint) {
	canvas.Save()
	defer canvas.Restore()
	// offset := s.SVG.OffsetToCenterWithinScaledSize(rect.Size)
	var offset Point
	canvas.Translate(rect.X+offset.X, rect.Y+offset.Y)
	// canvas.DrawPath(s.SVG.PathForSize(rect.Size), paint)

	for _, path := range s.SVG.paths {
		path.draw(canvas)
	}
}
