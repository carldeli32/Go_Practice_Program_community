<template>
  <div class="max-w-[650px] mx-auto">
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono', toast.type.value === 'error' ? 'bg-danger/20 border border-danger/40 text-danger' : 'bg-success/20 border border-success/40 text-success']">{{ toast.msg.value }}</div>

    <h2 class="text-xl font-bold text-fg font-heading tracking-wider mb-5">✉️ 私信</h2>

    <div v-if="conversations.length > 0" class="text-right mb-3">
      <button @click="markAllRead" class="px-3 py-1 bg-card border border-border rounded text-[11px] text-muted-fg font-mono cursor-pointer hover:text-primary hover:border-primary transition-all">全部已读</button>
    </div>

    <div v-if="conversations.length === 0" class="text-center text-muted-fg py-20">
      <p class="text-4xl mb-2">📭</p>
      <p class="text-sm">还没有私信</p>
      <p class="text-xs mt-1">去别人的主页点「发私信」开始聊天吧~</p>
    </div>

    <div v-for="conv in conversations" :key="conv.partner.id" @click="$router.push(`/messages/${conv.partner.id}`)"
      class="flex items-center gap-3.5 p-3.5 mb-2.5 bg-card border border-border rounded-lg cursor-pointer hover:border-border-strong transition-all">
      <span class="w-10 h-10 rounded-full bg-primary flex items-center justify-center text-sm font-bold text-bg flex-shrink-0">{{ conv.partner.username?.charAt(0).toUpperCase() }}</span>
      <div class="flex-1 min-w-0">
        <div class="flex justify-between items-center">
          <strong class="text-sm text-fg">{{ conv.partner.username }}</strong>
          <span class="text-[11px] text-muted-fg font-mono">{{ fmtDate(conv.last_message.created_at) }}</span>
        </div>
        <div class="flex justify-between items-center mt-1">
          <span class="text-xs text-muted-fg truncate">{{ truncate(conv.last_message.content, 50) }}</span>
          <span v-if="conv.unread_count > 0" class="min-w-[18px] h-[18px] rounded-full bg-danger text-white text-[10px] font-bold flex items-center justify-center px-1 ml-2">{{ conv.unread_count }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { useToast } from '../composables/toast'
const conversations = ref([])
const toast = useToast()
async function fetch() { try { const r = await api.get('/messages'); conversations.value = r.data.data.conversations } catch(e){} }
function markAllRead() { api.put('/messages/read-all').then(() => { conversations.value.forEach(c => c.unread_count = 0); toast.success('全部已读'); fetch() }).catch(() => {}) }
function fmtDate(s) { if(!s) return ''; const d = new Date(s), n = new Date(), diff = n - d; if(diff < 864e5) return d.toLocaleTimeString('zh-CN', {hour:'2-digit',minute:'2-digit'}); if(diff < 6048e5) return ['周日','周一','周二','周三','周四','周五','周六'][d.getDay()]; return d.toLocaleDateString('zh-CN') }
function truncate(s, n) { return s && s.length > n ? s.slice(0, n) + '...' : s }
onMounted(fetch)
</script>
