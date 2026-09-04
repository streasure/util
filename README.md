# util

 Go 通用工具库，提供泛型工具函数、组件生命周期管理、Prometheus 监控、etcd 服务治理等能力。所有功能模块支持插拔式接入，通过配置控制启用/禁用。

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
├── timeutil/         # 时间工具（偏移/格式化/天差值）
├── timeout_waitgroup/# 带超时的 WaitGroup
├── uuid/             # UUID V4 生成
├── id_allocator/     # ID 分配器
├── mathx/            # 2D 向量数学运算
├── component/        # 组件生命周期管理
├── version/          # 构建版本信息
├── gopool/           # goroutine 池（支持泛型返回值）
├── rand/             # 随机数工具（概率/权重/采样）
├── rand/wrand/       # 高性能加权随机选择（二分搜索）
├── gametime/         # 游戏时间工具（偏移/重置时间）
├── gametime/now/     # 游戏时间边界计算（日/周/月/年）
├── errors/           # 错误处理（堆栈/错误码/包装）
├── gevent/           # 反射事件分发器
├── sensitive/        # 敏感词过滤
├── prometheus/       # Prometheus 监控导出
├── etcd/             # etcd 服务注册/发现/配置中心
├── bucket/           # 令牌桶限流器
├── compressex/       # JSON+Gzip 压缩
├── config/           # 泛型配置加载器（YAML）
├── container/priority_queue/ # 优先队列（泛型）
├── backend/          # 通用 HTTP 客户端
├── config/           # 共用 YAML 配置文件
│   ├── etcd.yaml
│   ├── prometheus.yaml
│   └── grafana.yaml
└── examples/         # 配置示例
```

## 快速开始

```go
import "github.com/streasure/util/mathutil"
import "github.com/streasure/util/bitmask"
import "github.com/streasure/util/prometheus"
import "github.com/streasure/util/etcd"
import "github.com/streasure/util/bucket"
import "github.com/streasure/util/compressex"
import "github.com/streasure/util/config"
import "github.com/streasure/util/container/priority_queue"
import "github.com/streasure/util/rand/wrand"
import "github.com/streasure/util/backend"
import "github.com/streasure/util/gametime"
import "github.com/streasure/util/gametime/now"
import "github.com/streasure/util/timeutil"
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
| `SetOffset(d)` `GetOffset()` | 全局时间偏移 |
| `Now()` `Timestamp()` | 偏移感知的时间 |
| `Since(t)` `Until(t)` | 时间差计算 |
| `TimeStampToString` `TimeToString` | 时间格式化 |
| `DiffNatureDays(t1, t2)` | 自然天差值（高性能整除实现） |
| `DiffDays(end, start)` | 时间差天数 |
| `ZeroTimeOfDay(t)` | 截取到午夜 |
| `NormalizeTimeOfDay(t, hour)` | 归一化到每日重置时间 |
| `GetTomorrowStamp()` | 获取明天午夜 |
| `IsSameDay` `IsSameWeek` `IsSameMonth` `IsToday` | 时间周期判断 |
| `GetZeroTime` `GetTimeByHour` `GetDateKey` | 时间获取 |

```go
// 设置全局时间偏移（游戏服务器常用）
timeutil.SetOffset(5 * time.Hour)
now := timeutil.Now() // 返回当前时间 + 5小时

// 计算两个时间戳的自然天差
days := timeutil.DiffNatureDays(1609459200, 1609718400) // 3天

// 归一化到每日重置时间（如凌晨5点刷新）
resetTime := timeutil.NormalizeTimeOfDay(time.Now(), 5)
```

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
| `SetOffset(d)` `GetOffset()` | 全局时间偏移 |
| `Now()` `Timestamp()` | 偏移感知的时间/时间戳 |
| `Since(t)` `Until(t)` | 时间差计算 |
| `Unix(t, nsec)` | 时间戳转换 |

```go
// 设置全局时间偏移
gametime.SetOffset(-8 * time.Hour) // 时区补偿

now := gametime.Now()       // 当前游戏时间
ts := gametime.Timestamp()  // 当前游戏时间戳
d := gametime.Since(past)   // 距过去时间的差
```

### gametime/now/ - 游戏时间边界

| 函数 | 说明 |
|------|------|
| `BeginningOfDay()` | 当天开始（00:00:00） |
| `BeginningOfWeek()` | 当周开始（可配置周起始日） |
| `BeginningOfMonth()` | 当月开始 |
| `BeginningOfYear()` | 当年开始 |
| `EndOfDay()` | 当天结束（23:59:59） |
| `EndOfWeek()` | 当周结束 |
| `EndOfMonth()` | 当月结束 |
| `EndOfYear()` | 当年结束 |
| `SetWeekStartDay(day)` | 设置周起始日（默认周一） |
| `With(t)` `New(t)` | 为任意时间创建边界计算器 |

```go
// 获取本周开始时间
weekStart := now.BeginningOfWeek()

// 获取当月结束时间
monthEnd := now.EndOfMonth()

// 设置周日为起始日
now.SetWeekStartDay(time.Sunday)
weekStart := now.BeginningOfWeek()
```

### gametime/ref_time.go - 参考时间系统

| 函数 | 说明 |
|------|------|
| `NewRefTime(DailyTime{H,M,S})` | 创建参考时间（如每日5点刷新） |
| `NextNDayResetTime(t, days)` | 下N天重置时间 |
| `NextNWeeksResetTime(t, weeks)` | 下N周重置时间 |
| `NextNMonthsMonthdayResetTime(t, months, day)` | 下N月某日重置时间 |
| `IsSameDay` `IsSameWeek` `IsSameMonth` | 周期判断（支持多种时间类型） |
| `SubDay(a, b)` | 天数差 |

```go
// 每日凌晨5点重置
rt := gametime.NewRefTime(gametime.DailyTime{Hour: 5})
resetTime := rt.NextNDayResetTime(time.Now(), 1)

// 判断两个时间是否在同一周
sameWeek := gametime.IsSameWeek(time1, time2)
```

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
| `etcd.yaml` | `etcd.ComponentConfig` | etcd 连接、注册、发现、租约和动态配置 |
| `prometheus.yaml` | `prometheus.ExporterConfig` | Prometheus 指标导出 |
| `grafana.yaml` | `prometheus.GrafanaConfig` | Grafana 数据源 + Dashboard 导入 |

所有配置通过根级 `enabled` 和各功能级 `enabled` 字段控制启停，使用者只需加载 YAML 后创建一个 `etcd.Component`。

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

### etcd/ - etcd 服务治理

etcd 功能通过一个 `Component` 统一接入。注册、发现和动态配置可以单独启用，也可以同时启用。组件实现 `component.Component` 接口，直接加入现有 `component.Container` 即可。

| API | 说明 |
|------|------|
| `New(ComponentConfig)` | 创建单个 etcd Component |
| `Pick()` | 使用配置的负载均衡策略选择一个服务地址 |
| `Acquire()` | 选择地址并返回请求完成函数，供 P2C 统计并发负载 |
| `Services()` | 获取当前服务发现缓存 |
| `OnServiceChange` | 监听实例新增、地址更新和注销 |
| `OnConfigChange` | 监听原始动态配置字节变更 |
| `OnTypedConfigChange[T]` | 按 JSON/YAML/TOML 解析并校验后回调 |
| `BindConfig[T]` | 将当前配置快照绑定到指定类型 |
| `ConfigSnapshot()` | 获取当前有效配置快照 |
| `LastError()` | 获取后台 Watch/KeepAlive 最近一次错误 |

```go
etcdComponent := etcd.New(etcd.ComponentConfig{
    Enabled: true,
    Etcd: etcd.Config{
        Endpoints: []string{"http://etcd-1:2379", "http://etcd-2:2379"},
        DialTimeout: "5s",
        ServicePrefix: "/services",
        ConfigPrefix: "/config",
    },
    Registration: etcd.RegistrationConfig{
        Enabled: true,
        ServiceID: "my-svc",
        InstanceID: "instance-1",
        Address: "127.0.0.1:8080",
        LeaseTTL: "10s",
    },
    Discovery: etcd.DiscoveryConfig{
        Enabled: true,
        ServiceID: "upstream-svc",
        LoadBalance: etcd.LoadBalanceP2C,
    },
    Config: etcd.DynamicConfig{
        Enabled: true,
        Key: "/config/app.yaml",
        Format: "yaml",
    },
})

etcdComponent.OnServiceChange(func(event etcd.ServiceEvent) {
    fmt.Printf("service %s: %s\n", event.Type, event.Address)
})
container.Add(etcdComponent)
```

服务注册使用以下格式：

```text
{servicePrefix}/{serviceID}/{instanceID} = {address}
```

注册时会申请 Lease 并将 Key 绑定到 Lease，默认 TTL 为 `10s`。组件持续执行 `KeepAlive`，租约失效后会自动申请新 Lease 并重新写入同一个服务 Key。进程正常销毁时主动删除 Key；进程崩溃或无法续租时，由 etcd 在租约到期后自动删除 Key。

服务发现通过服务前缀读取实例，并使用 etcd `Watch` 实时同步新增、删除和地址更新。支持以下策略：

| 策略 | 配置值 | 说明 |
|------|--------|------|
| 轮询 | `round_robin` | 按稳定排序后的地址依次选择 |
| 随机 | `random` | 随机选择一个地址 |
| P2C | `p2c` | 随机抽取两个地址，选择当前并发负载较低者 |

```go
address, done, ok := etcdComponent.Acquire()
if ok {
    defer done()
    // 使用 address 发起请求
}
```

动态配置支持 JSON、YAML 和 TOML。配置更新只有在格式解析和业务校验通过后才会替换内存快照。

```go
type AppConfig struct {
    Name string `yaml:"name" json:"name" toml:"name"`
}

func (c AppConfig) Validate() error {
    if c.Name == "" { return fmt.Errorf("name is required") }
    return nil
}

etcd.OnTypedConfigChange(etcdComponent, func(config *AppConfig) error {
    return config.Validate()
})
```

安全连接配置支持用户名密码和双向 TLS：

| 字段 | 说明 |
|------|------|
| `username` / `password` | etcd 用户认证，可与 TLS 同时使用 |
| `certFile` | 客户端证书文件 |
| `certKeyFile` | 客户端私钥文件 |
| `caCertFile` | CA 根证书文件 |
| `tlsServerName` | TLS 服务端名称覆盖 |
| `endpoints` | 多个 etcd 节点地址，支持集群故障转移 |

TLS 配置要求 etcd 节点使用 `https` 地址；未配置 TLS 时使用普通 `http` 连接。当前 Go `1.22.5` 使用官方客户端 `go.etcd.io/etcd/client/v3 v3.5.18`。

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

# etcd Component（见 config/etcd.yaml）
enabled: true
etcd:
  endpoints:
    - "http://127.0.0.1:2379"
registration:
  enabled: false
  serviceID: "my-service"
  instanceID: "instance-1"
  address: "127.0.0.1:8080"
  leaseTTL: "10s"
discovery:
  enabled: false
  serviceID: "upstream-service"
  loadBalance: "round_robin"
config:
  enabled: false
  key: "/config/app.yaml"
  format: "yaml"
```

根级 `enabled: false` 时组件完全禁用，不创建 etcd 网络连接；`registration.enabled`、`discovery.enabled` 和 `config.enabled` 分别控制三项功能。

## 性能基准测试

运行压测：

```bash
go test -v -run='^$' -bench='.' -benchmem ./...
```

### 核心函数性能（i5-10400F）

| 包 | 函数 | 性能 | 内存分配 |
|----|------|------|---------|
| timeutil | `DiffNatureDays` | 0.25 ns/op | 0 B/op |
| timeutil | `IsSameDayUnix` | 0.25 ns/op | 0 B/op |
| timeutil | `Now` | 10.66 ns/op | 0 B/op |
| gametime | `Now` | 17.17 ns/op | 0 B/op |
| gametime | `GetOffset` | 0.31 ns/op | 0 B/op |
| gametime/now | `BeginningOfDay` | 73 ns/op | 0 B/op |
| gametime/now | `BeginningOfWeek` | 77 ns/op | 0 B/op |
| slice | `UniqueSlice` (小) | 84 ns/op | 64 B/op |
| slice | `StrToSlice` | 0.58 ns/op | 0 B/op |
| rand | `IntN` | 11.87 ns/op | 0 B/op |
| rand | `Float64` | 8.52 ns/op | 0 B/op |
| mathx | `V2.Add` | 2.53 ns/op | 0 B/op |
| mathx | `Normalize` | 11.33 ns/op | 0 B/op |
| mathutil | `Clamp` | 0.31 ns/op | 0 B/op |
| mathutil | `Abs` | 0.30 ns/op | 0 B/op |
| overflow | `CalcAddOverflow` | 0.32 ns/op | 0 B/op |
| bitmask | `SetBit` | 0.31 ns/op | 0 B/op |

## 测试

```bash
# 运行所有单元测试
go test ./... -v

# 运行压测
go test -v -run='^$' -bench='.' -benchmem ./...
```

## 特性

- **依赖可控**：Prometheus 使用标准库 HTTP；etcd 使用官方 v3 客户端
- **泛型优先**：`Clamp[T]` `Max[T]` `Abs[T]` `SetBit[T]` `WeightRandom[T]` `PriorityQueue` 等
- **组件化**：所有服务实现 `Component` 接口，统一生命周期管理
- **插拔式**：配置 `enabled: false` 即可禁用，不引入任何开销
- **配置共用**：`config/` 目录提供 etcd/Prometheus/Grafana 标准 YAML，项目直接引用
- **高性能**：核心函数通过整除运算、内存池、预排序缓存等优化
- **游戏时间支持**：全局时间偏移、每日重置时间、周期判断等游戏服务器常用功能
