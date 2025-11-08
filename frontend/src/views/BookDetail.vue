<template>
  <div class="book-detail-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <el-button @click="$router.back()">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <div class="header-actions" v-if="book">
            <el-button
              v-if="userStore.isAdmin"
              type="warning"
              @click="handleEdit"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button
              type="primary"
              :disabled="book.available_quantity === 0"
              @click="handleBorrow"
            >
              <el-icon><Reading /></el-icon>
              借阅
            </el-button>
          </div>
        </div>
      </template>

      <div v-if="book" class="book-detail">
        <div class="book-main-info">
          <h1 class="book-title">{{ book.title }}</h1>
          <div class="book-meta">
            <el-tag type="info">{{ book.category }}</el-tag>
            <el-tag :type="book.available_quantity > 0 ? 'success' : 'danger'">
              可借：{{ book.available_quantity }} / {{ book.total_quantity }}
            </el-tag>
          </div>
        </div>

        <el-descriptions :column="2" border>
          <el-descriptions-item label="作者">{{ book.author }}</el-descriptions-item>
          <el-descriptions-item label="ISBN">{{ book.isbn }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ book.category }}</el-descriptions-item>
          <el-descriptions-item label="总数量">{{ book.total_quantity }}</el-descriptions-item>
          <el-descriptions-item label="可借数量">{{ book.available_quantity }}</el-descriptions-item>
          <el-descriptions-item label="添加时间">{{ book.create_time }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ book.update_time }}</el-descriptions-item>
        </el-descriptions>

        <div class="book-description">
          <h3>图书描述</h3>
          <p>{{ book.description || '暂无描述' }}</p>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Edit, Reading } from '@element-plus/icons-vue'
import { getBookDetail } from '@/api/book'
import { borrowBook } from '@/api/borrow'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const book = ref(null)

const loadBookDetail = async () => {
  loading.value = true
  try {
    // 将路由参数转换为整数
    const bookId = parseInt(route.params.id, 10)
    if (isNaN(bookId)) {
      ElMessage.error('无效的图书ID')
      router.push({ name: 'Books' })
      return
    }
    
    const res = await getBookDetail({
      token: userStore.token,
      id: bookId
    })
    book.value = res.data.book
  } catch (error) {
    console.error('加载图书详情失败:', error)
  } finally {
    loading.value = false
  }
}

const handleEdit = () => {
  router.push({ name: 'BookEdit', params: { id: book.value.id } })
}

const handleBorrow = async () => {
  try {
    await ElMessageBox.confirm(`确认借阅《${book.value.title}》吗？`, '确认借阅', {
      type: 'warning'
    })
    
    await borrowBook({
      token: userStore.token,
      book_id: book.value.id
    })
    ElMessage.success('借阅成功')
    loadBookDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('借阅失败:', error)
    }
  }
}

onMounted(() => {
  loadBookDetail()
})
</script>

<style lang="scss" scoped>
.book-detail-container {
  max-width: 1200px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .header-actions {
    display: flex;
    gap: 10px;
  }
}

.book-detail {
  .book-main-info {
    margin-bottom: 30px;

    .book-title {
      font-size: 28px;
      font-weight: bold;
      margin-bottom: 15px;
      color: #303133;
    }

    .book-meta {
      display: flex;
      gap: 10px;
    }
  }

  .book-description {
    margin-top: 30px;

    h3 {
      font-size: 18px;
      margin-bottom: 15px;
      color: #303133;
    }

    p {
      font-size: 14px;
      line-height: 1.8;
      color: #606266;
    }
  }
}

@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;

    .header-actions {
      width: 100%;
      justify-content: flex-end;
    }
  }

  .book-detail .book-main-info .book-title {
    font-size: 22px;
  }
}
</style>
