<template>
  <el-container class="main-layout">
    <el-header class="header">
      <div class="header-content">
        <div class="logo">
          <el-icon><Reading /></el-icon>
          <span>图书管理系统</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          class="header-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="books">
            <el-icon><Collection /></el-icon>
            <span>图书列表</span>
          </el-menu-item>
          <el-menu-item index="my-borrows">
            <el-icon><Document /></el-icon>
            <span>我的借阅</span>
          </el-menu-item>
          <el-menu-item v-if="userStore.isAdmin" index="all-borrows">
            <el-icon><List /></el-icon>
            <span>全部借阅记录</span>
          </el-menu-item>
          <el-menu-item v-if="userStore.isAdmin" index="book-add">
            <el-icon><Plus /></el-icon>
            <span>添加图书</span>
          </el-menu-item>
        </el-menu>
        <div class="user-info">
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-icon><User /></el-icon>
              <span>{{ userStore.email }}</span>
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><Setting /></el-icon>
                  个人信息
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-header>
    <el-main class="main-content">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import {
  Reading,
  Collection,
  Document,
  List,
  Plus,
  User,
  ArrowDown,
  Setting,
  SwitchButton
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => {
  return route.name === 'BookDetail' || route.name === 'BookEdit' || route.name === 'BookAdd'
    ? 'books'
    : route.name?.toLowerCase().replace(/-/g, '-') || 'books'
})

const handleMenuSelect = (key) => {
  router.push({ name: key })
}

const handleCommand = (command) => {
  if (command === 'profile') {
    router.push({ name: 'Profile' })
  } else if (command === 'logout') {
    userStore.logout()
    ElMessage.success('已退出登录')
    router.push({ name: 'Login' })
  }
}
</script>

<style lang="scss" scoped>
.main-layout {
  min-height: 100vh;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  .header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100%;
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 20px;

    .logo {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 20px;
      font-weight: bold;
      color: #409eff;
      margin-right: 40px;

      .el-icon {
        font-size: 24px;
      }
    }

    .header-menu {
      flex: 1;
      border-bottom: none;
    }

    .user-info {
      margin-left: 20px;

      .user-dropdown {
        display: flex;
        align-items: center;
        gap: 8px;
        cursor: pointer;
        color: #606266;
        padding: 8px 12px;
        border-radius: 4px;
        transition: background-color 0.3s;

        &:hover {
          background-color: #f5f7fa;
        }
      }
    }
  }
}

.main-content {
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
  padding: 20px;

  @media (max-width: 768px) {
    padding: 15px;
  }
}

@media (max-width: 768px) {
  .header .header-content {
    flex-direction: column;
    height: auto;
    padding: 10px;

    .logo {
      margin-right: 0;
      margin-bottom: 10px;
    }

    .header-menu {
      width: 100%;
    }

    .user-info {
      margin-left: 0;
      margin-top: 10px;
    }
  }
}
</style>
