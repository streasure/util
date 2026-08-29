# util

Go 通用工具库。

## 包结构

```
github.com/streasure/util
├── util/           # 通用工具函数
├── mathx/          # 2D 向量数学运算
├── component/      # 组件生命周期管理
├── version/        # 构建版本信息
├── gopool/         # goroutine 池
├── rand/           # 随机数工具
├── gametime/       # 游戏时间工具
├── errors/         # 错误处理（带堆栈）
├── gevent/         # 事件分发器
├── sensitive/      # 敏感词过滤
├── monitor/        # Prometheus 监控导出
└── nacos/          # Nacos 服务注册/发现/配置中心
```

## 包功能说明

### util/ - 通用工具函数

| 文件 | 函数 | 说明 |
|------|------|------|
| math.go | `Clamp[T]` | 值域限制（泛型，支持所有可比较类型） |
| math.go | `Max[T]/Min[T]` | 最大值/最小值（泛型） |
| math.go | `Abs[T]` | 绝对值（泛型，支持所有有符号整数） |
| math.go | `Operator[T]` | 字符串运算符比较（泛型） |
| math.go | `RandInt/RandInt32/RandInt64` | [min, max] 范围随机数 |
| time.go | `TimeStampToString/TimeToString` | 时间格式化 |
| time.go | `DiffNatureDays` | 计算自然天差值 |
| time.go | `IsSameDay/IsSameWeek/IsSameMonth` | 时间周期判断 |
| time.go | `GetZeroTime/GetTimeByHour/GetDateKey` | 时间获取 |
| net.go | `DecodeUrlValues` | URL 参数解析到结构体 |
| net.go | `GetHttpIP/LocalIP` | 获取客户端/本机 IP |
| http.go | `CorsHandlerFunc/CorsHandler` | CORS 跨域中间件 |
| panic.go | `PanicCatcher` | 线程安全的 panic 捕获器 |
| slice.go | `IsNil` | 反射 nil 值检查 |
| slice.go | `StrToSlice` | 零拷贝 string → []byte |
| slice.go | `UniqueSlice` | 切片去重（泛型） |
| slice.go | `EqualSlice` | 无序切片相等判断（泛型） |
| bitmask.go | `SetBitSlice/ResetBitSlice/HasBitSlice` | 字节数组位操作 |
| bitmask.go | `SetBit[T]/ResetBit[T]/HasBit[T]` | 整数位操作（泛型） |
| overflow.go | `CalcAddOverflow/CalcMinusOverflow` | 溢出安全的加减法 |
| msg.go | `SyncMessage/AsyncMessage` | 同步/异步消息模式 |
| uuid.go | `NewUUID/NewUUIDBytes` | UUID V4 生成 |
| sys.go | `GoRoutineId` | 获取当前 goroutine ID |
| id_allocator.go | `Uint32IdAllocator` | 无锁原子 ID 分配器 |
| id_allocator.go | `GenerateSessionId` | 生成 32 位随机会话 ID |

### mathx/ - 2D 向量数学

| 函数 | 说明 |
|------|------|
| `NewV2(x, y)` | 创建 2D 向量 |
| `V2.Len/LenSqrt/Normalize` | 向量长度/归一化 |
| `V2.Add/Sub/Mul/Div/Dot` | 向量运算 |
| `Add/Sub/Mul/Div/Dot` | 包级别向量运算 |
| `IsFloatSame` | 浮点数近似相等比较 |
| `Clamp/IsClamped` | 浮点数域限制 |

### component/ - 组件生命周期

| 类型/函数 | 说明 |
|-----------|------|
| `Component` 接口 | Name/Order/Init/Start/Destroy |
| `BaseComponent` | 基础组件实现 |
| `Container` | 组件容器，按顺序初始化/启动，信号触发逆序销毁 |

### version/ - 版本信息

| 变量/函数 | 说明 |
|-----------|------|
| `Version/GoVersion/Built/GitCommit/OSArch` | 构建时注入的版本变量 |
| `VersionDetail(appName)` | 格式化版本详情字符串 |

### gopool/ - goroutine 池

| 类型/函数 | 说明 |
|-----------|------|
| `NewPool(count)` | 创建固定 worker 数量的池 |
| `Pool.Add(item)` | 提交任务 |
| `Pool.Stop()` | 停止并等待所有 worker 完成 |
| `NewTypedPool[T](count)` | 创建带返回值的泛型池 |
| `TypedPool.Add(fn)` | 提交任务，返回 `<-chan T` |

### rand/ - 随机数工具

| 函数 | 说明 |
|------|------|
| `InTenThousandsProbability(rate)` | 万分比概率判断 |
| `InRandomProbability(rate, total)` | 通用概率判断 |
| `RangeInt/RangeInt32/RangeInt64` | [min, max] 范围随机 |
| `RangeInts(min, max, n)` | 生成 n 个不重复随机数 |
| `SliceOne/SliceN` | 随机选取切片元素 |
| `RandWeightSlice` | 权重随机选择索引 |
| `WeightRandom[T]` | 泛型权重随机选择器 |

### gametime/ - 游戏时间工具

| 函数 | 说明 |
|------|------|
| `SetOffset(d)` | 设置全局时间偏移量 |
| `Now/Since/Until` | 偏移感知的时间获取 |
| `NewRefTime(DailyTime)` | 创建参考时间（如每日 5 点刷新） |
| `RefTime.IsSameDay/IsSameWeek/IsSameMonth` | 基于参考时间的周期判断 |
| `RefTime.NextNDayResetTime` | 计算 N 天后的重置时间 |
| `RefTime.SubDay` | 计算天数差 |

### errors/ - 错误处理

| 函数 | 说明 |
|------|------|
| `New/NewWithStack` | 创建错误（可带堆栈） |
| `WithMessage/WithMessagef` | 添加错误消息 |
| `WithStack` | 添加堆栈信息 |
| `Wrap/Wrapf` | 包装错误 |
| `Append/WithOverride` | 合并/覆盖错误 |
| `Code` | 提取错误码 |
| `Cause` | 获取根因错误 |
| `Is/As/Unwrap` | 标准错误操作 |
| `NewECode(code)` | 创建带错误码的错误 |

### gevent/ - 事件分发器

| 类型/函数 | 说明 |
|-----------|------|
| `NewDispatcher(opts...)` | 创建事件分发器 |
| `Register(event, handler)` | 注册事件处理函数 |
| `RegisterService(receiver)` | 注册结构体所有方法为事件 |
| `Dispatch(event, args...)` | 触发事件（多 handler） |
| `Call(event, args...)` | 调用事件（单 handler，返回结果） |

### sensitive/ - 敏感词过滤

| 函数 | 说明 |
|------|------|
| `InitWords(words)` | 初始化敏感词库 |
| `CensorIsPass(text)` | 检查文本是否通过审查 |
| `CensorAndReplace(text)` | 审查并替换敏感词为 `*` |

### monitor/ - Prometheus 监控导出

| 类型/函数 | 说明 |
|-----------|------|
| `NewExporter(cfg, provider)` | 创建 Prometheus 导出器组件 |
| `StatsProvider` 接口 | 自定义指标提供者 |
| `Stats` | 指标数据结构（连接/消息/延迟/安全/系统） |
| `RenderPrometheusText(s)` | 渲染 Prometheus 文本格式 |

所有组件通过 `Enabled` 配置控制是否接入：

```yaml
monitor:
  prometheus:
    enabled: true   # false 则不启动
    addr: ":9100"
    path: "/metrics"
```

### nacos/ - Nacos 服务注册/发现/配置中心

| 类型/函数 | 说明 |
|-----------|------|
| `NewRegistry(cfg)` | 服务注册组件（心跳保活） |
| `NewDiscovery(cfg)` | 服务发现组件（轮询拉取+变更回调） |
| `NewConfigCenter(cfg)` | 配置中心组件（轮询监听配置变更） |
| `Config` | Nacos 连接配置（支持 v1/v3 API） |
| `ServiceInfo` | 服务实例信息 |

所有组件通过 `Enabled` 配置控制是否接入：

```yaml
nacosRegistry:
  enabled: false    # false 则不注册
  nacos:
    endpoint: "http://127.0.0.1:8080"
    apiVersion: "v3"
  service:
    serviceName: "my-service"
    addr: "127.0.0.1:8080"

nacosDiscovery:
  enabled: false    # false 则不发现
  nacos:
    endpoint: "http://127.0.0.1:8080"
  service:
    serviceName: "upstream-service"
  scanInterval: "10s"

nacosConfigCenter:
  enabled: false    # false 则不拉取配置
  nacos:
    endpoint: "http://127.0.0.1:8080"
    dataID: "app.yaml"
```

完整配置示例见 `examples/config.yaml`。

## 使用

```go
import "github.com/streasure/util/util"
import "github.com/streasure/util/rand"
import "github.com/streasure/util/errors"
import "github.com/streasure/util/monitor"
import "github.com/streasure/util/nacos"
```

## 测试

```bash
go test ./...
```
