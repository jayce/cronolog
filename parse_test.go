package cronolog

import "testing"

func TestSizeString(t *testing.T) {
	data := []struct {
		b   Size
		str string
	}{
		{Size(0), "0b"}, {Byte, "1b"}, {Size(500), "500b"},
		{KiB + Size(500), "1.48kb"}, {-Size(1 * KiB), "-1kb"},
		{MiB, "1mb"}, {MiB, "1mb"},
		{25 * GiB, "25gb"},
	}

	for _, v := range data {
		if v.b.String() != v.str {
			t.Errorf("%s != %s", v.b, v.str)
		}
	}
}

func BenchmarkSizeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Size(20 * KiB).String()
	}
}

func TestParseSize(t *testing.T) {
	goods := []struct {
		s    string
		want string
	}{
		{"-1.2", "-1b"},
		{"-1.2b", "-1b"},
		{"1.2kb", "1.19kb"},
		{"1mb", "1mb"},
		{"1.3245gb", "1.32gb"},
	}

	for _, v := range goods {
		size, err := ParseSize(v.s)
		if err != nil {
			t.Error(err)
		}

		if size.String() != v.want {
			t.Errorf("%s want %s, but got %s", v.s, v.want, size)
		}
	}

	bads := []string{
		"-1..1b",
		"-1.b",
		"--1.1b",
		".000b",
		"asd",
	}

	for _, v := range bads {
		_, err := ParseSize(v)
		if err == nil {
			t.Errorf("%s is a bad 'Size', but passed", v)
		}
	}
}
