<template>
  <div class="conversation-page">
    <div class="chat-header">
      <el-button @click="$router.push('/messages')" text><el-icon><ArrowLeft /></el-icon> 返回</el-button>
      <router-link :to="`/user/${partner.id}`" class="partner-name">{{ partner.username }}</router-link>
      <el-button size="small" @click="showNewThread = true" style="margin-left:auto;">+ 新主题</el-button>
    </div>

    <!-- 主题选择 -->
    <div v-if="threads.length > 0" class="thread-tabs">
      <el-radio-group v-model="activeThread" size="small" @change="switchThread">
        <el-radio-button v-for="t in threads" :key="t.id" :value="t.id">
          {{ t.title }} ({{ t.message_count }})
        </el-radio-button>
      </el-radio-group>
      <el-popconfirm title="删除该主题及所有消息？" @confirm="deleteThread(activeThread)">
        <template #reference>
          <el-button text size="small" type="danger">删除此主题</el-button>
        </template>
      </el-popconfirm>
    </div>

    <!-- 无主题引导 -->
    <div v-if="threads.length === 0" style="text-align:center;color:#999;padding:40px;">
      <p>还没有对话主题</p>
      <el-button type="primary" @click="showNewThread = true">创建第一个主题</el-button>
    </div>

    <!-- 消息 -->
    <div class="chat-messages" ref="msgList" v-if="activeThread">
      <div v-if="messages.length === 0" style="text-align:center;color:#999;padding:40px;">发送第一条消息吧~</div>
      <div v-for="msg in messages" :key="msg.id" :class="['msg-bubble', msg.from_user_id === myID ? 'msg-mine' : 'msg-other']">
        <div class="msg-content">{{ msg.content }}</div>
        <div class="msg-time">{{ formatTime(msg.created_at) }}</div>
      </div>
    </div>

    <div class="chat-input" v-if="activeThread">
      <el-input v-model="text" placeholder="输入消息..." @keyup.enter="sendMsg" size="large">
        <template #append><el-button @click="sendMsg">发送</el-button></template>
      </el-input>
    </div>

    <el-dialog v-model="showNewThread" title="新建对话主题" width="380px">
      <el-input v-model="newThreadTitle" placeholder="比如：旅行计划、日常聊天..." />
      <template #footer>
        <el-button @click="showNewThread = false">取消</el-button>
        <el-button type="primary" @click="createThread">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '../stores/auth'
import api from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const auth = useAuth()
const myID = computed(() => auth.user?.id)
const partner = ref({})
const threads = ref([])
const activeThread = ref(0)
const messages = ref([])
const text = ref('')
const msgList = ref(null)
const showNewThread = ref(false)
const newThreadTitle = ref('')

async function loadThreads() {
  const res = await api.get('/threads', { params: { with: route.params.id } })
  threads.value = res.data.threads
  if (threads.value.length > 0 && !activeThread.value) {
    activeThread.value = threads.value[0].id
  }
  if (activeThread.value) await loadMessages()
}

async function loadMessages() {
  const res = await api.get(`/messages/${route.params.id}`, { params: { thread: activeThread.value } })
  partner.value = res.data.partner
  messages.value = res.data.messages
  await api.put(`/messages/${route.params.id}/read`)
  scrollToBottom()
}

async function switchThread() { await loadMessages() }

function sendMsg() {
  if (!text.value.trim()) return
  api.post('/messages', { to_user_id: Number(route.params.id), thread_id: activeThread.value, content: text.value })
    .then(() => { text.value = ''; loadMessages() })
    .catch(e => ElMessage.error(e.message))
}

function createThread() {
  if (!newThreadTitle.value.trim()) return
  api.post('/threads', { with_user_id: Number(route.params.id), title: newThreadTitle.value })
    .then(res => {
      showNewThread.value = false
      newThreadTitle.value = ''
      const tid = res.data.thread.id
      return api.get('/threads', { params: { with: route.params.id } })
    })
    .then(res => {
      threads.value = res.data.threads
      activeThread.value = threads.value[threads.value.length - 1].id
      return loadMessages()
    })
    .catch(e => ElMessage.error(e.message))
}

function deleteThread(id) {
  api.delete(`/threads/${id}`)
    .then(() => {
      threads.value = threads.value.filter(t => t.id !== id)
      if (activeThread.value === id) {
        activeThread.value = threads.value.length > 0 ? threads.value[0].id : 0
      }
      if (activeThread.value) loadMessages()
      ElMessage.success('已删除')
    })
    .catch(e => ElMessage.error(e.message))
}

function scrollToBottom() { nextTick(() => { if (msgList.value) msgList.value.scrollTop = msgList.value.scrollHeight }) }
function formatTime(s) { return s ? new Date(s).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }) : '' }

onMounted(loadThreads)
</script>

<style scoped>
.conversation-page { max-width: 650px; margin: 0 auto; display: flex; flex-direction: column; height: calc(100vh - 140px); }
.chat-header { display: flex; align-items: center; gap: 12px; padding: 8px 0; border-bottom: 1px solid #eee; margin-bottom: 8px; }
.partner-name { font-size: 16px; font-weight: 600; color: #409eff; text-decoration: none; }
.thread-tabs { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; flex-wrap: wrap; }
.chat-messages { flex: 1; overflow-y: auto; padding: 8px 0; }
.msg-bubble { margin-bottom: 14px; max-width: 75%; }
.msg-mine { margin-left: auto; text-align: right; }
.msg-mine .msg-content { background: #409eff; color: #fff; }
.msg-other .msg-content { background: #f0f0f0; color: #303133; }
.msg-content { display: inline-block; padding: 10px 14px; border-radius: 16px; font-size: 14px; line-height: 1.5; word-break: break-word; white-space: pre-wrap; }
.msg-time { font-size: 11px; color: #c0c4cc; margin-top: 4px; }
.chat-input { margin-top: 12px; padding-top: 12px; border-top: 1px solid #eee; }
</style>
