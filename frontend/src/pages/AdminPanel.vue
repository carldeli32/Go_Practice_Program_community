<template>
  <div class="admin-page">
    <h2 style="margin-bottom:20px;">👑 管理面板</h2>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="用户管理" name="users">
        <div style="margin-bottom:12px;display:flex;gap:8px;flex-wrap:wrap;">
          <el-input v-model="newUsername" placeholder="用户名" style="width:140px;" size="small" />
          <el-input v-model="newPassword" placeholder="密码" style="width:140px;" size="small" type="password" show-password />
          <el-button type="success" size="small" @click="createUser">创建用户</el-button>
        </div>
        <div style="margin-bottom:12px;display:flex;gap:8px;">
          <el-input v-model="userSearch" placeholder="搜索用户名..." style="width:250px;" clearable @clear="fetchUsers" @keyup.enter="fetchUsers" />
          <el-button type="primary" @click="fetchUsers">搜索</el-button>
        </div>
        <el-table :data="users" border stripe max-height="500">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="username" label="用户名" />
          <el-table-column label="管理员" width="80">
            <template #default="{row}"><el-tag v-if="row.is_admin" type="warning" size="small">管理员</el-tag></template>
          </el-table-column>
          <el-table-column label="状态" width="80">
            <template #default="{row}">
              <el-tag v-if="row.is_banned" type="danger" size="small">已封禁</el-tag>
              <el-tag v-else type="success" size="small">正常</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="motto" label="座右铭" show-overflow-tooltip />
          <el-table-column label="操作" width="180">
            <template #default="{row}">
              <template v-if="!row.is_admin">
                <el-button size="small" type="danger" @click="banUser(row.id)" :disabled="row.is_banned">封禁</el-button>
                <el-button size="small" type="success" @click="unbanUser(row.id)" :disabled="!row.is_banned">解禁</el-button>
              </template>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="站内公告" name="announce">
        <div v-if="announcement" style="margin-bottom:12px;">
          <el-alert :title="announcement" type="warning" :closable="false" show-icon />
        </div>
        <el-input v-model="newAnnounce" type="textarea" :rows="3" placeholder="输入公告内容..." />
        <div style="margin-top:8px;display:flex;gap:8px;">
          <el-button type="primary" @click="setAnnounce">发布</el-button>
          <el-button type="danger" plain @click="delAnnounce" v-if="announcement">删除</el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { ElMessage } from 'element-plus'

const activeTab = ref('users')
const users = ref([])
const userSearch = ref('')
const newUsername = ref('')
const newPassword = ref('')
const announcement = ref('')
const newAnnounce = ref('')

function fetchUsers() {
  api.get('/admin/users', { params: { q: userSearch.value } }).then(res => users.value = res.data.users)
}
function createUser() {
  if (!newUsername.value || !newPassword.value) return ElMessage.warning('请输入用户名和密码')
  api.post('/admin/users', { username: newUsername.value, password: newPassword.value })
    .then(() => { ElMessage.success('用户已创建'); newUsername.value = ''; newPassword.value = ''; fetchUsers() })
    .catch(e => ElMessage.error(e.message))
}
function banUser(id) { api.put(`/admin/users/${id}/ban`).then(() => { ElMessage.success('已封禁'); fetchUsers() }) }
function unbanUser(id) { api.put(`/admin/users/${id}/unban`).then(() => { ElMessage.success('已解禁'); fetchUsers() }) }
function fetchAnnounce() { api.get('/announcement').then(res => announcement.value = res.data.content) }
function setAnnounce() {
  if (!newAnnounce.value.trim()) return
  api.post('/admin/announcement', { content: newAnnounce.value })
    .then(() => { announcement.value = newAnnounce.value; newAnnounce.value = ''; ElMessage.success('已发布') })
}
function delAnnounce() {
  api.delete('/admin/announcement').then(() => { announcement.value = ''; ElMessage.success('已删除') })
}

onMounted(() => { fetchUsers(); fetchAnnounce() })
</script>
