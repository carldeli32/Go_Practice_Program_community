<template>
  <div class="max-w-[700px] mx-auto" v-if="user.username">
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono', toast.type.value === 'error' ? 'bg-danger/20 border border-danger/40 text-danger' : 'bg-success/20 border border-success/40 text-success']">{{ toast.msg.value }}</div>

    <!-- Profile Card -->
    <div class="bg-card border border-border rounded-lg p-6 mb-5">
      <div class="flex items-start gap-5">
        <span class="w-20 h-20 rounded-full bg-primary flex items-center justify-center text-3xl font-bold text-bg flex-shrink-0">{{ user.username?.charAt(0).toUpperCase() }}</span>
        <div class="flex-1">
          <div class="flex items-center gap-3 mb-2">
            <h2 class="text-xl font-bold text-fg font-heading tracking-wider">{{ user.username }}</h2>
            <span class="px-2.5 py-0.5 rounded-xl text-xs font-heading tracking-wide"
              :class="levelClass(user.level?.level)">{{ user.level?.badge }} {{ user.level?.name }}</span>
          </div>
          <div v-if="auth.token && auth.user?.id !== user.id" class="flex gap-2 mt-2">
            <button v-if="!user.is_following" @click="handleFollow" :disabled="followLoading" class="px-4 py-1.5 bg-primary text-bg rounded-md text-xs font-bold font-heading tracking-wide cursor-pointer hover:brightness-110 transition-all disabled:opacity-50">{{ followLoading ? '...' : '关注' }}</button>
            <button v-else @click="handleUnfollow" :disabled="followLoading" class="px-4 py-1.5 bg-card border border-border rounded-md text-xs text-fg font-heading tracking-wide cursor-pointer hover:border-primary transition-all">已关注</button>
            <button @click="$router.push(`/messages/${user.id}`)" class="px-4 py-1.5 bg-card border border-border rounded-md text-xs text-fg font-heading tracking-wide cursor-pointer hover:border-primary transition-all">发私信</button>
          </div>
        </div>
      </div>

      <hr class="border-border my-5">

      <div class="grid grid-cols-2 gap-3 text-sm">
        <div><span class="text-muted-fg">性别：</span>{{ user.gender || '未填写' }}</div>
        <div><span class="text-muted-fg">年龄：</span>{{ user.age > 0 ? user.age + '岁' : '未填写' }}</div>
        <div class="col-span-2"><span class="text-muted-fg">工作：</span>{{ user.job || '未填写' }}</div>
        <div class="col-span-2"><span class="text-muted-fg">座右铭：</span><em v-if="user.motto" class="text-fg">"{{ user.motto }}"</em><span v-else class="text-muted-fg">未填写</span></div>
        <div><span class="text-muted-fg">发帖：</span><span class="text-fg font-bold">{{ user.post_count }}</span></div>
        <div><span class="text-muted-fg">评论：</span><span class="text-fg font-bold">{{ user.comment_count }}</span></div>
        <div @click="$router.push('/friends?tab=following')" class="cursor-pointer hover:text-primary"><span class="text-muted-fg">关注：</span><span class="text-primary font-bold">{{ user.following_count }}</span></div>
        <div @click="$router.push('/friends?tab=followers')" class="cursor-pointer hover:text-primary"><span class="text-muted-fg">粉丝：</span><span class="text-primary font-bold">{{ user.follower_count }}</span></div>
        <div class="col-span-2"><span class="text-muted-fg">注册：</span><span class="font-mono text-xs">{{ fmtDate(user.created_at) }}</span></div>
      </div>
    </div>

    <!-- User Posts -->
    <div class="bg-card border border-border rounded-lg p-5">
      <h3 class="text-sm font-bold text-fg font-heading tracking-wide mb-3">📝 {{ user.username }} 的帖子</h3>
      <div v-if="posts.length === 0" class="text-center text-muted-fg py-6 text-sm">还没有发过帖子</div>
      <div v-for="post in posts" :key="post.id" @click="$router.push(`/post/${post.id}`)" class="flex justify-between py-2.5 border-b border-border last:border-b-0 cursor-pointer hover:text-primary transition-colors">
        <span class="text-sm text-fg">{{ post.title }}</span>
        <span class="text-xs text-muted-fg font-mono flex-shrink-0 ml-4">{{ fmtDate(post.created_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '../stores/auth'
import api from '../api'
import { useToast } from '../composables/toast'
const route = useRoute(); const auth = useAuth(); const toast = useToast()
const user = ref({}); const posts = ref([]); const followLoading = ref(false)
async function fetch() { try { const r = await api.get(`/users/${route.params.id}`); user.value = r.data.data.user; posts.value = r.data.data.posts } catch(e){} }
function handleFollow() { followLoading.value = true; api.post('/follow', { user_id: user.value.id }).then(() => { user.value.is_following = true; user.value.follower_count++; toast.success('关注成功 🤝') }).catch(e => toast.error(e.message)).finally(() => followLoading.value = false) }
function handleUnfollow() { followLoading.value = true; api.delete(`/follow/${user.value.id}`).then(() => { user.value.is_following = false; user.value.follower_count--; toast.success('已取消关注') }).catch(e => toast.error(e.message)).finally(() => followLoading.value = false) }
function levelClass(l) { return { 1: 'bg-success/20 text-success', 2: 'bg-warning/20 text-warning', 3: 'bg-primary/20 text-primary', 4: 'bg-[#f72585]/20 text-[#f72585]' }[l] || 'bg-muted-fg/20 text-muted-fg' }
function fmtDate(s) { return s ? new Date(s).toLocaleString('zh-CN') : '' }
onMounted(fetch)
</script>
