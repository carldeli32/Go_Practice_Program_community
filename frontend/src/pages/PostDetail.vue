<template>
  <div class="detail" v-if="post.id">
    <!-- 帖子主体 -->
    <el-card>
      <h2>{{ post.title }}</h2>
      <div class="post-meta">
        <span>@{{ post.user?.username || '匿名' }}</span>
        <span>·</span>
        <span>{{ formatDate(post.created_at) }}</span>
        <span v-if="post.updated_at !== post.created_at" style="color:#c0c4cc;">（已编辑）</span>
      </div>
      <div class="post-body">{{ post.content }}</div>
    </el-card>

    <!-- 评论区 -->
    <el-card style="margin-top:20px;">
      <template #header>
        <span>💬 评论 ({{ comments.length }})</span>
      </template>

      <div v-if="comments.length === 0" style="color:#999;text-align:center;padding:20px;">
        还没有评论，来抢沙发吧~
      </div>

      <div v-for="c in comments" :key="c.id" class="comment-item">
        <div class="comment-header">
          <strong>@{{ c.user?.username || '匿名' }}</strong>
          <span class="comment-time">{{ formatDate(c.created_at) }}</span>
        </div>
        <p class="comment-content">{{ c.content }}</p>
      </div>

      <!-- 发表评论 -->
      <div v-if="auth.token" class="comment-form">
        <el-input v-model="commentText" type="textarea" :rows="3" placeholder="写下你的评论..." />
        <el-button type="primary" size="small" style="margin-top:8px;" @click="submitComment" :loading="submitting">
          发表
        </el-button>
      </div>
      <div v-else style="text-align:center;margin-top:12px;">
        <el-link type="primary" @click="$router.push('/login')">登录后发表评论</el-link>
      </div>
    </el-card>

    <!-- 操作按钮（作者可见） -->
    <div v-if="isAuthor" style="margin-top:16px;display:flex;gap:8px;">
      <el-button @click="$router.push({path:'/create', query:{edit:post.id}})">编辑</el-button>
      <el-button type="danger" @click="handleDelete">删除</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../stores/auth'
import api from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const auth = useAuth()

const post = ref({})
const comments = ref([])
const commentText = ref('')
const submitting = ref(false)

const isAuthor = computed(() => auth.user && post.value.user_id === auth.user.id)

async function fetchPost() {
  const res = await api.get(`/posts/${route.params.id}`)
  post.value = res.data.post
}

async function fetchComments() {
  const res = await api.get(`/posts/${route.params.id}/comments`)
  comments.value = res.data.comments
}

function submitComment() {
  if (!commentText.value.trim()) return
  submitting.value = true
  api.post(`/posts/${route.params.id}/comments`, { content: commentText.value })
    .then(() => {
      commentText.value = ''
      ElMessage.success('评论成功')
      fetchComments()
    })
    .catch((e) => {
      ElMessage.error(e.message)
    })
    .finally(() => {
      submitting.value = false
    })
}

async function handleDelete() {
  try {
    await ElMessageBox.confirm('确定删除这篇帖子吗？', '提示', { type: 'warning' })
    await api.delete(`/posts/${route.params.id}`)
    ElMessage.success('删除成功')
    router.push('/')
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message)
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchPost()
  fetchComments()
})
</script>

<style scoped>
.post-meta { color: #909399; font-size: 13px; margin: 8px 0 16px; display: flex; gap: 8px; }
.post-body { white-space: pre-wrap; line-height: 1.8; color: #303133; }
.comment-item { padding: 12px 0; border-bottom: 1px solid #f0f0f0; }
.comment-item:last-child { border-bottom: none; }
.comment-header { display: flex; justify-content: space-between; margin-bottom: 4px; }
.comment-time { color: #c0c4cc; font-size: 12px; }
.comment-content { color: #606266; line-height: 1.5; }
.comment-form { margin-top: 16px; padding-top: 16px; border-top: 1px solid #f0f0f0; }
</style>
