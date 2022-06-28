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

type testConfig struct {
	tasksCount     int
	workersCount   int
	maxErrorsCount int
	errors         bool
}

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

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("1. if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tc := testConfig{50, 23, 10, true}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		_, tasks = prepareTasks(tc.tasksCount, tasks, &runTasksCount, tc.errors)

		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(tc.workersCount+tc.maxErrorsCount), "extra tasks were started")
	})

	t.Run("2. tasks without errors", func(t *testing.T) {
		tc := testConfig{50, 23, 1, false}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		sumTime, tasks := prepareTasks(tc.tasksCount, tasks, &runTasksCount, tc.errors)

		start := time.Now()
		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("3. m <= 0, with errors", func(t *testing.T) {
		tc := testConfig{50, 5, -1, true}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		sumTime, tasks := prepareTasks(tc.tasksCount, tasks, &runTasksCount, tc.errors)

		start := time.Now()
		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("4. n=1, without errors", func(t *testing.T) {
		tc := testConfig{10, 1, 1, false}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		sumTime, tasks := prepareTasks(tc.tasksCount, tasks, &runTasksCount, tc.errors)

		start := time.Now()
		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(sumTime), int64(elapsedTime), "tasks were not run sequentially?")
	})

	t.Run("5. Num tasks = m errors", func(t *testing.T) {
		tc := testConfig{10, 10, 10, true}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		_, tasks = prepareTasks(tc.tasksCount, tasks, &runTasksCount, tc.errors)

		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, runTasksCount, int32(tc.workersCount+tc.maxErrorsCount), "extra tasks were started")
	})

	t.Run("6. workers more than tasks", func(t *testing.T) {
		tc := testConfig{10, 15, 1, false}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		_, tasks = prepareTasks(tc.tasksCount, tasks, &runTasksCount, tc.errors)

		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, runTasksCount, int32(tc.workersCount+tc.maxErrorsCount), "extra tasks were started")
	})

	t.Run("7. workers less than tasks", func(t *testing.T) {
		tc := testConfig{50, 10, 25, true}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		_, tasks = prepareTasks(tc.tasksCount, tasks, &runTasksCount, tc.errors)

		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tc.workersCount))
	})

	t.Run("8. tasks with errors less then M", func(t *testing.T) {
		tc := testConfig{50, 10, 10, false}
		tasks := make([]Task, 0, tc.tasksCount)
		var runTasksCount int32

		errorsCount := 0
		for i := 0; i < tc.tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

			if errorsCount < (tc.maxErrorsCount - 5) {
				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			} else {
				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}
			errorsCount++
		}

		result := Run(tasks, tc.workersCount, tc.maxErrorsCount)

		require.NotEqual(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(tc.tasksCount), "all tasks were started")
	})
}
