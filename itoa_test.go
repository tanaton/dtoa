package dtoa

import (
	"math"
	"strconv"
	"testing"
)

func TestU32toa(t *testing.T) {
	buf := []byte{}
	table := []uint32{
		0,
		1,
		12,
		123,
		1234,
		12345,
		123456,
		1234567,
		12345678,
		123456789,
		math.MaxUint32,
	}

	for _, it := range table {
		buf = buf[:0]
		buf = U32toa(buf, it)
		s := strconv.FormatUint(uint64(it), 10)
		if string(buf) != s {
			t.Errorf("%d, %s, U32toa = %s\n", it, s, string(buf))
		}
	}
	buf = make([]byte, 100)
	for _, it := range table {
		start := len(buf)
		buf = U32toa(buf, it)
		s := strconv.FormatUint(uint64(it), 10)
		if string(buf[start:]) != s {
			t.Errorf("%d, %s, U32toa = %s\n", it, s, string(buf[start:]))
		}
	}
}

func TestI32toa(t *testing.T) {
	buf := []byte{}
	table := []int32{
		0,
		1,
		12,
		123,
		1234,
		12345,
		123456,
		1234567,
		12345678,
		123456789,
		-1,
		-12,
		-123,
		-1234,
		-12345,
		-123456,
		-1234567,
		-12345678,
		-123456789,
		math.MaxInt32,
		math.MinInt32,
	}

	for _, it := range table {
		buf = buf[:0]
		buf = I32toa(buf, it)
		s := strconv.FormatInt(int64(it), 10)
		if string(buf) != s {
			t.Errorf("%d, %s, I32toa = %s\n", it, s, string(buf))
		}
	}
	buf = make([]byte, 100)
	for _, it := range table {
		start := len(buf)
		buf = I32toa(buf, it)
		s := strconv.FormatInt(int64(it), 10)
		if string(buf[start:]) != s {
			t.Errorf("%d, %s, I32toa = %s\n", it, s, string(buf[start:]))
		}
	}
}

func TestU64toa(t *testing.T) {
	buf := []byte{}
	table := []uint64{
		0,
		1,
		12,
		123,
		1234,
		12345,
		123456,
		1234567,
		12345678,
		123456789,
		1234567891,
		12345678912,
		123456789123,
		1234567891234,
		12345678912345,
		123456789123456,
		1234567891234567,
		12345678912345678,
		123456789123456789,
		1234567891234567891,
		12345678912345678912,
		math.MaxUint64,
	}

	for _, it := range table {
		buf = buf[:0]
		buf = U64toa(buf, it)
		s := strconv.FormatUint(uint64(it), 10)
		if string(buf) != s {
			t.Errorf("%d, %s, U64toa = %s\n", it, s, string(buf))
		}
	}
	buf = make([]byte, 100)
	for _, it := range table {
		start := len(buf)
		buf = U64toa(buf, it)
		s := strconv.FormatUint(uint64(it), 10)
		if string(buf[start:]) != s {
			t.Errorf("%d, %s, U64toa = %s\n", it, s, string(buf[start:]))
		}
	}
}

func TestI64toa(t *testing.T) {
	buf := []byte{}
	table := []int64{
		0,
		1,
		12,
		123,
		1234,
		12345,
		123456,
		1234567,
		12345678,
		123456789,
		1234567891,
		12345678912,
		123456789123,
		1234567891234,
		12345678912345,
		123456789123456,
		1234567891234567,
		12345678912345678,
		123456789123456789,
		1234567891234567891,
		-1,
		-12,
		-123,
		-1234,
		-12345,
		-123456,
		-1234567,
		-12345678,
		-123456789,
		-1234567891,
		-12345678912,
		-123456789123,
		-1234567891234,
		-12345678912345,
		-123456789123456,
		-1234567891234567,
		-12345678912345678,
		-123456789123456789,
		-1234567891234567891,
		math.MaxInt64,
		math.MinInt64,
	}

	for _, it := range table {
		buf = buf[:0]
		buf = I64toa(buf, it)
		s := strconv.FormatInt(int64(it), 10)
		if string(buf) != s {
			t.Errorf("%d, %s, I64toa = %s\n", it, s, string(buf))
		}
	}
	buf = make([]byte, 100)
	for _, it := range table {
		start := len(buf)
		buf = I64toa(buf, it)
		s := strconv.FormatInt(int64(it), 10)
		if string(buf[start:]) != s {
			t.Errorf("%d, %s, I64toa = %s\n", it, s, string(buf[start:]))
		}
	}
}

func BenchmarkStrconvAppendInt1(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendInt(buf, -42, 10)
	}
}
func BenchmarkStrconvAppendInt2(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendInt(buf, math.MaxInt64, 10)
	}
}
func BenchmarkStrconvAppendInt3(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendInt(buf, math.MinInt64, 10)
	}
}

func BenchmarkDtoaI64toa1(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		I64toa(buf, -42)
	}
}
func BenchmarkDtoaI64toa2(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		I64toa(buf, math.MaxInt64)
	}
}
func BenchmarkDtoaI64toa3(b *testing.B) {
	buf := []byte{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		I64toa(buf, math.MinInt64)
	}
}
