package hw06pipelineexecution

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage
)

func sendIntData(in Bi, data []int) {
	for _, v := range data {
		fmt.Printf("[main] send %v\n", v)
		in <- v
	}
	close(in)
}

var g = func(s string, f func(v interface{}) interface{}) Stage {
	// Stage generator
	return func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for v := range in {
				// fmt.Printf("[stage %s] received %v\n", s, v)
				time.Sleep(sleepPerStage)
				out <- f(v)
				// fmt.Printf("\t[stage %s] send\n", s)
			}
		}()
		return out
	}
}

var stages = []Stage{
	g("Dummy", func(v interface{}) interface{} { return v }),
	g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
	g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
	g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
}

func TestPipeline(t *testing.T) {
	t.Run("1 stage", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go sendIntData(in, data)

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages[3]) {
			fmt.Printf("[main] received %s\n", s)
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"1", "2", "3", "4", "5"}, result)
		require.Less(t,
			int64(elapsed),
			int64(sleepPerStage)*int64(len(data))+int64(fault))
	})

	t.Run("0 stage", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go sendIntData(in, data)

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, nil, make([]Stage, 0)...) {
			fmt.Printf("[main] received %s\n", s)
			result = append(result, s.(string))
		}

		require.Equal(t, 0, len(result))
	})

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go sendIntData(in, data)

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			// fmt.Printf("[main] received %s\n", s)
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("empty data", func(t *testing.T) {
		in := make(Bi)
		data := []int{}

		go sendIntData(in, data)

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			fmt.Printf("[main] received %s\n", s)
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, 0, len(result))
		require.Less(t,
			int64(elapsed),
			int64(sleepPerStage)*int64(len(stages))+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go sendIntData(in, data)

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}
