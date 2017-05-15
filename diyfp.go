package dtoa

import (
	"math"
)

type DiyFp struct {
	f uint64
	e int
}

func DiyFpDouble(d float64) DiyFp {
	//u64 := *(*uint64)(unsafe.Pointer(&d))
	u64 := math.Float64bits(d)

	biased_e := int((u64 & kDpExponentMask) >> uint64(kDpSignificandSize))
	significand := (u64 & kDpSignificandMask)
	if biased_e != 0 {
		return DiyFp{significand + kDpHiddenBit, biased_e - kDpExponentBias}
	}
	return DiyFp{significand, kDpMinExponent + 1}
}
func (df DiyFp) Minus(rhs DiyFp) DiyFp {
	return DiyFp{df.f - rhs.f, df.e}
}
func (df DiyFp) Multiplication(rhs DiyFp) DiyFp {
	const M32 uint64 = 0xFFFFFFFF
	a := df.f >> 32
	b := df.f & M32
	c := rhs.f >> 32
	d := rhs.f & M32
	ac := a * c
	bc := b * c
	ad := a * d
	bd := b * d
	tmp := (bd >> 32) + (ad & M32) + (bc & M32)
	tmp += uint64(1) << 31 /// mult_round
	return DiyFp{ac + (ad >> 32) + (bc >> 32) + (tmp >> 32), df.e + rhs.e + 64}
}
func (df DiyFp) Normalize() DiyFp {
	for (df.f & (uint64(1) << 63)) == 0 {
		df.f <<= 1
		df.e--
	}
	return df
}
func (df DiyFp) NormalizeBoundary() DiyFp {
	for (df.f & (kDpHiddenBit << 1)) == 0 {
		df.f <<= 1
		df.e--
	}
	df.f <<= (kDiySignificandSize - kDpSignificandSize - 2)
	df.e = df.e - (kDiySignificandSize - kDpSignificandSize - 2)
	return df
}
func (df DiyFp) NormalizedBoundaries() (DiyFp, DiyFp) {
	pl := DiyFp{(df.f << 1) + 1, df.e - 1}.NormalizeBoundary()
	var mi DiyFp
	if df.f == kDpHiddenBit {
		mi = DiyFp{(df.f << 2) - 1, df.e - 2}
	} else {
		mi = DiyFp{(df.f << 1) - 1, df.e - 1}
	}
	mi.f <<= uint64(mi.e - pl.e)
	mi.e = pl.e
	return mi, pl
}
func (df DiyFp) ToDouble() float64 {
	var be uint64
	if df.e == kDpDenormalExponent && (df.f&kDpHiddenBit) == 0 {
		be = 0
	} else {
		be = uint64(df.e + kDpExponentBias)
	}
	be = (df.f & kDpSignificandMask) | (be << kDpSignificandSize)
	return math.Float64frombits(be)
}

func GetCachedPowerByIndex(index uint) DiyFp {
	return DiyFp{kCachedPowers_F[index], kCachedPowers_E[index]}
}

func GetCachedPower(e int) (DiyFp, int) {
	//int k = static_cast<int>(ceil((-61 - e) * 0.30102999566398114)) + 374;
	var dk float64 = (-61-float64(e))*0.30102999566398114 + 347 // dk must be positive, so can do ceiling in positive
	k := int(dk)
	if dk-float64(k) > 0.0 {
		k++
	}
	index := uint((k >> 3) + 1)

	return GetCachedPowerByIndex(index), -(-348 + int(index<<3)) // decimal exponent no need lookup table
}
func GetCachedPower10(exp int, outExp *int) DiyFp {
	var index uint = uint(exp) + 348/8
	*outExp = -348 + int(index)*8
	return GetCachedPowerByIndex(index)
}
