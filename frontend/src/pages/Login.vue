<template>
  <div class="flex justify-center pt-5">
    <div class="w-[420px] bg-card border border-border rounded-lg p-8">
      <h2 class="text-center text-xl font-bold text-fg font-heading tracking-wider mb-6">{{ isRegister ? '注册' : '登录' }}</h2>

      <div class="space-y-4">
        <input v-model="form.username" placeholder="用户名" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors" />
        <input v-model="form.password" type="password" placeholder="密码" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors" />

        <template v-if="isRegister">
          <input v-model="form.email" placeholder="邮箱（选填）" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors" />
          <select v-model="form.gender" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg font-mono outline-none focus:border-primary transition-colors">
            <option value="">性别（选填）</option>
            <option value="男">男</option><option value="女">女</option><option value="其他">其他</option>
          </select>
          <input v-model.number="form.age" type="number" min="1" max="150" placeholder="年龄（选填）" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors" />
          <input v-model="form.job" placeholder="当前工作（选填）" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors" />
          <input v-model="form.motto" placeholder="座右铭（选填）" maxlength="200" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors" />
        </template>

        <button @click="handleSubmit" :disabled="loading" class="w-full h-10 bg-primary text-bg border-0 rounded-md text-sm font-bold font-heading tracking-wider cursor-pointer hover:brightness-110 transition-all disabled:opacity-50">
          {{ loading ? '处理中...' : (isRegister ? '注册' : '登录') }}
        </button>
      </div>

      <p v-if="errorMsg" class="text-center text-danger text-xs mt-3">{{ errorMsg }}</p>

      <div class="text-center mt-4">
        <button @click="isRegister = !isRegister; errorMsg = ''" class="text-xs text-primary hover:underline bg-transparent border-0 cursor-pointer font-mono">
          {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import { login } from '../stores/auth'

const router = useRouter()
const isRegister = ref(false); const loading = ref(false); const errorMsg = ref('')
const form = reactive({ username: '', password: '', email: '', gender: '', age: null, job: '', motto: '' })

async function handleSubmit() {
  if (!form.username || !form.password) { errorMsg.value = '请填写用户名和密码'; return }
  loading.value = true; errorMsg.value = ''
  try {
    const url = isRegister.value ? '/register' : '/login'
    const res = await api.post(url, form)
    if (isRegister.value) {
      const lr = await api.post('/login', { username: form.username, password: form.password })
      login(lr.data.data.token, lr.data.data.user)
    } else {
      login(res.data.data.token, res.data.data.user)
    }
    router.push('/')
  } catch (e) { errorMsg.value = e.message } finally { loading.value = false }
}
</script>
