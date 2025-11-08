<template>
  <div class="all-borrows-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>全部借阅记录</span>
          <div class="header-filters">
            <el-input
              v-model="filters.user_email"
              placeholder="用户邮箱"
              clearable
              style="width: 200px; margin-right: 10px"
              @keyup.enter="handleSearch"
            />
            <el-input
              v-model="filters.book_title"
              placeholder="图书名称"
              clearable
              style="width: 200px; margin-right: 10px"
              @keyup.enter="handleSearch"
            />
            <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px; margin-right: 10px">
              <el-option label="全部" value="" />
              <el-option label="未归还" value="borrowed" />
              <el-option label="已归还" value="returned" />
            </el-select>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="records"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="user_email" label="用户邮箱" width="200" />
        <el-table-column prop="book_title" label="图书名称" width="200" />
        <el-table-column prop="borrow_date" label="借阅日期" width="180" />
        <el-table-column prop="due_date" label="应还日期" width="180" />
        <el-table-column prop="return_date" label="归还日期" width="180">
          <template #default="{ row }">
            {{ row.return_date || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'borrowed' ? 'warning' : 'success'">
              {{ row.status === 'borrowed' ? '未归还' : '已归还' }}
            </el-tag>
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
import { getAllRecords } from '@/api/borrow'

const userStore = useUserStore()

const loading = ref(false)
const records = ref([])
const total = ref(0)

const filters = reactive({
  user_email: '',
  book_title: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  limit: 10
})

const loadRecords = async () => {
  loading.value = true
  try {
    const res = await getAllRecords({
      token: userStore.token,
      user_email: filters.user_email || undefined,
      book_title: filters.book_title || undefined,
      status: filters.status || undefined,
      page: pagination.page,
      limit: pagination.limit
    })
    records.value = res.data.records || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载借阅记录失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadRecords()
}

const handleReset = () => {
  filters.user_email = ''
  filters.book_title = ''
  filters.status = ''
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
})
</script>

<style lang="scss" scoped>
.all-borrows-container {
  max-width: 1400px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;

  .header-filters {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
  }
}

.pagination {
  margin-top: 20px;
  justify-content: center;
}

@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;

    .header-filters {
      width: 100%;
      flex-direction: column;

      .el-input,
      .el-select {
        width: 100% !important;
        margin-right: 0 !important;
        margin-bottom: 10px;
      }
    }
  }
}
</style>
