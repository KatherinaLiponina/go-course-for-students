package storage

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount     int
	maxSet              bool
	cwcMutex            sync.Mutex
	currentWorkersCount int64
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

var childrenAmount int64

func (a *sizer) Size(ctx context.Context, d Dir) (res Result, err error) {
	collector := make(chan Result)
	errorChannel := make(chan error)
	childrenAmount = 1
	go a.exploreDir(ctx, d, collector, errorChannel)
	for childrenAmount != 0 {
		select {
		case r, ok := <-collector:
			if !ok {
				goto done
			}
			res.Count += r.Count
			res.Size += r.Size
		case e := <-errorChannel:
			return Result{}, errors.Join(errors.New("error occured in one of goroutines"), e)
		case <-ctx.Done():
			return Result{}, errors.New("context was canceled")
		}
	}
done:
	close(errorChannel)
	return res, nil
}

func (a *sizer) exploreDir(ctx context.Context, d Dir, resultChannel chan Result, errorChannel chan error) {
	var res Result
	directories, files, err := d.Ls(ctx)
	defer atomic.AddInt64(&childrenAmount, -1)
	defer func() {
		if childrenAmount-1 == 0 {
			close(resultChannel)
		}
	}()
	if err != nil {
		errorChannel <- err
		return
	}
	for _, f := range files {
		res.Count++
		size, err := f.Stat(ctx)
		if err != nil {
			errorChannel <- err
			return
		}
		res.Size += size
	}
	if !a.maxSet {
		for _, dir := range directories {
			atomic.AddInt64(&childrenAmount, 1)
			go a.exploreDir(ctx, dir, resultChannel, errorChannel)
		}
		resultChannel <- res
		return
	}

	for _, dir := range directories {
		//for atomic operation
		a.cwcMutex.Lock()
		if a.currentWorkersCount < int64(a.maxWorkersCount) {
			//possible to create goroutine
			a.currentWorkersCount += 1
			a.cwcMutex.Unlock()
			atomic.AddInt64(&childrenAmount, 1)
			go a.exploreDir(ctx, dir, resultChannel, errorChannel)
		} else {
			a.cwcMutex.Unlock()
			//just continue on this goroutine
			//use add 'cause function is ending althougt goroutine is not
			atomic.AddInt64(&childrenAmount, 1)
			a.exploreDir(ctx, dir, resultChannel, errorChannel)
		}
	}
	resultChannel <- res
}
