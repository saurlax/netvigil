<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { ElAside, ElContainer, ElMain, ElMenu, ElMenuItem, ElMessage, ElMessageBox, ElScrollbar, } from 'element-plus'
import { RouterView, useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { user, netstats, threats } from './utils'

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
    axios.get('/api/netstats', {
      headers: {
        Authorization: `Bearer ${user.value.token}`
      }
    }).then(res => {
      netstats.value = res.data
    }).catch(e => {
      if (e.response.status === 401) {
        router.push('/login')
      } else {
        ElMessage.error(e.response.data.error ?? e.message)
      }
    })

    axios.get('/api/threats', {
      headers: {
        Authorization: `Bearer ${user.value.token}`
      }
    }).then(res => {
      threats.value = res.data
    }).catch(e => {
      if (e.response.status === 401) {
        router.push('/login')
      } else {
        ElMessage.error(e.response.data.error ?? e.message)
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
          <ElMenuItem index="home">可视大屏</ElMenuItem>
          <ElMenuItem index="threat">威胁情报</ElMenuItem>
          <ElMenuItem index="netstat">网络流量</ElMenuItem>
          <ElMenuItem index="config">配置文件</ElMenuItem>
          <ElMenuItem v-if="user" @click="ElMessageBox.confirm('确定要退出登录吗？').then(() => { router.push('/login') })">退出登录
          </ElMenuItem>
        </ElMenu>
      </ElScrollbar>
    </ElAside>
    <ElMain :class="route.name">
      <RouterView />
    </ElMain>
  </ElContainer>
</template>


<style scoped>
.wrapper {
  height: 100vh;
}

.el-aside {
  --el-aside-width: 180px;
  --el-menu-bg-color: transnparent;
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

.el-main.home {
  padding: 0;
}
</style>