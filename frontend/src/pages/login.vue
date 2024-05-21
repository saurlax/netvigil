<script setup lang="ts">
import axios from 'axios'
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus'
import { reactive } from 'vue'
import { user } from '../store'
import { useRouter } from 'vue-router'

const router = useRouter()
const form = reactive({
  username: '',
  password: '',
})

const submit = async () => {
  axios.postForm('/api/login', form).then(res => {
    user.value = {
      username: form.username,
      token: res.data.token
    }
    localStorage.setItem('user', JSON.stringify(user.value))
    ElMessage.success('登录成功')
    router.push('/')
  }).catch(err => {
    ElMessage.error(err.response.data.error)
  })
}
</script>

<template>
  <div class="wrapper">
    <ElForm @submit.prevent="submit">
      <ElFormItem>
        <img src="/logo.webp" alt="logo" />
      </ElFormItem>
      <ElFormItem>
        <ElInput v-model="form.username" placeholder="账号" />
      </ElFormItem>
      <ElFormItem>
        <ElInput v-model="form.password" placeholder="密码" type="password" />
      </ElFormItem>
      <ElButton type="primary" native-type="submit">登录</ElButton>
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