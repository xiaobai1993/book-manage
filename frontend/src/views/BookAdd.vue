<template>
  <div class="book-add-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>添加图书</span>
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
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">提交</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { addBook } from '@/api/book'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const loading = ref(false)

const form = reactive({
  title: '',
  author: '',
  isbn: '',
  category: '',
  total_quantity: 1,
  description: ''
})

const rules = {
  title: [{ required: true, message: '请输入书名', trigger: 'blur' }],
  author: [{ required: true, message: '请输入作者', trigger: 'blur' }],
  isbn: [{ required: true, message: '请输入ISBN编号', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  total_quantity: [{ required: true, message: '请输入总数量', trigger: 'blur' }]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await addBook({
          token: userStore.token,
          ...form
        })
        ElMessage.success('添加成功')
        router.push({ name: 'Books' })
      } catch (error) {
        console.error('添加失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.book-add-container {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
