<template>
  <div class="min-h-screen bg-bg text-fg font-body overflow-x-hidden">
    <!-- INFO BAR -->
    <div class="sticky top-0 z-100 bg-bg-top border-b border-border">
      <div class="max-w-[1440px] mx-auto px-3 sm:px-5 h-9 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="w-1.5 h-1.5 rounded-full bg-primary animate-pulse"></span>
          <span class="text-[10px] text-muted-fg font-mono tracking-[2px] uppercase hidden sm:inline">Thrum · v3.0</span>
        </div>
        <div class="flex items-center gap-1.5">
          <template v-if="auth.token">
            <button @click="$router.push('/messages')" class="flex items-center gap-1 px-1.5 sm:px-2 h-9 text-muted-fg hover:text-primary text-[11px] font-mono transition-colors">
              <Bell :size="12" />
              <span v-if="unreadCount" class="min-w-[16px] h-[16px] rounded-full bg-danger text-white text-[9px] font-bold flex items-center justify-center px-1">{{ unreadCount }}</span>
              <span class="hidden sm:inline">消息</span>
            </button>
            <span class="w-px h-3.5 bg-border hidden sm:block"></span>
            <button @click="$router.push('/friends')" class="hidden sm:flex items-center gap-1 px-2 h-9 text-muted-fg hover:text-primary text-[11px] font-mono transition-colors">
              <Users :size="12" />好友
            </button>
            <span class="w-px h-3.5 bg-border hidden sm:block"></span>
            <button @click="$router.push('/friends?tab=following')" class="hidden sm:flex items-center gap-1 px-2 h-9 text-muted-fg hover:text-success text-[11px] font-mono transition-colors">
              <Rss :size="12" />关注
            </button>
            <span class="w-px h-3.5 bg-border hidden sm:block"></span>
            <router-link :to="`/user/${auth.user?.id}`" class="flex items-center gap-1.5 px-2 h-9 no-underline hover:bg-secondary rounded-full transition-colors">
              <span class="w-[22px] h-[22px] rounded-full bg-gradient-to-br from-primary to-accent flex items-center justify-center text-[10px] font-bold text-primary-foreground">{{ auth.user?.username?.charAt(0) }}</span>
              <span class="text-xs text-fg font-medium hidden sm:inline">{{ auth.user?.username }}</span>
            </router-link>
            <button @click="handleLogout" class="flex items-center gap-1 px-2 h-9 text-muted-fg hover:text-danger text-[11px] font-mono transition-colors"><LogOut :size="12" /></button>
            <!-- Mobile menu toggle -->
            <button @click="mobileMenu = !mobileMenu" class="sm:hidden flex items-center px-2 h-9 text-muted-fg text-lg">☰</button>
          </template>
          <template v-else>
            <button @click="$router.push('/login')" class="text-[11px] text-primary hover:underline font-mono">登录</button>
          </template>
        </div>
      </div>
    </div>

    <!-- HEADER -->
    <header class="bg-bg-top/95 backdrop-blur-sm border-b border-border">
      <div class="max-w-[1440px] mx-auto px-3 sm:px-5 py-3 flex items-center gap-3 sm:gap-5">
        <div @click="$router.push('/')" class="flex items-center gap-2.5 cursor-pointer flex-shrink-0">
          <div class="w-8 sm:w-10 h-8 sm:h-10 border-2 border-primary flex items-center justify-center text-base sm:text-lg" style="clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)">⚡</div>
          <div>
            <h1 class="text-lg sm:text-xl font-black text-fg font-heading tracking-[3px] leading-none" style="text-shadow: 0 0 15px rgba(0,200,255,0.5)">Thrum</h1>
            <p class="text-[8px] sm:text-[9px] text-muted-fg font-mono tracking-[3px] uppercase hidden sm:block">连接每一个思想节点</p>
          </div>
        </div>
        <nav class="hidden md:flex items-center gap-1 flex-1 justify-center">
          <button v-for="item in navItems" :key="item.path" @click="$router.push(item.path)"
            :class="['px-3 sm:px-4 py-2 text-sm font-semibold font-heading tracking-wide transition-all border-b-2', $route.path === item.path ? 'text-primary border-primary' : 'text-muted-fg border-transparent hover:text-fg']">
            {{ item.label }}
          </button>
        </nav>
        <div class="flex items-center gap-2 sm:gap-4 flex-shrink-0">
          <div class="hidden xl:flex items-center gap-4">
            <div class="text-right"><div class="text-[10px] text-muted-fg font-mono">在线</div><div class="text-sm font-bold text-primary font-heading">--</div></div>
            <span class="w-px h-8 bg-border"></span>
            <div class="text-right"><div class="text-[10px] text-muted-fg font-mono">今日帖</div><div class="text-sm font-bold text-success font-heading">--</div></div>
          </div>
          <template v-if="auth.token">
            <button @click="goCreate" class="bg-primary text-bg border-0 rounded-md px-3 sm:px-4 py-1.5 text-xs font-bold font-heading tracking-wide cursor-pointer hover:brightness-110 transition-all">✏️ <span class="hidden sm:inline">发帖</span></button>
            <button v-if="auth.user?.is_admin" @click="$router.push('/admin')" class="border border-warning/40 text-warning bg-transparent rounded-md px-2 sm:px-3 py-1.5 text-xs font-bold font-heading tracking-wide cursor-pointer hover:bg-warning/10 transition-all">⚙️ <span class="hidden sm:inline">管理</span></button>
          </template>
        </div>
      </div>

      <!-- Mobile nav dropdown -->
      <div v-if="mobileMenu" class="md:hidden border-t border-border bg-bg-top px-3 py-2">
        <nav class="flex flex-col gap-1">
          <button v-for="item in navItems" :key="item.path" @click="$router.push(item.path); mobileMenu = false"
            :class="['px-4 py-2.5 text-sm font-semibold font-heading tracking-wide text-left rounded', $route.path === item.path ? 'text-primary bg-primary/10' : 'text-muted-fg hover:text-fg']">
            {{ item.label }}
          </button>
          <button @click="$router.push('/friends'); mobileMenu = false" class="px-4 py-2.5 text-sm font-semibold text-muted-fg hover:text-fg font-heading tracking-wide text-left rounded sm:hidden">好友</button>
          <button @click="$router.push('/friends?tab=following'); mobileMenu = false" class="px-4 py-2.5 text-sm font-semibold text-muted-fg hover:text-fg font-heading tracking-wide text-left rounded sm:hidden">关注</button>
        </nav>
      </div>
    </header>

    <!-- ANNOUNCEMENT -->
    <div v-if="announcement" class="bg-gradient-to-r from-[#001428] via-[#002040] to-[#001428] border-b border-primary/30">
      <div class="max-w-[1440px] mx-auto px-3 sm:px-5 h-[38px] flex items-center gap-3">
        <span class="px-2.5 py-0.5 bg-primary/10 border border-primary/30 text-[10px] text-primary font-bold font-mono tracking-[2px] flex-shrink-0">📢 公告</span>
        <span class="text-xs text-secondary-fg overflow-hidden whitespace-nowrap text-ellipsis">{{ announcement }}</span>
      </div>
    </div>

    <!-- MAIN -->
    <main class="max-w-[1440px] mx-auto px-2 sm:px-5 py-3 sm:py-5 overflow-y-auto" style="height: calc(100vh - 140px)">
      <router-view v-slot="{ Component }">
        <keep-alive :include="['Home', 'Messages', 'Friends']">
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </main>

    <!-- FOOTER -->
    <footer class="border-t border-border bg-bg-top">
      <div class="max-w-[1440px] mx-auto px-5 py-4 flex items-center justify-between text-xs text-muted-fg font-mono">
        <span class="font-heading text-sm font-bold text-fg tracking-wider">⚡ Thrum</span>
        <span class="hidden sm:inline">© 2026 · 社区项目</span>
        <div class="flex gap-5"><span class="cursor-pointer hover:text-primary transition-colors">关于</span><span class="cursor-pointer hover:text-primary transition-colors">帮助</span></div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, provide } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth, logout } from './stores/auth'
import api from './api'
import { Bell, Users, LogOut, Rss } from 'lucide-vue-next'

const auth = useAuth()
const router = useRouter()
const unreadCount = ref(0)
const announcement = ref('')
const mobileMenu = ref(false)
const sharedCategories = ref([])
provide('sharedCategories', sharedCategories)

const navItems = [
  { path: '/', label: '首页' },
  { path: '/messages', label: '私信' },
  { path: '/friends', label: '好友' },
]

function goCreate() { router.push('/create') }
function handleLogout() { logout(); router.push('/') }

async function fetchUnread() {
  if (!auth.token) return
  try { const r = await api.get('/messages/unread-count'); unreadCount.value = r.data.data.count } catch (e) {}
}
async function fetchAnnouncement() {
  try { const r = await api.get('/announcement'); announcement.value = r.data.data.content } catch (e) {}
}

watch(() => router.currentRoute.value, () => { fetchUnread(); mobileMenu.value = false })
onMounted(() => { fetchUnread(); fetchAnnouncement(); fetchSharedCategories() })
const timer = setInterval(fetchUnread, 30000)
onUnmounted(() => clearInterval(timer))

async function fetchSharedCategories() {
  try { const r = await api.get('/categories'); sharedCategories.value = r.data.data.categories } catch (e) {}
}
</script>
