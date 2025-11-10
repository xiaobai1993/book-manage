<template>
  <div class="image-upload">
    <div 
      class="upload-area" 
      :class="{ dragover: isDragging }"
      @click="triggerFileInput"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
      @drop.prevent="handleDrop"
    >
      <input
        ref="fileInputRef"
        type="file"
        accept="image/jpeg,image/png,image/webp,image/gif"
        style="display: none"
        @change="handleFileChange"
      />
      <img v-if="imageUrl" :src="imageUrl" alt="预览" class="preview-image" />
      <div v-else class="upload-placeholder">
        <el-icon class="upload-icon"><Plus /></el-icon>
        <div class="upload-text">
          <p>点击或拖拽上传图片</p>
          <p class="hint">支持 JPG、PNG、WebP、GIF，最大 5MB</p>
        </div>
      </div>
    </div>
    
    <div v-if="imageUrl" class="image-actions">
      <el-button size="small" @click="handleRemove">删除图片</el-button>
    </div>
    
    <div v-if="uploading" class="upload-progress">
      <el-progress :percentage="uploadProgress" :status="uploadStatus" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { uploadCover } from '@/api/book'

const props = defineProps({
  bookId: {
    type: Number,
    required: true
  },
  currentImageUrl: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['uploaded', 'removed'])

const userStore = useUserStore()
const fileInputRef = ref()
const imageUrl = ref(props.currentImageUrl || '')
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadStatus = ref('')
const isDragging = ref(false)

// 监听外部图片URL变化
watch(() => props.currentImageUrl, (newUrl) => {
  imageUrl.value = newUrl || ''
}, { immediate: true })

// 触发文件选择
const triggerFileInput = () => {
  fileInputRef.value?.click()
}

// 拖拽处理
const handleDragOver = (e) => {
  isDragging.value = true
  e.dataTransfer.dropEffect = 'copy'
}

const handleDragLeave = () => {
  isDragging.value = false
}

const handleDrop = (e) => {
  isDragging.value = false
  const file = e.dataTransfer.files?.[0]
  if (file) {
    processFile(file)
  }
}

// 处理文件选择
const handleFileChange = (event) => {
  const file = event.target.files?.[0]
  if (file) {
    processFile(file)
  }
  // 清空文件输入，允许重复选择同一文件
  event.target.value = ''
}

// 处理文件（验证和上传）
const processFile = (file) => {
  if (!file) return
  
  // 验证文件类型
  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.error('仅支持 JPG、PNG、WebP、GIF 格式')
    return
  }
  
  // 验证文件大小（5MB）
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB')
    return
  }
  
  // 预览图片
  const reader = new FileReader()
  reader.onload = (e) => {
    imageUrl.value = e.target.result
  }
  reader.readAsDataURL(file)
  
  // 上传文件
  uploadFile(file)
}

// 上传文件
const uploadFile = async (file) => {
  uploading.value = true
  uploadProgress.value = 0
  uploadStatus.value = ''
  
  try {
    const formData = new FormData()
    formData.append('token', userStore.token)
    formData.append('book_id', props.bookId)
    formData.append('image', file)
    
    // 模拟上传进度
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += 10
      }
    }, 200)
    
    const res = await uploadCover(formData)
    
    clearInterval(progressInterval)
    uploadProgress.value = 100
    uploadStatus.value = 'success'
    
    if (res.code === 0) {
      imageUrl.value = res.data.image_url
      ElMessage.success('上传成功')
      emit('uploaded', res.data.image_url)
      
      // 延迟隐藏进度条
      setTimeout(() => {
        uploading.value = false
        uploadProgress.value = 0
      }, 1000)
    }
  } catch (error) {
    uploading.value = false
    uploadProgress.value = 0
    uploadStatus.value = 'exception'
    ElMessage.error(error.message || '上传失败')
    
    // 恢复预览
    if (props.currentImageUrl) {
      imageUrl.value = props.currentImageUrl
    } else {
      imageUrl.value = ''
    }
  }
}

// 删除图片
const handleRemove = () => {
  imageUrl.value = ''
  emit('removed')
}

// 更新图片URL（供外部调用）
const updateImageUrl = (url) => {
  imageUrl.value = url || ''
}

defineExpose({
  updateImageUrl
})
</script>

<style lang="scss" scoped>
.image-upload {
  .upload-area {
    width: 100%;
    aspect-ratio: 2 / 3; /* 书籍封面标准比例 2:3 */
    min-height: 200px;
    border: 2px dashed #dcdfe6;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #fafafa;
    cursor: pointer;
    transition: border-color 0.3s;
    position: relative;
    
    &:hover {
      border-color: #409eff;
    }
    
    &.dragover {
      border-color: #409eff;
      background-color: #ecf5ff;
    }
    
    .preview-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
      border-radius: 4px;
    }
    
    .upload-placeholder {
      text-align: center;
      padding: 40px 20px;
      
      .upload-icon {
        font-size: 48px;
        color: #8c939d;
        margin-bottom: 16px;
      }
      
      .upload-text {
        p {
          margin: 8px 0;
          color: #606266;
          font-size: 14px;
          
          &.hint {
            color: #909399;
            font-size: 12px;
          }
        }
      }
    }
  }
  
  .image-actions {
    margin-top: 10px;
    text-align: center;
  }
  
  .upload-progress {
    margin-top: 10px;
  }
}
</style>

