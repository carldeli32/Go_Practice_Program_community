<template>
  <div id="app-layout">
    <!-- 顶栏 -->
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
          <span class="username">{{ auth.user?.username }}</span>
          <el-button @click="handleLogout" size="small">退出</el-button>
        </template>
        <template v-else>
          <el-button @click="$router.push('/login')" type="primary" size="small">登录</el-button>
        </template>
      </div>
    </el-menu>

    <!-- 页面内容 -->
    <div class="main-container">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuth, logout } from './stores/auth'

const auth = useAuth()
const router = useRouter()

function handleLogout() {
  logout()
  router.push('/')
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { background: #f5f7fa; font-family: 'Helvetica Neue', Arial, sans-serif; }

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 56px;
}

.navbar .el-menu-item { border-bottom: none !important; }

.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username { color: #606266; font-size: 14px; }

.main-container {
  max-width: 800px;
  margin: 24px auto;
  padding: 0 16px;
}
</style>
