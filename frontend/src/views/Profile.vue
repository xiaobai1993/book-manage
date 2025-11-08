<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>个人信息</span>
          </template>
          <el-descriptions :column="1" border v-loading="loading">
            <el-descriptions-item label="邮箱">{{ userInfo?.email || '-' }}</el-descriptions-item>
            <el-descriptions-item label="角色">
              <el-tag :type="userInfo?.role === 'admin' ? 'danger' : 'primary'">
                {{ userInfo?.role === 'admin' ? '管理员' : '普通用户' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="注册时间">{{ userInfo?.register_time || '-' }}</el-descriptions-item>
            <el-descriptions-item label="当前借阅数量">{{ currentBorrowCount || 0 }}</el-descriptions-item>
            <el-descriptions-item label="账户状态">
              <el-tag :type="userInfo?.status === 'normal' ? 'success' : 'danger'">
                {{ userInfo?.status === 'normal' ? '正常' : '已禁用' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <span>修改密码</span>
          </template>
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
          >
            <el-form-item label="原密码" prop="old_password">
              <el-input
                v-model="passwordForm.old_password"
                type="password"
                show-password
                placeholder="请输入原密码"
              />
            </el-form-item>
            <el-form-item label="新密码" prop="new_password">
              <el-input
                v-model="passwordForm.new_password"
                type="password"
                show-password
                placeholder="请输入新密码（至少8位）"
              />
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirm_new_password">
              <el-input
                v-model="passwordForm.confirm_new_password"
                type="password"
                show-password
                placeholder="请确认新密码"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="passwordLoading" @click="handleChangePassword">
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { getProfile, changePassword } from '@/api/user'

const userStore = useUserStore()

const loading = ref(false)
const passwordLoading = ref(false)
const userInfo = ref(null)
const currentBorrowCount = ref(0)

const passwordFormRef = ref()
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_new_password: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度至少8位', trigger: 'blur' }
  ],
  confirm_new_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const loadProfile = async () => {
  loading.value = true
  try {
    const res = await getProfile({ token: userStore.token })
    userInfo.value = res.data.user_info
    currentBorrowCount.value = res.data.current_borrow_count || 0
  } catch (error) {
    console.error('加载个人信息失败:', error)
  } finally {
    loading.value = false
  }
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      passwordLoading.value = true
      try {
        await changePassword({
          token: userStore.token,
          ...passwordForm
        })
        ElMessage.success('密码修改成功')
        passwordForm.old_password = ''
        passwordForm.new_password = ''
        passwordForm.confirm_new_password = ''
        passwordFormRef.value.resetFields()
      } catch (error) {
        console.error('修改密码失败:', error)
      } finally {
        passwordLoading.value = false
      }
    }
  })
}

onMounted(() => {
  loadProfile()
})
</script>

<style lang="scss" scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
}
</style>
