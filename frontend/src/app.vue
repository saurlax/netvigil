<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { ElAside, ElContainer, ElMain, ElMenu, ElMenuItem, ElMessage, ElMessageBox, ElScrollbar, } from 'element-plus'
import { RouterView, useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { user, records } from './utils'

const route = useRoute()
const router = useRouter()


onMounted(() => {
  const userStorage = localStorage.getItem('user')
  if (userStorage) {
    user.value = JSON.parse(userStorage)
  } else {
    router.push('/login')
  }
})

watch(user, () => {
  if (user.value) {
    axios.get('/api/records', {
      headers: {
        Authorization: `Bearer ${user.value.token}`
      }
    }).then(res => {
      records.value = res.data
    }).catch(e => {
      if (e.response.status === 401) {
        router.push('/login')
      } else {
        ElMessage.error(e.response.data.error)
      }
    })
  }
})

const navigate = (name: string) => {
  router.push({ name })
}
</script>

<template>
  <RouterView v-if="route.name == 'login'" />
  <ElContainer v-else class="wrapper">
    <ElAside>
      <ElScrollbar>
        <div class="logo"></div>
        <ElMenu :default-active="route.name?.toString()" @select="navigate">
          <ElMenuItem index="home">统计数据</ElMenuItem>
          <ElMenuItem index="records">情报记录</ElMenuItem>
          <ElMenuItem index="tix">情报中心</ElMenuItem>
          <ElMenuItem v-if="user" @click="ElMessageBox.confirm('确定要退出登录吗？').then(() => { router.push('/login') })">退出登录
          </ElMenuItem>
        </ElMenu>
      </ElScrollbar>
    </ElAside>
    <ElMain>
      <ElScrollbar>
        <RouterView />
      </ElScrollbar>
    </ElMain>
  </ElContainer>
</template>


<style scoped>
.wrapper {
  height: 100vh;
}

.el-aside {
  --el-aside-width: 180px;
  --el-menu-bg-color: trasnparent;
  background-image: linear-gradient(to bottom right, #e2eefb, #b9cfe5);
}

.logo {
  margin: 10px 20px;
  height: 50px;
  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  background-image: url('/logo.webp');
}

.el-menu {
  border: none;
  user-select: none;
}

.el-menu-item {
  margin: 8px;
  height: 50px;
  line-height: initial;
  border-radius: 4px;
  box-shadow: 0 0 1px #e3e3e3;
  border: 2px solid transparent;
}

.el-menu-item:not(:hover) {
  background-color: #ffffff88;
}

.el-menu-item.is-active {
  background-color: white;
  border-color: var(--el-menu-active-color);
}
</style>