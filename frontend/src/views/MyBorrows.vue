<template>
  <div class="my-borrows-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的借阅记录</span>
          <el-radio-group v-model="statusFilter" @change="handleStatusChange">
            <el-radio-button label="all">全部</el-radio-button>
            <el-radio-button label="borrowed">未归还</el-radio-button>
            <el-radio-button label="returned">已归还</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="records"
        stripe
        style="width: 100%"
      >
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
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'borrowed'"
              type="primary"
              size="small"
              @click="handleReturn(row)"
            >
              归还
            </el-button>
            <span v-else>-</span>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBorrowRecords } from '@/api/borrow'
import { returnBook } from '@/api/borrow'

const userStore = useUserStore()

const loading = ref(false)
const records = ref([])
const total = ref(0)
const statusFilter = ref('all')

const pagination = reactive({
  page: 1,
  limit: 10
})

const loadRecords = async () => {
  loading.value = true
  try {
    const res = await getBorrowRecords({
      token: userStore.token,
      status: statusFilter.value === 'all' ? undefined : statusFilter.value,
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

const handleStatusChange = () => {
  pagination.page = 1
  loadRecords()
}

const handleReturn = async (row) => {
  try {
    await ElMessageBox.confirm(`确认归还《${row.book_title}》吗？`, '确认归还', {
      type: 'warning'
    })
    
    await returnBook({
      token: userStore.token,
      book_id: row.book_id
    })
    ElMessage.success('归还成功')
    loadRecords()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('归还失败:', error)
    }
  }
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
.my-borrows-container {
  max-width: 1400px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  justify-content: center;
}

@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
