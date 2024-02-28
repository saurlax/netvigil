<script setup lang="ts">
import { ElAside, ElContainer, ElMain, ElMenu, ElMenuItem, ElScrollbar, } from 'element-plus'
import { RouterView, useRouter } from 'vue-router'

const router = useRouter()

const navigate = (dest: string) => {
  router.push({ name: dest })
}
</script>

<template>
  <ElContainer class="container">
    <ElAside class="aside">
      <ElScrollbar>
        <div class="logo"></div>
        <ElMenu :default-active="router.currentRoute.value.name?.toString()" @select="navigate">
          <ElMenuItem index="home">Home</ElMenuItem>
          <ElMenuItem index="iprecords">IP Records</ElMenuItem>
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
.container {
  height: 100vh;
}

.aside {
  --el-aside-width: 180px;
  --el-menu-bg-color: trasnparent;
  background-color: #e9eef8;
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
  box-shadow: 0 0 1px #fdfdfd;
  border: 2px solid transparent;
}

.el-menu-item:not(:hover) {
  background-color: #fafbff;
}

.el-menu-item.is-active {
  background-color: white;
  border-color: var(--el-menu-active-color);
}
</style>