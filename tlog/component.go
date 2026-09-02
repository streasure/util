package tlog

import (
	"github.com/streasure/util/component"
)

// LogComponent treasure-slog 日志组件，实现 component.Component 接口
type LogComponent struct {
}

// 确保 LogComponent 实现 component.Component 接口
var _ component.Component = (*LogComponent)(nil)

// NewLogComponent 创建日志组件
func NewLogComponent() *LogComponent {
	return &LogComponent{}
}

// Name 组件名称
func (l *LogComponent) Name() string {
	return "TLog"
}

// Order 组件顺序，按照数字从小到达排序，可以为负数，默认为 0
// 日志组件需要最早初始化，返回最小整数
func (l *LogComponent) Order() int {
	return 0
}

// Init 组件初始化
func (l *LogComponent) Init() error {
	return nil
}

// Start 组件启动（可选，tlog 无需额外启动逻辑）
func (l *LogComponent) Start() error {
	return nil
}

// Destroy 组件销毁，同步并关闭 logger
func (l *LogComponent) Destroy() {
	Sync()
}
