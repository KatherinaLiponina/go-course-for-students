package storage

import (
	"context"
	"errors"
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
	maxWorkersCount int

	// TODO: add other fields as you wish
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
	go exploreDir(context.TODO(), d, collector, errorChannel)
	for childrenAmount != 0 {
		select {
		case r:= <-collector:
			res.Count += r.Count
			res.Size += r.Size
		case e := <-errorChannel:
			return Result{}, errors.Join(errors.New("error occured in one of goroutines"), e)
		}
	}
	close(collector)
	close(errorChannel)
	return res, nil
}

func exploreDir(ctx context.Context, d Dir, resultChannel chan Result, errorChannel chan error) () {
	var res Result
	directories, files, err := d.Ls(ctx)
	defer atomic.AddInt64(&childrenAmount, -1)
	if err != nil {
		errorChannel <- err
		return
	}
	if directories != nil {
		for _, dir := range directories {
			res.Count ++
			atomic.AddInt64(&childrenAmount, 1)
			go exploreDir(context.TODO(), dir, resultChannel, errorChannel)
		}
	}
	if files != nil {
		for _, f := range files {
			res.Count++
			size, err := f.Stat(context.TODO())
			if err != nil {
				errorChannel <- err
				return
			}
			res.Size += size
		}
	}
	resultChannel <- res
}
