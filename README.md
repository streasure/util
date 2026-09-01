# util

Go 通用工具库，提供泛型工具函数、组件生命周期管理、Prometheus 监控、Nacos 服务治理等能力。所有功能模块支持插拔式接入，通过配置控制启用/禁用。

## 目录

```
github.com/streasure/util
├── bitmask/          # 泛型位操作
├── httputil/         # HTTP 工具（CORS 中间件）
├── mathutil/         # 泛型数学运算
├── msg/              # 同步/异步消息
├── netutil/          # 网络工具（URL 解析/IP 获取）
├── overflow/         # 溢出安全加减
├── panicutil/        # 线程安全 panic 捕获
├── rangeable/        # 区间重叠检测
├── slice/            # 切片操作
├── sys/              # 系统工具（goroutine ID）
├── timeutil/         # 时间工具
├── timeout_waitgroup/# 带超时的 WaitGroup
├── uuid/             # UUID V4 生成
├── id_allocator/     # ID 分配器
├── mathx/            # 2D 向量数学运算
├── component/        # 组件生命周期管理
├── version/          # 构建版本信息
├── gopool/           # goroutine 池（支持泛型返回值）
├── rand/             # 随机数工具（概率/权重/采样）
├── rand/wrand/       # 高性能加权随机选择（二分搜索）
├── gametime/         # 游戏时间工具（重置时间/偏移）
├── errors/           # 错误处理（堆栈/错误码/包装）
├── gevent/           # 反射事件分发器
├── sensitive/        # 敏感词过滤
├── prometheus/       # Prometheus 监控导出
├── nacos/            # Nacos 服务注册/发现/配置中心
├── bucket/           # 令牌桶限流器
├── compressex/       # JSON+Gzip 压缩
├── config/           # 泛型配置加载器（YAML）
├── container/priority_queue/ # 优先队列（泛型）
├── backend/          # 通用 HTTP 客户端
├── config/           # 共用 YAML 配置文件
│   ├── nacos.yaml
│   ├── prometheus.yaml
│   └── grafana.yaml
└── examples/         # 配置示例
```

## 快速开始

```go
import "github.com/streasure/util/mathutil"
import "github.com/streasure/util/bitmask"
import "github.com/streasure/util/prometheus"
import "github.com/streasure/util/nacos"
import "github.com/streasure/util/bucket"
import "github.com/streasure/util/compressex"
import "github.com/streasure/util/config"
import "github.com/streasure/util/container/priority_queue"
import "github.com/streasure/util/rand/wrand"
import "github.com/streasure/util/backend"
```

## 包说明

### bitmask/ - 位操作

| 函数 | 说明 |
|------|------|
| `SetBit[T]` `ResetBit[T]` `HasBit[T]` | 泛型位操作 |
| `SetBitSlice` `ResetBitSlice` `HasBitSlice` | 字节数组位操作 |

### httputil/ - HTTP 工具

| 函数 | 说明 |
|------|------|
| `CorsHandlerFunc` `CorsHandler` | CORS 跨域中间件 |

### mathutil/ - 数学运算

| 函数 | 说明 |
|------|------|
| `Clamp[T]` `Max[T]` `Min[T]` `Abs[T]` `Operator[T]` | 泛型数学运算 |
| `RandInt` `RandInt32` `RandInt64` | [min, max] 范围随机 |

### msg/ - 消息处理

| 类型 | 说明 |
|------|------|
| `SyncMessage` | 同步消息 |
| `AsyncMessage` | 异步消息 |

### netutil/ - 网络工具

| 函数 | 说明 |
|------|------|
| `DecodeUrlValues` | URL 参数解析到 struct |
| `GetHttpIP` `LocalIP` | 获取客户端/本机 IP |

### overflow/ - 溢出处理

| 函数 | 说明 |
|------|------|
| `CalcAddOverflow` `CalcMinusOverflow` | 溢出安全加减 |

### panicutil/ - panic 捕获

| 类型 | 说明 |
|------|------|
| `PanicCatcher` | 线程安全 panic 捕获 |

### rangeable/ - 区间操作

| 函数 | 说明 |
|------|------|
| `CheckRangeIntersect[T]` | 区间重叠检测 |

### slice/ - 切片操作

| 函数 | 说明 |
|------|------|
| `IsNil` `StrToSlice` `UniqueSlice[T]` `EqualSlice[T]` | 切片操作 |

### sys/ - 系统工具

| 函数 | 说明 |
|------|------|
| `GoRoutineId` | 当前 goroutine ID |

### timeutil/ - 时间工具

| 函数 | 说明 |
|------|------|
| `TimeStampToString` `TimeToString` | 时间格式化 |
| `DiffNatureDays` | 自然天差值 |
| `IsSameDay` `IsSameWeek` `IsSameMonth` `IsToday` | 时间周期判断 |
| `GetZeroTime` `GetTimeByHour` `GetDateKey` | 时间获取 |

### timeout_waitgroup/ - 超时等待组

| 函数 | 说明 |
|------|------|
| `NewTimeoutWaitGroup(n)` | 带超时的 WaitGroup |

### uuid/ - UUID 生成

| 函数 | 说明 |
|------|------|
| `NewUUID` `NewUUIDBytes` | UUID V4 生成 |

### id_allocator/ - ID 分配器

| 类型 | 说明 |
|------|------|
| `Uint32IdAllocator` | ID 分配器 |
| `GenerateSessionId` | 会话 ID 生成 |

### mathx/ - 2D 向量数学

| 函数 | 说明 |
|------|------|
| `NewV2(x, y)` | 创建向量 |
| `V2.Len` `LenSqrt` `Normalize` | 长度/归一化 |
| `V2.Add` `Sub` `Mul` `Div` `Dot` | 向量运算 |
| `IsFloatSame` | 浮点数近似比较 |
| `Clamp` `IsClamped` | 浮点数域限制 |

### component/ - 组件生命周期

| 类型 | 说明 |
|------|------|
| `Component` 接口 | `Name` `Order` `Init` `Start` `Destroy` |
| `BaseComponent` | 可嵌入的基础实现 |
| `Container` | 容器：按 Order 排序初始化/启动，信号触发逆序销毁 |

```go
c := component.NewContainer()
c.Add(myComponent)
c.Serve() // 阻塞，等待 SIGINT/SIGTERM 后逆序 Destroy
```

### gopool/ - goroutine 池

| 函数 | 说明 |
|------|------|
| `NewPool(count)` | 无返回值任务池 |
| `NewTypedPool[T](count)` | 带返回值的泛型任务池 |

```go
pool := gopool.NewTypedPool[int](4)
ch := pool.Add(func() int { return 42 })
result := <-ch
pool.Stop()
```

### rand/ - 随机数工具

| 函数 | 说明 |
|------|------|
| `InTenThousandsProbability(rate)` | 万分比概率 |
| `InRandomProbability(rate, total)` | 通用概率 |
| `RangeInt` `RangeInt32` `RangeInt64` | 范围随机 |
| `RangeInts(min, max, n)` | n 个不重复随机数 |
| `SliceOne` `SliceN` | 切片随机采样 |
| `RandWeightSlice` | 权重随机索引 |
| `WeightRandom[T]` | 泛型权重选择器 |

### rand/wrand/ - 高性能加权随机

| 函数 | 说明 |
|------|------|
| `NewRandChooser[T, W](choices...)` | 创建加权选择器 |
| `chooser.Pick()` | 随机选择一个 |
| `chooser.PickSource(rs)` | 指定随机源选择 |
| `chooser.PickN(n)` | 随机选择 n 个不重复 |
| `NewRandChoices[T]([][]T)` | 二维数组创建 choices |

基于二分搜索的预排序缓存，大规模数据下性能显著优于线性扫描。

```go
chooser, _ := wrand.NewRandChooser(
    wrand.NewRandChoice("rare", 1),
    wrand.NewRandChoice("common", 10),
)
item := chooser.Pick() // "common" 被选中的概率是 "rare" 的 10 倍
```

### gametime/ - 游戏时间

| 函数 | 说明 |
|------|------|
| `SetOffset(d)` | 全局时间偏移 |
| `Now` `Since` `Until` | 偏移感知的时间 |
| `NewRefTime(DailyTime)` | 参考时间（如每日 5 点刷新） |
| `RefTime.IsSameDay` `IsSameWeek` `IsSameMonth` | 周期判断 |
| `RefTime.NextNDayResetTime` `NextNWeeksResetTime` | 重置时间计算 |
| `RefTime.SubDay` | 天数差 |

### errors/ - 错误处理

| 函数 | 说明 |
|------|------|
| `New` `NewWithStack` | 创建错误 |
| `WithMessage` `WithMessagef` | 添加消息 |
| `WithStack` | 添加堆栈 |
| `Wrap` `Wrapf` | 包装错误 |
| `Append` `WithOverride` | 合并/覆盖 |
| `Code` | 提取错误码 |
| `Cause` | 根因 |
| `Is` `As` `Unwrap` | 标准操作 |
| `NewECode(code)` | 带码错误 |

### gevent/ - 事件分发器

| 函数 | 说明 |
|------|------|
| `NewDispatcher(opts...)` | 创建分发器 |
| `Register(event, handler)` | 注册事件 |
| `RegisterService(receiver)` | 注册 struct 方法 |
| `Dispatch(event, args...)` | 触发（多 handler） |
| `Call(event, args...)` | 调用（单 handler，返回值） |

### sensitive/ - 敏感词过滤

| 函数 | 说明 |
|------|------|
| `InitWords(words)` | 初始化词库 |
| `CensorIsPass(text)` | 检查是否通过 |
| `CensorAndReplace(text)` | 替换敏感词为 `*` |

### bucket/ - 令牌桶限流器

| 函数 | 说明 |
|------|------|
| `NewBucket(interval, capacity, quantum)` | 创建令牌桶 |
| `TakeAvailable(now, count)` | 加锁获取令牌 |
| `TakeAvailableNoLock(now, count)` | 无锁获取令牌 |

```go
b := bucket.NewBucket(100*time.Millisecond, 10, 2)
taken := b.TakeAvailable(time.Now(), 5) // 获取 5 个令牌
```

### compressex/ - JSON+Gzip 压缩

| 函数 | 说明 |
|------|------|
| `ProtoMarshal(v)` | JSON 序列化 + Gzip 压缩 |
| `ProtoUnmarshal(data, v)` | Gzip 解压 + JSON 反序列化 |

```go
compressed, _ := compressex/proto.Marshal(myData)
var result MyData
compressex.ProtoUnmarshal(compressed, &result)
```

### config/ - 泛型配置加载器

| 函数 | 说明 |
|------|------|
| `Load[T](paths...)` | 加载 YAML 配置到泛型结构体 |

支持多文件合并，自动校验必填字段。

```go
type MyConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

cfg, err := config.Load[MyConfig]("config.yaml", "override.yaml")
```

### container/priority_queue/ - 优先队列

| 函数 | 说明 |
|------|------|
| `New(opts...)` | 创建优先队列（默认最大堆） |
| `WithMin(true)` | 切换为最小堆 |
| `Push(x, priority)` | 入队 |
| `Pop()` | 出队（返回值和优先级） |
| `Peek()` | 查看队首 |
| `Len()` | 长度 |
| `Clear()` | 清空 |

```go
pq := priority_queue.New()
pq.Push("task1", 10)
pq.Push("task2", 5)
v, _ := pq.Pop() // 返回 "task1"（优先级最高）
```

### backend/ - 通用 HTTP 客户端

| 函数 | 说明 |
|------|------|
| `HttpCommonGet[T](ctx, url)` | GET 请求，解析 CommonAck 响应 |
| `HttpCommonPost[T](ctx, url, body)` | POST 请求，解析 CommonAck 响应 |
| `ActivationCode(ctx, url)` | 激活码请求 |
| `CommonAck[T]` | 通用响应结构（Code/Msg/Data） |

```go
result, code, err := backend.HttpCommonGet[ServerInfo](ctx, "http://api/servers")
if code == backend.CodeSuccess {
    fmt.Println(result.Name)
}
```

### prometheus/ - Prometheus 监控

| 类型 | 说明 |
|------|------|
| `NewExporter(cfg, provider)` | 创建导出器组件 |
| `StatsProvider` 接口 | 实现 `Stats() Stats` 提供自定义指标 |
| `Stats` | 指标结构（连接/消息/延迟/安全/系统/自定义） |
| `ExporterConfig` | 导出器配置（enabled/addr/path/prefix） |
| `GrafanaConfig` | Grafana 数据源/Dashboard 配置 |
| `RenderPrometheusText(s)` | 渲染 Prometheus 文本格式 |
| `RenderPrometheusTextWithPrefix(s, prefix)` | 带前缀渲染 |

```go
exp := prometheus.NewExporter(prometheus.ExporterConfig{
    Enabled: true, Addr: ":9100", Path: "/metrics", Prefix: "app",
}, myProvider)
container.Add(exp)
```

### config/ - 共用 YAML 配置

| 文件 | 配置结构 | 用途 |
|------|---------|------|
| `nacos.yaml` | `nacos.RegistryConfig` + `nacos.DiscoveryConfig` + `nacos.ConfigCenterConfig` | Nacos 连接/注册/发现/配置中心 |
| `prometheus.yaml` | `prometheus.ExporterConfig` | Prometheus 指标导出 |
| `grafana.yaml` | `prometheus.GrafanaConfig` | Grafana 数据源 + Dashboard 导入 |

所有配置通过 `enabled` 字段控制启停，项目引用时直接加载对应 YAML 即可。

支持的内置指标：

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `app_connections_total` | counter | 总连接数 |
| `app_connections_active` | gauge | 活跃连接数 |
| `app_messages_received_total` | counter | 接收消息数 |
| `app_messages_forwarded_total` | counter | 转发消息数 |
| `app_messages_dropped_*_total` | counter | 丢弃消息（多种原因） |
| `app_latency_p50_us` | gauge | P50 延迟 |
| `app_latency_p95_us` | gauge | P95 延迟 |
| `app_latency_p99_us` | gauge | P99 延迟 |
| `app_goroutines` | gauge | goroutine 数 |
| `app_memory_alloc_bytes` | gauge | 内存分配 |
| `app_gc_count` | counter | GC 次数 |
| `app_custom_*` | gauge | 自定义指标 |

### nacos/ - Nacos 服务治理

| 类型 | 说明 |
|------|------|
| `NewRegistry(cfg)` | 服务注册（心跳保活） |
| `NewDiscovery(cfg)` | 服务发现（轮询+变更回调） |
| `NewConfigCenter(cfg)` | 配置中心（轮询监听） |
| `Config` | 连接配置（支持 v1/v3 API） |

```go
// 注册
reg := nacos.NewRegistry(nacos.RegistryConfig{
    Enabled: true,
    Nacos:   nacos.Config{Endpoint: "http://nacos:8080"},
    Service: nacos.NamingConfig{ServiceName: "my-svc", Addr: ":8080"},
})
container.Add(reg)

// 发现
disc := nacos.NewDiscovery(nacos.DiscoveryConfig{
    Enabled: true,
    Nacos:   nacos.Config{Endpoint: "http://nacos:8080"},
    Service: nacos.NamingConfig{ServiceName: "upstream-svc"},
})
disc.OnServiceChange(func(e nacos.ServiceEvent) {
    fmt.Printf("service %s: %s\n", e.Type, e.Service.Address)
})
container.Add(disc)

// 配置中心
cc := nacos.NewConfigCenter(nacos.ConfigCenterConfig{
    Enabled: true,
    Nacos:   nacos.Config{Endpoint: "http://nacos:8080", DataID: "app.yaml"},
})
cc.OnConfigChange(func(data []byte) {
    fmt.Printf("config changed: %d bytes\n", len(data))
})
container.Add(cc)
```

## 配置示例

完整配置见 `examples/config.yaml` 和 `config/` 目录：

```yaml
# Prometheus 监控（见 config/prometheus.yaml）
monitor:
  prometheus:
    enabled: true
    addr: ":9100"
    path: "/metrics"
    prefix: "app"

# Grafana（见 config/grafana.yaml）
grafana:
  datasources:
    - name: Prometheus
      type: prometheus
      url: "http://localhost:9090"
  dashboards:
    - name: "Game Gateway Overview"
      file: "dashboards/gateway-overview.json"
      folder: "sgate"

# Nacos 服务注册（见 config/nacos.yaml）
nacosRegistry:
  enabled: false
  nacos:
    endpoint: "http://127.0.0.1:8080"
    apiVersion: "v3"
  service:
    serviceName: "my-service"
    addr: "127.0.0.1:8080"

# Nacos 服务发现
nacosDiscovery:
  enabled: false
  nacos:
    endpoint: "http://127.0.0.1:8080"
  service:
    serviceName: "upstream-service"
  scanInterval: "10s"

# Nacos 配置中心
nacosConfigCenter:
  enabled: false
  nacos:
    endpoint: "http://127.0.0.1:8080"
    dataID: "app.yaml"
    pollInterval: "5s"
```

所有组件通过 `enabled` 控制是否接入，`false` 时不启动任何网络连接。

## 测试

```bash
go test ./... -v
```

## 特性

- **零外部依赖**：monitor/nacos 全部使用标准库 `net/http` 实现
- **泛型优先**：`Clamp[T]` `Max[T]` `Abs[T]` `SetBit[T]` `WeightRandom[T]` `PriorityQueue` 等
- **组件化**：所有服务实现 `Component` 接口，统一生命周期管理
- **插拔式**：配置 `enabled: false` 即可禁用，不引入任何开销
- **配置共用**：`config/` 目录提供 Nacos/Prometheus/Grafana 标准 YAML，项目直接引用
