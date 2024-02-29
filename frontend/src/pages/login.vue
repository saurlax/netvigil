<script setup lang="ts">
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus'
import { reactive } from 'vue';

const form = reactive({
  username: '',
  password: '',
})

const login = () => {
  ElMessage.info(`username: ${form.username}, password: ${form.password}`)
}
</script>
<template>
  <div class="wrapper">
    <ElForm>
      <ElFormItem>
        <img src="/logo.webp" alt="logo" />
      </ElFormItem>
      <ElFormItem>
        <ElInput v-model="form.username" placeholder="username" />
      </ElFormItem>
      <ElFormItem>
        <ElInput v-model="form.password" placeholder="password" />
      </ElFormItem>
      <ElButton type="primary" @click="login">Login</ElButton>
    </ElForm>
  </div>
</template>

<style scoped>
.wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  background-image: linear-gradient(to bottom right, #e2eefb, #b9cfe5);
}

img {
  width: 60%;
  margin: auto;
}

.el-form {
  width: 240px;
  padding: 20px;
  border-radius: 4px;
  background-color: #ffffff88;
}

.el-button {
  width: 100%;
}
</style>