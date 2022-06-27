package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func prepareTasks(tasksCount int, tasks []Task, runTasksCount *int32, errors bool) (time.Duration, []Task) {
	var sumTime time.Duration

	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)

		taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
		sumTime += taskSleep

		tasks = append(tasks, func() error {
			time.Sleep(taskSleep)
			atomic.AddInt32(runTasksCount, 1)
			if errors {
				return err
			}
			return nil
		})
	}
	return sumTime, tasks
}

type testConfig struct {
	tasksCount     int
	workersCount   int
	maxErrorsCount int
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tc := testConfig{50, 23, 10}
		tasks := make([]Task, 0, tc.tasksCount)

		var runTasksCount int32

		_, tasks = prepareTasks(tc.tasksCount, tasks, &runTasksCount, true)

		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(tc.workersCount+tc.maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tc := testConfig{50, 5, 1}
		tasks := make([]Task, 0, tc.tasksCount)

		var runTasksCount int32

		sumTime, tasks := prepareTasks(tc.tasksCount, tasks, &runTasksCount, false)

		start := time.Now()
		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("m <= 0, with errors", func(t *testing.T) {
		tc := testConfig{50, 5, -1}
		tasks := make([]Task, 0, tc.tasksCount)

		var runTasksCount int32

		sumTime, tasks := prepareTasks(tc.tasksCount, tasks, &runTasksCount, true)

		start := time.Now()
		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

}
