<template>
  <div class="test-component">
    <div class="comp-header">
      <div class="title-wrap">
        <h4>📡 请求回显与网关链路检测 (Echo & Headers Inspector)</h4>
        <span class="endpoint-badge">POST /api/echo</span>
      </div>
      <span class="tip-text">测试自定义请求头、反代代理头及 Body 结构反解</span>
    </div>

    <div class="form-grid">
      <div class="form-col">
        <label class="field-label">自定义请求头 (X-Custom-Header)</label>
        <input v-model="customHeader" type="text" class="input-dark" placeholder="例如: Jenkins-CI-Pipeline-101" />

        <label class="field-label mt-2">测试 JSON Body</label>
        <textarea v-model="jsonBody" class="textarea-dark" rows="4"></textarea>

        <button class="btn btn-primary mt-2" @click="sendEcho" :disabled="sending">
          {{ sending ? '发送中...' : '发送请求测试' }}
        </button>
      </div>

      <div class="form-col">
        <label class="field-label">后端 Gin 解析结果 (Echo 回显)</label>
        <div class="echo-result">
          <div v-if="echoResponse">
            <pre class="json-preview">{{ JSON.stringify(echoResponse, null, 2) }}</pre>
          </div>
          <div v-else class="placeholder">点击“发送请求测试”查看回显结果...</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const customHeader = ref('Jenkins-CI-Token-Test')
const jsonBody = ref('{\n  "module": "payment",\n  "action": "refund",\n  "amount": 99.8\n}')
const echoResponse = ref(null)
const sending = ref(false)

const sendEcho = async () => {
  sending.value = true
  try {
    let parsedBody = {}
    try {
      parsedBody = JSON.parse(jsonBody.value)
    } catch (e) {
      alert('JSON 语法错误，请检查输入')
      sending.value = false
      return
    }

    const res = await fetch('/api/echo?timestamp=' + Date.now(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Custom-Header': customHeader.value
      },
      body: JSON.stringify(parsedBody)
    })
    echoResponse.value = await res.json()
  } catch (err) {
    echoResponse.value = { error: err.message }
  } finally {
    sending.value = false
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
  background: rgba(168, 85, 247, 0.15);
  color: #c084fc;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-family: monospace;
}
.tip-text {
  font-size: 0.8rem;
  color: #94a3b8;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 1.25rem;
}
@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
.field-label {
  display: block;
  font-size: 0.8rem;
  color: #94a3b8;
  margin-bottom: 0.35rem;
}
.mt-2 {
  margin-top: 0.65rem;
}
.input-dark, .textarea-dark {
  width: 100%;
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
  border-radius: 6px;
  padding: 0.5rem 0.75rem;
  font-family: monospace;
  font-size: 0.85rem;
  box-sizing: border-box;
}
.input-dark:focus, .textarea-dark:focus {
  outline: none;
  border-color: #38bdf8;
}
.btn {
  padding: 0.55rem 1rem;
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
.btn-primary:hover {
  background: #2563eb;
}
.echo-result {
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  min-height: 155px;
  padding: 0.6rem;
  display: flex;
  flex-direction: column;
}
.placeholder {
  color: #64748b;
  font-size: 0.85rem;
  margin: auto;
  text-align: center;
}
.json-preview {
  margin: 0;
  font-family: monospace;
  font-size: 0.78rem;
  color: #a5f3fc;
  max-height: 220px;
  overflow-y: auto;
}
</style>
