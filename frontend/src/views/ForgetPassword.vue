<template>
  <div class="forget-password-container">
    <div class="forget-password-box">
      <div class="forget-password-header">
        <el-icon><Reading /></el-icon>
        <h2>找回密码</h2>
      </div>
      <el-form
        ref="forgetPasswordFormRef"
        :model="forgetPasswordForm"
        :rules="forgetPasswordRules"
        class="forget-password-form"
      >
        <el-form-item prop="email">
          <el-input
            v-model="forgetPasswordForm.email"
            placeholder="请输入注册邮箱"
            size="large"
            :prefix-icon="Message"
          />
        </el-form-item>
        <el-form-item prop="code">
          <div class="code-input">
            <el-input
              v-model="forgetPasswordForm.code"
              placeholder="请输入验证码"
              size="large"
              :prefix-icon="Key"
            />
            <el-button
              :disabled="codeDisabled"
              :loading="codeLoading"
              @click="handleSendCode"
            >
              {{ codeDisabled ? `${countdown}秒后重试` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item prop="newPassword">
          <el-input
            v-model="forgetPasswordForm.newPassword"
            type="password"
            placeholder="请输入新密码（至少8位）"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item prop="confirmNewPassword">
          <el-input
            v-model="forgetPasswordForm.confirmNewPassword"
            type="password"
            placeholder="请确认新密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleForgetPassword"
            style="width: 100%"
          >
            重置密码
          </el-button>
        </el-form-item>
        <el-form-item>
          <div class="forget-password-links">
            <el-link type="primary" @click="$router.push('/login')">返回登录</el-link>
          </div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Reading, Message, Lock, Key } from '@element-plus/icons-vue'
import { forgetPassword, sendEmailCode } from '@/api/user'

const router = useRouter()

const forgetPasswordFormRef = ref()
const loading = ref(false)
const codeLoading = ref(false)
const codeDisabled = ref(false)
const countdown = ref(0)

const forgetPasswordForm = reactive({
  email: '',
  code: '',
  newPassword: '',
  confirmNewPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== forgetPasswordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const forgetPasswordRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度至少8位', trigger: 'blur' }
  ],
  confirmNewPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleSendCode = async () => {
  if (!forgetPasswordForm.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(forgetPasswordForm.email)) {
    ElMessage.warning('请输入正确的邮箱格式')
    return
  }

  codeLoading.value = true
  try {
    await sendEmailCode({ email: forgetPasswordForm.email, action: 'forget' })
    ElMessage.success('验证码已发送，请查收邮箱')
    codeDisabled.value = true
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        codeDisabled.value = false
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    console.error('发送验证码失败:', error)
  } finally {
    codeLoading.value = false
  }
}

const handleForgetPassword = async () => {
  if (!forgetPasswordFormRef.value) return

  await forgetPasswordFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await forgetPassword(forgetPasswordForm)
        ElMessage.success('密码重置成功，请登录')
        router.push('/login')
      } catch (error) {
        console.error('重置密码失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.forget-password-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.forget-password-box {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);

  .forget-password-header {
    text-align: center;
    margin-bottom: 30px;

    .el-icon {
      font-size: 48px;
      color: #409eff;
      margin-bottom: 10px;
    }

    h2 {
      color: #303133;
      margin: 0;
    }
  }

  .forget-password-form {
    .code-input {
      display: flex;
      gap: 10px;

      .el-input {
        flex: 1;
      }
    }

    .forget-password-links {
      width: 100%;
      text-align: center;
    }
  }
}

@media (max-width: 768px) {
  .forget-password-box {
    padding: 30px 20px;
  }
}
</style>
