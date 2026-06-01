<template>
  <div class="friends-page">
    <h2 style="margin-bottom:20px;">👥 好友</h2>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="我关注的人" name="following">
        <div v-if="following.length === 0" style="text-align:center;color:#999;padding:40px;">
          还没有关注任何人~
        </div>
        <div v-for="f in following" :key="f.id" class="friend-row">
          <router-link :to="`/user/${f.id}`" class="friend-name">@{{ f.username }}</router-link>
          <span v-if="f.motto" class="friend-motto">"{{ f.motto }}"</span>
          <div style="margin-left:auto;display:flex;gap:8px;">
            <el-button size="small" @click="$router.push(`/messages/${f.id}`)">
              <el-icon><Message /></el-icon> 私信
            </el-button>
            <el-button size="small" type="danger" plain @click="unfollow(f.id)">取关</el-button>
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="关注我的人" name="followers">
        <div v-if="followers.length === 0" style="text-align:center;color:#999;padding:40px;">
          还没有粉丝~
        </div>
        <div v-for="f in followers" :key="f.id" class="friend-row">
          <router-link :to="`/user/${f.id}`" class="friend-name">@{{ f.username }}</router-link>
          <span v-if="f.motto" class="friend-motto">"{{ f.motto }}"</span>
          <el-button size="small" style="margin-left:auto;" @click="$router.push(`/messages/${f.id}`)">
            <el-icon><Message /></el-icon> 私信
          </el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import { ElMessage } from 'element-plus'

const activeTab = ref('following')
const following = ref([])
const followers = ref([])

async function fetchData() {
  try {
    const [fingRes, ferRes] = await Promise.all([
      api.get('/following'),
      api.get('/followers'),
    ])
    following.value = fingRes.data.users
    followers.value = ferRes.data.users
  } catch (e) {
    console.error(e)
  }
}

function unfollow(userID) {
  api.delete(`/follow/${userID}`)
    .then(() => {
      ElMessage.success('已取消关注')
      following.value = following.value.filter(f => f.id !== userID)
    })
    .catch(e => ElMessage.error(e.message))
}

onMounted(fetchData)
</script>

<style scoped>
.friends-page { max-width: 600px; margin: 0 auto; }
.friend-row {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  gap: 12px;
}
.friend-row:last-child { border-bottom: none; }
.friend-name { color: #409eff; text-decoration: none; font-weight: 600; font-size: 15px; }
.friend-name:hover { text-decoration: underline; }
.friend-motto { color: #909399; font-size: 13px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 200px; }
</style>
