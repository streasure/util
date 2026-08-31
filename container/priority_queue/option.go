package priority_queue

type Option func(*priorityQueue)

func WithMin(v bool) Option {
	return func(p *priorityQueue) {
		p.items.isMin = v
	}
}
