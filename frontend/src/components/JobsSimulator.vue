<template>
  <div class="test-component">
    <div class="comp-header">
      <div class="title-wrap">
        <h4>🚀 异步构建流水线任务测试 (CI/CD Jobs Simulator)</h4>
        <span class="endpoint-badge">GET/POST /api/jobs</span>
      </div>
      <button class="btn btn-xs" @click="fetchJobs" :disabled="loading">
        {{ loading ? '刷新中...' : '刷新任务列表' }}
      </button>
    </div>

    <!-- Quick Dispatch Form -->
    <form @submit.prevent="dispatchJob" class="job-form">
      <input
        v-model="newJobName"
        type="text"
        placeholder="流水线任务名称 (如: sonar-scan, docker-build)"
        class="input-dark flex-1"
        required
      />
      <input
        v-model="newJobBranch"
        type="text"
        placeholder="分支 (默认: main)"
        class="input-dark branch-input"
      />
      <button type="submit" class="btn btn-primary" :disabled="submitting">
        {{ submitting ? '创建中...' : '+ 派发模拟构建' }}
      </button>
    </form>

    <!-- Jobs Table / List -->
    <div class="jobs-list">
      <div v-for="job in jobs" :key="job.id" class="job-item">
        <div class="job-info">
          <span class="job-id">#{{ job.id }}</span>
          <strong class="job-name">{{ job.name }}</strong>
          <span class="job-branch">git:{{ job.branch }}</span>
        </div>
        <div class="job-meta">
          <span :class="['status-tag', 'tag-' + job.status]">
            {{ job.status }}
          </span>
          <span class="job-duration">{{ job.duration_sec }}s</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const jobs = ref([])
const loading = ref(false)
const submitting = ref(false)
const newJobName = ref('')
const newJobBranch = ref('')

const fetchJobs = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/jobs')
    if (res.ok) {
      const data = await res.json()
      jobs.value = data.data || []
    }
  } catch (e) {
    console.error('Fetch jobs failed:', e)
  } finally {
    loading.value = false
  }
}

const dispatchJob = async () => {
  if (!newJobName.value.trim()) return
  submitting.value = true
  try {
    const res = await fetch('/api/jobs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: newJobName.value,
        branch: newJobBranch.value || 'main'
      })
    })
    if (res.ok) {
      newJobName.value = ''
      newJobBranch.value = ''
      await fetchJobs()
    }
  } catch (e) {
    alert('派发任务失败: ' + e.message)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchJobs()
})
</script>

<style scoped>
.test-component {
  background: rgba(15, 23, 42, 0.45);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 1.25rem;
  margin-top: 1rem;
}
.comp-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}
.title-wrap {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}
.title-wrap h4 {
  font-size: 1.05rem;
  color: #e2e8f0;
}
.endpoint-badge {
  font-size: 0.75rem;
  background: rgba(234, 179, 8, 0.15);
  color: #facc15;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-family: monospace;
}
.btn-xs {
  padding: 0.3rem 0.65rem;
  font-size: 0.8rem;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.btn-xs:hover {
  background: rgba(255, 255, 255, 0.2);
}
.job-form {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.85rem;
}
.flex-1 { flex: 1; }
.branch-input { width: 140px; }
.input-dark {
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
  border-radius: 6px;
  padding: 0.5rem 0.75rem;
  font-size: 0.85rem;
}
.btn {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.85rem;
  cursor: pointer;
  border: none;
  font-weight: 500;
}
.btn-primary {
  background: #3b82f6;
  color: white;
}
.jobs-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 200px;
  overflow-y: auto;
}
.job-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.05);
  padding: 0.6rem 0.85rem;
  border-radius: 6px;
}
.job-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.job-id {
  font-family: monospace;
  color: #94a3b8;
  font-size: 0.8rem;
}
.job-name {
  color: #f8fafc;
  font-size: 0.9rem;
}
.job-branch {
  font-size: 0.75rem;
  background: rgba(255, 255, 255, 0.08);
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  color: #cbd5e1;
}
.job-meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.status-tag {
  font-size: 0.75rem;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
  text-transform: uppercase;
}
.tag-success {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}
.tag-running {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
}
.tag-failed {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
}
.job-duration {
  font-size: 0.75rem;
  color: #94a3b8;
  font-family: monospace;
}
</style>
