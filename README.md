# Kbase

Kbase 是一个本地优先的 TUI 命令知识库工具，使用 YAML 持久化常用命令，并提供毫秒级启动与模糊搜索体验。

## 目标
- 终端内快速检索与记忆常用命令
- 本地文件存储，用户可自行通过 Git 同步
- 通过 $EDITOR 直接编辑数据源并即时重载

## 数据文件
默认路径：~/.config/kbase/commands.yaml（可通过参数指定）

## 项目结构
```
Kbase/
├── cmd/
│   └── kbase/
│       └── main.go           # 程序入口 (Entry Point)
├── internal/                 # 内部业务逻辑 (不对外暴露库)
│   ├── config/               # 配置加载、YAML读写、文件路径处理
│   │   └── config.go
│   ├── model/                # 数据结构定义 (Command, Tag 等)
│   │   └── command.go
│   └── tui/                  # Bubble Tea 的 UI 逻辑 (核心部分)
│       ├── model.go          # TUI 的状态定义 (Tea Model)
│       ├── update.go         # 键盘事件处理 (Ctrl+C, e, Enter)
│       ├── view.go           # 界面渲染逻辑
│       └── styles.go         # Lipgloss 样式定义 (颜色、边框)
├── pkg/                      # (可选) 如果有可以被外部项目复用的通用工具
│   └── utils/                # 如剪贴板封装等
├── assets/                   # 静态资源
│   └── default.yaml          # 默认的示例配置文件 (可嵌入二进制)
├── go.mod
├── go.sum
├── PRD.md
└── README.md
```

## 关键特性
- 启动时读取 YAML，按 platform 过滤
- 实时模糊搜索（cmd/desc/tags）
- 复制命令到剪贴板（Ctrl+c）
- 使用 $EDITOR 打开并编辑数据文件（e）

## 规格说明
详细 PRD 请见 PRD.md。
