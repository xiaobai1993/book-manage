<template>
  <div class="books-container">
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入书名或作者"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="请选择分类" clearable>
            <el-option label="全部" value="" />
            <el-option label="文学" value="文学" />
            <el-option label="科幻" value="科幻" />
            <el-option label="历史" value="历史" />
            <el-option label="童话" value="童话" />
            <el-option label="科技" value="科技" />
            <el-option label="教育" value="教育" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>图书列表</span>
          <el-button
            v-if="userStore.isAdmin"
            type="primary"
            @click="$router.push({ name: 'BookAdd' })"
          >
            <el-icon><Plus /></el-icon>
            添加图书
          </el-button>
        </div>
      </template>

      <el-empty v-if="books.length === 0 && !loading" description="暂无数据" />
      
      <el-row :gutter="20" v-loading="loading">
        <el-col
          v-for="book in books"
          :key="book.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
          class="book-col"
        >
          <el-card class="book-card" shadow="hover" @click="handleBookClick(book.id)">
            <div class="book-cover-wrapper" v-if="book.cover_image_url">
              <img :src="book.cover_image_url" :alt="book.title" class="book-cover" />
            </div>
            <div class="book-cover-placeholder" v-else>
              <el-icon><Picture /></el-icon>
              <span>暂无封面</span>
            </div>
            <div class="book-info">
              <h3 class="book-title">{{ book.title }}</h3>
              <p class="book-author">作者：{{ book.author }}</p>
              <p class="book-isbn">ISBN：{{ book.isbn }}</p>
              <p class="book-category">分类：{{ book.category }}</p>
              <div class="book-quantity">
                <el-tag :type="book.available_quantity > 0 ? 'success' : 'danger'">
                  可借：{{ book.available_quantity }} / {{ book.total_quantity }}
                </el-tag>
              </div>
            </div>
            <div class="book-actions">
              <el-button
                type="primary"
                size="small"
                :disabled="book.available_quantity === 0"
                @click.stop="handleBorrow(book)"
              >
                借阅
              </el-button>
              <el-button
                v-if="userStore.isAdmin"
                type="warning"
                size="small"
                @click.stop="handleEdit(book.id)"
              >
                编辑
              </el-button>
              <el-button
                v-if="userStore.isAdmin"
                type="danger"
                size="small"
                @click.stop="handleDelete(book.id)"
              >
                删除
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

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
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Picture } from '@element-plus/icons-vue'
import { searchBooks, deleteBook } from '@/api/book'
import { borrowBook } from '@/api/borrow'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const books = ref([])
const total = ref(0)

const searchForm = reactive({
  keyword: '',
  category: ''
})

const pagination = reactive({
  page: 1,
  limit: 20
})

const loadBooks = async () => {
  loading.value = true
  try {
    const res = await searchBooks({
      token: userStore.token,
      keyword: searchForm.keyword || undefined,
      category: searchForm.category || undefined,
      page: pagination.page,
      limit: pagination.limit
    })
    books.value = res.data.books || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('加载图书列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadBooks()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.category = ''
  pagination.page = 1
  loadBooks()
}

const handleBookClick = (id) => {
  router.push({ name: 'BookDetail', params: { id } })
}

const handleBorrow = async (book) => {
  try {
    await ElMessageBox.confirm(`确认借阅《${book.title}》吗？`, '确认借阅', {
      type: 'warning'
    })
    
    await borrowBook({
      token: userStore.token,
      book_id: book.id
    })
    ElMessage.success('借阅成功')
    loadBooks()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('借阅失败:', error)
    }
  }
}

const handleEdit = (id) => {
  router.push({ name: 'BookEdit', params: { id } })
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确认删除该图书吗？', '确认删除', {
      type: 'warning'
    })
    
    await deleteBook({
      token: userStore.token,
      id
    })
    ElMessage.success('删除成功')
    loadBooks()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const handleSizeChange = () => {
  loadBooks()
}

const handlePageChange = () => {
  loadBooks()
}

onMounted(() => {
  loadBooks()
})
</script>

<style lang="scss" scoped>
.books-container {
  max-width: 1400px;
  margin: 0 auto;
}

.search-card {
  margin-bottom: 20px;

  .search-form {
    margin: 0;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.book-col {
  margin-bottom: 20px;
}

.book-card {
  cursor: pointer;
  transition: transform 0.3s;

  &:hover {
    transform: translateY(-5px);
  }

  .book-cover-wrapper {
    width: 100%;
    aspect-ratio: 2 / 3; /* 书籍封面标准比例 2:3 */
    margin-bottom: 15px;
    border-radius: 4px;
    overflow: hidden;
    background-color: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;

    .book-cover {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .book-cover-placeholder {
    width: 100%;
    aspect-ratio: 2 / 3; /* 书籍封面标准比例 2:3 */
    margin-bottom: 15px;
    border-radius: 4px;
    background-color: #f5f5f5;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #909399;
    font-size: 14px;

    .el-icon {
      font-size: 48px;
      margin-bottom: 8px;
    }
  }

  .book-info {
    margin-bottom: 15px;

    .book-title {
      font-size: 18px;
      font-weight: bold;
      margin-bottom: 10px;
      color: #303133;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .book-author,
    .book-isbn,
    .book-category {
      font-size: 14px;
      color: #606266;
      margin-bottom: 8px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .book-quantity {
      margin-top: 10px;
    }
  }

  .book-actions {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
  }
}

.pagination {
  margin-top: 20px;
  justify-content: center;
}

@media (max-width: 768px) {
  .search-form {
    :deep(.el-form-item) {
      margin-bottom: 10px;
    }
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
