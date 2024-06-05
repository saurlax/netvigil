<template>
  <div class="tix-config">
    <ElForm label-position="top" :model="config">
      <ElFormItem label="用户名">
        <ElInput v-model="config.username" readonly></ElInput>
      </ElFormItem>
      <ElFormItem label="管理地址">
        <ElInput v-model="config.web" readonly></ElInput>
      </ElFormItem>
      <ElFormItem label="抓包间隔">
        <ElInput v-model="config.capture_interval" readonly></ElInput>
      </ElFormItem>
      <ElFormItem label="检测间隔">
        <ElInput v-model="config.check_interval" readonly></ElInput>
      </ElFormItem>
    </ElForm>

    <div v-for="t in config.tix" :key="t.type">
      <h3>{{ t.type }}</h3>
      <ElForm v-for="[k, v] in Object.entries(t).filter(([k, _]) => k != 'type')" :key="k" label-position="top">
        <ElFormItem :label="k">
          <ElInput :value="v" readonly></ElInput>
        </ElFormItem>
        <ElButton type="primary" @click="goToEdit">修改资料</ElButton>
      </ElForm>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus';
import { user } from '../utils';

interface Config {
  capture_interval: string;
  check_interval: string;
  username: string;
  web: string;
  tix: Record<string, any>[];
}

const router = useRouter();
const config = ref<Config>({
  capture_interval: '',
  check_interval: '',
  username: '',
  web: '',
  tix: []
});

onMounted(() => {
  if (user.value) {
    axios.get('/api/config', {
      headers: {
        Authorization: `Bearer ${user.value.token}`
      }
    }).then(res => {
      config.value = res.data;
    }).catch(e => {
      if (e.response.status === 401) {
        router.push('/login');
      } else {
        ElMessage.error(e.response.data.error);
      }
    });
  }
});

const goToEdit = () => {
  router.push('/edit');
};
</script>

<style scoped>
.tix-config {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

h3 {
  margin-top: 20px;
  font-weight: bold;
  color: #409EFF;
}
</style>
