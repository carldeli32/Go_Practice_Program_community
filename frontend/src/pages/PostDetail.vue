<template>
  <div class="max-w-[800px] mx-auto">
    <!-- Toast -->
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono shadow-lg transition-all', toast.type.value === 'success' ? 'bg-success/20 border border-success/40 text-success' : toast.type.value === 'error' ? 'bg-danger/20 border border-danger/40 text-danger' : 'bg-warning/20 border border-warning/40 text-warning']">{{ toast.msg.value }}</div>

    <div v-if="post.id">
      <!-- Post Card -->
      <div class="bg-card border border-border rounded-lg p-5 mb-5">
        <div class="flex items-center gap-2 mb-2 text-xs text-muted-fg">
          <router-link :to="`/user/${post.user?.id}`" class="font-semibold text-primary hover:underline">@{{ post.user?.username || '匿名' }}</router-link>
          <span>·</span>
          <span class="font-mono">{{ fmtDate(post.created_at) }}</span>
          <span>·</span>
          <span class="text-[11px] text-primary font-mono bg-primary/10 border border-primary/20 rounded px-1.5">{{ getIcon(post.category?.name) }} {{ post.category?.name }}</span>
        </div>
        <h1 class="text-xl font-bold text-fg font-heading tracking-wide mb-3">{{ post.title }}</h1>
        <p class="text-sm text-fg/90 leading-relaxed whitespace-pre-wrap">{{ post.content }}</p>

        <div v-if="canModifyPost" class="flex gap-2 mt-4 pt-3 border-t border-border">
          <button @click="$router.push({path:'/create', query:{edit:post.id}})" class="px-3 py-1.5 bg-card border border-border rounded-md text-xs text-fg font-mono cursor-pointer hover:border-primary transition-colors">编辑</button>
          <button @click="handleDeletePost" class="px-3 py-1.5 bg-card border border-danger/40 rounded-md text-xs text-danger font-mono cursor-pointer hover:bg-danger/10 transition-colors">删除</button>
        </div>
      </div>

      <!-- Comments -->
      <div class="bg-card border border-border rounded-lg p-5">
        <h3 class="text-sm font-bold text-fg font-heading tracking-wide mb-4">💬 评论 ({{ comments.length }})</h3>

        <div v-if="comments.length === 0" class="text-center text-muted-fg py-8 text-sm">还没有评论，来抢沙发吧~</div>

        <div v-for="c in comments" :key="c.id" class="py-3 border-b border-border last:border-b-0">
          <div class="flex justify-between mb-1.5">
            <router-link :to="`/user/${c.user?.id}`" class="text-xs font-semibold text-primary hover:underline">@{{ c.user?.username || '匿名' }}</router-link>
            <span class="text-[11px] text-muted-fg font-mono">{{ fmtDate(c.created_at) }}</span>
          </div>
          <p v-if="editingId !== c.id" class="text-xs text-fg/85 leading-relaxed mb-1.5">{{ c.content }}</p>
          <div v-else class="mb-1.5">
            <textarea v-model="editText" rows="3" class="w-full bg-input border border-border rounded-md p-2 text-xs text-fg placeholder-muted-fg outline-none focus:border-primary resize-none"></textarea>
            <div class="flex gap-2 mt-1.5">
              <button @click="saveEditComment(c.id)" class="px-3 py-1 bg-primary text-bg rounded text-[11px] font-bold font-heading tracking-wide cursor-pointer">保存</button>
              <button @click="cancelEdit" class="px-3 py-1 bg-card border border-border rounded text-[11px] text-muted-fg font-mono cursor-pointer">取消</button>
            </div>
          </div>
          <div v-if="c.can_edit || c.can_delete" class="flex gap-2 mt-1">
            <button v-if="c.can_edit" @click="startEdit(c)" class="text-[10px] text-muted-fg hover:text-primary font-mono cursor-pointer bg-transparent border-0">编辑</button>
            <button v-if="c.can_delete" @click="delComment(c.id)" class="text-[10px] text-muted-fg hover:text-danger font-mono cursor-pointer bg-transparent border-0">删除</button>
          </div>
        </div>

        <!-- New Comment -->
        <div v-if="isLoggedIn" class="mt-4 pt-4 border-t border-border">
          <textarea v-model="commentText" rows="3" placeholder="写下你的评论..." class="w-full bg-input border border-border rounded-md p-3 text-sm text-fg placeholder-muted-fg outline-none focus:border-primary resize-none"></textarea>
          <button @click="doComment" class="mt-2 px-5 py-2 bg-primary text-bg rounded-md text-xs font-bold font-heading tracking-wide cursor-pointer hover:brightness-110 transition-all">发表</button>
        </div>
        <div v-else class="text-center mt-4 pt-4 border-t border-border">
          <router-link to="/login" class="text-xs text-primary hover:underline">登录后发表评论</router-link>
        </div>
      </div>

      <!-- Confirm Dialog -->
      <div v-if="confirmOpen" class="fixed inset-0 z-500 flex items-center justify-center bg-black/50" @click.self="confirmOpen = false">
        <div class="bg-card border border-border rounded-lg p-6 max-w-[360px] w-full mx-4 shadow-2xl">
          <p class="text-sm text-fg mb-4">{{ confirmMsg }}</p>
          <div class="flex justify-end gap-2">
            <button @click="confirmOpen = false" class="px-4 py-1.5 bg-card border border-border rounded text-xs text-muted-fg font-mono cursor-pointer">取消</button>
            <button @click="onConfirm(); confirmOpen = false" class="px-4 py-1.5 bg-danger text-white rounded text-xs font-bold font-heading tracking-wide cursor-pointer">确定</button>
          </div>
        </div>
      </div>
    </div>
    <div v-else-if="!loading" class="text-center text-muted-fg py-20 text-sm">帖子不存在</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth, isLoggedIn } from '../stores/auth'
import api from '../api'
import { useToast } from '../composables/toast'

const route = useRoute(); const router = useRouter(); const auth = useAuth()
const loading = ref(true); const post = ref({}); const comments = ref([])
const commentText = ref(''); const editingId = ref(0); const editText = ref('')
const confirmOpen = ref(false); const confirmMsg = ref(''); let confirmCb = null
const toast = useToast()

const canModifyPost = computed(() => auth.user && post.value.id && (auth.user.id === post.value.user_id || auth.user.is_admin))

async function load() {
  loading.value = true
  try {
    const [p, c] = await Promise.all([
      api.get(`/posts/${route.params.id}`),
      api.get(`/posts/${route.params.id}/comments`),
    ])
    post.value = p.data.data.post
    comments.value = c.data.data.comments || []
  } catch (e) {} finally { loading.value = false }
}

function doComment() {
  const t = commentText.value.trim(); if (!t) return
  api.post(`/posts/${route.params.id}/comments`, { content: t })
    .then(() => { commentText.value = ''; toast.success('评论成功'); load() })
    .catch(e => toast.error(e.message))
}

function startEdit(c) { editingId.value = c.id; editText.value = c.content }
function cancelEdit() { editingId.value = 0; editText.value = '' }
function saveEditComment(id) {
  const t = editText.value.trim(); if (!t) return
  api.put(`/comments/${id}`, { content: t })
    .then(() => { cancelEdit(); toast.success('已更新'); load() })
    .catch(e => toast.error(e.message))
}
function delComment(id) { confirm('确定删除评论？', () => api.delete(`/comments/${id}`).then(() => { toast.success('已删除'); load() }).catch(e => toast.error(e.message))) }
function handleDeletePost() { confirm('确定删除帖子？', () => api.delete(`/posts/${route.params.id}`).then(() => { toast.success('已删除'); router.push('/') }).catch(e => toast.error(e.message))) }
function confirm(msg, cb) { confirmMsg.value = msg; confirmCb = cb; confirmOpen.value = true }
function onConfirm() { if (confirmCb) confirmCb() }

function getIcon(n) { const m = {'综合讨论':'💬','技术交流':'💻','军事纵横':'⚔️','历史长廊':'📜','文学艺术':'🎨','生活杂谈':'🌻'}; return m[n] || '📌' }
function fmtDate(s) { return s ? new Date(s).toLocaleString('zh-CN') : '' }

onMounted(load)
</script>
