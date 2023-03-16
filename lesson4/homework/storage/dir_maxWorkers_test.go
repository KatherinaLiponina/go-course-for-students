package storage

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_MaxWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok, dummy storage", func(t *testing.T) {
		var s sizer = sizer{}
		s.maxWorkersCount = 5
		s.maxSet = true
		ch := make(chan int64)
		var b bool = true

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go maxWorkerCount(&s, &b, ch)
		result, err := s.Size(ctx, getDummySet())
		b = false
		max := <-ch
		assert.LessOrEqual(t, max, int64(s.maxWorkersCount))
		assert.NoError(t, err)
		assert.Equal(t, int64(14), result.Count)
		assert.Equal(t, int64(37254162), result.Size)
	})
}

func maxWorkerCount(s * sizer, b * bool, ch chan int64) () {
	var max int64 = -1
	for *b {
		if s.currentWorkersCount > max {
			max = s.currentWorkersCount
		}
		time.Sleep(10 * time.Millisecond)
	}
	ch <- max
}
