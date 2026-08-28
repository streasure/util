package component

type Component interface {
	Name() string
	Order() int
	Init() error
	Start() error
	Destroy()
}

type BaseComponent struct {
	order int
}

func (c *BaseComponent) Name() string {
	return "base"
}

func (c *BaseComponent) Order() int {
	return c.order
}

func (c *BaseComponent) Init() error {
	return nil
}

func (c *BaseComponent) Start() error {
	return nil
}

func (c *BaseComponent) Destroy() {}

func (c *BaseComponent) SetOrder(order int) {
	c.order = order
}
