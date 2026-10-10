# Vue + Go Jenkins CI/CD 自动化部署演示项目

本项目是一个专门用于演练 **Jenkins 持续集成与持续部署 (CI/CD)** 的前后端分离 Demo。
包含了完整的 **Vue 3 前端**、**Go 1.22+ 后端**、**多阶段 Docker 容器化打包配置**、**Docker Compose 编排文件** 以及生产级标准的 **Jenkinsfile (Pipeline as Code)**。

---

## 目录结构

```text
.
├── backend/                  # Go 后端工程
│   ├── Dockerfile            # 后端多阶段构建 Dockerfile
│   ├── go.mod                # Go 模块定义 (零第三方依赖，纯标准库，极速构建)
│   ├── main.go               # REST API 服务 (含 ldflags 动态注入构建元数据)
│   ├── main_test.go          # 单元测试 (用于 Jenkins 自动化测试阶段)
│   └── .dockerignore
├── frontend/                 # Vue 3 前端工程 (Vite)
│   ├── Dockerfile            # 前端多阶段构建 Dockerfile (Node 构建 + Nginx 托管)
│   ├── nginx.conf            # Nginx 生产反代配置 (解决 SPA 路由与 /api/ 跨域)
│   ├── package.json          # 前端依赖配置
│   ├── vite.config.js        # Vite 配置 (开发代理配置)
│   ├── src/                  # 页面源码 (动态读取后端版本、提交记录、健康状态)
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

## 核心特性

1. **动态版本注入 (DevOps 实战技巧)**:
   - Go 编译时利用 `-ldflags "-X main.Version=... -X main.GitCommit=... -X main.BuildTime=..."`，自动将 Jenkins 的 `${BUILD_NUMBER}`、Git 短 Hash 和构建时间注入二进制中。
   - 打开前端页面即可实时看到本次构建的流水线元数据，直观验证部署是否真正生效。
2. **自动化测试守门 (Quality Gate)**:
   - 后端自带 `go test` 单元测试，流水线可设置测试不通过即时熔断阻断部署。
3. **极简轻量，高可靠**:
   - 后端基于 Go 1.22+ 原生强大路由与并发模型，无冗余三方包，彻底避免 CI 机器因为网络拉取包超时失败。
   - 前后端均采用多阶段构建 (Multi-stage build)，生成的镜像体积小、安全度高。

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
