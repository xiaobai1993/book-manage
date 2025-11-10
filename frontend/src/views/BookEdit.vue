<template>
  <div class="book-edit-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>编辑图书</span>
          <el-button @click="$router.back()">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        v-if="!loading"
      >
        <el-form-item label="书名" prop="title">
          <el-input v-model="form.title" placeholder="请输入书名" />
        </el-form-item>
        <el-form-item label="作者" prop="author">
          <el-input v-model="form.author" placeholder="请输入作者" />
        </el-form-item>
        <el-form-item label="ISBN" prop="isbn">
          <el-input v-model="form.isbn" placeholder="请输入ISBN编号" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
            <el-option label="文学" value="文学" />
            <el-option label="科幻" value="科幻" />
            <el-option label="历史" value="历史" />
            <el-option label="童话" value="童话" />
            <el-option label="科技" value="科技" />
            <el-option label="教育" value="教育" />
          </el-select>
        </el-form-item>
        <el-form-item label="总数量" prop="total_quantity">
          <el-input-number v-model="form.total_quantity" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="图书描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="5"
            placeholder="请输入图书描述"
          />
        </el-form-item>
        <el-form-item label="封面图片">
          <div v-if="form.id === 0" class="upload-hint">
            <el-alert
              type="info"
              :closable="false"
              show-icon
            >
              <template #default>
                <span>正在加载图书信息...</span>
              </template>
            </el-alert>
          </div>
          <ImageUpload
            v-else
            :book-id="form.id"
            :current-image-url="form.cover_image_url"
            @uploaded="handleImageUploaded"
            @removed="handleImageRemoved"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">提交</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getBookDetail, editBook, deleteCover } from '@/api/book'
import ImageUpload from '@/components/ImageUpload.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const loading = ref(false)
const submitting = ref(false)

const form = reactive({
  id: 0,
  title: '',
  author: '',
  isbn: '',
  category: '',
  total_quantity: 1,
  description: '',
  cover_image_url: ''
})

const rules = {
  title: [{ required: true, message: '请输入书名', trigger: 'blur' }],
  author: [{ required: true, message: '请输入作者', trigger: 'blur' }],
  isbn: [{ required: true, message: '请输入ISBN编号', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  total_quantity: [{ required: true, message: '请输入总数量', trigger: 'blur' }]
}

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
    const book = res.data.book
    form.id = book.id
    form.title = book.title
    form.author = book.author
    form.isbn = book.isbn
    form.category = book.category
    form.total_quantity = book.total_quantity
    form.description = book.description || ''
    form.cover_image_url = book.cover_image_url || ''
  } catch (error) {
    console.error('加载图书详情失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        await editBook({
          token: userStore.token,
          ...form
        })
        ElMessage.success('编辑成功')
        router.push({ name: 'Books' })
      } catch (error) {
        console.error('编辑失败:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleImageUploaded = (imageUrl) => {
  form.cover_image_url = imageUrl
  ElMessage.success('封面图片上传成功')
}

const handleImageRemoved = async () => {
  try {
    await deleteCover({
      token: userStore.token,
      book_id: form.id
    })
    form.cover_image_url = ''
    ElMessage.success('封面图片已删除')
  } catch (error) {
    console.error('删除封面失败:', error)
  }
}

onMounted(() => {
  loadBookDetail()
})
</script>

<style lang="scss" scoped>
.book-edit-container {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
