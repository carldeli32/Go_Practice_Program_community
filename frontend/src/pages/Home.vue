<template>
  <div class="home">
    <!-- 公告横幅 -->
    <el-alert
      v-if="announcement"
      :title="announcement"
      type="warning"
      show-icon
      :closable="false"
      style="margin-bottom:20px;" />

    <h2 style="margin-bottom:20px;">📋 帖子列表</h2>

    <div v-if="posts.length === 0 && !loading" style="text-align:center;color:#999;padding:40px;">
      还没有帖子，快去发一篇吧~
    </div>

    <el-card v-for="post in posts" :key="post.id" class="post-card" shadow="hover"
      @click="$router.push(`/post/${post.id}`)">
      <template #header>
        <div class="post-header">
          <span class="post-title">{{ post.title }}</span>
          <span class="post-author">
            <router-link :to="`/user/${post.user?.id}`" class="user-link">
              @{{ post.user?.username || '匿名' }}
            </router-link>
          </span>
        </div>
      </template>
      <p class="post-content">{{ post.content.slice(0, 200) }}{{ post.content.length > 200 ? '...' : '' }}</p>
      <div class="post-meta">
        <span>{{ formatDate(post.created_at) }}</span>
      </div>
    </el-card>

    <div v-if="total > 0" style="margin-top:20px;text-align:center;">
      <el-pagination background layout="prev, pager, next" :total="total" :page-size="pageSize"
        v-model:current-page="page" @current-change="fetchPosts" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const posts = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const announcement = ref('')

async function fetchPosts() {
  loading.value = true
  try {
    const res = await api.get('/posts', { params: { page: page.value, page_size: pageSize.value } })
    posts.value = res.data.posts
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function fetchAnnouncement() {
  try {
    const res = await api.get('/announcement')
    announcement.value = res.data.content
  } catch (e) { /* ignore */ }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchPosts()
  fetchAnnouncement()
})
</script>

<style scoped>
.post-card { margin-bottom: 16px; cursor: pointer; transition: transform 0.15s; }
.post-card:hover { transform: translateY(-2px); }
.post-header { display: flex; justify-content: space-between; align-items: center; }
.post-title { font-size: 17px; font-weight: 600; }
.post-author { color: #909399; font-size: 13px; }
.user-link { color: #409eff; text-decoration: none; }
.user-link:hover { text-decoration: underline; }
.post-content { color: #606266; line-height: 1.6; }
.post-meta { color: #c0c4cc; font-size: 12px; margin-top: 8px; }
</style>
