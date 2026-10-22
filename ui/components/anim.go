package components

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op"
)

type FloatAnim struct {
	Value    float32
	Target   float32
	Duration time.Duration
	lastTime time.Time
}

func NewFloatAnim(initial float32, duration time.Duration) *FloatAnim {
	return &FloatAnim{
		Value:    initial,
		Target:   initial,
		Duration: duration,
	}
}

func (a *FloatAnim) SetTarget(target float32) {
	if a.Target != target {
		a.Target = target
		a.lastTime = time.Time{}
	}
}

func (a *FloatAnim) IsAnimating() bool {
	return a.Value != a.Target
}

func (a *FloatAnim) Update(gtx layout.Context) float32 {
	if a.Value == a.Target {
		a.lastTime = time.Time{}
		return a.Value
	}

	now := gtx.Now
	if a.lastTime.IsZero() {
		a.lastTime = now
	}
	dt := float32(now.Sub(a.lastTime).Seconds())
	a.lastTime = now

	// Cap dt in case frame was paused or took long
	if dt > 0.1 {
		dt = 0.016
	}

	if a.Duration <= 0 {
		a.Value = a.Target
		return a.Value
	}

	speed := 1.0 / float32(a.Duration.Seconds())
	delta := speed * dt

	if a.Value < a.Target {
		a.Value += delta
		if a.Value > a.Target {
			a.Value = a.Target
		}
	} else if a.Value > a.Target {
		a.Value -= delta
		if a.Value < a.Target {
			a.Value = a.Target
		}
	}

	if a.Value != a.Target {
		gtx.Execute(op.InvalidateCmd{})
	}
	return a.Value
}

// EaseOutCubic applies cubic ease out curve: 1 - (1-t)^3
func EaseOutCubic(t float32) float32 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	f := 1.0 - t
	return 1.0 - f*f*f
}
