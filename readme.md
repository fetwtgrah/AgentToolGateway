# AgentToolGateway

一个基于 Gin 的 Agent 工具调用网关，提供工具注册、风险评估、审批流程和审计日志。

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

### 1. Client 和 GatewayServer 分别负责什么？

- **Client**（`client/`）：代表 Agent 侧，只负责把工具调用请求（`call_id`、`tool_name`、`args`）通过 HTTP 发送给网关。它不做权限判断，也不直接执行工具。
- **GatewayServer**（`server/gateway/` + `server/main.go`）：所有工具执行的唯一入口，负责：
  - 工具注册（`ToolRepo`）
  - 风险策略判定（`Policy`）
  - 审批管理（`ApprovalManager`）
  - 审计日志（`AuditLogger`）
  - 实际调用工具执行（`HandleToolCall` / `ExecuteApproval`）

### 2. 为什么 Agent 输出不能被视为用户授权？

Agent（LLM）是不可信组件，可能被提示注入、产生幻觉，或被恶意控制。它的输出只能算“请求”，不能算“授权”。真正的授权必须来自用户通过独立通道的明确确认。

因此本项目对中风险操作会生成审批，由用户调用 `/approvals/:id` 明确同意后才执行，而不是相信 Agent 说“用户已经同意了”。

### 3. 为什么确认必须绑定具体参数，而不能只绑定工具名称？

如果审批只绑定工具名，就存在“审批后参数被替换”的风险。例如用户批准了 `send_email`，但真正执行时把收件人、内容换成攻击者想要的内容。

本项目在审批对象中保存完整的 `ToolCall`（包含 `Args`），`ExecuteApproval` 只执行审批时保存的那份参数，保证“批准什么，就执行什么”。

### 4. 为什么工具的风险策略必须保存在 Gateway，而不能由 Agent 提交？

Agent 不可信。如果风险等级由 Agent 提交，Agent（或被注入的指令）可以把高风险工具标记为 `low`，从而绕过审批或拒绝机制。

风险等级在服务端注册工具时确定（`Tool.Risk`），由 `Policy.Check` 在网关侧判定，Agent 无法影响决策。

### 5. 如果一次工具执行超时，系统如何判断它能否安全重试？

能否安全重试，取决于工具是否幂等、以及是否可能已经产生了副作用：

- **幂等工具**（重复执行结果相同）可以安全自动重试，例如 `read_note`。
- **有副作用的工具**（如 `send_email`、`delete_note`）在超时时可能已经部分执行或已经生效，不能盲目重试，否则可能重复发信、重复删除。

正确做法是：

1. 在 Gateway 中记录每次执行的状态（`pending` / `success` / `fail`）。
2. 重试前先查询上一次执行是否已生效，或使用幂等键/去重机制。
3. 只有确认“未产生副作用”或“工具本身幂等”时才自动重试；否则需要用户重新确认。

> 注：当前代码尚未实现超时重试机制，以上是设计上应遵循的原则。

## 注意

- 服务端 `ToolCall.Args` 的 JSON 标签为 `args`，客户端调用时字段名需保持一致。
- 测试功能
```bash
go test ./test -v 
```