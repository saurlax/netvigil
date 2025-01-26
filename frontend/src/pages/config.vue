<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage, ElInputNumber, ElSwitch, ElSkeleton, ElSpace, ElMessageBox } from 'element-plus'
import { user } from '../utils'
import { useRouter } from 'vue-router';


const router = useRouter()
const config = ref()

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
        ElMessage.error(e.response.data.error ?? e.message)
      }
    })
  }
})

const update = () => {
  ElMessageBox.confirm('确定要更新配置文件吗？').then(() => {
    axios.post('/api/config', config.value, {
      headers: {
        Authorization: `Bearer ${user.value?.token}`
      }
    }).then(() => {
      ElMessage.success('更新成功，相关配置将会在重启后生效')
    }).catch(e => {
      if (e.response.status === 401) {
        router.push('/login')
      } else {
        ElMessage.error(e.response.data.error ?? e.message)
      }
    })
  })
}
</script>

<template>
  <ElForm label-position="left" label-width="auto">
    <div v-if="config">
      <ElFormItem label="检测间隔">
        <ElInput v-model="config.check_period" placeholder="5s" />
      </ElFormItem>
      <ElFormItem label="管理地址">
        <ElInput v-model="config.web" />
      </ElFormItem>
      <ElFormItem label="情报中心" v-for="t in config.tic">
        <ElSpace class="tics" direction="vertical" fill>
          <ElFormItem v-for="(v, k, _) in t" :label="k.toString()">
            <ElInput disabled v-if="typeof v === 'string'" :value="v" />
            <ElInputNumber disabled v-else-if="typeof v === 'number'" :value="v" />
            <ElSwitch disabled v-else-if="typeof v === 'boolean'" :value="v" />
            <ElInput disabled v-else :value="v" />
          </ElFormItem>
        </ElSpace>
      </ElFormItem>
    </div>
    <div v-else>
      <ElSkeleton />
    </div>
    <ElFormItem>
      <ElButton type="primary" @click="update">更新配置文件</ElButton>
    </ElFormItem>
  </ElForm>
</template>

<style scoped>
.tics {
  flex-grow: 1;
}
</style>
