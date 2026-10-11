# AgentToolGateway

一个基于 Gin 的简单 Agent 工具调用网关，提供工具注册、风险评估、审批流程和审计日志。

## 项目结构

```
AgentToolGateway/
├── client/
│   ├── controller/service.go   # 简单的 HTTP 客户端
│   └── main.go                 # 客户端示例
├── server/
│   ├── gateway/
│   │   ├── model.go            # ToolCall、Tool、风险等级等模型
│   │   ├── decision.go         # 策略决策（allow / required_approval / deny）
│   │   ├── approval.go         # 审批管理
│   │   ├── audit.go            # 审计日志
│   │   └── service.go          # 核心网关服务
│   ├── tools/
│   │   ├── read.go             # read_note 工具
│   │   ├── email.go            # send_email 工具
│   │   └── delete.go           # delete_note 工具
│   └── main.go                 # HTTP 服务入口
├── sandbox/
│   ├── notes/                  # 笔记文件（read_note / delete_note 使用）
│   └── outbox/                 # 发送邮件后的输出目录
├── go.mod
└── go.sum
```

## 运行

```bash
# 启动服务端（默认监听 :8080）
go run ./server

# 运行客户端示例（调用 send_email）
go run ./client
```

## API

### 1. 调用工具

```
POST /toolCalls
```

请求体示例：

```json
{
  "call_id": "001",
  "tool_name": "send_email",
  "args": {
    "to": "test@qq.com",
    "subject": "hello",
    "body": "this is a test"
  }
}
```

### 2. 审批

```
POST /approvals/:id
```

审批通过后执行对应工具。

### 3. 审计日志

```
GET /logs
```

返回已记录的审计事件列表。

## 已注册工具

| 工具名       | 风险等级 | 说明                         |
| ------------ | -------- | ---------------------------- |
| `read_note`  | low      | 读取 `sandbox/notes` 下的文件 |
| `send_email` | low      | 将邮件写入 `sandbox/outbox`   |
| `delete_note`| low      | 删除 `sandbox/notes` 下的文件 |

风险等级与决策：

- `low` -> `allow`（直接执行）
- `medium` -> `required_approval`（需要审批）
- `high` -> `deny`（拒绝执行）

## 注意

- 服务端 `ToolCall.Args` 的 JSON 标签为 `args`，客户端调用时字段名需保持一致。
