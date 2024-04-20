<script setup lang="ts">
  import axios from 'axios';
  import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus'
  import { reactive } from 'vue';
  import { useRouter } from 'vue-router';

  const form = reactive({
    username: '',
    password: '',
  })

  const router = useRouter();

  const handleSubmit = async () => {
    try {
      // const res = await axios.postForm('/api/login', form);
      // const response = res.data;
      const response = await axios({
        method: 'post',
        url: '/api/login',
        data: form,
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
      });

      if (response.status === 200) {
        // const data = await response.json();
        sessionStorage.setItem('user', JSON.stringify(response.data))
        ElMessage.success('Login Successfully!');
        router.push('/');
      } else {
        throw new Error('Login Failed!');
      }
    } catch (error) {
      console.log(error);
      const message = (error instanceof Error) ? error.message : 'Login failed';
      ElMessage.error(message);
    }
  }
</script>
<template>
  <div class="wrapper">
    <ElForm @submit.prevent="handleSubmit">
      <ElFormItem>
        <img src="/logo.webp" alt="logo" />
      </ElFormItem>
      <ElFormItem>
        <ElInput v-model="form.username" placeholder="username" />
      </ElFormItem>
      <ElFormItem>
        <ElInput v-model="form.password" placeholder="password" />
      </ElFormItem>
      <ElButton type="primary" native-type="submit">Login</ElButton>
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