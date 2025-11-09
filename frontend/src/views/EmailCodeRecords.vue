<template>
  <div class="email-code-records-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>验证码记录管理</span>
          <div class="header-filters">
            <el-input
              v-model="filters.keyword"
              placeholder="邮箱/验证码"
              clearable
              style="width: 200px; margin-right: 10px"
              @keyup.enter="handleSearch"
            />
            <el-select v-model="filters.action" placeholder="用途" clearable style="width: 120px; margin-right: 10px">
              <el-option label="全部" value="" />
              <el-option label="注册" value="register" />
              <el-option label="忘记密码" value="forget" />
            </el-select>
            <el-select v-model="filters.is_used" placeholder="状态" clearable style="width: 120px; margin-right: 10px">
              <el-option label="全部" value="" />
              <el-option label="已使用" :value="true" />
              <el-option label="未使用" :value="false" />
            </el-select>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
            <el-button type="info" @click="loadStats">刷新统计</el-button>
          </div>
        </div>
      </template>

      <!-- 统计信息 -->
      <el-row :gutter="20" style="margin-bottom: 20px">
        <el-col :span="6">
          <el-statistic title="总验证码数" :value="stats.total_count" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="已使用" :value="stats.used_count" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="未使用" :value="stats.unused_count" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="已过期" :value="stats.expired_count" />
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 20px">
        <el-col :span="12">
          <el-statistic title="注册验证码" :value="stats.register_count" />
        </el-col>
        <el-col :span="12">
          <el-statistic title="忘记密码验证码" :value="stats.forget_count" />
        </el-col>
      </el-row>

      <el-table
        v-loading="loading"
        :data="records"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="code" label="验证码" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="用途" width="120">
          <template #default="{ row }">
            <el-tag :type="row.action === 'register' ? 'primary' : 'warning'">
              {{ row.action === 'register' ? '注册' : '忘记密码' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="expires_at" label="过期时间" width="180">
          <template #default="{ row }">
            <span :class="{ 'expired': isExpired(row.expires_at) }">
              {{ formatDateTime(row.expires_at) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="is_used" label="使用状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.is_used ? 'success' : 'info'">
              {{ row.is_used ? '已使用' : '未使用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="used_at" label="使用时间" width="180">
          <template #default="{ row }">
            {{ row.used_at ? formatDateTime(row.used_at) : '-' }}
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > 0"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="pagination"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getEmailCodeList, getEmailCodeStats } from '@/api/admin'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

const loading = ref(false)
const records = ref([])
const total = ref(0)

const stats = reactive({
  total_count: 0,
  used_count: 0,
  unused_count: 0,
  expired_count: 0,
  register_count: 0,
  forget_count: 0
})

const filters = reactive({
  keyword: '',
  action: '',
  is_used: null
})

const pagination = reactive({
  page: 1,
  limit: 20
})

const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'
  const date = new Date(dateTime)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const isExpired = (expiresAt) => {
  if (!expiresAt) return false
  return new Date(expiresAt) < new Date()
}

const loadRecords = async () => {
  loading.value = true
  try {
    const params = {
      token: userStore.token,
      page: pagination.page,
      limit: pagination.limit
    }
    
    if (filters.keyword) {
      params.keyword = filters.keyword
    }
    if (filters.action) {
      params.action = filters.action
    }
    if (filters.is_used !== null && filters.is_used !== '') {
      params.is_used = filters.is_used
    }

    const res = await getEmailCodeList(params)
    if (res.code === 0) {
      records.value = res.data.list
      total.value = res.data.total
    } else {
      ElMessage.error(res.msg || '获取验证码记录失败')
    }
  } catch (error) {
    console.error('获取验证码记录失败:', error)
    ElMessage.error('获取验证码记录失败')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getEmailCodeStats()
    if (res.code === 0) {
      Object.assign(stats, res.data)
    } else {
      ElMessage.error(res.msg || '获取统计信息失败')
    }
  } catch (error) {
    console.error('获取统计信息失败:', error)
    ElMessage.error('获取统计信息失败')
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadRecords()
}

const handleReset = () => {
  filters.keyword = ''
  filters.action = ''
  filters.is_used = null
  pagination.page = 1
  loadRecords()
}

const handleSizeChange = () => {
  loadRecords()
}

const handlePageChange = () => {
  loadRecords()
}

onMounted(() => {
  loadRecords()
  loadStats()
})
</script>

<style scoped>
.email-code-records-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-filters {
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  justify-content: center;
}

.expired {
  color: #f56c6c;
  font-weight: bold;
}
</style>

