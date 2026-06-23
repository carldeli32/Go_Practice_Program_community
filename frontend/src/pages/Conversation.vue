<template>
  <div class="max-w-[650px] mx-auto flex flex-col" style="height: calc(100vh - 160px);">
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono', toast.type.value === 'error' ? 'bg-danger/20 border border-danger/40 text-danger' : 'bg-success/20 border border-success/40 text-success']">{{ toast.msg.value }}</div>

    <!-- Header -->
    <div class="flex items-center gap-3 pb-2 border-b border-border mb-2">
      <button @click="$router.push('/messages')" class="text-xs text-muted-fg font-mono cursor-pointer bg-transparent border-0 hover:text-primary">← 返回</button>
      <router-link :to="`/user/${partner.id}`" class="text-base font-bold text-primary font-heading tracking-wide hover:underline">{{ partner.username }}</router-link>
      <button @click="showNewThread = true" class="ml-auto px-3 py-1 bg-card border border-border rounded text-[11px] text-fg font-mono cursor-pointer hover:border-primary transition-all">+ 新主题</button>
    </div>

    <!-- Threads -->
    <div v-if="threads.length > 0" class="flex items-center gap-2 mb-2 flex-wrap">
      <button v-for="t in threads" :key="t.id" @click="switchThread(t.id)"
        :class="['px-3 py-1 rounded text-[11px] font-mono cursor-pointer transition-all', activeThread === t.id ? 'bg-primary text-bg border border-primary' : 'bg-card text-muted-fg border border-border hover:border-primary']">
        {{ t.title }} ({{ t.message_count }})
      </button>
      <button @click="confirmDelete" class="px-2 py-1 text-[10px] text-muted-fg hover:text-danger font-mono cursor-pointer bg-transparent border-0">删除</button>
    </div>

    <!-- No threads -->
    <div v-if="threads.length === 0" class="text-center text-muted-fg py-10">
      <p class="text-sm mb-3">还没有对话主题</p>
      <button @click="showNewThread = true" class="px-4 py-2 bg-primary text-bg rounded-md text-xs font-bold font-heading tracking-wide cursor-pointer">创建第一个主题</button>
    </div>

    <!-- Messages -->
    <div v-if="activeThread" ref="msgList" class="flex-1 overflow-y-auto py-2 space-y-3.5">
      <div v-if="messages.length === 0" class="text-center text-muted-fg py-10 text-sm">发送第一条消息吧~</div>
      <div v-for="msg in messages" :key="msg.id" :class="['max-w-[75%]', msg.from_user_id === myID ? 'ml-auto text-right' : '']"
        @contextmenu.prevent="openMenu(msg, $event)">
        <template v-if="msg.is_recalled">
          <div class="inline-block px-3.5 py-2.5 rounded-2xl text-sm italic bg-card border border-border text-muted-fg">
            {{ msg.from_user_id === myID ? '你撤回了一条消息' : '对方撤回了一条消息' }}
          </div>
        </template>
        <template v-else>
          <div :class="['inline-block px-3.5 py-2.5 rounded-2xl text-sm leading-relaxed whitespace-pre-wrap break-words', msg.from_user_id === myID ? 'bg-primary text-bg' : 'bg-secondary text-fg']">{{ msg.content }}</div>
          <div class="text-[10px] text-muted-fg font-mono mt-1 px-1">{{ fmtTime(msg.created_at) }}</div>
        </template>
      </div>
    </div>

    <!-- Input -->
    <div v-if="activeThread" class="mt-3 pt-3 border-t border-border">
      <div class="flex gap-2">
        <input v-model="text" @keyup.enter="sendMsg" placeholder="输入消息..." class="flex-1 h-9 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg outline-none focus:border-primary transition-colors">
        <button @click="sendMsg" class="px-4 bg-primary text-bg border-0 rounded-md text-xs font-bold font-heading tracking-wide cursor-pointer hover:brightness-110 transition-all">发送</button>
      </div>
    </div>

    <!-- Context Menu Backdrop -->
    <div v-if="ctxMenu.show" class="fixed inset-0 z-599" @click="closeMenu" @contextmenu.prevent="closeMenu"></div>

    <!-- Context Menu -->
    <div v-if="ctxMenu.show" :style="{ left: ctxMenu.x + 'px', top: ctxMenu.y + 'px' }" class="fixed z-600 bg-card border border-border rounded-md py-1 shadow-lg min-w-[100px]">
      <button @click="recallMsg" class="w-full text-left px-3 py-1.5 text-xs text-danger font-mono cursor-pointer bg-transparent border-0 hover:bg-secondary transition-colors">撤回</button>
    </div>

    <!-- New Thread Dialog -->
    <div v-if="showNewThread" class="fixed inset-0 z-500 flex items-center justify-center bg-black/50" @click.self="showNewThread = false">
      <div class="bg-card border border-border rounded-lg p-6 max-w-[380px] w-full mx-4">
        <h3 class="text-sm font-bold text-fg font-heading tracking-wide mb-3">新建对话主题</h3>
        <input v-model="newThreadTitle" placeholder="比如：旅行计划、日常聊天..." class="w-full h-9 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg outline-none focus:border-primary mb-4">
        <div class="flex justify-end gap-2">
          <button @click="showNewThread = false" class="px-4 py-1.5 bg-card border border-border rounded text-xs text-muted-fg font-mono cursor-pointer">取消</button>
          <button @click="createThread" class="px-4 py-1.5 bg-primary text-bg rounded text-xs font-bold font-heading tracking-wide cursor-pointer">创建</button>
        </div>
      </div>
    </div>

    <!-- Delete confirm -->
    <div v-if="deleteOpen" class="fixed inset-0 z-500 flex items-center justify-center bg-black/50" @click.self="deleteOpen = false">
      <div class="bg-card border border-border rounded-lg p-6 max-w-[360px] w-full mx-4">
        <p class="text-sm text-fg mb-4">删除该主题及所有消息？</p>
        <div class="flex justify-end gap-2">
          <button @click="deleteOpen = false" class="px-4 py-1.5 bg-card border border-border rounded text-xs text-muted-fg font-mono cursor-pointer">取消</button>
          <button @click="doDelete" class="px-4 py-1.5 bg-danger text-white rounded text-xs font-bold font-heading tracking-wide cursor-pointer">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '../stores/auth'
import api from '../api'
import { useToast } from '../composables/toast'

const route = useRoute(); const auth = useAuth(); const toast = useToast()
const myID = computed(() => auth.user?.id)
const partner = ref({}); const threads = ref([]); const activeThread = ref(0)
const messages = ref([]); const text = ref(''); const msgList = ref(null)
const showNewThread = ref(false); const newThreadTitle = ref('')
const deleteOpen = ref(false)
const stream = ref(null)
const partnerId = computed(() => Number(route.params.id))

const ctxMenu = reactive({ show: false, x: 0, y: 0, msg: null })

function openMenu(msg, e) {
  if (!msg.is_recalled && msg.from_user_id === myID.value) {
    ctxMenu.msg = msg
    ctxMenu.x = e.clientX
    ctxMenu.y = e.clientY
    ctxMenu.show = true
  }
}
function closeMenu() { ctxMenu.show = false }
function recallMsg() {
  if (!ctxMenu.msg) return
  api.put(`/messages/recall/${ctxMenu.msg.id}`)
    .then(() => { ctxMenu.msg.is_recalled = true; closeMenu(); toast.success('已撤回') })
    .catch(e => { toast.error(e.message); closeMenu() })
}

async function loadThreads() {
  const r = await api.get('/threads', { params: { with: route.params.id } })
  threads.value = r.data.data.threads
  if (threads.value.length && !activeThread.value) activeThread.value = threads.value[0].id
  if (activeThread.value) await loadMessages()
}
async function loadMessages() {
  const r = await api.get(`/messages/${route.params.id}`, { params: { thread: activeThread.value } })
  partner.value = r.data.data.partner; messages.value = r.data.data.messages
  api.put(`/messages/${route.params.id}/read`).catch(() => {})
  nextTick(() => { if(msgList.value) msgList.value.scrollTop = msgList.value.scrollHeight })
}
async function switchThread(tid) { activeThread.value = tid; await loadMessages() }
function sendMsg() {
  if (!text.value.trim()) return
  api.post('/messages', { to_user_id: Number(route.params.id), thread_id: activeThread.value, content: text.value })
    .then(() => { text.value = ''; loadMessages() }).catch(e => toast.error(e.message))
}
function createThread() {
  if (!newThreadTitle.value.trim()) return
  api.post('/threads', { with_user_id: Number(route.params.id), title: newThreadTitle.value })
    .then(res => { showNewThread.value = false; newThreadTitle.value = ''; return api.get('/threads', { params: { with: route.params.id } }) })
    .then(res => { threads.value = res.data.data.threads; activeThread.value = threads.value[threads.value.length - 1].id; return loadMessages() })
    .catch(e => toast.error(e.message))
}
function confirmDelete() { deleteOpen.value = true }
function doDelete() {
  api.delete(`/threads/${activeThread.value}`).then(() => {
    threads.value = threads.value.filter(t => t.id !== activeThread.value)
    activeThread.value = threads.value.length ? threads.value[0].id : 0
    deleteOpen.value = false
    if (activeThread.value) loadMessages(); else messages.value = []
    toast.success('已删除')
  }).catch(e => toast.error(e.message))
}
function fmtTime(s) { return s ? new Date(s).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }) : '' }

function connectSSE() {
  const token = localStorage.getItem('token')
  stream.value = new EventSource('/api/messages/stream?token=' + encodeURIComponent(token))
  stream.value.onmessage = (e) => {
    const d = JSON.parse(e.data)
    if (d.type === 'new_message' && d.data.from_user_id === partnerId.value) {
      messages.value.push(d.data)
      api.put(`/messages/${partnerId.value}/read`).catch(() => {})
      nextTick(() => { if (msgList.value) msgList.value.scrollTop = msgList.value.scrollHeight })
    } else if (d.type === 'recall') {
      const msg = messages.value.find(m => m.id === d.data.message_id)
      if (msg) msg.is_recalled = true
    }
  }
}

onMounted(() => { loadThreads(); connectSSE() })
onBeforeUnmount(() => { closeMenu(); stream.value?.close() })
</script>
