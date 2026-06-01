<template>
  <div class="detail" v-if="post.id">
    <el-card>
      <h2>{{ post.title }}</h2>
      <div class="post-meta">
        <router-link :to="`/user/${post.user?.id}`" class="user-link">@{{ post.user?.username || '匿名' }}</router-link>
        <span>·</span>
        <span>{{ formatDate(post.created_at) }}</span>
      </div>
      <div class="post-body">{{ post.content }}</div>
    </el-card>

    <div v-if="canModifyPost" style="margin-top:12px;display:flex;gap:8px;">
      <el-button @click="$router.push({path:'/create', query:{edit:post.id}})">编辑</el-button>
      <el-button type="danger" @click="handleDeletePost">删除</el-button>
    </div>

    <!-- 评论区 -->
    <el-card style="margin-top:20px;">
      <template #header><span>💬 评论 ({{ comments.length }})</span></template>
      <div v-if="comments.length === 0" style="color:#999;text-align:center;padding:20px;">还没有评论，来抢沙发吧~</div>

      <div v-for="c in comments" :key="c.id" class="comment-item">
        <div class="comment-header">
          <router-link :to="`/user/${c.user?.id}`" class="user-link">@{{ c.user?.username || '匿名' }}</router-link>
          <span class="comment-time">{{ formatDate(c.created_at) }}</span>
        </div>
        <p class="comment-content" v-if="editingComment !== c.id">{{ c.content }}</p>
        <div v-else class="edit-comment-area">
          <el-input v-model="editCommentText" type="textarea" :rows="3" />
          <div style="margin-top:6px;display:flex;gap:6px;">
            <el-button size="small" type="primary" @click="saveEditComment(c.id)">保存</el-button>
            <el-button size="small" @click="cancelEdit">取消</el-button>
          </div>
        </div>
        <div v-if="c.can_edit || c.can_delete" class="comment-actions">
          <el-button text size="small" @click="startEdit(c)" v-if="c.can_edit">编辑</el-button>
          <el-button text size="small" type="danger" @click="delComment(c.id)" v-if="c.can_delete">删除</el-button>
        </div>
      </div>

      <!-- 发表评论 -->
      <div v-if="isLoggedIn" class="comment-form">
        <el-input v-model="commentText" type="textarea" :rows="3" placeholder="写下你的评论..." />
        <el-button type="primary" size="small" style="margin-top:8px;" @click="doComment">发表</el-button>
      </div>
      <div v-else style="text-align:center;margin-top:12px;">
        <el-link type="primary" @click="$router.push('/login')">登录后发表评论</el-link>
      </div>
    </el-card>
  </div>
  <div v-else-if="!loading" style="text-align:center;color:#999;padding:40px;">帖子不存在</div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth, isLoggedIn } from '../stores/auth'
import api from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const loading = ref(true)
const post = ref({})
const comments = ref([])
const commentText = ref('')
const editingComment = ref(0)
const editCommentText = ref('')

const canModifyPost = computed(() => {
  if (!auth.user || !post.value.id) return false
  return auth.user.id === post.value.user_id || auth.user.is_admin
})

async function load() {
  loading.value = true
  try {
    const [p, c] = await Promise.all([
      api.get(`/posts/${route.params.id}`),
      api.get(`/posts/${route.params.id}/comments`),
    ])
    post.value = p.data.post
    comments.value = c.data.comments || []
  } catch (e) {
    console.error(e)
  }
  loading.value = false
}

function doComment() {
  const txt = commentText.value.trim()
  if (!txt) return
  api.post(`/posts/${route.params.id}/comments`, { content: txt })
    .then(() => { commentText.value = ''; ElMessage.success('评论成功'); load() })
    .catch(e => ElMessage.error(e.message))
}

function startEdit(c) { editingComment.value = c.id; editCommentText.value = c.content }
function cancelEdit() { editingComment.value = 0; editCommentText.value = '' }
function saveEditComment(id) {
  const txt = editCommentText.value.trim()
  if (!txt) return
  api.put(`/comments/${id}`, { content: txt })
    .then(() => { cancelEdit(); ElMessage.success('已更新'); load() })
    .catch(e => ElMessage.error(e.message))
}
function delComment(id) {
  ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' })
    .then(() => api.delete(`/comments/${id}`))
    .then(() => { ElMessage.success('已删除'); load() })
    .catch(() => {})
}
function handleDeletePost() {
  ElMessageBox.confirm('确定删除？', '提示', { type: 'warning' })
    .then(() => api.delete(`/posts/${route.params.id}`))
    .then(() => { ElMessage.success('已删除'); router.push('/') })
    .catch(() => {})
}
function formatDate(s) { return s ? new Date(s).toLocaleString('zh-CN') : '' }

onMounted(load)
</script>

<style scoped>
.post-meta { color: #909399; font-size: 13px; margin: 8px 0 16px; display: flex; gap: 8px; }
.post-body { white-space: pre-wrap; line-height: 1.8; color: #303133; }
.comment-item { padding: 10px 0; border-bottom: 1px solid #f0f0f0; }
.comment-item:last-child { border-bottom: none; }
.comment-header { display: flex; justify-content: space-between; margin-bottom: 4px; }
.comment-time { color: #c0c4cc; font-size: 12px; }
.comment-content { color: #606266; line-height: 1.5; margin: 0; }
.comment-actions { margin-top: 6px; display: flex; gap: 4px; }
.edit-comment-area { margin: 6px 0; }
.user-link { color: #409eff; text-decoration: none; font-weight: 600; }
.user-link:hover { text-decoration: underline; }
.comment-form { margin-top: 16px; padding-top: 16px; border-top: 1px solid #f0f0f0; }
</style>
