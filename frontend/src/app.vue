<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElAside, ElContainer, ElMain, ElMenu, ElMenuItem, ElScrollbar, } from 'element-plus'
import { RouterView, useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const isLoggedIn = ref(false)

const updateloginStatus = () => {
  const user = sessionStorage.getItem('user');
  isLoggedIn.value = !!user;
  if(!isLoggedIn.value){
    router.push({ path: '/login' })
  }
}

const logout = () => {
  sessionStorage.removeItem('user');
  isLoggedIn.value = false
  router.push({ name: 'login' }) // Navigate to login route on logout
}

const navigate = (dest: string) => {
  router.push({ name: dest })
}

onMounted(() => {
  updateloginStatus();
});
</script>

<template>
  <RouterView v-if="route.name == 'login'" />
  <ElContainer v-else class="wrapper">
    <ElAside>
      <ElScrollbar>
        <div class="logo"></div>
        <ElMenu :default-active="route.name?.toString()" @select="navigate">
          <ElMenuItem index="home">Home</ElMenuItem>
          <ElMenuItem index="iprecords">IP Records</ElMenuItem>
          <ElMenuItem index="tix">TIX</ElMenuItem>
          <!-- <ElMenuItem v-if="!isLoggedIn" index="login" @click="login">Login</ElMenuItem> -->
          <ElMenuItem v-if="isLoggedIn" index="logout" @click="logout">Logout</ElMenuItem>
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