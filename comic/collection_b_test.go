package comic

import "testing"

var c2k = Comics{comics: setupComics(20000, false)}
var c2krand = Comics{comics: setupComics(20000, true), sorted: false}

func benchmarkLoading(comics []XKCD, count int, b *testing.B) {
	c := Comics{}
	for i := 0; i < count; i++ {
		c.Load(comics)
	}
}

func benchmarkAdding(comic *XKCD, count int, b *testing.B) {
	c := Comics{}

	for i := 1; i <= count; i++ {
		c.Add(comic)
	}
}

func BenchmarkComicsLoad10(b *testing.B) {
	count := 10
	comics := setupComics(count, false)

	benchmarkLoading(comics, count, b)
}

func BenchmarkComicsLoad100(b *testing.B) {
	count := 100
	comics := setupComics(count, false)

	benchmarkLoading(comics, count, b)
}

func BenchmarkComicsLoad1000(b *testing.B) {
	count := 1000
	comics := setupComics(count, false)

	benchmarkLoading(comics, count, b)
}

func BenchmarkComicAdd10(b *testing.B) {
	count := 10
	comic := &XKCD{Number: 1, Title: "Unit test", ImageURL: "https://xkcd.com/10"}

	benchmarkAdding(comic, count, b)
}

func BenchmarkComicAdd100(b *testing.B) {
	count := 100
	comic := &XKCD{Number: 1, Title: "Unit test", ImageURL: "https://xkcd.com/100"}

	benchmarkAdding(comic, count, b)
}

func BenchmarkComicAdd1000(b *testing.B) {
	count := 1000
	comic := &XKCD{Number: 1, Title: "Unit test", ImageURL: "https://xkcd.com/1000"}

	benchmarkAdding(comic, count, b)
}

func BenchmarkComicGet1234(b *testing.B) {
	for i := 0; i < len(c2k.comics); i++ {
		c2k.Get(1234)
	}
}

func BenchmarkComicGetFirst(b *testing.B) {
	for i := 0; i < len(c2k.comics); i++ {
		c2k.Get(1)
	}
}

func BenchmarkComicGetLast(b *testing.B) {

	for i := 0; i < len(c2k.comics); i++ {
		c2k.Get(2000)
	}
}

func BenchmarkComicGetFromRandomSource(b *testing.B) {
	for i := 0; i < len(c2krand.comics); i++ {
		c2krand.Get(12534)
	}
}
