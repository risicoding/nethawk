package cache

import (
	"sync"
	"time"
)

type Cache[T any] struct {
	mu       sync.RWMutex
	function func() (T, error)
	value    T
	err      error
	stop     chan struct{}
	close    sync.Once
}

func New[T any](function func() (T, error)) *Cache[T] {
	c := &Cache[T]{
		function: function,
		stop:     make(chan struct{}),
	}

	c.Refresh()

	go c.run()

	return c
}

func (c *Cache[T]) Refresh() {
	value, err := c.function()

	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil {
		c.err = err
		return // Preserve the last successful value.
	}

	c.value = value
	c.err = nil
}

func (c *Cache[T]) run() {
	timer := time.NewTimer(10 * time.Second)

	defer timer.Stop()
	for {
		select {
		case <-c.stop:
			return

		case <-timer.C:

			c.Refresh()
			timer = time.NewTimer(10 * time.Second)

		}
	}
}

func (c *Cache[T]) Get() (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.value, c.err
}

func (c *Cache[T]) Close() {
	c.close.Do(func() {
		close(c.stop)
	})
}
