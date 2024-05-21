<script setup lang="ts">
import axios from 'axios';
import { ElDescriptions, ElDescriptionsItem, ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { user } from '../utils'

interface Config {
  capture_interval: string
  check_interval: string
  username: string
  web: string
  tix: Record<string, any>[]
}

const router = useRouter()
const config = ref<Config>()

onMounted(() => {
  if (user.value) {
    axios.get('/api/config', {
      headers: {
        Authorization: `Bearer ${user.value.token}`
      }
    }).then(res => {
      config.value = res.data
    }).catch(e => {
      if (e.response.status === 401) {
        router.push('/login')
      } else {
        ElMessage.error(e.response.data.error)
      }
    })
  }
})
</script>



<template>
  <div class="config">
    <ElDescriptions title="基本设置">
      <ElDescriptionsItem label="账户">{{ config?.username }} </ElDescriptionsItem>
      <ElDescriptionsItem label="管理地址">{{ config?.web }} </ElDescriptionsItem>
      <ElDescriptionsItem label="抓包间隔">{{ config?.capture_interval }} </ElDescriptionsItem>
      <ElDescriptionsItem label="检测间隔">{{ config?.check_interval }} </ElDescriptionsItem>
    </ElDescriptions>
    <ElDescriptions v-for="t in config?.tix" :title="t.type">
      <ElDescriptionsItem v-for="[k, v] in Object.entries(t).filter(([k, _]) => k != 'type')" :label="k">{{ v }}
      </ElDescriptionsItem>
    </ElDescriptions>
  </div>
</template>

<style scoped>
.config {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>