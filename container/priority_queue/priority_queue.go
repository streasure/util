package priority_queue

type priorityItem struct {
	Value    any
	Priority int64
	Index    int
}

type priorityItems struct {
	items []*priorityItem
	isMin bool
}

func (p *priorityItems) Len() int {
	return len(p.items)
}

func (p *priorityItems) Less(i, j int) bool {
	if p.isMin {
		return p.items[i].Priority < p.items[j].Priority
	}
	return p.items[i].Priority > p.items[j].Priority
}

func (p *priorityItems) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
	p.items[i].Index = i
	p.items[j].Index = j
}

func (p *priorityItems) Push(item *priorityItem) {
	item.Index = p.Len()
	p.items = append(p.items, item)
	p.up()
}

func (p *priorityItems) up() {
	j := p.Len() - 1
	for {
		i := (j - 1) / 2
		if i == j || !p.Less(j, i) {
			break
		}
		p.Swap(i, j)
		j = i
	}
}

func (p *priorityItems) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 < n && p.Less(j2, j1) {
			j = j2
		}
		if !p.Less(j, i) {
			break
		}
		p.Swap(i, j)
		i = j
	}
	return i > i0
}

func (p *priorityItems) Pop() *priorityItem {
	n := p.Len() - 1
	p.Swap(0, n)
	p.down(0, n)

	old := p.items
	n = len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	p.items = old[0 : n-1]

	return item
}

func (p *priorityItems) Peek() *priorityItem {
	return p.items[0]
}

type priorityQueue struct {
	items *priorityItems
}

func (p *priorityQueue) Push(x any, priority int64) {
	item := &priorityItem{
		Value:    x,
		Priority: priority,
		Index:    -1,
	}
	p.items.Push(item)
}

func (p *priorityQueue) Pop() (any, int64) {
	if p.items.Len() == 0 {
		return nil, 0
	}

	v := p.items.Pop()
	return v.Value, v.Priority
}

func (p *priorityQueue) Peek() (any, int64) {
	if p.items.Len() == 0 {
		return nil, 0
	}

	v := p.items.Peek()
	return v.Value, v.Priority
}

func (p *priorityQueue) Len() int {
	return p.items.Len()
}

func (p *priorityQueue) Clear() {
	if p.items.Len() == 0 {
		return
	}
	p.items.items = nil
}

type PriorityQueue interface {
	Push(x any, priority int64)
	Pop() (x any, priority int64)
	Peek() (x any, priority int64)
	Len() int
	Clear()
}

func New(opt ...Option) PriorityQueue {
	ret := &priorityQueue{
		items: &priorityItems{},
	}

	for _, o := range opt {
		o(ret)
	}
	return ret
}
