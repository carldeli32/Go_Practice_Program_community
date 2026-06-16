<template>
  <div class="create-post">
    <h2 style="margin-bottom:20px;">{{ isEdit ? '编辑帖子' : '✍️ 发表新帖' }}</h2>

    <el-card>
      <el-form :model="form" label-width="0">
        <el-form-item>
          <el-select v-model="form.category_id" placeholder="选择分类" filterable style="width:100%">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="getCategoryIcon(cat.name) + ' ' + cat.name"
              :value="cat.id" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-input v-model="form.title" placeholder="标题" size="large" maxlength="200" show-word-limit />
        </el-form-item>

        <el-form-item>
          <el-input v-model="form.content" type="textarea" :rows="10" placeholder="写下你的想法..." />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" @click="handleSubmit" :loading="submitting" style="width:100%">
            {{ isEdit ? '保存修改' : '发布帖子' }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const submitting = ref(false)
const categories = ref([])
const form = reactive({ title: '', content: '', category_id: null })

const isEdit = computed(() => !!route.query.edit)

function getCategoryIcon(name) {
  const icons = {
    '综合讨论': '💬', '技术交流': '💻', '军事纵横': '⚔️',
    '历史长廊': '📜', '文学艺术': '🎨', '生活杂谈': '🌻'
  }
  return icons[name] || '📌'
}

onMounted(async () => {
  // 加载分类列表
  api.get('/categories').then(res => {
    categories.value = res.data.data.categories
    // 从 URL query 预选分类
    const catId = route.query.category_id
    if (catId) {
      form.category_id = Number(catId)
    } else {
      form.category_id = categories.value[0]?.id || 1
    }
  }).catch(e => ElMessage.error('加载分类失败'))

  // 编辑模式：加载帖子
  if (isEdit.value) {
    api.get(`/posts/${route.query.edit}`).then(res => {
      form.title = res.data.data.post.title
      form.content = res.data.data.post.content
      form.category_id = res.data.data.post.category_id
    }).catch(e => {
      ElMessage.error('无法加载帖子')
      router.push('/')
    })
  }
})

function handleSubmit() {
  if (!form.title.trim() || !form.content.trim()) {
    ElMessage.warning('请填写标题和内容')
    return
  }
  submitting.value = true
  const payload = {
    title: form.title,
    content: form.content,
    category_id: form.category_id
  }

  if (isEdit.value) {
    api.put(`/posts/${route.query.edit}`, payload)
      .then(() => {
        ElMessage.success('更新成功')
        router.push(`/post/${route.query.edit}`)
      })
      .catch(e => ElMessage.error(e.message))
      .finally(() => submitting.value = false)
  } else {
    api.post('/posts', payload)
      .then(res => {
        ElMessage.success('发布成功')
        router.push(`/post/${res.data.data.post.id}`)
      })
      .catch(e => ElMessage.error(e.message))
      .finally(() => submitting.value = false)
  }
}
</script>
