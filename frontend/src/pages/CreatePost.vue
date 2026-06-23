<template>
  <div class="max-w-[700px] mx-auto">
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono', toast.type.value === 'success' ? 'bg-success/20 border border-success/40 text-success' : 'bg-danger/20 border border-danger/40 text-danger']">{{ toast.msg.value }}</div>

    <h2 class="text-xl font-bold text-fg font-heading tracking-wider mb-5">{{ isEdit ? '编辑帖子' : '✍️ 发表新帖' }}</h2>

    <div class="bg-card border border-border rounded-lg p-5 space-y-4">
      <select v-model="categoryId" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg font-mono outline-none focus:border-primary transition-colors">
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ getIcon(cat.name) }} {{ cat.name }}</option>
      </select>

      <input v-model="form.title" placeholder="标题" maxlength="200" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg outline-none focus:border-primary transition-colors" />

      <Editor v-model="form.content" :post-id="postId" placeholder="分享你的想法，支持拖入或粘贴图片..." />

      <button @click="handleSubmit" :disabled="submitting" class="w-full h-10 bg-primary text-bg border-0 rounded-md text-sm font-bold font-heading tracking-wider cursor-pointer hover:brightness-110 transition-all disabled:opacity-50">
        {{ submitting ? '提交中...' : (isEdit ? '保存修改' : '发布帖子') }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import { useToast } from '../composables/toast'
import Editor from '../components/Editor.vue'
import TurndownService from 'turndown'

const turndown = new TurndownService({ headingStyle: 'atx' })

const route = useRoute(); const router = useRouter()
const submitting = ref(false)
const categoryId = ref(null)
const form = reactive({ title: '', content: '' })
const isEdit = computed(() => !!route.query.edit)
const toast = useToast()
const postId = ref(0)
const published = ref(false)
const pageReady = ref(0)

// 尝试用 App.vue 预热的分类，秒显下拉框
const sharedCategories = inject('sharedCategories', null)
const categories = sharedCategories?.value?.length ? sharedCategories : ref([])

function getIcon(n) { const m = {'综合讨论':'💬','技术交流':'💻','军事纵横':'⚔️','历史长廊':'📜','文学艺术':'🎨','生活杂谈':'🌻'}; return m[n] || '📌' }

onMounted(() => {
  pageReady.value = Date.now()

  // 如果没有预热，自己拉分类
  if (!categories.value.length) {
    api.get('/categories').then(r => {
      categories.value = r.data.data.categories
      const cid = route.query.category_id
      categoryId.value = cid ? Number(cid) : (categories.value[0]?.id ?? null)
    }).catch(() => toast.error('加载分类失败'))
  } else {
    const cid = route.query.category_id
    categoryId.value = cid ? Number(cid) : (categories.value[0]?.id ?? null)
  }

  if (isEdit.value) {
    api.get(`/posts/${route.query.edit}`).then(r => {
      const p = r.data.data.post
      postId.value = p.id
      form.title = p.title
      form.content = p.content
      categoryId.value = p.category_id
    }).catch(() => { toast.error('无法加载帖子'); router.push('/') })
  } else {
    // 后台创建草稿，不阻塞页面渲染
    const cid = route.query.category_id || 1
    api.post('/posts', { title: '', content: '', category_id: Number(cid), status: 'draft' })
      .then(r => { postId.value = r.data.data.post.id })
      .catch(() => toast.error('初始化失败'))
  }
})

onBeforeUnmount(() => {
  if (!isEdit.value && !published.value && postId.value) {
    api.delete(`/posts/${postId.value}`).catch(() => {})
  }
})

function handleSubmit() {
  if (!form.title.trim() || !form.content.trim()) { toast.warning('请填写标题和内容'); return }
  if (!categoryId.value) { toast.warning('请选择分类'); return }
  if (Date.now() - pageReady.value < 1000) { toast.warning('编辑时间不得少于 1 秒'); return }
  submitting.value = true
  const md = turndown.turndown(form.content)
  const payload = { title: form.title, content: md, category_id: Number(categoryId.value) }

  if (isEdit.value) {
    api.put(`/posts/${route.query.edit}`, payload)
      .then(() => { published.value = true; toast.success('更新成功'); router.push(`/post/${route.query.edit}`) })
      .catch(e => toast.error(e.message)).finally(() => submitting.value = false)
  } else {
    api.put(`/posts/${postId.value}`, { ...payload, status: 'published' })
      .then(() => { published.value = true; toast.success('发布成功'); router.push(`/post/${postId.value}`) })
      .catch(e => toast.error(e.message)).finally(() => submitting.value = false)
  }
}
</script>
