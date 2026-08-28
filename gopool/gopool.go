package gopool

import "sync"

type Item interface {
	Execute()
}

type Worker func()

func (w Worker) Execute() {
	w()
}

type Pool struct {
	items chan Item
	once  sync.Once
	wg    sync.WaitGroup
}

func NewPool(count int) *Pool {
	pool := &Pool{
		items: make(chan Item, count*10),
	}
	pool.wg.Add(count)
	for i := 0; i < count; i++ {
		go pool.run()
	}
	return pool
}

func (p *Pool) Add(item Item) {
	p.items <- item
}

func (p *Pool) run() {
	defer p.wg.Done()
	for item := range p.items {
		item.Execute()
	}
}

func (p *Pool) Stop() {
	p.once.Do(func() {
		close(p.items)
	})
	p.wg.Wait()
}

type Task[T any] struct {
	Fn   func() T
	Ch   chan<- T
}

func (t *Task[T]) Execute() {
	t.Ch <- t.Fn()
}

type TypedPool[T any] struct {
	items chan *Task[T]
	once  sync.Once
	wg    sync.WaitGroup
}

func NewTypedPool[T any](count int) *TypedPool[T] {
	pool := &TypedPool[T]{
		items: make(chan *Task[T], count*10),
	}
	pool.wg.Add(count)
	for i := 0; i < count; i++ {
		go pool.run()
	}
	return pool
}

func (p *TypedPool[T]) Add(fn func() T) <-chan T {
	ch := make(chan T, 1)
	p.items <- &Task[T]{Fn: fn, Ch: ch}
	return ch
}

func (p *TypedPool[T]) run() {
	defer p.wg.Done()
	for task := range p.items {
		task.Execute()
	}
}

func (p *TypedPool[T]) Stop() {
	p.once.Do(func() {
		close(p.items)
	})
	p.wg.Wait()
}
