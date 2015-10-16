package clui

import (
	term "github.com/nsf/termbox-go"
	"time"
)

/* Push button control
onClick function is called every time a user clicks the Button. Sent event has only one valid field sender ID(Ctrl field of Event structure).
Button can be clicked with mouse or using space on keyboard when the button is active.
*/
type Button struct {
	ControlBase
	pressed     bool
	shadowColor term.Attribute

	onClick func(Event)
}

func NewButton(view View, parent Control, width, height int, title string, scale int) *Button {
	b := new(Button)

	b.view = view
	b.parent = parent
	b.align = AlignCenter

	if height < 4 {
		height = 4
	}
	if width < 6 {
		width = 6
	}

	b.SetTitle(title)
	b.SetSize(width, height)
	b.SetConstraints(width, height)

	b.fg = ColorBlue
	b.bg = ColorYellow
	b.shadowColor = ColorBlue

	if parent != nil {
		parent.AddChild(b, scale)
	}

	return b
}

func (b *Button) Repaint() {
	x, y := b.Pos()
	w, h := b.Size()
	canvas := b.view.Canvas()
	tm := b.view.Screen().Theme()

	fg, bg := b.fg, b.bg
	shadow := RealColor(tm, b.shadowColor, ColorControlShadow)
	if !b.Enabled() {
		fg, bg = RealColor(tm, fg, ColorControlDisabledText), RealColor(tm, bg, ColorControlDisabledBack)
	} else if b.Active() {
		fg, bg = RealColor(tm, fg, ColorControlActiveText), RealColor(tm, bg, ColorControlActiveBack)
	} else {
		fg, bg = RealColor(tm, fg, ColorControlText), RealColor(tm, bg, ColorControlBack)
	}

	dy := int((h - 1) / 2)
	shift, text := AlignText(b.title, w-1, b.align)
	if !b.pressed {
		canvas.FillRect(x+1, y+1, w-1, h-1, term.Cell{Ch: ' ', Bg: shadow})
		canvas.FillRect(x, y, w-1, h-1, term.Cell{Ch: ' ', Bg: bg})
		canvas.PutText(x+shift, y+dy, text, fg, bg)
	} else {
		canvas.FillRect(x+1, y+1, w-1, h-1, term.Cell{Ch: ' ', Bg: bg})
		canvas.PutText(x+1+shift, y+1+dy, b.title, fg, bg)
	}
}

func (b *Button) ProcessEvent(event Event) bool {
	if (!b.active && event.Type == EventKey) || !b.Enabled() || b.pressed {
		return false
	}

	if (event.Type == EventKey && event.Key == term.KeySpace) || event.Type == EventMouse {
		b.pressed = true
		timer := time.NewTimer(time.Millisecond * 150)
		go func() {
			<-timer.C
			b.pressed = false
			// generate ButtonClickEvent
			if b.parent != nil {
				if b.onClick != nil {
					ev := Event{Sender: b}
					b.onClick(ev)
				}

				ev := Event{Type: EventRedraw, Sender: b}
				b.view.Screen().PutEvent(ev)
			}
		}()
		return true
	}

	return false
}

func (b *Button) OnClick(fn func(Event)) {
	b.onClick = fn
}