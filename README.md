# Flomo Go Tools

这是一个用 Go 语言开发的 Flomo 工具集，提供命令行工具和 MCP 服务器两种方式来发送笔记到 Flomo。

## 功能特点

- 📝 支持通过命令行快速发送笔记
- 🏷️ 支持添加标签
- 🔌 提供 MCP 服务器集成
- 📊 详细的日志记录
- 🌈 友好的命令行界面

## 安装

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/mcp-server-flomo-go.git
cd mcp-server-flomo-go
```

2. 设置环境变量：
```bash
cp .env.example .env
```
编辑 `.env` 文件，添加你的 Flomo API URL：
```
FLOMO_API_URL=https://flomoapp.com/iwh/xxx/xxx
```

3. 编译 CLI 工具：
```bash
go build -o flomo cmd/flomo/main.go
```

## 使用方法

### CLI 工具

1. 基本使用：
```bash
./flomo -c "你的笔记内容"
```

2. 添加标签：
```bash
./flomo -c "笔记内容" -t "标签1,标签2"
```

3. 从管道输入：
```bash
echo "笔记内容" | ./flomo
```

4. 显示详细信息：
```bash
./flomo -c "笔记内容" -v
```

### 命令行参数

- `-c, --content`: 笔记内容（必需）
- `-t, --tags`: 标签列表，用逗号分隔（可选）
- `-v, --verbose`: 显示详细信息（可选）
- `-h, --help`: 显示帮助信息

### MCP 服务器

1. 启动服务器：
```bash
go run server.go
```

2. 服务器提供的工具：
- `write_note`: 写入笔记
  - 参数：`content` (string) - 笔记内容，支持 Markdown 格式

## 项目结构

```
.
├── cmd
│   └── flomo
│       └── main.go    # CLI 工具实现
├── pkg
│   └── flomo
│       └── client.go  # Flomo API 客户端
├── server.go          # MCP 服务器实现
├── .env              # 环境变量配置
└── README.md         # 本文档
```

## 开发

1. 代码风格遵循 Go 标准
2. 使用 `go fmt` 格式化代码
3. 确保所有日志信息清晰可读

## 示例输出

```
Note sent successfully! 🎉
Created at: 2025-04-16 00:35:06
Tags: 测试
View at: https://v.flomoapp.com/mine/?memo_id=xxx

Detailed information:
- Source: incoming_webhook
- Creator ID: xxx
- Response code: 0
- Response message: 已记录
- Total time: 896.754667ms
```

## 注意事项

1. 请妥善保管你的 Flomo API URL，不要分享给他人
2. 建议在发送大量笔记时适当控制频率
3. 如果遇到问题，可以使用 `-v` 参数查看详细日志

## 许可证

MIT License 