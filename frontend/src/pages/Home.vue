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

    <!-- ========== 模式 1：分类大厅 ========== -->
    <template v-if="!activeCategory">
      <div class="category-header">
        <h2>🏠 分类大厅</h2>
        <el-input
          v-model="categorySearch"
          placeholder="搜索分类..."
          :prefix-icon="Search"
          clearable
          style="width:280px;"
          @input="filterCategories" />
      </div>

      <div v-if="filteredCategories.length === 0 && !loadingCats" style="text-align:center;color:#999;padding:40px;">
        没有找到匹配的分类
      </div>

      <div class="category-grid">
        <el-card
          v-for="cat in filteredCategories"
          :key="cat.id"
          class="category-card"
          shadow="hover"
          @click="enterCategory(cat)">
          <div class="category-card-content">
            <span class="category-icon">{{ getCategoryIcon(cat.name) }}</span>
            <div class="category-info">
              <span class="category-name">{{ cat.name }}</span>
              <span class="category-desc">{{ cat.description }}</span>
            </div>
            <el-icon class="category-arrow"><ArrowRight /></el-icon>
          </div>
        </el-card>
      </div>
    </template>

    <!-- ========== 模式 2：帖子列表 ========== -->
    <template v-else>
      <div class="post-list-header">
        <el-button @click="backToCategories" :icon="ArrowLeft" text>返回分类</el-button>
        <h2>{{ getCategoryIcon(activeCategory.name) }} {{ activeCategory.name }}</h2>
        <el-input
          v-model="postSearch"
          placeholder="搜索帖子标题..."
          :prefix-icon="Search"
          clearable
          style="width:260px;"
          @keyup.enter="searchPosts"
          @clear="searchPosts" />
      </div>

      <div v-if="posts.length === 0 && !loading" style="text-align:center;color:#999;padding:40px;">
        这个分类还没有帖子，快去发一篇吧~
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
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import api from '../api'

const route = useRoute()
const router = useRouter()

// ===== 分类数据 =====
const categories = ref([])
const categorySearch = ref('')
const loadingCats = ref(false)

const filteredCategories = computed(() => {
  if (!categorySearch.value) return categories.value
  const q = categorySearch.value.toLowerCase()
  return categories.value.filter(c => c.name.toLowerCase().includes(q) || c.description.toLowerCase().includes(q))
})

// ===== 当前激活的分类 =====
const activeCategory = ref(null)

// ===== 帖子数据 =====
const posts = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const postSearch = ref('')
const announcement = ref('')

// ===== 初始化 =====
async function fetchCategories() {
  loadingCats.value = true
  try {
    const res = await api.get('/categories')
    categories.value = res.data.data.categories
  } catch (e) {
    console.error(e)
  } finally {
    loadingCats.value = false
  }
}

async function fetchAnnouncement() {
  try {
    const res = await api.get('/announcement')
    announcement.value = res.data.data.content
  } catch (e) { /* ignore */ }
}

async function fetchPosts() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (activeCategory.value) {
      params.category_id = activeCategory.value.id
    }
    if (postSearch.value) {
      params.q = postSearch.value
    }
    const res = await api.get('/posts', { params })
    posts.value = res.data.data.posts
    total.value = res.data.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

// ===== 分类操作 =====
function enterCategory(cat) {
  activeCategory.value = cat
  postSearch.value = ''
  page.value = 1
  router.replace({ query: { category_id: cat.id } })
  fetchPosts()
}

function backToCategories() {
  activeCategory.value = null
  router.replace({ query: {} })
}

function searchPosts() {
  page.value = 1
  fetchPosts()
}

// ===== 工具函数 =====
function getCategoryIcon(name) {
  const icons = {
    '综合讨论': '💬', '技术交流': '💻', '军事纵横': '⚔️',
    '历史长廊': '📜', '文学艺术': '🎨', '生活杂谈': '🌻'
  }
  return icons[name] || '📌'
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

// ===== 监听路由变化 =====
watch(() => route.query.category_id, (newVal) => {
  if (newVal) {
    const cat = categories.value.find(c => c.id === Number(newVal))
    if (cat) {
      activeCategory.value = cat
      page.value = 1
      postSearch.value = ''
      fetchPosts()
    }
  } else {
    activeCategory.value = null
  }
})

onMounted(async () => {
  await fetchCategories()
  fetchAnnouncement()

  // 如果 URL 带 category_id，自动进入
  const catId = route.query.category_id
  if (catId) {
    const cat = categories.value.find(c => c.id === Number(catId))
    if (cat) {
      activeCategory.value = cat
      fetchPosts()
    }
  }
})
</script>

<style scoped>
/* ===== 分类卡片 ===== */
.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.category-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.category-card {
  cursor: pointer;
  transition: transform 0.15s;
}
.category-card:hover {
  transform: translateY(-3px);
}
.category-card-content {
  display: flex;
  align-items: center;
  gap: 14px;
}
.category-icon {
  font-size: 32px;
  flex-shrink: 0;
}
.category-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}
.category-name {
  font-size: 16px;
  font-weight: 600;
}
.category-desc {
  font-size: 13px;
  color: #909399;
}
.category-arrow {
  color: #c0c4cc;
  flex-shrink: 0;
}

/* ===== 帖子列表 ===== */
.post-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.post-list-header h2 {
  margin: 0;
  flex: 1;
  margin-left: 12px;
}
.post-card {
  margin-bottom: 16px;
  cursor: pointer;
  transition: transform 0.15s;
}
.post-card:hover {
  transform: translateY(-2px);
}
.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.post-title {
  font-size: 17px;
  font-weight: 600;
}
.post-author {
  color: #909399;
  font-size: 13px;
}
.user-link {
  color: #409eff;
  text-decoration: none;
}
.user-link:hover {
  text-decoration: underline;
}
.post-content {
  color: #606266;
  line-height: 1.6;
}
.post-meta {
  color: #c0c4cc;
  font-size: 12px;
  margin-top: 8px;
}
</style>
