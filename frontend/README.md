# Frontend

基于 React 19 + TypeScript + Vite 7 + Tailwind CSS 4 + shadcn/ui 的前端工程。

## 常用命令

```bash
npm install     # 安装依赖
npm run dev     # 开发服务器（http://localhost:5173，/api 自动代理到 8080）
npm run lint    # ESLint 检查
npm run build   # 生产构建（输出到 dist/）
```

## 目录结构

```
src/
├── components/ui/   # shadcn/ui 组件（Button、Card 等）
├── lib/             # 工具函数（cn 等）
├── App.tsx          # 根组件
├── main.tsx         # 入口
└── index.css        # Tailwind 与主题变量
```

## 添加 shadcn/ui 组件

```bash
npx shadcn@latest add [component-name]
```

更多说明见根目录 [README.md](../README.md)。
