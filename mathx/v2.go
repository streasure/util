package mathx

import "math"

func NewV2(x float64, y float64) V2 {
	return V2{X: x, Y: y}
}

type V2 struct {
	X float64
	Y float64
}

func (v *V2) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v *V2) LenSqrt() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *V2) Normalize() *V2 {
	l := v.Len()
	if l < Epsilon {
		v.X = 0
		v.Y = 0
		return v
	}
	v.X /= l
	v.Y /= l
	return v
}

func (v *V2) Add(other V2) *V2 {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v *V2) Sub(other V2) *V2 {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v *V2) Mul(f float64) *V2 {
	v.X *= f
	v.Y *= f
	return v
}

func (v *V2) Div(f float64) *V2 {
	if f == 0 {
		return v
	}
	v.X /= f
	v.Y /= f
	return v
}

func (v *V2) Dot(other *V2) float64 {
	return v.X*other.X + v.Y*other.Y
}

func Add(v1 V2, v2 V2) V2 {
	return NewV2(v1.X+v2.X, v1.Y+v2.Y)
}

func Sub(v1 V2, v2 V2) V2 {
	return NewV2(v1.X-v2.X, v1.Y-v2.Y)
}

func Mul(v V2, f float64) V2 {
	return NewV2(v.X*f, v.Y*f)
}

func Div(v V2, f float64) V2 {
	if f == 0 {
		return v
	}
	return NewV2(v.X/f, v.Y/f)
}

func Dot(v1 V2, v2 V2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}
