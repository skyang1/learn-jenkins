# Vue + Go (Gin) Jenkins CI/CD 自动化部署演示项目

本项目是一个专门用于演练 **Jenkins 持续集成与持续部署 (CI/CD)** 的前后端分离 Demo。
包含了完整的 **Vue 3 前端**、**Go (Gin Web 框架) 后端**、**多阶段 Docker 容器化打包配置**、**Docker Compose 编排文件** 以及生产级标准的 **Jenkinsfile (Pipeline as Code)**。

---

## 目录结构

```text
.
├── backend/                  # Go 后端工程 (Gin 框架)
│   ├── Dockerfile            # 后端多阶段构建 Dockerfile (配置 GOPROXY 加速与缓存)
│   ├── go.mod / go.sum       # Go 依赖 (Gin Web Framework)
│   ├── main.go               # Gin REST API 服务 (含 ldflags 动态注入构建元数据)
│   ├── main_test.go          # 单元测试 (测试 Gin 路由与接口，用于 Jenkins 门禁测试)
│   └── .dockerignore
├── frontend/                 # Vue 3 前端工程 (Vite)
│   ├── Dockerfile            # 前端多阶段构建 Dockerfile (Node 构建 + Nginx 托管)
│   ├── nginx.conf            # Nginx 生产反代配置 (解决 SPA 路由与 /api/ 跨域)
│   ├── package.json          # 前端依赖配置
│   ├── vite.config.js        # Vite 配置 (开发代理配置)
│   ├── src/
│   │   ├── components/       # 接口与数据结构测试套件组件
│   │   │   ├── MetricsCard.vue    # Go 运行时指标 (Goroutines、堆内存、CPU、GC)
│   │   │   ├── StatusTester.vue   # 状态码与异常熔断测试 (200, 400, 403, 404, 500)
│   │   │   ├── EchoInspector.vue  # 请求回显与自定义 Header / JSON Body 探测
│   │   │   └── JobsSimulator.vue  # 异步 CI/CD 流水线任务流模拟
│   │   ├── App.vue           # 页面主框架 (Tab 切换、健康检查、流水线图示)
│   │   └── main.js
│   └── .dockerignore
├── docker-compose.yml        # Docker 编排文件 (一键启动/部署)
├── Jenkinsfile               # 标准 Jenkins 声明式流水线 (推荐，基于 Docker)
├── Jenkinsfile.direct        # 传统宿主机流水线 (备用，无 Docker 时直接运行二进制)
├── scripts/
│   ├── build.sh              # 离线/手动构建脚本
│   └── deploy.sh             # 快捷部署脚本
└── README.md
```

---

## 核心特性与测试套件

1. **主流高性能 Web 框架 (Gin)**:
   - 基于 Gin 路由引擎，提供优雅的 JSON 序列化、中间件机制与参数校验。
2. **丰富的后端接口测试套件 (Vue 交互组件)**:
   - ⚡ **运行时指标 (GET /api/metrics)**: 实时探测 Go 协程数、堆分配、CPU 核数和 GC 次数。
   - 🚨 **状态码与异常熔断 (GET /api/mock/status)**: 模拟不同 HTTP 状态，验证前端错误拦截与探针熔断。
   - 📡 **请求回显测试 (POST /api/echo)**: 发送自定义 Header（如 CI 令牌）和复杂 JSON Body，检验网关反代和 Gin 参数绑定。
   - 🚀 **异步任务模拟 (GET/POST /api/jobs)**: 模拟提交 CI/CD 扫描、打包等流水线子任务。
   - 📝 **留言板数据测试 (GET/POST/DELETE /api/messages)**: 经典 CRUD 接口验证连通性。
3. **动态版本注入 (DevOps 实战技巧)**:
   - Go 编译时利用 `-ldflags "-X main.Version=... -X main.GitCommit=... -X main.BuildTime=..."`，自动将 Jenkins 的 `${BUILD_NUMBER}`、Git 短 Hash 和构建时间注入二进制中。
   - 页面实时展示流水线元数据，直观检验部署生效状态。
4. **自动化测试守门 (Quality Gate)**:
   - 自带针对 Gin 路由的自动化单元测试，若测试失败流水线将自动熔断阻止发布。

---

## 本地快速体验与调试

### 1. 启动 Go 后端 (终端 1)
```bash
cd backend
go run main.go
# 服务默认运行在 http://localhost:8080
# 测试接口: curl http://localhost:8080/api/health
```

### 2. 启动 Vue 前端 (终端 2)
```bash
cd frontend
npm install
npm run dev
# 访问 http://localhost:3000
```

### 3. 使用 Docker Compose 一键启动验证
```bash
docker compose up -d --build
# 访问 http://localhost (前端已自动反代 /api 至 Go 后端)
```

---

## Jenkins 自动化部署配置教程

### 步骤一：环境准备 (在 Jenkins 所在服务器上)

1. **安装必要工具**:
   - 安装 Docker 和 Docker Compose。
2. **配置 Jenkins 用户权限** (关键避坑点):
   Jenkins 默认使用 `jenkins` 用户执行 shell，必须将其加入 `docker` 用户组，否则会报错 `permission denied`:
   ```bash
   sudo usermod -aG docker jenkins
   sudo systemctl restart jenkins
   ```
3. **Jenkins 推荐插件**:
   - `Git` (核心插件)
   - `Pipeline` (工作流插件)
   - `Docker Pipeline` (可选)

---

### 步骤二：新建流水线任务 (Pipeline)

1. 登录 Jenkins 后台，点击 **新建任务 (New Item)**。
2. 输入任务名称（如 `vue-go-ci-cd`），选择 **流水线 (Pipeline)**，点击确定。
3. 进入任务配置页：
   - 勾选 **“This project is parameterized” (参数化构建过程)**：
     - 添加 **选项参数 (Choice Parameter)**:
       - 名称: `DEPLOY_ENV`
       - 选项: `test\nprod`
     - 添加 **布尔值参数 (Boolean Parameter)**:
       - 名称: `SKIP_TESTS`
       - 默认值: 保持不勾选
   - 滑动到最下方 **流水线 (Pipeline)** 区域：
     - 定义: 选择 **Pipeline script from SCM**
     - SCM: 选择 **Git**
     - Repository URL: 填入你的代码仓库地址 (GitHub / Gitee / GitLab)
     - Credentials: 添加并选择你的代码访问凭据 (SSH Key 或账号密码)
     - 指定分支: `*/main` 或 `*/master`
     - 脚本路径: `Jenkinsfile`
4. 点击 **保存 (Save)**。

---

### 步骤三：触发构建与验证

1. 在任务页面点击 **Build with Parameters (带参数构建)**。
2. 选择环境后点击 **开始构建 (Build)**。
3. 观察流水线阶段视图 (Stage View)：
   - **准备与提取元数据** &rarr;
   - **自动化测试 (Go 单元测试)** &rarr;
   - **镜像构建 (Docker Build)** &rarr;
   - **部署上线 (Docker Compose)** &rarr;
   - **健康检查 (curl /api/health)**
4. 构建成功后，在浏览器访问目标服务器 `http://<服务器IP>`，可以看到页面上展示了：
   - 当前 Jenkins 的 **构建版本号** (如 `1.0.1`)
   - 代码的 **Git Commit Hash**
   - 构建时间与 Go 运行状态

---

### 步骤四：配置 Webhook 自动触发 (实现 Push 代码自动部署)

1. **Jenkins 开启触发器**:
   - 在任务配置的 **构建触发器** 中勾选 **GitHub hook trigger for GITScm polling** (或 Gitee WebHook)。
2. **代码托管平台配置 Webhook**:
   - 进入 GitHub / Gitee 仓库的 `Settings` &rarr; `Webhooks`。
   - 添加 Webhook，Payload URL 填写：
     - GitHub 格式: `http://<Jenkins外网IP:端口>/github-webhook/`
     - Gitee 格式: `http://<Jenkins外网IP:端口>/gitee-project-hook/`
   - 选择触发事件: `Push events`。
3. 配置完成后，本地修改代码并 `git push`，Jenkins 将会自动开始流水线构建！

---

## 备用方案：宿主机直接部署 (无需 Docker)

如果你的 Jenkins 环境或测试机没有安装 Docker，可以使用本项目提供的 `Jenkinsfile.direct`:
1. 在 Jenkins 任务配置中，将脚本路径改为 `Jenkinsfile.direct`。
2. 确保 Jenkins 节点安装了 `go` 和 `node` 环境。
3. 流水线将直接在宿主机编译 Go 可执行文件并运行，编译前端后输出到指定目录。
