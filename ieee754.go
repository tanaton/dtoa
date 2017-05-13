package dtoa

import (
	"math"
)

type Double struct {
	u_ uint64
}

func NewDouble(d float64) Double {
	return Double{math.Float64bits(d)}
}
func NewDoubleUint64(u uint64) Double {
	return Double{u}
}

func (d Double) Value() float64 {
	return math.Float64frombits(d.u_)
}
func (d Double) Uint64Value() uint64 {
	return d.u_
}
func (d Double) NextPositiveDouble() float64 {
	return math.Float64frombits(d.u_ + 1)
}
func (d Double) Sign() bool {
	return (d.u_ & kSignMask) != 0
}
func (d Double) Significand() uint64 {
	return d.u_ & kSignificandMask
}
func (d Double) Exponent() int {
	return int(((d.u_ & kExponentMask) >> kSignificandSize) - kExponentBias)
}
func (d Double) IsNan() bool {
	return (d.u_&kExponentMask) == kExponentMask && d.Significand() != 0
}
func (d Double) IsInf() bool {
	return (d.u_&kExponentMask) == kExponentMask && d.Significand() == 0
}
func (d Double) IsNanOrInf() bool {
	return (d.u_ & kExponentMask) == kExponentMask
}
func (d Double) IsNormal() bool {
	return (d.u_&kExponentMask) != 0 || d.Significand() == 0
}
func (d Double) IsZero() bool {
	return (d.u_ & (kExponentMask | kSignificandMask)) == 0
}
func (d Double) IntegerSignificand() uint64 {
	if d.IsNormal() {
		return d.Significand() | kHiddenBit
	} else {
		return d.Significand()
	}
}
func (d Double) IntegerExponent() int {
	var tmp int
	if d.IsNormal() {
		tmp = d.Exponent()
	} else {
		tmp = kDenormalExponent
	}
	return tmp - kSignificandSize
}
func (d Double) ToBias() uint64 {
	if (d.u_ & kSignMask) != 0 {
		return ^d.u_ + 1
	} else {
		return d.u_ | kSignMask
	}
}
func (d Double) EffectiveSignificandSize(order int) int {
	if order >= -1021 {
		return 53
	} else if order <= -1074 {
		return 0
	} else {
		return order + 1074
	}
}
