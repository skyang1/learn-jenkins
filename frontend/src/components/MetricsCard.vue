<template>
  <div class="test-component">
    <div class="comp-header">
      <div class="title-wrap">
        <h4>⚡ 后端运行时指标 (Go Runtime Metrics)</h4>
        <span class="endpoint-badge">GET /api/metrics</span>
      </div>
      <div class="actions">
        <label class="auto-refresh-label">
          <input type="checkbox" v-model="autoRefresh" @change="toggleAutoRefresh" />
          自动刷新 (3s)
        </label>
        <button class="btn btn-xs" @click="fetchMetrics" :disabled="loading">
          {{ loading ? '刷新中...' : '手动刷新' }}
        </button>
      </div>
    </div>

    <div class="metrics-grid">
      <div class="metric-box">
        <span class="label">Goroutines 协程数</span>
        <span class="value text-cyan">{{ metrics.num_goroutine ?? '-' }}</span>
      </div>
      <div class="metric-box">
        <span class="label">当前堆内存 (Alloc)</span>
        <span class="value text-emerald">{{ metrics.alloc_mb ?? '-' }}</span>
      </div>
      <div class="metric-box">
        <span class="label">累计堆分配 (TotalAlloc)</span>
        <span class="value text-blue">{{ metrics.total_alloc_mb ?? '-' }}</span>
      </div>
      <div class="metric-box">
        <span class="label">系统占用内存 (Sys)</span>
        <span class="value text-purple">{{ metrics.sys_mb ?? '-' }}</span>
      </div>
      <div class="metric-box">
        <span class="label">CPU 逻辑核数</span>
        <span class="value text-amber">{{ metrics.num_cpu ?? '-' }} Cores</span>
      </div>
      <div class="metric-box">
        <span class="label">GC 垃圾回收次数</span>
        <span class="value text-rose">{{ metrics.num_gc ?? '-' }} 次</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const metrics = ref({})
const loading = ref(false)
const autoRefresh = ref(false)
let timer = null

const fetchMetrics = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/metrics')
    if (res.ok) {
      metrics.value = await res.json()
    }
  } catch (e) {
    console.error('Fetch metrics failed:', e)
  } finally {
    loading.value = false
  }
}

const toggleAutoRefresh = () => {
  if (autoRefresh.value) {
    timer = setInterval(fetchMetrics, 3000)
  } else if (timer) {
    clearInterval(timer)
    timer = null
  }
}

onMounted(() => {
  fetchMetrics()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
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
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-family: monospace;
}
.actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.auto-refresh-label {
  font-size: 0.85rem;
  color: #94a3b8;
  display: flex;
  align-items: center;
  gap: 0.35rem;
  cursor: pointer;
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
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 0.85rem;
}
.metric-box {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  padding: 0.75rem;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
.metric-box .label {
  font-size: 0.75rem;
  color: #94a3b8;
}
.metric-box .value {
  font-size: 1.15rem;
  font-weight: 700;
  font-family: monospace;
}
.text-cyan { color: #22d3ee; }
.text-emerald { color: #34d399; }
.text-blue { color: #60a5fa; }
.text-purple { color: #c084fc; }
.text-amber { color: #fbbf24; }
.text-rose { color: #fb7185; }
</style>
