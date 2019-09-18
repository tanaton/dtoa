package dtoa

import (
	"math"
	"strconv"
	"testing"
)

func TestDtoa(t *testing.T) {
	buf := []byte{}
	table := []struct {
		i int
		f float64
		s string
	}{
		{324, 0.0, "0.0"},
		//{324, -0.0, "-0.0"},
		{324, math.Float64frombits(1 << 63), "-0.0"},
		{324, 1.0, "1.0"},
		{324, -1.0, "-1.0"},
		{324, 1.2345, "1.2345"},
		{324, 1.2345678, "1.2345678"},
		{324, 0.123456789012, "0.123456789012"},
		{324, 1234567.8, "1234567.8"},
		{324, -79.39773355813419, "-79.39773355813419"},
		{324, 0.000001, "0.000001"},
		{324, 0.0000001, "1e-7"},
		{324, 1e30, "1e30"},
		{324, 1.234567890123456e30, "1.234567890123456e30"},
		{324, 5e-324, "5e-324"},                                   // Min subnormal positive double
		{324, 2.225073858507201e-308, "2.225073858507201e-308"},   // Max subnormal positive double
		{324, 2.2250738585072014e-308, "2.2250738585072014e-308"}, // Min normal positive double
		{324, 1.7976931348623157e308, "1.7976931348623157e308"},   // Max double
		{3, 0.0, "0.0"},
		{1, 0.0, "0.0"},
		//{3, -0.0, "-0.0"},
		{3, math.Float64frombits(1 << 63), "-0.0"},
		{3, 1.0, "1.0"},
		{3, -1.0, "-1.0"},
		{3, 1.2345, "1.234"},
		{2, 1.2345, "1.23"},
		{1, 1.2345, "1.2"},
		{3, 1.2345678, "1.234"},
		{3, 1.0001, "1.0"},
		{2, 1.0001, "1.0"},
		{1, 1.0001, "1.0"},
		{3, 0.123456789012, "0.123"},
		{2, 0.123456789012, "0.12"},
		{1, 0.123456789012, "0.1"},
		{4, 0.0001, "0.0001"},
		{3, 0.0001, "0.0"},
		{2, 0.0001, "0.0"},
		{1, 0.0001, "0.0"},
		{3, 1234567.8, "1234567.8"},
		{3, 1e30, "1e30"},
		{3, 5e-324, "0.0"},                  // Min subnormal positive double
		{3, 2.225073858507201e-308, "0.0"},  // Max subnormal positive double
		{3, 2.2250738585072014e-308, "0.0"}, // Min normal positive double
		{3, 1.7976931348623157e308, "1.7976931348623157e308"}, // Max double
		{5, -0.14000000000000001, "-0.14"},
		{4, -0.14000000000000001, "-0.14"},
		{3, -0.14000000000000001, "-0.14"},
		{3, -0.10000000000000001, "-0.1"},
		{2, -0.10000000000000001, "-0.1"},
		{1, -0.10000000000000001, "-0.1"},
	}

	for _, it := range table {
		buf = buf[:0]
		buf = Dtoa(buf, it.f, it.i)
		if string(buf) != it.s {
			t.Errorf("%f, %d, %s, Dtoa = %s\n", it.f, it.i, it.s, string(buf))
		}
	}
	buf = make([]byte, 100)
	for _, it := range table {
		start := len(buf)
		buf = Dtoa(buf, it.f, it.i)
		if string(buf[start:]) != it.s {
			t.Errorf("%f, %d, %s, Dtoa = %s\n", it.f, it.i, it.s, string(buf[start:]))
		}
	}
}

func BenchmarkStrconvAppendFloat1(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendFloat(buf, 1.1, 'f', -1, 64)
	}
}
func BenchmarkStrconvAppendFloat2(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendFloat(buf, 3.1415926535, 'f', -1, 64)
	}
}
func BenchmarkStrconvAppendFloat3(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendFloat(buf, 2.225073858507201e-308, 'f', 1, 64)
	}
}

func BenchmarkDtoaDtoa1(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dtoa(buf, 1.1, -1)
	}
}
func BenchmarkDtoaDtoa2(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dtoa(buf, 3.1415926535, -1)
	}
}
func BenchmarkDtoaDtoa3(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dtoa(buf, 2.225073858507201e-308, 1)
	}
}
