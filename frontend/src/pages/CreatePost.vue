<template>
  <div class="create-post">
    <h2 style="margin-bottom:20px;">{{ isEdit ? '编辑帖子' : '✍️ 发表新帖' }}</h2>

    <el-card>
      <el-form :model="form" label-width="0">
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
const form = reactive({ title: '', content: '' })

const isEdit = computed(() => !!route.query.edit)

onMounted(async () => {
  if (isEdit.value) {
    try {
      const res = await api.get(`/posts/${route.query.edit}`)
      form.title = res.data.post.title
      form.content = res.data.post.content
    } catch (e) {
      ElMessage.error('无法加载帖子')
      router.push('/')
    }
  }
})

async function handleSubmit() {
  if (!form.title.trim() || !form.content.trim()) {
    ElMessage.warning('请填写标题和内容')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await api.put(`/posts/${route.query.edit}`, form)
      ElMessage.success('更新成功')
      router.push(`/post/${route.query.edit}`)
    } else {
      const res = await api.post('/posts', form)
      ElMessage.success('发布成功')
      router.push(`/post/${res.data.post.id}`)
    }
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}
</script>
