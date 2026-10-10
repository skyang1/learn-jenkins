pipeline {
    agent any

    parameters {
        choice(name: 'DEPLOY_ENV', choices: ['test', 'prod'], description: '部署环境')
        booleanParam(name: 'SKIP_TESTS', defaultValue: false, description: '是否跳过单元测试')
    }

    environment {
        APP_NAME         = 'vue-go-demo'
        BUILD_VERSION    = "1.0.${BUILD_NUMBER}"
        BUILD_TIME       = ""
        GIT_COMMIT_SHORT = ""
    }

    stages {
        stage('准备与提取元数据') {
            steps {
                script {
                    // 获取当前构建时间和短 Commit Hash
                    env.BUILD_TIME = sh(script: "date '+%Y-%m-%d %H:%M:%S'", returnStdout: true).trim()
                    env.GIT_COMMIT_SHORT = sh(script: "git rev-parse --short HEAD || echo 'local'", returnStdout: true).trim()
                    echo "=========================================="
                    echo "流水线构建开始: ${env.APP_NAME}"
                    echo "版本号:        ${env.BUILD_VERSION}"
                    echo "Git Commit:    ${env.GIT_COMMIT_SHORT}"
                    echo "构建时间:      ${env.BUILD_TIME}"
                    echo "目标环境:      ${params.DEPLOY_ENV}"
                    echo "=========================================="
                }
            }
        }

        stage('自动化测试') {
            when {
                expression { return !params.SKIP_TESTS }
            }
            parallel {
                stage('Go 后端单元测试') {
                    steps {
                        dir('backend') {
                            echo "运行 Go 单元测试..."
                            sh 'go test -v ./...'
                        }
                    }
                }
                stage('Vue 前端语法检查') {
                    steps {
                        dir('frontend') {
                            echo "前端依赖预校验..."
                            sh 'test -f package.json'
                        }
                    }
                }
            }
        }

        stage('镜像构建 (Docker Build)') {
            steps {
                script {
                    echo "正在构建 Docker 镜像..."
                    // 使用 docker-compose 进行镜像构建，注入流水线元数据
                    sh """
                        BUILD_VERSION=${env.BUILD_VERSION} \
                        GIT_COMMIT=${env.GIT_COMMIT_SHORT} \
                        BUILD_TIME="${env.BUILD_TIME}" \
                        docker compose build
                    """
                }
            }
        }

        stage('部署上线 (Docker Compose)') {
            steps {
                script {
                    echo "正在执行部署升级..."
                    // 重新创建并后台启动容器
                    sh """
                        BUILD_VERSION=${env.BUILD_VERSION} \
                        GIT_COMMIT=${env.GIT_COMMIT_SHORT} \
                        BUILD_TIME="${env.BUILD_TIME}" \
                        docker compose up -d --remove-orphans
                    """
                }
            }
        }

        stage('健康检查 (Smoke Test)') {
            steps {
                script {
                    echo "等待服务启动并在进行冒烟测试..."
                    // 轮询检查后端 /api/health 接口，最多重试 10 次
                    sh '''
                        for i in $(seq 1 10); do
                            echo "检查健康状态 (第 $i 次)..."
                            STATUS=$(curl -s http://localhost:8080/api/health | grep '"status":"UP"' || true)
                            if [ -n "$STATUS" ]; then
                                echo "后端服务健康检查通过！"
                                exit 0
                            fi
                            sleep 2
                        done
                        echo "健康检查超时失败！"
                        exit 1
                    '''
                }
            }
        }
    }

    post {
        success {
            echo "=========================================="
            echo "🎉 Jenkins 流水线执行成功！"
            echo "访问地址: 前端 http://<服务器IP> | 后端 API http://<服务器IP>:8080/api/health"
            echo "=========================================="
        }
        failure {
            echo "=========================================="
            echo "❌ Jenkins 流水线执行失败，请检查构建日志！"
            echo "=========================================="
        }
        always {
            script {
                // 清理悬空镜像，避免磁盘占满
                sh 'docker image prune -f || true'
            }
        }
    }
}
