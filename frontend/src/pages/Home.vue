<template>
  <div class="grid grid-cols-1 lg:grid-cols-[240px_1fr_240px] gap-5 items-start">
    <!-- LEFT SIDEBAR -->
    <aside class="space-y-5">
      <div>
        <h3 class="text-xs font-bold font-heading tracking-wide uppercase mb-2.5 pl-2.5 border-l-[3px] border-success text-success">主题频道</h3>
        <input v-model="channelSearch" placeholder="搜索频道..." class="w-full h-8 bg-input border border-border rounded-md px-2.5 text-xs text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors mb-2">
        <div class="bg-card border border-border rounded-lg overflow-hidden">
          <div v-for="cat in filteredCategories" :key="cat.id"
            :class="['flex items-center gap-2.5 px-3 py-2.5 cursor-pointer transition-colors border-b border-border last:border-b-0 hover:bg-card-hover', { 'bg-primary/8': activeCategory?.id === cat.id }]"
            @click="enterCategory(cat)">
            <span class="text-lg flex-shrink-0">{{ getIcon(cat.name) }}</span>
            <div class="flex-1 min-w-0">
              <span class="text-[13px] font-semibold text-fg block">{{ cat.name }}</span>
              <span class="text-[10px] text-muted-fg font-mono block">{{ cat.description }}</span>
            </div>
            <span v-if="activeCategory?.id === cat.id" class="w-1.5 h-1.5 rounded-full bg-primary shadow-[0_0_6px_var(--color-primary)] flex-shrink-0"></span>
          </div>
        </div>
      </div>
      <div>
        <h3 class="text-xs font-bold font-heading tracking-wide uppercase mb-2.5 pl-2.5 border-l-[3px] border-primary text-primary">热门标签</h3>
        <div class="flex flex-wrap gap-1.5">
          <button v-for="cat in categories" :key="cat.id" @click="enterCategory(cat)"
            class="bg-card border border-border rounded px-2.5 py-1 text-[10px] text-muted-fg font-mono cursor-pointer hover:border-primary hover:text-primary transition-all">#{{ cat.name }}</button>
        </div>
      </div>
    </aside>

    <!-- CENTER FEED -->
    <main class="min-w-0 space-y-3.5">
      <!-- Compose -->
      <div v-if="auth.token" @click="$router.push('/create')" class="flex items-center gap-3 px-4 py-3 bg-card border border-border rounded-lg cursor-pointer hover:border-border-strong transition-colors">
        <span class="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-accent flex items-center justify-center text-xs font-bold text-bg flex-shrink-0">{{ auth.user?.username?.charAt(0) }}</span>
        <span class="flex-1 text-[13px] text-muted-fg">分享你的想法...</span>
        <span class="text-[11px] text-primary font-heading tracking-wide">✏️ 发帖</span>
      </div>

      <!-- Header -->
      <div class="flex items-center justify-between pb-2.5 border-b border-border">
        <div v-if="activeCategory" class="flex items-center gap-3">
          <span @click="backToCategories" class="text-xs text-muted-fg font-mono cursor-pointer hover:text-primary">← 返回</span>
          <span class="text-base font-bold text-fg">{{ getIcon(activeCategory.name) }} {{ activeCategory.name }}</span>
        </div>
        <div v-else class="flex gap-1">
          <span class="text-[13px] font-semibold font-heading tracking-wide text-primary border-b-2 border-primary px-3 py-1.5">📋 最新帖子</span>
        </div>
        <input v-if="activeCategory" v-model="postSearch" @change="searchPosts" placeholder="搜索帖子..." class="w-[220px] h-8 bg-input border border-border rounded-md px-3 text-xs text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors">
      </div>

      <!-- Empty -->
      <div v-if="posts.length === 0 && !loading" class="text-center text-muted-fg py-15 text-sm">还没有帖子，快来发第一篇吧~</div>

      <!-- Post Cards -->
      <article v-for="post in posts" :key="post.id" @click="$router.push(`/post/${post.id}`)"
        class="bg-card border border-border rounded-lg p-4 cursor-pointer transition-all hover:border-border-strong hover:shadow-[0_0_20px_var(--color-primary-glow)]">
        <div class="flex items-center justify-between mb-2.5">
          <div class="flex items-center gap-2">
            <span class="w-7 h-7 rounded-full flex items-center justify-center text-[11px] font-bold text-bg flex-shrink-0" :style="{ background: `linear-gradient(135deg, ${avatarColor(post.user?.username)}, ${avatarColor(post.user?.username)}88)` }">{{ post.user?.username?.charAt(0) || '?' }}</span>
            <div>
              <span class="text-xs font-semibold text-fg">{{ post.user?.username || '匿名' }}</span>
              <span class="text-[9px] text-muted-fg font-mono bg-primary/10 border border-primary/20 rounded px-1 ml-1">Lv.1</span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-[9px] text-primary font-mono bg-primary/10 border border-primary/20 rounded px-1.5 py-0.5">{{ post.category?.name || '综合讨论' }}</span>
            <span class="text-[10px] text-muted-fg font-mono">{{ fmtTime(post.created_at) }}</span>
          </div>
        </div>
        <h3 class="text-[15px] font-semibold text-fg font-heading tracking-wide mb-1.5 group-hover:text-primary transition-colors">{{ post.title }}</h3>
        <p class="text-xs text-muted-fg leading-relaxed line-clamp-2" v-html="previewMD(post.content)" />
      </article>

      <!-- Pagination -->
      <div v-if="total > pageSize" class="flex justify-center mt-5 gap-1">
        <button :disabled="page <= 1" @click="page--; fetchPosts()" class="px-3 py-1.5 bg-card border border-border rounded text-xs text-fg font-mono cursor-pointer disabled:opacity-40 hover:border-primary transition-colors">‹</button>
        <button v-for="p in pageRange" :key="p" @click="page = p; fetchPosts()"
          :class="['px-3 py-1.5 rounded text-xs font-mono cursor-pointer transition-all', p === page ? 'bg-primary text-bg border border-primary' : 'bg-card border border-border text-fg hover:border-primary']">{{ p }}</button>
        <button :disabled="page >= maxPage" @click="page++; fetchPosts()" class="px-3 py-1.5 bg-card border border-border rounded text-xs text-fg font-mono cursor-pointer disabled:opacity-40 hover:border-primary transition-colors">›</button>
      </div>
    </main>

    <!-- RIGHT SIDEBAR -->
    <aside class="space-y-5">
      <!-- User Search -->
      <div class="relative">
        <div class="relative">
          <Search :size="12" class="absolute left-2.5 top-1/2 -translate-y-1/2 text-muted-fg" />
          <input v-model="userSearch" @focus="onUserFocus" @blur="onUserBlur" placeholder="搜索用户..."
            class="w-full h-9 bg-input border border-border rounded-md pl-8 pr-3 text-xs text-fg placeholder-muted-fg font-mono outline-none focus:border-primary transition-colors" />
        </div>
        <div v-if="userSearch && showUserDropdown" class="absolute z-50 left-0 right-0 mt-1 bg-card border border-border rounded-lg overflow-hidden shadow-xl max-h-[240px] overflow-y-auto">
          <div v-if="userLoading" class="px-3 py-2 text-[10px] text-muted-fg font-mono text-center">搜索中...</div>
          <div v-else-if="userResults.length === 0" class="px-3 py-2 text-[10px] text-muted-fg font-mono text-center">无匹配用户</div>
          <div v-for="u in userResults" :key="u.id"
            @mousedown.prevent="$router.push(`/user/${u.id}`); userSearch = ''; userResults = []"
            class="flex items-center gap-2 px-3 py-2 cursor-pointer hover:bg-card-hover transition-colors border-b border-border last:border-b-0">
            <span class="w-6 h-6 rounded-full bg-primary flex items-center justify-center text-[10px] font-bold text-bg flex-shrink-0">{{ u.username?.charAt(0) }}</span>
            <div class="min-w-0">
              <span class="text-xs font-semibold text-fg block">{{ u.username }}</span>
              <span class="text-[10px] text-muted-fg font-mono truncate block">{{ u.motto || '这个人很懒，什么都没写' }}</span>
            </div>
          </div>
        </div>
      </div>
      <div>
        <h3 class="text-xs font-bold font-heading tracking-wide uppercase mb-2.5 pl-2.5 border-l-[3px] border-primary text-primary">社区统计</h3>
        <div class="grid grid-cols-3 gap-2">
          <div class="bg-card border border-border rounded-md p-3 text-center"><div class="text-lg font-bold text-primary font-heading">{{ total }}</div><div class="text-[9px] text-muted-fg font-mono mt-0.5">总帖数</div></div>
          <div class="bg-card border border-border rounded-md p-3 text-center"><div class="text-lg font-bold text-success font-heading">--</div><div class="text-[9px] text-muted-fg font-mono mt-0.5">在线</div></div>
          <div class="bg-card border border-border rounded-md p-3 text-center"><div class="text-lg font-bold text-warning font-heading">{{ categories.length }}</div><div class="text-[9px] text-muted-fg font-mono mt-0.5">频道</div></div>
        </div>
      </div>
      <div>
        <h3 class="text-xs font-bold font-heading tracking-wide uppercase mb-2.5 pl-2.5 border-l-[3px] border-[#f72585] text-[#f72585]">热门推荐</h3>
        <div class="space-y-1.5">
          <div v-if="hotPosts.length === 0" class="text-xs text-muted-fg text-center py-5">暂无数据</div>
          <div v-for="(p, i) in hotPosts" :key="p.id" @click="$router.push(`/post/${p.id}`)"
            class="flex items-start gap-2.5 p-2.5 bg-card border border-border rounded-md cursor-pointer hover:border-border-strong transition-colors">
            <span :class="['text-sm font-bold font-heading flex-shrink-0 w-[22px] text-center', i === 0 ? 'text-warning' : i === 1 ? 'text-secondary-fg' : 'text-[#ff6b35]']">{{ i + 1 }}</span>
            <div class="flex-1 min-w-0">
              <span class="text-xs text-fg block truncate">{{ p.title }}</span>
              <span class="text-[10px] text-muted-fg">{{ p.user?.username || '匿名' }}</span>
            </div>
          </div>
        </div>
      </div>
      <div>
        <h3 class="text-xs font-bold font-heading tracking-wide uppercase mb-2.5 pl-2.5 border-l-[3px] border-warning text-warning">快速入口</h3>
        <div class="space-y-1.5">
          <button v-if="auth.token" @click="$router.push('/create')" class="w-full text-left bg-card border border-border rounded-md px-3.5 py-2.5 text-xs text-fg font-heading tracking-wide cursor-pointer hover:border-primary hover:bg-primary/5 transition-all">✏️ 发布新帖</button>
          <button v-if="auth.token" @click="$router.push('/messages')" class="w-full text-left bg-card border border-border rounded-md px-3.5 py-2.5 text-xs text-fg font-heading tracking-wide cursor-pointer hover:border-primary hover:bg-primary/5 transition-all">💬 我的私信</button>
          <button v-if="auth.token" @click="$router.push('/friends')" class="w-full text-left bg-card border border-border rounded-md px-3.5 py-2.5 text-xs text-fg font-heading tracking-wide cursor-pointer hover:border-primary hover:bg-primary/5 transition-all">👥 好友列表</button>
          <button v-if="auth.user?.is_admin" @click="$router.push('/admin')" class="w-full text-left border border-warning/30 rounded-md px-3.5 py-2.5 text-xs text-warning font-heading tracking-wide cursor-pointer hover:bg-warning/5 transition-all">⚙️ 管理后台</button>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
import api from '../api'
import { Search } from 'lucide-vue-next'
import { marked } from 'marked'

marked.setOptions({ breaks: true, gfm: true })

const route = useRoute(); const router = useRouter(); const auth = useAuth()
const categories = ref([]); const posts = ref([]); const hotPosts = ref([])
const total = ref(0); const page = ref(1); const pageSize = ref(10)
const loading = ref(false); const postSearch = ref(''); const activeCategory = ref(null)
const channelSearch = ref('')
const userSearch = ref(''); const userResults = ref([]); const userLoading = ref(false); const showUserDropdown = ref(false)
let userDebounce = null

const filteredCategories = computed(() => {
  if (!channelSearch.value) return categories.value
  const q = channelSearch.value.toLowerCase()
  return categories.value.filter(c => c.name.toLowerCase().includes(q) || c.description.toLowerCase().includes(q))
})

const maxPage = computed(() => Math.ceil(total.value / pageSize.value) || 1)
const pageRange = computed(() => { const r = []; const max = maxPage.value; const p = page.value; for (let i = Math.max(1, p - 2); i <= Math.min(max, p + 2); i++) r.push(i); return r })

async function fetchCategories() {
  try { const r = await api.get('/categories'); categories.value = r.data.data.categories } catch (e) {}
}
async function fetchPosts() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (activeCategory.value) params.category_id = activeCategory.value.id
    if (postSearch.value) params.q = postSearch.value
    const r = await api.get('/posts', { params }); posts.value = r.data.data.posts; total.value = r.data.data.total
  } catch (e) {} finally { loading.value = false }
}
async function fetchHotPosts() {
  try { const r = await api.get('/posts', { params: { page: 1, page_size: 5 } }); hotPosts.value = r.data.data.posts || [] } catch (e) {}
}
function onUserFocus() { showUserDropdown.value = true; if (userSearch.value) doUserSearch() }
function onUserBlur() { setTimeout(() => showUserDropdown.value = false, 200) }
function doUserSearch() {
  const q = userSearch.value.trim()
  if (!q) { userResults.value = []; return }
  userLoading.value = true
  api.get('/users', { params: { q } })
    .then(r => userResults.value = r.data.data.users || [])
    .catch(() => userResults.value = [])
    .finally(() => userLoading.value = false)
}

function enterCategory(cat) { activeCategory.value = cat; postSearch.value = ''; page.value = 1; router.replace({ query: { category_id: cat.id } }); fetchPosts() }
function backToCategories() { activeCategory.value = null; router.replace({ query: {} }); fetchPosts() }
function searchPosts() { page.value = 1; fetchPosts() }

function getIcon(n) { const m = { '综合讨论':'💬','技术交流':'💻','军事纵横':'⚔️','历史长廊':'📜','文学艺术':'🎨','生活杂谈':'🌻' }; return m[n] || '📌' }
function avatarColor(n) { const c = ['#00c8ff','#7c3aed','#00ff94','#ff6b35','#f72585','#ffd60a']; let h = 0; for (let ch of (n || '?')) h = ch.charCodeAt(0) + ((h << 5) - h); return c[Math.abs(h) % c.length] }
function fmtTime(s) { if (!s) return ''; const d = new Date(s), n = new Date(), diff = n - d; if (diff < 6e4) return '刚刚'; if (diff < 36e5) return Math.floor(diff/6e4)+'分钟前'; if (diff < 864e5) return Math.floor(diff/36e5)+'小时前'; return d.toLocaleString('zh-CN') }
function truncate(s, n) { return s && s.length > n ? s.slice(0, n) + '...' : s }
function previewMD(text) {
  if (!text) return ''
  const html = marked.parse(text)
  // Strip all HTML tags for plain text preview
  return html.replace(/<[^>]*>/g, '').slice(0, 200)
}

watch(() => route.query.category_id, (v) => { if (v) { const c = categories.value.find(x => x.id === Number(v)); if (c) { activeCategory.value = c; fetchPosts() } } else activeCategory.value = null })
watch(userSearch, () => {
  clearTimeout(userDebounce)
  if (!userSearch.value.trim()) { userResults.value = []; return }
  userDebounce = setTimeout(doUserSearch, 300)
})
onMounted(async () => { await fetchCategories(); fetchPosts(); fetchHotPosts(); const cid = route.query.category_id; if (cid) { const c = categories.value.find(x => x.id === Number(cid)); if (c) { activeCategory.value = c; fetchPosts() } } })
</script>
