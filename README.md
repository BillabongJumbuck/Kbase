# Kbase

Kbase 是一个本地优先的 TUI 命令知识库工具，使用 YAML 持久化常用命令，并提供毫秒级启动与模糊搜索体验。

## 目标
- 终端内快速检索与记忆常用命令
- 本地文件存储，用户可自行通过 Git 同步
- 通过 $EDITOR 直接编辑数据源并即时重载

## 数据文件
默认路径：~/.config/kbase/commands.yaml（可通过参数指定）

## 关键特性
- 启动时读取 YAML，按 platform 过滤
- 实时模糊搜索（cmd/desc/tags）
- 复制命令到剪贴板（Ctrl+c）
- 使用 $EDITOR 打开并编辑数据文件（e）

## 规格说明
详细 PRD 请见 PRD.md。
