package ds

import "fmt"

type Middleware[T any] func(io T, next func(error))

type Chain[T any] struct {
	stack []Middleware[T]
}

func (p *Chain[T]) Use(middleware Middleware[T]) {
	p.stack = append(p.stack, middleware)
}

func (p *Chain[T]) Execute(io T) {
	handler := func(_io T) {
		p.handle(_io, func(err error) {
			if err != nil {
				panic("Handler Error")
			}
		})
	}
	handler(io)
}

func (p *Chain[T]) handle(io T, callback func(error)) {
	var idx int

	var next func(err error)
	next = func(err error) {
		if err != nil {
			callback(err)
			return
		}
		if idx >= len(p.stack) {
			callback(nil)
			return
		}
		layer := p.stack[idx]
		idx++

		func() {
			defer func() {
				if r := recover(); r != nil {
					next(fmt.Errorf("%v", r))
				}
			}()

			layer(io, next)
		}()
	}

	next(nil)
}
