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
	"fmt"
	"runtime"

	"github.com/lafriks/go-svg"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/strokecap"
	"github.com/richardwilkes/unison/enums/strokejoin"
	"github.com/richardwilkes/unison/internal/skia"
)

var strokeCaps = map[svg.CapMode]strokecap.Enum{
	svg.NilCap:       strokecap.Butt, // default to something that won't break
	svg.ButtCap:      strokecap.Butt,
	svg.SquareCap:    strokecap.Square,
	svg.RoundCap:     strokecap.Round,
	svg.CubicCap:     strokecap.Butt, // default to something that won't break
	svg.QuadraticCap: strokecap.Butt, // default to something that won't break
}

var strokeJoins = map[svg.JoinMode]strokejoin.Enum{
	svg.Arc:       strokejoin.Miter, // default to something that won't break
	svg.Round:     strokejoin.Round,
	svg.Bevel:     strokejoin.Bevel,
	svg.Miter:     strokejoin.Miter,
	svg.MiterClip: strokejoin.Miter, // default to something that won't break
	svg.ArcClip:   strokejoin.Miter, // default to something that won't break
}

// Path2 holds geometry.
type Path2 struct {
	svgPath     svg.SvgPath
	path        skia.Path
	fillPaint   *Paint // if nil don't draw
	strokePaint *Paint // if nil don't draw
}

func newPath2(svgPath svg.SvgPath) (*Path2, error) {
	p := &Path2{
		svgPath: svgPath,
	}

	p.createPath()
	if err := p.createFillPaint(); err != nil {
		return nil, err
	}
	if err := p.createStrokePaint(); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(p, func(obj *Path2) {
		ReleaseOnUIThread(func() {
			skia.PathDelete(obj.path)
		})
	})
	return p, nil
}

// createPath creates a skia path from an svg.Path.
//
// The coordinates used in svg.Operation are of type fixed.Int26_6, which has a fractional
// part of 6 bits. When converting to a float the values are divided by 64.
func (p *Path2) createPath() {
	p.path = skia.PathNew()

	for _, op := range p.svgPath.Path {
		switch op := op.(type) {
		case svg.OpMoveTo:
			skia.PathMoveTo(p.path, float32(op.X)/64.0, float32(op.Y)/64.0)
		case svg.OpLineTo:
			skia.PathLineTo(p.path, float32(op.X)/64.0, float32(op.Y)/64.0)
		case svg.OpQuadTo:
			skia.PathQuadTo(p.path,
				float32(op[0].X)/64.0, float32(op[0].Y)/64.0,
				float32(op[1].X)/64.0, float32(op[1].Y)/64.0,
			)
		case svg.OpCubicTo:
			skia.PathCubicTo(p.path,
				float32(op[0].X)/64.0, float32(op[0].Y)/64.0,
				float32(op[1].X)/64.0, float32(op[1].Y)/64.0,
				float32(op[2].X)/64.0, float32(op[2].Y)/64.0,
			)
		case svg.OpClose:
			skia.PathClose(p.path)
		}
	}

	if p.svgPath.Style.Transform != svg.Identity {
		skia.PathTransform(p.path, geom.Matrix[float32]{
			ScaleX: float32(p.svgPath.Style.Transform.A),
			SkewX:  float32(p.svgPath.Style.Transform.C),
			TransX: float32(p.svgPath.Style.Transform.E),
			SkewY:  float32(p.svgPath.Style.Transform.B),
			ScaleY: float32(p.svgPath.Style.Transform.D),
			TransY: float32(p.svgPath.Style.Transform.F),
		})
	}
}

func (p *Path2) createFillPaint() error {
	if p.svgPath.Style.FillerColor == nil ||
		p.svgPath.Style.FillOpacity == 0 {
		return nil
	}

	p.fillPaint = NewPaint()
	p.fillPaint.SetAntialias(true)
	p.fillPaint.SetStyle(paintstyle.Fill)

	if color, ok := p.svgPath.Style.FillerColor.(svg.PlainColor); ok {
		alpha := float32(color.A) / 0xFF
		alpha *= float32(p.svgPath.Style.FillOpacity)
		p.fillPaint.SetColor(ARGB(alpha, int(color.R), int(color.G), int(color.B)))
	} else {
		return fmt.Errorf("unsupported path fill style %T", p.svgPath.Style.FillerColor)
	}

	return nil
}

func (p *Path2) createStrokePaint() error {
	if p.svgPath.Style.LinerColor == nil ||
		p.svgPath.Style.LineOpacity == 0 ||
		p.svgPath.Style.LineWidth == 0 {
		return nil
	}

	p.strokePaint = NewPaint()
	p.strokePaint.SetAntialias(true)
	p.strokePaint.SetStrokeCap(strokeCaps[p.svgPath.Style.Join.TrailLineCap])
	p.strokePaint.SetStrokeJoin(strokeJoins[p.svgPath.Style.Join.LineJoin])
	p.strokePaint.SetStrokeWidth(float32(p.svgPath.Style.LineWidth))
	p.strokePaint.SetStyle(paintstyle.Stroke)

	if color, ok := p.svgPath.Style.LinerColor.(svg.PlainColor); ok {
		alpha := float32(color.A) / 0xFF
		alpha *= float32(p.svgPath.Style.LineOpacity)
		p.strokePaint.SetColor(ARGB(alpha, int(color.R), int(color.G), int(color.B)))
	} else {
		return fmt.Errorf("unsupported path stroke style %T", p.svgPath.Style.LinerColor)
	}

	return nil
}

func (p *Path2) draw(canvas *Canvas) {
	if p.fillPaint != nil {
		canvas.DrawPath2(p, p.fillPaint)
	}
	if p.strokePaint != nil {
		canvas.DrawPath2(p, p.strokePaint)
	}
}
