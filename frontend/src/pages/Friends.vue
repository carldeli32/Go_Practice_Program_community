<template>
  <div class="max-w-[600px] mx-auto">
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono', toast.type.value === 'error' ? 'bg-danger/20 border border-danger/40 text-danger' : 'bg-success/20 border border-success/40 text-success']">{{ toast.msg.value }}</div>
    <h2 class="text-xl font-bold text-fg font-heading tracking-wider mb-5">👥 好友</h2>

    <div class="flex gap-0 mb-4 border-b border-border">
      <button :class="['px-4 py-2 text-xs font-semibold font-heading tracking-wide border-b-2 transition-all', activeTab === 'following' ? 'text-primary border-primary' : 'text-muted-fg border-transparent hover:text-fg']" @click="activeTab = 'following'">我关注的人</button>
      <button :class="['px-4 py-2 text-xs font-semibold font-heading tracking-wide border-b-2 transition-all', activeTab === 'followers' ? 'text-primary border-primary' : 'text-muted-fg border-transparent hover:text-fg']" @click="activeTab = 'followers'">关注我的人</button>
    </div>

    <template v-if="activeTab === 'following'">
      <div v-if="following.length === 0" class="text-center text-muted-fg py-10 text-sm">还没有关注任何人~</div>
      <div v-for="f in following" :key="f.id" class="flex items-center gap-3 py-3 border-b border-border last:border-b-0">
        <router-link :to="`/user/${f.id}`" class="text-sm font-semibold text-primary hover:underline">@{{ f.username }}</router-link>
        <span v-if="f.motto" class="text-xs text-muted-fg truncate flex-1">"{{ f.motto }}"</span>
        <div class="flex gap-2 ml-auto">
          <button @click="$router.push(`/messages/${f.id}`)" class="px-2.5 py-1 bg-card border border-border rounded text-[11px] text-fg font-mono cursor-pointer hover:border-primary transition-all">私信</button>
          <button @click="unfollow(f.id)" class="px-2.5 py-1 bg-card border border-danger/30 rounded text-[11px] text-danger font-mono cursor-pointer hover:bg-danger/10 transition-all">取关</button>
        </div>
      </div>
    </template>

    <template v-if="activeTab === 'followers'">
      <div v-if="followers.length === 0" class="text-center text-muted-fg py-10 text-sm">还没有粉丝~</div>
      <div v-for="f in followers" :key="f.id" class="flex items-center gap-3 py-3 border-b border-border last:border-b-0">
        <router-link :to="`/user/${f.id}`" class="text-sm font-semibold text-primary hover:underline">@{{ f.username }}</router-link>
        <span v-if="f.motto" class="text-xs text-muted-fg truncate flex-1">"{{ f.motto }}"</span>
        <button @click="$router.push(`/messages/${f.id}`)" class="px-2.5 py-1 bg-card border border-border rounded text-[11px] text-fg font-mono cursor-pointer hover:border-primary transition-all ml-auto">私信</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { useToast } from '../composables/toast'
const toast = useToast()
const activeTab = ref('following')
const following = ref([]); const followers = ref([])
async function fetch() { try { const [a,b] = await Promise.all([api.get('/following'), api.get('/followers')]); following.value = a.data.data.users; followers.value = b.data.data.users } catch(e){} }
function unfollow(id) { api.delete(`/follow/${id}`).then(() => { toast.success('已取消关注'); following.value = following.value.filter(f => f.id !== id) }).catch(e => toast.error(e.message)) }
onMounted(() => {
  const tab = new URLSearchParams(window.location.search).get('tab')
  if (tab === 'followers') activeTab.value = 'followers'
  fetch()
})
</script>
