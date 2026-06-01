<template>
  <div class="messages-page">
    <h2 style="margin-bottom:20px;">✉️ 私信</h2>

    <div v-if="conversations.length > 0" style="text-align:right;margin-bottom:12px;">
      <el-button size="small" @click="markAllRead">全部已读</el-button>
    </div>

    <div v-if="conversations.length === 0" style="text-align:center;color:#999;padding:60px;">
      <p style="font-size:40px;">📭</p>
      <p>还没有私信</p>
      <p style="font-size:13px;">去别人的主页点「发私信」开始聊天吧~</p>
    </div>

    <el-card v-for="conv in conversations" :key="conv.partner.id" class="conv-card" shadow="hover"
      @click="$router.push(`/messages/${conv.partner.id}`)">
      <div class="conv-row">
        <el-avatar :size="40" style="background:#409eff;">
          {{ conv.partner.username?.charAt(0).toUpperCase() }}
        </el-avatar>
        <div class="conv-info">
          <div class="conv-top">
            <strong>{{ conv.partner.username }}</strong>
            <span class="conv-time">{{ formatDate(conv.last_message.created_at) }}</span>
          </div>
          <div class="conv-bottom">
            <span class="conv-preview">{{ conv.last_message.content.slice(0, 50) }}</span>
            <el-badge v-if="conv.unread_count > 0" :value="conv.unread_count" type="danger" />
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { ElMessage } from 'element-plus'

const conversations = ref([])

async function fetchConversations() {
  try {
    const res = await api.get('/messages')
    conversations.value = res.data.conversations
  } catch (e) {
    console.error(e)
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now - d
  if (diff < 86400000) return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  if (diff < 604800000) return ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][d.getDay()]
  return d.toLocaleDateString('zh-CN')
}

function markAllRead() {
  api.put('/messages/read-all')
    .then(() => {
      conversations.value.forEach(c => c.unread_count = 0)
      ElMessage.success('全部已读')
      fetchConversations()  // 刷新列表
    })
    .catch(e => console.error(e))
}

onMounted(fetchConversations)
</script>

<style scoped>
.messages-page { max-width: 650px; margin: 0 auto; }
.conv-card { margin-bottom: 10px; cursor: pointer; }
.conv-card:hover { background: #fafafa; }
.conv-row { display: flex; align-items: center; gap: 14px; }
.conv-info { flex: 1; min-width: 0; }
.conv-top { display: flex; justify-content: space-between; align-items: center; }
.conv-bottom { display: flex; justify-content: space-between; align-items: center; margin-top: 4px; }
.conv-time { color: #c0c4cc; font-size: 12px; }
.conv-preview { color: #909399; font-size: 13px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
