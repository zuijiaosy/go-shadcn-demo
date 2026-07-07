# 贡献指南

感谢你对本项目的关注！欢迎任何形式的贡献。

## 提交 Issue

- 报告 Bug 时请附上复现步骤、期望行为和实际行为
- 提出新功能建议前，请先搜索是否已有相关讨论

## 提交 Pull Request

1. Fork 本仓库并创建特性分支：`git checkout -b feat/my-feature`
2. 完成开发并确保通过本地检查：

   ```bash
   # Go 检查
   gofmt -l .          # 应无输出
   go vet ./...
   go test ./...

   # 前端检查
   cd frontend
   npm run lint
   npm run build
   ```

3. 提交信息建议遵循 [Conventional Commits](https://www.conventionalcommits.org/zh-hans/) 规范（如 `feat: 添加用户认证`、`fix: 修复静态资源 404`）
4. 推送分支并发起 Pull Request，描述清楚变更内容和动机

## 代码规范

- Go 代码使用 `gofmt` 格式化，注释遵循 godoc 格式
- 前端代码需通过 ESLint 检查，组件文件使用 PascalCase 命名
- 保持提交原子化：一个 PR 只做一件事

## 许可证

提交贡献即表示你同意你的代码以 [MIT 许可证](LICENSE) 发布。
