# Go Redis MQ 项目

这是一个使用 Go 语言和 Redis Stream 实现的消息队列项目。

## 项目概览

本项目旨在提供一个简单可靠的消息队列解决方案，利用 Redis Stream 的特性进行消息的生产、消费、以及失败重试处理。同时提供一个简单的仪表盘来监控队列状态。

## 如何运行

确保你已经安装了 Go 和 Redis。

1.  **初始化 Redis 客户端**
    所有组件在启动时都会调用 `redisclient.InitRedis()` 来初始化 Redis 连接。

2.  **启动生产者 (Producer)**
    生产者用于向 Redis Stream (`task_stream`) 发送任务消息。
    ```bash
    go run cmd/producer/main.go "your-task-payload"
    ```
    如果不提供参数，将使用默认的 `"default-task"` 作为 payload。

3.  **启动消费者 (Consumer)**
    消费者用于从 Redis Stream (`task_stream`) 读取并处理任务消息。
    ```bash
    go run cmd/consumer/main.go
    ```

4.  **启动仪表盘 (Dashboard)**
    仪表盘提供一个 Web 界面来查看队列的统计信息。
    ```bash
    go run cmd/dashboard/main.go
    ```
    启动后，可以通过浏览器访问 `http://localhost:8080` 查看仪表盘。
    API 端点 `/api/stats` 返回 JSON 格式的统计数据。

5.  **启动重试处理器 (Retry Worker)**
    重试处理器用于处理之前失败的任务。
    ```bash
    go run cmd/retry_worker/main.go
    ```

### 多节点协同
水平扩展：新增更多 Producer/Consumer，只要指向同一集群即可。
高可用：Redis 集群／哨兵自动故障转移，业务感知最小化。

### 环境变量配置
    ```
    export REDIS_MODE=sentinel
    export REDIS_ADDRS="192.168.1.10:26379,192.168.1.11:26379,192.168.1.12:26379"
    export REDIS_MASTER_NAME=mymaster
    export REDIS_PASSWORD=yourpass  # 如果你设置了密码
    ```
然后在任意节点同时启动生产者／消费者／重试／Dashboard，都会连到同一套 Redis 后端，实现分布式生产和消费。

### Docker运行

1. 启动 Redis 主从与 Sentinel 集群：
docker-compose up -d

验证：
    ```
    # 查看 sentinel 是否识别了主节点
    docker exec -it sentinel1 redis-cli -p 26379
    > SENTINEL get-master-addr-by-name mymaster
    ```
输出类似：
    ```
    1) "172.18.0.2"
    2) "6379"
    ```




## 项目结构

```
.
├── cmd
│   ├── consumer         # 消费者应用
│   │   └── main.go
│   ├── dashboard        # 仪表盘应用
│   │   └── main.go
│   ├── producer         # 生产者应用
│   │   └── main.go
│   └── retry_worker     # 重试处理器应用
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── dashboard        # 仪表盘内部逻辑
│   │   └── stats.go
│   ├── redisclient      # Redis 客户端封装
│   │   └── client.go
│   ├── retry            # 重试逻辑处理
│   │   └── handler.go
│   └── task             # 任务定义及相关处理
│       ├── retry_checker.go
│       └── task.go
└── web
    └── static           # 仪表盘静态文件
        └── index.html
```

## 主要功能

*   **消息生产**：通过生产者将任务发送到 Redis Stream。
*   **消息消费**：通过消费者从 Redis Stream 读取并处理任务。
*   **任务定义**：统一定义任务结构，包含 `task_id`, `payload`, `created_at`。
*   **失败重试**：(推测功能，基于 `retry_worker` 和 `internal/retry`) 对于处理失败的任务，提供重试机制。
*   **仪表盘监控**：通过 Web 界面实时监控队列状态和统计数据。
*   **Redis Stream 应用**：充分利用 Redis Stream 作为消息队列的底层实现。

## 依赖

-   [github.com/redis/go-redis/v9](https://github.com/redis/go-redis)
-   [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
-   (其他依赖请查看 `go.mod` 文件)

## 后续可能的计划

### 基础功能增强
1. **任务幂等性设计**：通过任务ID或业务唯一标识实现幂等处理，避免重复消费。可使用Redis的SET NX命令或事务机制实现。
2. **任务持久日志系统**：将任务执行日志记录至MongoDB或结构化文件，支持查询和分析。可集成标准日志库如zap或logrus。
3. **任务优先级支持**：实现多级优先级队列，通过多个Stream结合ZSet进行调度，确保高优先级任务优先处理。
4. **配置管理**：使用viper库将Redis连接信息、队列名称、端口号等配置外部化，支持配置文件、环境变量和命令行参数。
5. **优雅停机**：实现信号处理（如SIGTERM、SIGINT），确保消费者和重试处理器在收到终止信号时完成当前任务再退出。

### 性能与可靠性
6. **批量处理与并发控制**：实现工作池模式，允许消费者批量获取和并发处理任务，提供配置项控制并发数量。
7. **死信队列 (DLQ)**：为多次重试失败的任务建立专门的死信队列，支持手动检查和重新入队。
8. **限流机制**：基于令牌桶或漏桶算法实现生产者和消费者的限流，防止系统过载。
9. **熔断机制**：当下游服务异常时自动熔断，避免持续失败消耗资源，可集成hystrix-go等库。
10. **任务超时控制**：为任务处理设置最大执行时间，超时自动中断并进入重试流程。

### 监控与运维
11. **增强仪表盘功能**：提供任务详情查看、手动重试特定任务、清理已处理任务等功能，支持按状态筛选和搜索。
12. **指标收集与导出**：集成Prometheus，暴露关键指标如队列长度、处理速率、错误率等，便于监控和告警。
13. **分布式追踪**：集成OpenTelemetry，跟踪任务在不同组件间的流转，便于问题排查。
14. **报警系统**：实现任务积压或失败率异常时的告警通知，支持钉钉/企业微信等多渠道。
15. **健康检查API**：提供各组件健康状态的API接口，便于负载均衡和自动恢复。

### 部署与扩展
16. **Docker Compose打包**：提供完整的Docker Compose配置，实现一键部署整个系统（包括Redis）。
17. **水平扩展支持**：通过消费者组实现多实例并行消费，提高处理能力，支持动态扩缩容。
18. **多租户支持**：通过命名空间隔离不同租户的任务队列，实现资源隔离。
19. **插件系统**：设计插件接口，支持自定义任务处理器、存储后端和通知渠道。
20. **更完善的测试**：增加单元测试、集成测试和基准测试，确保各组件的稳定性、正确性和性能。

### 安全性增强
21. **访问控制**：为仪表盘和API添加基本的认证和授权机制，防止未授权访问。
22. **数据加密**：支持敏感任务数据的加密存储和传输，保护数据安全。
23. **审计日志**：记录关键操作的审计日志，包括手动重试、清理任务等管理操作。

以上扩展计划均可基于现有架构稳定实现，并可根据实际需求分阶段逐步集成。