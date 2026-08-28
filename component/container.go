package component

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"syscall"
)

type Container struct {
	components []Component
	sigch      chan os.Signal
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Add(component Component) {
	c.components = append(c.components, component)
}

func (c *Container) Serve() {
	sort.SliceStable(c.components, func(i, j int) bool {
		return c.components[i].Order() < c.components[j].Order()
	})

	fmt.Println("[component] all component init-ing.")

	for _, comp := range c.components {
		if err := comp.Init(); err != nil {
			fmt.Printf("[component] component[%s] init error: %s\n", comp.Name(), err.Error())
			panic(err)
		}
	}

	fmt.Println("[component] all component starting.")

	for index, comp := range c.components {
		if err := comp.Start(); err != nil {
			fmt.Printf("[component] component[%s] start error: %s\n", comp.Name(), err.Error())
			for i := index; i > 0; i-- {
				c.components[i-1].Destroy()
			}
			panic(err)
		}
	}

	fmt.Println("[component] all component start success.")

	c.sigch = make(chan os.Signal, 1)

	switch runtime.GOOS {
	case "windows":
		signal.Notify(c.sigch, syscall.SIGINT, syscall.SIGTERM)
	default:
		signal.Notify(c.sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	}

	<-c.sigch

	signal.Stop(c.sigch)

	fmt.Println("[component] all component destroying.")

	for i := len(c.components); i > 0; i-- {
		c.components[i-1].Destroy()
	}

	fmt.Println("[component] all components destroyed.")
}

func (c *Container) Destroy() {
	if c.sigch != nil {
		c.sigch <- syscall.SIGTERM
	}
}
