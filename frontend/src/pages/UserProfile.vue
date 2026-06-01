<template>
  <div class="profile-page" v-if="user.username">
    <!-- 用户信息卡片 -->
    <el-card class="info-card">
      <div class="profile-header">
        <el-avatar :size="80" style="font-size:32px;background:#409eff;">
          {{ user.username?.charAt(0).toUpperCase() }}
        </el-avatar>
        <div class="profile-title">
          <h2>{{ user.username }}</h2>
          <span class="level-badge" :class="'level-' + user.level.level">
            {{ user.level.badge }} {{ user.level.name }}
          </span>
        </div>
        <div v-if="auth.token && auth.user?.id !== user.id" class="profile-actions">
          <el-button
            v-if="!user.is_following"
            type="primary"
            @click="handleFollow"
            :loading="followLoading">
            <el-icon><Plus /></el-icon> 关注
          </el-button>
          <el-button
            v-else
            @click="handleUnfollow"
            :loading="followLoading">
            已关注
          </el-button>
          <el-button @click="$router.push(`/messages/${user.id}`)">
            <el-icon><Message /></el-icon> 发私信
          </el-button>
        </div>
      </div>

      <el-divider />

      <el-descriptions :column="2" border>
        <el-descriptions-item label="性别">
          {{ user.gender || '未填写' }}
        </el-descriptions-item>
        <el-descriptions-item label="年龄">
          {{ user.age > 0 ? user.age + '岁' : '未填写' }}
        </el-descriptions-item>
        <el-descriptions-item label="当前工作" :span="2">
          {{ user.job || '未填写' }}
        </el-descriptions-item>
        <el-descriptions-item label="座右铭" :span="2">
          <em v-if="user.motto">"{{ user.motto }}"</em>
          <span v-else style="color:#c0c4cc;">未填写</span>
        </el-descriptions-item>
        <el-descriptions-item label="发帖数">
          {{ user.post_count }}
        </el-descriptions-item>
        <el-descriptions-item label="评论数">
          {{ user.comment_count }}
        </el-descriptions-item>
        <el-descriptions-item label="关注">
          <router-link :to="`/friends?tab=following`" class="count-link" @click.stop>
            {{ user.following_count }}
          </router-link>
        </el-descriptions-item>
        <el-descriptions-item label="粉丝">
          <router-link :to="`/friends?tab=followers`" class="count-link" @click.stop>
            {{ user.follower_count }}
          </router-link>
        </el-descriptions-item>
        <el-descriptions-item label="注册时间" :span="2">
          {{ formatDate(user.created_at) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 该用户的帖子列表 -->
    <el-card style="margin-top:20px;">
      <template #header>
        <span>📝 {{ user.username }} 的帖子</span>
      </template>
      <div v-if="posts.length === 0" style="color:#999;text-align:center;padding:20px;">
        还没有发过帖子
      </div>
      <div v-for="post in posts" :key="post.id" class="user-post" @click="$router.push(`/post/${post.id}`)">
        <span class="post-title">{{ post.title }}</span>
        <span class="post-time">{{ formatDate(post.created_at) }}</span>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '../stores/auth'
import api from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const auth = useAuth()
const user = ref({})
const posts = ref([])
const followLoading = ref(false)

async function fetchProfile() {
  try {
    const res = await api.get(`/users/${route.params.id}`)
    user.value = res.data.user
    posts.value = res.data.posts
  } catch (e) {
    console.error(e)
  }
}

function handleFollow() {
  followLoading.value = true
  api.post('/follow', { user_id: user.value.id })
    .then(() => {
      user.value.is_following = true
      user.value.follower_count++
      ElMessage.success('关注成功 🤝')
    })
    .catch(e => ElMessage.error(e.message))
    .finally(() => followLoading.value = false)
}

function handleUnfollow() {
  followLoading.value = true
  api.delete(`/follow/${user.value.id}`)
    .then(() => {
      user.value.is_following = false
      user.value.follower_count--
      ElMessage.success('已取消关注')
    })
    .catch(e => ElMessage.error(e.message))
    .finally(() => followLoading.value = false)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(fetchProfile)
</script>

<style scoped>
.profile-page { max-width: 700px; margin: 0 auto; }
.profile-header { display: flex; align-items: center; gap: 20px; }
.profile-title h2 { margin: 0 0 4px 0; }
.profile-actions { display: flex; gap: 8px; margin-left: auto; }
.level-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 13px;
}
.level-1 { background: #e8f5e9; color: #4caf50; }
.level-2 { background: #fff3e0; color: #ff9800; }
.level-3 { background: #e3f2fd; color: #2196f3; }
.level-4 { background: #fce4ec; color: #e91e63; }
.count-link { color: #409eff; text-decoration: none; }
.count-link:hover { text-decoration: underline; }
.user-post {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
}
.user-post:last-child { border-bottom: none; }
.user-post:hover .post-title { color: #409eff; }
.post-title { font-size: 14px; }
.post-time { color: #c0c4cc; font-size: 12px; }
</style>
