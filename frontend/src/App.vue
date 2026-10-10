<template>
  <div class="container">
    <!-- Header -->
    <header class="header">
      <div class="logo-group">
        <span class="badge vue-badge">Vue 3</span>
        <span class="plus">+</span>
        <span class="badge go-badge">Go 1.22+</span>
        <span class="plus">+</span>
        <span class="badge jenkins-badge">Jenkins CI/CD</span>
      </div>
      <h1>自动化部署演示项目</h1>
      <p class="subtitle">用于演练 Jenkins 流水线 (Checkout &rarr; Test &rarr; Build &rarr; Deploy &rarr; HealthCheck)</p>
    </header>

    <!-- Status Cards Grid -->
    <div class="grid">
      <!-- Service Health -->
      <div class="card">
        <div class="card-header">
          <h3>后端健康状态</h3>
          <span :class="['status-pill', healthStatus === 'UP' ? 'online' : 'offline']">
            {{ healthStatus === 'UP' ? '服务正常 (UP)' : '离线 (OFFLINE)' }}
          </span>
        </div>
        <div class="card-body">
          <p><strong>运行时长:</strong> {{ healthData.uptime || 'N/A' }}</p>
          <p><strong>心跳时间:</strong> {{ healthData.timestamp || 'N/A' }}</p>
          <button class="btn btn-secondary" @click="fetchHealth">
            刷新健康状态
          </button>
        </div>
      </div>

      <!-- Build & CI/CD Info -->
      <div class="card">
        <div class="card-header">
          <h3>CI/CD 构建元数据</h3>
          <span class="tag">动态注入 ldflags</span>
        </div>
        <div class="card-body">
          <p><strong>应用名称:</strong> {{ appInfo.app || '加载中...' }}</p>
          <p><strong>部署版本:</strong> <span class="highlight">{{ appInfo.version || 'v1.0.0' }}</span></p>
          <p><strong>Git Commit:</strong> <code>{{ appInfo.git_commit || 'unspecified' }}</code></p>
          <p><strong>构建时间:</strong> {{ appInfo.build_time || 'unspecified' }}</p>
          <p><strong>Go 版本:</strong> {{ appInfo.go_version || 'N/A' }}</p>
          <p><strong>运行环境:</strong> <span class="env-badge">{{ appInfo.environment || 'production' }}</span></p>
        </div>
      </div>
    </div>

    <!-- Messages / CRUD Test -->
    <div class="card mt-4">
      <div class="card-header">
        <h3>交互验证：实时留言列表</h3>
        <span class="subtext">验证前后端联调及数据读写接口</span>
      </div>
      <div class="card-body">
        <form @submit.prevent="addMessage" class="form-row">
          <input
            v-model="newAuthor"
            type="text"
            placeholder="发布人 (如: Jenkins Bot)"
            class="input input-small"
          />
          <input
            v-model="newContent"
            type="text"
            placeholder="输入发布内容，测试接口连通性..."
            class="input input-large"
            required
          />
          <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
            {{ isSubmitting ? '发布中...' : '提交测试' }}
          </button>
        </form>

        <div v-if="loadingMessages" class="loading">正在加载数据...</div>

        <ul v-else-if="messages.length" class="message-list">
          <li v-for="item in messages" :key="item.id" class="message-item">
            <div class="message-meta">
              <span class="author">{{ item.author }}</span>
              <span class="time">{{ formatTime(item.created_at) }}</span>
            </div>
            <div class="message-content">{{ item.content }}</div>
            <button class="btn-del" @click="deleteMessage(item.id)" title="删除">&times;</button>
          </li>
        </ul>
        <div v-else class="empty">暂无留言数据</div>
      </div>
    </div>

    <!-- Jenkins Pipeline Explanation -->
    <div class="card mt-4 pipeline-guide">
      <div class="card-header">
        <h3>Jenkins 流水线设计与执行过程</h3>
      </div>
      <div class="card-body">
        <div class="pipeline-steps">
          <div class="step-card">
            <div class="step-num">1</div>
            <h4>拉取代码</h4>
            <p>Git Checkout 检出最新分支或 Tag 代码</p>
          </div>
          <div class="arrow">&rarr;</div>
          <div class="step-card">
            <div class="step-num">2</div>
            <h4>自动化测试</h4>
            <p>执行 <code>go test</code> 单元测试与语法检查</p>
          </div>
          <div class="arrow">&rarr;</div>
          <div class="step-card">
            <div class="step-num">3</div>
            <h4>编译构建</h4>
            <p>前端 <code>npm run build</code> + 容器镜像打包</p>
          </div>
          <div class="arrow">&rarr;</div>
          <div class="step-card">
            <div class="step-num">4</div>
            <h4>部署上线</h4>
            <p>Docker Compose 容器滚动升级或重启生效</p>
          </div>
          <div class="arrow">&rarr;</div>
          <div class="step-card">
            <div class="step-num">5</div>
            <h4>冒烟验证</h4>
            <p>探测 <code>/api/health</code> 接口并发送通知</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="footer">
      <p>Vue + Go Demo for Jenkins Automation Practice &middot; 快速验证 CI/CD 流水线状态</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const healthStatus = ref('OFFLINE')
const healthData = ref({})
const appInfo = ref({})
const messages = ref([])
const loadingMessages = ref(false)
const isSubmitting = ref(false)

const newAuthor = ref('')
const newContent = ref('')

const fetchHealth = async () => {
  try {
    const res = await fetch('/api/health')
    if (res.ok) {
      healthData.value = await res.json()
      healthStatus.value = 'UP'
    } else {
      healthStatus.value = 'DOWN'
    }
  } catch (e) {
    healthStatus.value = 'OFFLINE'
  }
}

const fetchInfo = async () => {
  try {
    const res = await fetch('/api/info')
    if (res.ok) {
      appInfo.value = await res.json()
    }
  } catch (e) {
    console.error('Failed to fetch info:', e)
  }
}

const fetchMessages = async () => {
  loadingMessages.value = true
  try {
    const res = await fetch('/api/messages')
    if (res.ok) {
      const data = await res.json()
      messages.value = data.data || []
    }
  } catch (e) {
    console.error('Failed to fetch messages:', e)
  } finally {
    loadingMessages.value = false
  }
}

const addMessage = async () => {
  if (!newContent.value.trim()) return
  isSubmitting.value = true
  try {
    const res = await fetch('/api/messages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        author: newAuthor.value || 'Anonymous',
        content: newContent.value
      })
    })
    if (res.ok) {
      newContent.value = ''
      await fetchMessages()
    }
  } catch (e) {
    alert('提交失败: ' + e.message)
  } finally {
    isSubmitting.value = false
  }
}

const deleteMessage = async (id) => {
  try {
    const res = await fetch(`/api/messages?id=${id}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      await fetchMessages()
    }
  } catch (e) {
    alert('删除失败: ' + e.message)
  }
}

const formatTime = (isoString) => {
  if (!isoString) return ''
  const d = new Date(isoString)
  return d.toLocaleString()
}

onMounted(() => {
  fetchHealth()
  fetchInfo()
  fetchMessages()
})
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.header {
  text-align: center;
  padding: 2rem 0;
}

.logo-group {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.badge {
  padding: 0.35rem 0.85rem;
  border-radius: 9999px;
  font-weight: 600;
  font-size: 0.9rem;
}

.vue-badge {
  background: rgba(66, 184, 131, 0.2);
  color: #42b883;
  border: 1px solid rgba(66, 184, 131, 0.4);
}

.go-badge {
  background: rgba(0, 172, 215, 0.2);
  color: #00acd7;
  border: 1px solid rgba(0, 172, 215, 0.4);
}

.jenkins-badge {
  background: rgba(211, 56, 51, 0.2);
  color: #f87171;
  border: 1px solid rgba(211, 56, 51, 0.4);
}

.plus {
  color: var(--text-secondary);
  font-weight: bold;
}

h1 {
  font-size: 2.2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(to right, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 1.05rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.5rem;
}

.card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 12px;
  padding: 1.5rem;
  backdrop-filter: blur(12px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.mt-4 {
  margin-top: 1rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.2rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--card-border);
}

.card-header h3 {
  font-size: 1.2rem;
  font-weight: 600;
}

.card-body p {
  margin: 0.6rem 0;
  font-size: 0.95rem;
  color: #cbd5e1;
}

.status-pill {
  font-size: 0.8rem;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-weight: 600;
}

.online {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border: 1px solid #10b981;
}

.offline {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
  border: 1px solid #ef4444;
}

.highlight {
  color: #38bdf8;
  font-weight: 600;
}

code {
  background: rgba(0, 0, 0, 0.3);
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  color: #f472b6;
  font-family: monospace;
}

.env-badge {
  background: #334155;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
  text-transform: uppercase;
}

.tag {
  font-size: 0.75rem;
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}

.subtext {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.btn {
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: white;
  margin-top: 0.5rem;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.2);
}

.form-row {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1.25rem;
}

.input {
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid var(--card-border);
  color: white;
  padding: 0.6rem 1rem;
  border-radius: 8px;
  outline: none;
}

.input:focus {
  border-color: #3b82f6;
}

.input-small {
  width: 180px;
}

.input-large {
  flex: 1;
}

.message-list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.message-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  padding: 0.85rem 1rem;
  border-radius: 8px;
}

.message-meta {
  display: flex;
  flex-direction: column;
  min-width: 140px;
}

.author {
  font-weight: 600;
  color: #38bdf8;
  font-size: 0.9rem;
}

.time {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.message-content {
  flex: 1;
  padding: 0 1rem;
  color: #e2e8f0;
}

.btn-del {
  background: transparent;
  color: #ef4444;
  border: none;
  font-size: 1.4rem;
  cursor: pointer;
  padding: 0 0.5rem;
  line-height: 1;
}

.btn-del:hover {
  color: #f87171;
}

.empty, .loading {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
}

/* Pipeline Guide */
.pipeline-steps {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  overflow-x: auto;
  padding: 1rem 0;
}

.step-card {
  flex: 1;
  min-width: 140px;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  padding: 1rem 0.8rem;
  text-align: center;
}

.step-num {
  width: 24px;
  height: 24px;
  background: #3b82f6;
  color: white;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.step-card h4 {
  font-size: 0.95rem;
  margin-bottom: 0.35rem;
  color: #f1f5f9;
}

.step-card p {
  font-size: 0.75rem;
  color: #94a3b8;
  line-height: 1.4;
  margin: 0;
}

.arrow {
  color: #64748b;
  font-size: 1.2rem;
  font-weight: bold;
}

.footer {
  text-align: center;
  padding: 2rem 0;
  color: var(--text-secondary);
  font-size: 0.85rem;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

@media (max-width: 768px) {
  .pipeline-steps {
    flex-direction: column;
  }
  .arrow {
    transform: rotate(90deg);
  }
  .form-row {
    flex-direction: column;
  }
  .input-small {
    width: 100%;
  }
}
</style>
