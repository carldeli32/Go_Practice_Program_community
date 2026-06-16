<template>
  <div class="admin-page">
    <h2 style="margin-bottom:20px;">👑 管理面板</h2>
    <el-tabs v-model="activeTab">
      <!-- ====== 用户管理 ====== -->
      <el-tab-pane label="用户管理" name="users">
        <!-- 创建用户（仅超管） -->
        <div v-if="isSuperAdmin" style="margin-bottom:12px;display:flex;gap:8px;flex-wrap:wrap;">
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
          <el-table-column label="角色" width="200">
            <template #default="{row}">
              <el-tag v-for="r in row.roles" :key="r" :type="roleTagType(r)" size="small" style="margin-right:4px;">
                {{ roleLabel(r) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="80">
            <template #default="{row}">
              <el-tag v-if="row.is_banned" type="danger" size="small">已封禁</el-tag>
              <el-tag v-else type="success" size="small">正常</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="motto" label="座右铭" show-overflow-tooltip />
          <el-table-column label="操作" width="240">
            <template #default="{row}">
              <template v-if="!row.is_admin">
                <el-button size="small" type="danger" @click="banUser(row.id)" :disabled="row.is_banned">封禁</el-button>
                <el-button size="small" type="success" @click="unbanUser(row.id)" :disabled="!row.is_banned">解禁</el-button>
              </template>
              <el-button v-if="isSuperAdmin && !row.is_admin" size="small" type="danger" plain @click="deleteUser(row.id, row.username)">删除</el-button>
              <el-button v-if="isSuperAdmin" size="small" @click="openRoleDialog(row)">角色</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- ====== 分类管理（仅超管） ====== -->
      <el-tab-pane v-if="isSuperAdmin" label="分类管理" name="categories">
        <div style="margin-bottom:12px;display:flex;gap:8px;">
          <el-input v-model="newCatName" placeholder="分类名称" style="width:160px;" size="small" />
          <el-input v-model="newCatDesc" placeholder="描述（可选）" style="width:220px;" size="small" />
          <el-button type="primary" size="small" @click="createCategory">创建分类</el-button>
        </div>
        <el-table :data="categories" border stripe>
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="description" label="描述" />
          <el-table-column label="操作" width="120">
            <template #default="{row}">
              <el-popconfirm v-if="row.id !== 1" title="确定删除？帖子会迁移至综合讨论" @confirm="deleteCategory(row.id)">
                <template #reference>
                  <el-button size="small" type="danger" plain>删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- ====== 站内公告 ====== -->
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

    <!-- 角色编辑弹窗 -->
    <el-dialog v-model="roleDialogVisible" title="编辑角色" width="400px">
      <el-checkbox-group v-model="selectedRoles">
        <el-checkbox v-for="r in allRoles" :key="r" :label="r" :value="r">
          {{ roleLabel(r) }}
        </el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRoles">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuth } from '../stores/auth'
import api from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const auth = useAuth()

const activeTab = ref('users')
const users = ref([])
const userSearch = ref('')
const newUsername = ref('')
const newPassword = ref('')
const announcement = ref('')
const newAnnounce = ref('')
const categories = ref([])
const newCatName = ref('')
const newCatDesc = ref('')

// 兼容旧 session（无 roles 数组）+ 新 session
const isSuperAdmin = computed(() => {
  if (!auth.user) return false
  if (auth.user.roles?.includes('super_admin')) return true
  // 回退：旧 session 数据只有 is_admin，且用户名是 root
  if (auth.user.is_admin && auth.user.username === 'root') return true
  return false
})

// 角色弹窗
const roleDialogVisible = ref(false)
const allRoles = ['super_admin', 'admin', 'moderator', 'user']
const selectedRoles = ref([])
const editingUserId = ref(0)

function roleLabel(r) {
  const map = { super_admin: '超管', admin: '管理员', moderator: '版主', user: '用户' }
  return map[r] || r
}
function roleTagType(r) {
  const map = { super_admin: 'danger', admin: 'warning', moderator: '', user: 'info' }
  return map[r] || 'info'
}

// ===== 用户管理 =====
function fetchUsers() {
  api.get('/admin/users', { params: { q: userSearch.value } })
    .then(res => users.value = res.data.data.users)
}
function createUser() {
  if (!newUsername.value || !newPassword.value) return ElMessage.warning('请输入用户名和密码')
  api.post('/admin/users', { username: newUsername.value, password: newPassword.value })
    .then(() => { ElMessage.success('用户已创建'); newUsername.value = ''; newPassword.value = ''; fetchUsers() })
    .catch(e => ElMessage.error(e.message))
}
function banUser(id) { api.put(`/admin/users/${id}/ban`).then(() => { ElMessage.success('已封禁'); fetchUsers() }) }
function unbanUser(id) { api.put(`/admin/users/${id}/unban`).then(() => { ElMessage.success('已解禁'); fetchUsers() }) }
function deleteUser(id, name) {
  ElMessageBox.confirm(`确定删除用户「${name}」？该操作不可逆。`, '危险操作', { type: 'warning', confirmButtonText: '确认删除' })
    .then(() => api.delete(`/admin/users/${id}`))
    .then(() => { ElMessage.success('已删除'); fetchUsers() })
    .catch(() => {})
}
function openRoleDialog(row) {
  editingUserId.value = row.id
  selectedRoles.value = [...row.roles]
  roleDialogVisible.value = true
}
function saveRoles() {
  api.put(`/admin/users/${editingUserId.value}/roles`, { roles: selectedRoles.value })
    .then(() => { ElMessage.success('角色已更新'); roleDialogVisible.value = false; fetchUsers() })
    .catch(e => ElMessage.error(e.message))
}

// ===== 分类管理 =====
function fetchCategories() {
  api.get('/categories').then(res => categories.value = res.data.data.categories)
}
function createCategory() {
  if (!newCatName.value.trim()) return ElMessage.warning('请输入分类名称')
  api.post('/admin/categories', { name: newCatName.value.trim(), description: newCatDesc.value.trim() })
    .then(() => { ElMessage.success('分类已创建'); newCatName.value = ''; newCatDesc.value = ''; fetchCategories() })
    .catch(e => ElMessage.error(e.message))
}
function deleteCategory(id) {
  api.delete(`/admin/categories/${id}`)
    .then(() => { ElMessage.success('已删除'); fetchCategories() })
    .catch(e => ElMessage.error(e.message))
}

// ===== 公告 =====
function fetchAnnounce() { api.get('/announcement').then(res => announcement.value = res.data.data.content) }
function setAnnounce() {
  if (!newAnnounce.value.trim()) return
  api.post('/admin/announcement', { content: newAnnounce.value })
    .then(() => { announcement.value = newAnnounce.value; newAnnounce.value = ''; ElMessage.success('已发布') })
}
function delAnnounce() {
  api.delete('/admin/announcement').then(() => { announcement.value = ''; ElMessage.success('已删除') })
}

onMounted(() => { fetchUsers(); fetchAnnounce(); fetchCategories() })
</script>
