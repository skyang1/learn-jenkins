<template>
  <div class="test-component">
    <div class="comp-header">
      <div class="title-wrap">
        <h4>🚨 状态码与异常熔断探测 (Status Code Simulator)</h4>
        <span class="endpoint-badge">GET /api/mock/status?code={code}</span>
      </div>
      <span class="tip-text">测试前端错误兜底与 Jenkins 探针熔断逻辑</span>
    </div>

    <div class="buttons-row">
      <button
        v-for="c in codes"
        :key="c.code"
        :class="['btn-status', c.cls, selectedCode === c.code ? 'active' : '']"
        @click="testStatus(c.code)"
      >
        HTTP {{ c.code }} ({{ c.label }})
      </button>
    </div>

    <div v-if="result" class="result-box">
      <div class="result-header">
        <span>响应代码: <strong :class="result.status >= 400 ? 'text-rose' : 'text-emerald'">{{ result.status }}</strong></span>
        <span>耗时: {{ result.duration }}ms</span>
      </div>
      <pre class="json-preview">{{ JSON.stringify(result.data, null, 2) }}</pre>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const selectedCode = ref(null)
const result = ref(null)

const codes = [
  { code: 200, label: '正常成功', cls: 'btn-ok' },
  { code: 400, label: '参数异常', cls: 'btn-warn' },
  { code: 403, label: '拒绝访问', cls: 'btn-warn' },
  { code: 404, label: '资源缺失', cls: 'btn-warn' },
  { code: 500, label: '服务崩溃', cls: 'btn-err' }
]

const testStatus = async (code) => {
  selectedCode.value = code
  const start = performance.now()
  try {
    const res = await fetch(`/api/mock/status?code=${code}`)
    const duration = Math.round(performance.now() - start)
    const data = await res.json()
    result.value = {
      status: res.status,
      duration,
      data
    }
  } catch (err) {
    result.value = {
      status: '网络错误',
      duration: Math.round(performance.now() - start),
      data: { error: err.message }
    }
  }
}
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
  background: rgba(244, 63, 94, 0.15);
  color: #fb7185;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-family: monospace;
}
.tip-text {
  font-size: 0.8rem;
  color: #94a3b8;
}
.buttons-row {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 0.85rem;
}
.btn-status {
  padding: 0.45rem 0.85rem;
  font-size: 0.82rem;
  border-radius: 6px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
  color: white;
}
.btn-ok {
  background: rgba(16, 185, 129, 0.2);
  border-color: rgba(16, 185, 129, 0.4);
}
.btn-ok:hover, .btn-ok.active {
  background: #10b981;
}
.btn-warn {
  background: rgba(245, 158, 11, 0.2);
  border-color: rgba(245, 158, 11, 0.4);
}
.btn-warn:hover, .btn-warn.active {
  background: #f59e0b;
}
.btn-err {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.4);
}
.btn-err:hover, .btn-err.active {
  background: #ef4444;
}
.result-box {
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  padding: 0.75rem;
}
.result-header {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  margin-bottom: 0.4rem;
  color: #cbd5e1;
}
.json-preview {
  margin: 0;
  font-family: monospace;
  font-size: 0.8rem;
  color: #f1f5f9;
  max-height: 140px;
  overflow-y: auto;
}
.text-rose { color: #f43f5e; }
.text-emerald { color: #10b981; }
</style>
