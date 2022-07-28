//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	mb          uint64 = 1 << 20
	memoryLimit uint64 = 30 * mb

	timeLimit = 300 * time.Millisecond
)

var data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

// go test -benchmem -v -run=^$ -tags bench -bench .
// old BenchmarkStats-4              76          14461933 ns/op         1769869 B/op      22376 allocs/op
// new
func BenchmarkStats(b *testing.B) {
	infoLog.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := bytes.NewBufferString(data)
		b.StartTimer()
		GetDomainStat(r, "biz")
		b.StopTimer()
	}
}

// old BenchmarkStats-4              1        1963230600 ns/op        310943680 B/op   3045380 allocs/op
func BenchmarkStatsZip(b *testing.B) {
	infoLog.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		r, _ := zip.OpenReader("testdata/users.dat.zip")
		defer r.Close()

		data, _ := r.File[0].Open()

		b.StartTimer()
		GetDomainStat(data, "biz")
		b.StopTimer()
	}
}

// go test -v -count=1 -timeout=30s -tags bench .
func TestGetDomainStat_Time_And_Memory(t *testing.T) {
	bench := func(b *testing.B) {
		b.Helper()
		b.StopTimer()

		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(t, err)

		b.StartTimer()
		stat, err := GetDomainStat(data, "biz")
		b.StopTimer()
		require.NoError(t, err)

		// add check from file
		require.Equal(t, expectedBizStat, stat)
	}

	result := testing.Benchmark(bench)
	mem := result.MemBytes
	t.Logf("time used: %s / %s", result.T, timeLimit)
	t.Logf("memory used: %dMb / %dMb", mem/mb, memoryLimit/mb)

	require.Less(t, int64(result.T), int64(timeLimit), "the program is too slow")
	require.Less(t, mem, memoryLimit, "the program is too greedy")

}

var expectedBizStat = DomainStat{}
