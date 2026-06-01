<template>
  <div class="login-page">
    <el-card class="login-card">
      <h2 style="text-align:center;margin-bottom:20px;">{{ isRegister ? '注册' : '登录' }}</h2>

      <el-form :model="form" label-width="0" @submit.prevent="handleSubmit">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" size="large" clearable />
        </el-form-item>

        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password />
        </el-form-item>

        <el-form-item v-if="isRegister">
          <el-input v-model="form.email" placeholder="邮箱（选填）" size="large" clearable />
        </el-form-item>

        <!-- 注册时填写更多个人信息 -->
        <template v-if="isRegister">
          <el-form-item>
            <el-select v-model="form.gender" placeholder="性别（选填）" size="large" style="width:100%" clearable>
              <el-option label="男" value="男" />
              <el-option label="女" value="女" />
              <el-option label="其他" value="其他" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-input-number v-model="form.age" :min="1" :max="150" placeholder="年龄（选填）" size="large" style="width:100%" controls-position="right" />
          </el-form-item>
          <el-form-item>
            <el-input v-model="form.job" placeholder="当前工作（选填）" size="large" clearable />
          </el-form-item>
          <el-form-item>
            <el-input v-model="form.motto" placeholder="座右铭（选填）" size="large" maxlength="200" show-word-limit clearable />
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" size="large" style="width:100%" :loading="loading" @click="handleSubmit">
            {{ isRegister ? '注册' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <p v-if="errorMsg" style="color:#f56c6c;text-align:center;">{{ errorMsg }}</p>

      <div style="text-align:center;margin-top:12px;">
        <el-link type="primary" @click="isRegister = !isRegister; errorMsg = ''">
          {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
        </el-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import { login } from '../stores/auth'

const router = useRouter()
const isRegister = ref(false)
const loading = ref(false)
const errorMsg = ref('')
const form = reactive({
  username: '',
  password: '',
  email: '',
  gender: '',
  age: null,
  job: '',
  motto: '',
})

async function handleSubmit() {
  if (!form.username || !form.password) {
    errorMsg.value = '请填写用户名和密码'
    return
  }
  loading.value = true
  errorMsg.value = ''

  try {
    const url = isRegister.value ? '/register' : '/login'
    const res = await api.post(url, form)
    const data = res.data

    if (isRegister.value) {
      // 注册成功后自动登录
      const loginRes = await api.post('/login', { username: form.username, password: form.password })
      login(loginRes.data.token, loginRes.data.user)
    } else {
      login(data.token, data.user)
    }
    router.push('/')
  } catch (e) {
    errorMsg.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  padding-top: 40px;
}
.login-card {
  width: 440px;
}
</style>
