<template>
  <div class="max-w-[900px] mx-auto">
    <div v-if="toast.msg.value" :class="['fixed top-20 right-5 z-200 px-4 py-2 rounded-lg text-xs font-mono', toast.type.value === 'error' ? 'bg-danger/20 border border-danger/40 text-danger' : 'bg-success/20 border border-success/40 text-success']">{{ toast.msg.value }}</div>
    <h2 class="text-xl font-bold text-fg font-heading tracking-wider mb-5">👑 管理面板</h2>

    <div class="flex gap-0 mb-4 border-b border-border">
      <button v-for="t in tabs" :key="t" :class="['px-4 py-2 text-xs font-semibold font-heading tracking-wide border-b-2 transition-all', activeTab === t ? 'text-primary border-primary' : 'text-muted-fg border-transparent hover:text-fg']" @click="activeTab = t">{{ t }}</button>
    </div>

    <!-- Users Tab -->
    <div v-if="activeTab === '用户管理'">
      <div v-if="isSuperAdmin" class="flex gap-2 mb-3 flex-wrap">
        <input v-model="newUsername" placeholder="用户名" class="w-[140px] h-8 bg-input border border-border rounded px-2 text-xs text-fg placeholder-muted-fg font-mono outline-none focus:border-primary">
        <input v-model="newPassword" type="password" placeholder="密码" class="w-[140px] h-8 bg-input border border-border rounded px-2 text-xs text-fg placeholder-muted-fg font-mono outline-none focus:border-primary">
        <button @click="createUser" class="px-3 h-8 bg-success text-bg rounded text-xs font-bold font-heading tracking-wide cursor-pointer hover:brightness-110 transition-all">创建用户</button>
      </div>
      <div class="flex gap-2 mb-3">
        <input v-model="userSearch" @keyup.enter="fetchUsers" placeholder="搜索用户名..." class="w-[250px] h-8 bg-input border border-border rounded px-2 text-xs text-fg placeholder-muted-fg font-mono outline-none focus:border-primary">
        <button @click="fetchUsers" class="px-3 h-8 bg-primary text-bg rounded text-xs font-bold font-heading tracking-wide cursor-pointer">搜索</button>
      </div>

      <div class="bg-card border border-border rounded-lg overflow-hidden">
        <table class="w-full text-xs">
          <thead class="bg-secondary">
            <tr>
              <th class="p-2 text-left text-muted-fg font-heading tracking-wide">ID</th>
              <th class="p-2 text-left text-muted-fg font-heading tracking-wide">用户名</th>
              <th class="p-2 text-left text-muted-fg font-heading tracking-wide">角色</th>
              <th class="p-2 text-left text-muted-fg font-heading tracking-wide">状态</th>
              <th class="p-2 text-left text-muted-fg font-heading tracking-wide">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id" class="border-t border-border hover:bg-card-hover">
              <td class="p-2 font-mono text-muted-fg">{{ u.id }}</td>
              <td class="p-2 text-fg font-medium">{{ u.username }}</td>
              <td class="p-2">
                <span v-for="r in u.roles" :key="r" class="inline-block px-1.5 py-0.5 rounded text-[10px] font-mono mr-1"
                  :class="{ 'bg-danger/20 text-danger': r === 'super_admin', 'bg-warning/20 text-warning': r === 'admin', 'bg-primary/20 text-primary': r === 'moderator', 'bg-muted-fg/20 text-muted-fg': r === 'user' }">{{ roleLabel(r) }}</span>
              </td>
              <td class="p-2">
                <span :class="['inline-block px-1.5 py-0.5 rounded text-[10px] font-mono', u.is_banned ? 'bg-danger/20 text-danger' : 'bg-success/20 text-success']">{{ u.is_banned ? '已封禁' : '正常' }}</span>
              </td>
              <td class="p-2">
                <div class="flex gap-1 flex-wrap">
                  <template v-if="!u.is_admin">
                    <button @click="banUser(u.id)" :disabled="u.is_banned" class="px-2 py-0.5 bg-card border border-danger/40 rounded text-[10px] text-danger font-mono cursor-pointer hover:bg-danger/10 transition-all disabled:opacity-40">封禁</button>
                    <button @click="unbanUser(u.id)" :disabled="!u.is_banned" class="px-2 py-0.5 bg-card border border-success/40 rounded text-[10px] text-success font-mono cursor-pointer hover:bg-success/10 transition-all disabled:opacity-40">解禁</button>
                  </template>
                  <button v-if="isSuperAdmin && !u.is_admin" @click="confirm('确定删除？', () => doDelete(u.id))" class="px-2 py-0.5 bg-card border border-border rounded text-[10px] text-danger font-mono cursor-pointer hover:bg-danger/10 transition-all">删除</button>
                  <button v-if="isSuperAdmin" @click="openRoleDialog(u)" class="px-2 py-0.5 bg-card border border-border rounded text-[10px] text-muted-fg font-mono cursor-pointer hover:border-primary transition-all">角色</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Categories Tab -->
    <div v-if="activeTab === '分类管理' && isSuperAdmin">
      <div class="flex gap-2 mb-3">
        <input v-model="newCatName" placeholder="分类名称" class="w-[160px] h-8 bg-input border border-border rounded px-2 text-xs text-fg placeholder-muted-fg outline-none focus:border-primary">
        <input v-model="newCatDesc" placeholder="描述" class="w-[220px] h-8 bg-input border border-border rounded px-2 text-xs text-fg placeholder-muted-fg outline-none focus:border-primary">
        <button @click="createCategory" class="px-3 h-8 bg-primary text-bg rounded text-xs font-bold font-heading tracking-wide cursor-pointer">创建</button>
      </div>
      <div class="bg-card border border-border rounded-lg overflow-hidden">
        <table class="w-full text-xs">
          <thead class="bg-secondary"><tr><th class="p-2 text-left text-muted-fg font-heading">ID</th><th class="p-2 text-left text-muted-fg font-heading">名称</th><th class="p-2 text-left text-muted-fg font-heading">描述</th><th class="p-2 text-left text-muted-fg font-heading">操作</th></tr></thead>
          <tbody>
            <tr v-for="c in categories" :key="c.id" class="border-t border-border">
              <td class="p-2 font-mono text-muted-fg">{{ c.id }}</td><td class="p-2 text-fg">{{ c.name }}</td><td class="p-2 text-muted-fg">{{ c.description }}</td>
              <td class="p-2"><button v-if="c.id !== 1" @click="confirm('确定删除？帖子会迁移至综合讨论', () => deleteCategory(c.id))" class="px-2 py-0.5 bg-card border border-danger/40 rounded text-[10px] text-danger font-mono cursor-pointer hover:bg-danger/10">删除</button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Announce Tab -->
    <div v-if="activeTab === '站内公告'">
      <div v-if="announcement" class="p-3 bg-warning/10 border border-warning/30 rounded-lg mb-3 text-xs text-warning">{{ announcement }}</div>
      <textarea v-model="newAnnounce" rows="3" placeholder="输入公告内容..." class="w-full bg-input border border-border rounded p-3 text-sm text-fg placeholder-muted-fg outline-none focus:border-primary resize-none mb-2"></textarea>
      <div class="flex gap-2">
        <button @click="setAnnounce" class="px-4 py-1.5 bg-primary text-bg rounded text-xs font-bold font-heading tracking-wide cursor-pointer">发布</button>
        <button v-if="announcement" @click="delAnnounce" class="px-4 py-1.5 bg-card border border-danger/40 rounded text-xs text-danger font-mono cursor-pointer">删除</button>
      </div>
    </div>

    <!-- Role Dialog -->
    <div v-if="roleOpen" class="fixed inset-0 z-500 flex items-center justify-center bg-black/50" @click.self="roleOpen = false">
      <div class="bg-card border border-border rounded-lg p-6 max-w-[400px] w-full mx-4">
        <h3 class="text-sm font-bold text-fg font-heading tracking-wide mb-4">编辑角色</h3>
        <div class="space-y-2 mb-4">
          <label v-for="r in allRoles" :key="r" class="flex items-center gap-2 text-xs text-fg cursor-pointer">
            <input type="checkbox" :value="r" v-model="selectedRoles" class="accent-primary"> {{ roleLabel(r) }}
          </label>
        </div>
        <div class="flex justify-end gap-2">
          <button @click="roleOpen = false" class="px-4 py-1.5 bg-card border border-border rounded text-xs text-muted-fg font-mono cursor-pointer">取消</button>
          <button @click="saveRoles" class="px-4 py-1.5 bg-primary text-bg rounded text-xs font-bold font-heading tracking-wide cursor-pointer">保存</button>
        </div>
      </div>
    </div>

    <!-- Confirm Dialog -->
    <div v-if="confirmOpen" class="fixed inset-0 z-500 flex items-center justify-center bg-black/50" @click.self="confirmOpen = false">
      <div class="bg-card border border-border rounded-lg p-6 max-w-[360px] w-full mx-4">
        <p class="text-sm text-fg mb-4">{{ confirmMsg }}</p>
        <div class="flex justify-end gap-2">
          <button @click="confirmOpen = false" class="px-4 py-1.5 bg-card border border-border rounded text-xs text-muted-fg font-mono cursor-pointer">取消</button>
          <button @click="onConfirm(); confirmOpen = false" class="px-4 py-1.5 bg-danger text-white rounded text-xs font-bold font-heading tracking-wide cursor-pointer">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuth } from '../stores/auth'
import api from '../api'
import { useToast } from '../composables/toast'
const auth = useAuth(); const toast = useToast()
const activeTab = ref('用户管理'); const tabs = ['用户管理', '分类管理', '站内公告']
const users = ref([]); const userSearch = ref(''); const newUsername = ref(''); const newPassword = ref('')
const announcement = ref(''); const newAnnounce = ref('')
const categories = ref([]); const newCatName = ref(''); const newCatDesc = ref('')
const roleOpen = ref(false); const allRoles = ['super_admin', 'admin', 'moderator', 'user']
const selectedRoles = ref([]); const editingUserId = ref(0)
const confirmOpen = ref(false); const confirmMsg = ref(''); let confirmCb = null
const isSuperAdmin = computed(() => auth.user && (auth.user.roles?.includes('super_admin') || (auth.user.is_admin && auth.user.username === 'root')))

function roleLabel(r) { return { super_admin: '超管', admin: '管理员', moderator: '版主', user: '用户' }[r] || r }
function confirm(msg, cb) { confirmMsg.value = msg; confirmCb = cb; confirmOpen.value = true }
function onConfirm() { if (confirmCb) confirmCb() }
function fetchUsers() { api.get('/admin/users', { params: { q: userSearch.value } }).then(r => users.value = r.data.data.users) }
function createUser() { if (!newUsername.value || !newPassword.value) return toast.warning('请输入用户名和密码'); api.post('/admin/users', { username: newUsername.value, password: newPassword.value }).then(() => { toast.success('已创建'); newUsername.value = ''; newPassword.value = ''; fetchUsers() }).catch(e => toast.error(e.message)) }
function banUser(id) { api.put(`/admin/users/${id}/ban`).then(() => { toast.success('已封禁'); fetchUsers() }) }
function unbanUser(id) { api.put(`/admin/users/${id}/unban`).then(() => { toast.success('已解禁'); fetchUsers() }) }
function doDelete(id) { api.delete(`/admin/users/${id}`).then(() => { toast.success('已删除'); fetchUsers() }).catch(e => toast.error(e.message)) }
function openRoleDialog(u) { editingUserId.value = u.id; selectedRoles.value = [...u.roles]; roleOpen.value = true }
function saveRoles() { api.put(`/admin/users/${editingUserId.value}/roles`, { roles: selectedRoles.value }).then(() => { toast.success('角色已更新'); roleOpen.value = false; fetchUsers() }).catch(e => toast.error(e.message)) }
function fetchCategories() { api.get('/categories').then(r => categories.value = r.data.data.categories) }
function createCategory() { if (!newCatName.value.trim()) return toast.warning('请输入名称'); api.post('/admin/categories', { name: newCatName.value.trim(), description: newCatDesc.value.trim() }).then(() => { toast.success('已创建'); newCatName.value = ''; newCatDesc.value = ''; fetchCategories() }).catch(e => toast.error(e.message)) }
function deleteCategory(id) { api.delete(`/admin/categories/${id}`).then(() => { toast.success('已删除'); fetchCategories() }).catch(e => toast.error(e.message)) }
function fetchAnnounce() { api.get('/announcement').then(r => announcement.value = r.data.data.content) }
function setAnnounce() { if (!newAnnounce.value.trim()) return; api.post('/admin/announcement', { content: newAnnounce.value }).then(() => { announcement.value = newAnnounce.value; newAnnounce.value = ''; toast.success('已发布') }) }
function delAnnounce() { api.delete('/admin/announcement').then(() => { announcement.value = ''; toast.success('已删除') }) }
onMounted(() => { fetchUsers(); fetchAnnounce(); fetchCategories() })
</script>
