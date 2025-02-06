<script setup lang="tsx">
import { ElButton, ElInput, ElMessage, ElSpace, ElTableV2 } from 'element-plus';
import { Client, user } from '../utils';
import { onMounted, ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

const router = useRouter();
const name = ref('');

const getAPIKey = () => {
  axios.get('/api/clients', {
    headers: {
      Authorization: `Bearer ${user.value?.token}`
    }
  }).then(res => {
    data.value = res.data;
  }).catch(e => {
    if (e.response.status === 401) {
      router.push('/login')
    } else {
      ElMessage.error(e.response.data.error ?? e.message)
    }
  })
}

const createAPIKey = (name: string) => {
  axios.post('/api/clients', { name }, {
    headers: {
      Authorization: `Bearer ${user.value?.token}`
    }
  }).then(() => {
    ElMessage.success('创建成功')
    getAPIKey();
  }).catch(e => {
    ElMessage.error(e.response.data.error ?? e.message)
  })
}

const deleteAPIKey = (apikey: string) => {
  axios.delete(`/api/clients/${apikey}`, {
    headers: {
      Authorization: `Bearer ${user.value?.token}`
    }
  }).then(() => {
    ElMessage.success('删除成功')
    getAPIKey();
  }).catch(e => {
    ElMessage.error(e.response.data.error ?? e.message)
  })
}

const data = ref<Client[]>([]);
const columns = [{
  key: 'name',
  title: '名称',
  dataKey: 'name',
  width: 300
}, {
  key: 'apikey',
  title: 'APIKey',
  dataKey: 'apikey',
  width: 400
}, {
  key: 'action',
  title: '操作',
  dataKey: 'apikey',
  width: 100,
  cellRenderer: ({ cellData: apikey }: { cellData: string }) =>
  (<ElButton type="danger" size="small" onClick={() => deleteAPIKey(apikey)}
  >删除</ElButton>)
}];

onMounted(getAPIKey);
</script>

<template>
  <ElAutoResizer>
    <template #default="{ height, width }">
      <ElSpace>
        <ElInput v-model="name" placeholder="APIKey 名称" />
        <ElButton type="primary" @click="() => createAPIKey(name)">创建</ElButton>
      </ElSpace>
      <ElTableV2 :columns :data :height="height - 50" :width="width" />
    </template>
  </ElAutoResizer>
</template>