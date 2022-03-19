// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
)

// DefaultPopupMenuTheme holds the default PopupMenuTheme values for PopupMenus. Modifying this data will not alter
// existing PopupMenus, but will alter any PopupMenus created in the future.
var DefaultPopupMenuTheme = PopupMenuTheme{
	Font:            SystemFont,
	BackgroundInk:   ControlColor,
	OnBackgroundInk: OnControlColor,
	EdgeInk:         ControlEdgeColor,
	SelectionInk:    SelectionColor,
	OnSelectionInk:  OnSelectionColor,
	CornerRadius:    4,
	HMargin:         8,
	VMargin:         1,
}

// PopupMenuTheme holds theming data for a PopupMenu.
type PopupMenuTheme struct {
	Font            Font
	BackgroundInk   Ink
	OnBackgroundInk Ink
	EdgeInk         Ink
	SelectionInk    Ink
	OnSelectionInk  Ink
	CornerRadius    float32
	HMargin         float32
	VMargin         float32
}

type popupMenuItem struct {
	item      interface{}
	textCache TextCache
	enabled   bool
}

// PopupMenu represents a clickable button that displays a menu of choices.
type PopupMenu struct {
	Panel
	PopupMenuTheme
	MenuFactory       MenuFactory
	SelectionCallback func()
	items             []*popupMenuItem
	selectedIndex     int
	textCache         TextCache
	Pressed           bool
}

type popupSeparationMarker struct{}

// NewPopupMenu creates a new PopupMenu.
func NewPopupMenu() *PopupMenu {
	p := &PopupMenu{PopupMenuTheme: DefaultPopupMenuTheme}
	p.Self = p
	p.SetFocusable(true)
	p.SetSizer(p.DefaultSizes)
	p.MenuFactory = DefaultMenuFactory()
	p.DrawCallback = p.DefaultDraw
	p.GainedFocusCallback = p.MarkForRedraw
	p.LostFocusCallback = p.MarkForRedraw
	p.MouseDownCallback = p.DefaultMouseDown
	p.KeyDownCallback = p.DefaultKeyDown
	return p
}

// DefaultSizes provides the default sizing.
func (p *PopupMenu) DefaultSizes(hint geom32.Size) (min, pref, max geom32.Size) {
	pref = LabelSize(p.textCache.Text("M", p.Font), nil, 0, 0)
	for _, one := range p.items {
		switch one.item.(type) {
		case *popupSeparationMarker:
		default:
			size := LabelSize(one.textCache.Text(fmt.Sprintf("%v", one.item), p.Font), nil, 0, 0)
			if pref.Width < size.Width {
				pref.Width = size.Width
			}
			if pref.Height < size.Height {
				pref.Height = size.Height
			}
		}
	}
	if border := p.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.Height += p.VMargin*2 + 2
	pref.Width += p.HMargin*2 + 2 + pref.Height*0.75
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	max.Width = mathf32.Max(DefaultMaxSize, pref.Width)
	max.Height = pref.Height
	return pref, pref, max
}

// DefaultDraw provides the default drawing.
func (p *PopupMenu) DefaultDraw(canvas *Canvas, dirty geom32.Rect) {
	var fg, bg Ink
	switch {
	case p.Pressed:
		bg = p.SelectionInk
		fg = p.OnSelectionInk
	default:
		bg = p.BackgroundInk
		fg = p.OnBackgroundInk
	}
	thickness := float32(1)
	if p.Focused() {
		thickness++
	}
	rect := p.ContentRect(false)
	DrawRoundedRectBase(canvas, rect, p.CornerRadius, thickness, bg, p.EdgeInk)
	rect.InsetUniform(1.5)
	rect.X += p.HMargin
	rect.Y += p.VMargin
	rect.Width -= p.HMargin * 2
	rect.Height -= p.VMargin * 2
	triWidth := rect.Height * 0.75
	triHeight := triWidth / 2
	rect.Width -= triWidth
	DrawLabel(canvas, rect, StartAlignment, MiddleAlignment, p.textObj(), fg, nil, 0, 0, !p.Enabled())
	rect.Width += triWidth + p.HMargin/2
	path := NewPath()
	path.MoveTo(rect.Right(), rect.Y+(rect.Height-triHeight)/2)
	path.LineTo(rect.Right()-triWidth, rect.Y+(rect.Height-triHeight)/2)
	path.LineTo(rect.Right()-triWidth/2, rect.Y+(rect.Height-triHeight)/2+triHeight)
	path.Close()
	paint := fg.Paint(canvas, rect, Fill)
	if !p.Enabled() {
		paint.SetColorFilter(Grayscale30PercentFilter())
	}
	canvas.DrawPath(path, paint)
}

func (p *PopupMenu) textObj() *Text {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.items) {
		one := p.items[p.selectedIndex].item
		switch one.(type) {
		case *popupSeparationMarker:
		default:
			return p.items[p.selectedIndex].textCache.Text(fmt.Sprintf("%v", one), p.Font)
		}
	}
	return nil
}

// Text the currently shown text.
func (p *PopupMenu) Text() string {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.items) {
		one := p.items[p.selectedIndex].item
		switch one.(type) {
		case *popupSeparationMarker:
		default:
			return fmt.Sprintf("%v", one)
		}
	}
	return ""
}

// Click performs any animation associated with a click and triggers the popup menu to appear.
func (p *PopupMenu) Click() {
	hasItem := false //nolint:ifshort // Cannot collapse this into the if statement, despite what the linter says
	m := p.MenuFactory.NewMenu(PopupMenuTemporaryBaseID, "", nil)
	defer m.Dispose()
	for i, one := range p.items {
		if _, ok := one.item.(*popupSeparationMarker); ok {
			m.InsertSeparator(-1, false)
		} else {
			hasItem = true
			m.InsertItem(-1, p.createMenuItem(m, i, one))
		}
	}
	if hasItem {
		m.Popup(p.RectToRoot(p.ContentRect(true)), p.selectedIndex)
	}
}

func (p *PopupMenu) createMenuItem(m Menu, index int, entry *popupMenuItem) MenuItem {
	return m.Factory().NewItem(PopupMenuTemporaryBaseID+index+1,
		fmt.Sprintf("%v", entry.item), KeyBinding{}, func(mi MenuItem) bool {
			return entry.enabled
		}, func(mi MenuItem) {
			if index != p.SelectedIndex() {
				p.SelectIndex(index)
				mi.SetCheckState(OnCheckState)
			}
		})
}

// AddItem appends a menu item to the end of the PopupMenu.
func (p *PopupMenu) AddItem(item interface{}) *PopupMenu {
	p.items = append(p.items, &popupMenuItem{
		item:    item,
		enabled: true,
	})
	return p
}

// AddDisabledItem appends a disabled menu item to the end of the PopupMenu.
func (p *PopupMenu) AddDisabledItem(item interface{}) *PopupMenu {
	p.items = append(p.items, &popupMenuItem{item: item})
	return p
}

// AddSeparator adds a separator to the end of the PopupMenu.
func (p *PopupMenu) AddSeparator() *PopupMenu {
	return p.AddDisabledItem(&popupSeparationMarker{})
}

// IndexOfItem returns the index of the specified menu item. -1 will be returned if the menu item isn't present.
func (p *PopupMenu) IndexOfItem(item interface{}) int {
	for i, one := range p.items {
		if one.item == item {
			return i
		}
	}
	return -1
}

// RemoveAllItems removes all items from the PopupMenu.
func (p *PopupMenu) RemoveAllItems() *PopupMenu {
	p.selectedIndex = 0
	p.items = nil
	p.MarkForRedraw()
	return p
}

// RemoveItem from the PopupMenu.
func (p *PopupMenu) RemoveItem(item interface{}) *PopupMenu {
	p.RemoveItemAt(p.IndexOfItem(item))
	return p
}

// RemoveItemAt the specified index from the PopupMenu.
func (p *PopupMenu) RemoveItemAt(index int) *PopupMenu {
	if index >= 0 {
		length := len(p.items)
		if index < length {
			if p.selectedIndex == index {
				if p.selectedIndex > length-2 {
					p.selectedIndex = length - 2
					if p.selectedIndex < 0 {
						p.selectedIndex = 0
					}
				}
				p.MarkForRedraw()
			} else if p.selectedIndex > index {
				p.selectedIndex--
			}
			copy(p.items[index:], p.items[index+1:])
			length--
			p.items[length] = nil
			p.items = p.items[:length]
		}
	}
	return p
}

// ItemCount returns the number of items in this PopupMenu.
func (p *PopupMenu) ItemCount() int {
	return len(p.items)
}

// ItemAt returns the menuItem at the specified index or nil.
func (p *PopupMenu) ItemAt(index int) interface{} {
	if index >= 0 && index < len(p.items) {
		return p.items[index].item
	}
	return nil
}

// SetItemAt sets the menuItem at the specified index.
func (p *PopupMenu) SetItemAt(index int, item interface{}, enabled bool) *PopupMenu {
	if index >= 0 && index < len(p.items) {
		if p.items[index].item != item {
			p.items[index].item = item
			p.items[index].enabled = enabled
			p.MarkForRedraw()
		}
	}
	return p
}

// Selected returns the currently selected menuItem or nil.
func (p *PopupMenu) Selected() interface{} {
	return p.ItemAt(p.selectedIndex)
}

// SelectedIndex returns the currently selected menuItem index.
func (p *PopupMenu) SelectedIndex() int {
	return p.selectedIndex
}

// Select an menuItem.
func (p *PopupMenu) Select(item interface{}) *PopupMenu {
	p.SelectIndex(p.IndexOfItem(item))
	return p
}

// SelectIndex selects an menuItem by its index.
func (p *PopupMenu) SelectIndex(index int) *PopupMenu {
	if index != p.selectedIndex && index >= 0 && index < len(p.items) {
		p.selectedIndex = index
		p.MarkForRedraw()
		if p.SelectionCallback != nil {
			p.SelectionCallback()
		}
	}
	return p
}

// DefaultMouseDown provides the default mouse down handling.
func (p *PopupMenu) DefaultMouseDown(where geom32.Point, button, clickCount int, mod Modifiers) bool {
	p.Click()
	return true
}

// DefaultKeyDown provides the default key down handling.
func (p *PopupMenu) DefaultKeyDown(keyCode KeyCode, mod Modifiers, repeat bool) bool {
	if IsControlAction(keyCode, mod) {
		p.Click()
		return true
	}
	return false
}
