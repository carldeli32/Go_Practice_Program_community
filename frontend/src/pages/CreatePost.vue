<template>
  <div class="max-w-[700px] mx-auto">
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono', toast.type.value === 'success' ? 'bg-success/20 border border-success/40 text-success' : 'bg-danger/20 border border-danger/40 text-danger']">{{ toast.msg.value }}</div>

    <h2 class="text-xl font-bold text-fg font-heading tracking-wider mb-5">{{ isEdit ? '编辑帖子' : '✍️ 发表新帖' }}</h2>

    <div class="bg-card border border-border rounded-lg p-5 space-y-4">
      <select v-model="form.category_id" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg font-mono outline-none focus:border-primary transition-colors">
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ getIcon(cat.name) }} {{ cat.name }}</option>
      </select>

      <input v-model="form.title" placeholder="标题" maxlength="200" class="w-full h-10 bg-input border border-border rounded-md px-3 text-sm text-fg placeholder-muted-fg outline-none focus:border-primary transition-colors" />

      <textarea v-model="form.content" rows="10" placeholder="写下你的想法..." class="w-full bg-input border border-border rounded-md p-3 text-sm text-fg placeholder-muted-fg outline-none focus:border-primary transition-colors resize-none"></textarea>

      <button @click="handleSubmit" :disabled="submitting" class="w-full h-10 bg-primary text-bg border-0 rounded-md text-sm font-bold font-heading tracking-wider cursor-pointer hover:brightness-110 transition-all disabled:opacity-50">
        {{ submitting ? '提交中...' : (isEdit ? '保存修改' : '发布帖子') }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import { useToast } from '../composables/toast'

const route = useRoute(); const router = useRouter()
const submitting = ref(false); const categories = ref([])
const form = reactive({ title: '', content: '', category_id: null })
const isEdit = computed(() => !!route.query.edit)
const toast = useToast()

function getIcon(n) { const m = {'综合讨论':'💬','技术交流':'💻','军事纵横':'⚔️','历史长廊':'📜','文学艺术':'🎨','生活杂谈':'🌻'}; return m[n] || '📌' }

onMounted(async () => {
  api.get('/categories').then(r => {
    categories.value = r.data.data.categories
    const cid = route.query.category_id
    form.category_id = cid ? Number(cid) : (categories.value[0]?.id || 1)
  }).catch(() => toast.error('加载分类失败'))

  if (isEdit.value) {
    api.get(`/posts/${route.query.edit}`).then(r => {
      form.title = r.data.data.post.title
      form.content = r.data.data.post.content
      form.category_id = r.data.data.post.category_id
    }).catch(() => { toast.error('无法加载帖子'); router.push('/') })
  }
})

function handleSubmit() {
  if (!form.title.trim() || !form.content.trim()) { toast.warning('请填写标题和内容'); return }
  submitting.value = true
  const payload = { title: form.title, content: form.content, category_id: form.category_id }
  const req = isEdit.value ? api.put(`/posts/${route.query.edit}`, payload) : api.post('/posts', payload)
  req.then(res => {
    toast.success(isEdit.value ? '更新成功' : '发布成功')
    router.push(isEdit.value ? `/post/${route.query.edit}` : `/post/${res.data.data.post.id}`)
  }).catch(e => toast.error(e.message)).finally(() => submitting.value = false)
}
</script>
