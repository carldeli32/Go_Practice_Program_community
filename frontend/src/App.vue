<template>
  <div id="app-layout">
    <el-menu mode="horizontal" :router="true" class="navbar">
      <el-menu-item index="/">
        <el-icon><ChatDotRound /></el-icon>
        <span>小墨社区</span>
      </el-menu-item>
      <div class="nav-right">
        <template v-if="auth.token">
          <el-button @click="$router.push('/create')" type="primary" size="small">
            <el-icon><Edit /></el-icon> 发帖
          </el-button>
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" :max="99">
            <el-button @click="$router.push('/messages')" size="small">
              <el-icon><Message /></el-icon> 私信
            </el-button>
          </el-badge>
          <el-button @click="$router.push('/friends')" size="small">
            <el-icon><UserFilled /></el-icon> 好友
          </el-button>
          <el-button v-if="auth.user?.is_admin" @click="$router.push('/admin')" type="warning" size="small">
            <el-icon><Setting /></el-icon> 管理
          </el-button>
          <router-link :to="`/user/${auth.user?.id}`" class="username">
            {{ auth.user?.username }}
          </router-link>
          <el-button @click="handleLogout" size="small">退出</el-button>
        </template>
        <template v-else>
          <el-button @click="$router.push('/login')" type="primary" size="small">登录</el-button>
        </template>
      </div>
    </el-menu>

    <div class="main-container">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth, logout } from './stores/auth'
import api from './api'

const auth = useAuth()
const router = useRouter()
const unreadCount = ref(0)

function handleLogout() {
  logout()
  router.push('/')
}

async function fetchUnread() {
  if (!auth.token) return
  try {
    const res = await api.get('/messages/unread-count')
    unreadCount.value = res.data.count
  } catch (e) { /* ignore */ }
}

// 每次路由跳转时刷新未读数（解决进站红点不显示的问题）
watch(() => router.currentRoute.value, () => {
  fetchUnread()
})

onMounted(() => {
  fetchUnread()
})

// 每 30 秒定时检查（降低频率避免频繁请求）
const timer = setInterval(fetchUnread, 30000)
onUnmounted(() => clearInterval(timer))
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { background: #f5f7fa; font-family: 'Helvetica Neue', Arial, sans-serif; }
.navbar { display: flex; justify-content: space-between; align-items: center; padding: 0 20px; height: 56px; }
.navbar .el-menu-item { border-bottom: none !important; }
.nav-right { display: flex; align-items: center; gap: 12px; }
.username { color: #409eff; font-size: 14px; text-decoration: none; }
.username:hover { text-decoration: underline; }
.main-container { max-width: 800px; margin: 24px auto; padding: 0 16px; }
</style>
