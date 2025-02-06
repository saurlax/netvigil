<script setup lang="tsx">
import { ElAutoResizer, ElButton, ElTableV2, ElMessageBox, ElMessage } from 'element-plus'
import { Threat } from '../utils'
import { computed, onMounted, ref } from 'vue'
import { user } from '../utils'
import dayjs from 'dayjs'
import axios from 'axios'
import { useRouter } from 'vue-router'

const riskLevel = ["未知", "安全", "正常", "可疑", "恶意"];
const credibilityLevel = ["低", "中", "高"];

const threats = ref<Threat[]>([])
const router = useRouter()

onMounted(() => {
  axios.get('/api/threats', {
    headers: { Authorization: `Bearer ${user.value?.token}` }
  }).then(res => {
    threats.value = res.data
  }).catch(e => {
    if (e.response.status === 401) {
      router.push('/login')
    } else {
      ElMessage.error(e.response.data.error ?? e.message)
    }
  })
})

const deleteFireWallRule = async (ip: string) => {
  ElMessageBox.confirm('确定删除该威胁记录吗?', '提示', {
    type: 'warning'
  }).then(() => {
    axios.delete(`/api/threats/${ip}`, {
      headers: { Authorization: `Bearer ${user.value?.token}` },
    }).then(() => {
      ElMessage.success('删除成功')
      threats.value = threats.value.filter(n => n.ip !== ip)
    }).catch(e => {
      if (e.response.status === 401) {
        router.push('/login')
      } else {
        ElMessage.error(e.response.data.error ?? e.message)
      }
    })
  })
}

const columns = [{
  key: 'id',
  title: 'ID',
  dataKey: 'id',
  width: 100
}, {
  key: 'time',
  title: '记录时间',
  dataKey: 'time',
  width: 200
}, {
  key: 'ip',
  title: 'IP',
  dataKey: 'ip',
  width: 200
}, {
  key: 'tic',
  title: '情报来源',
  dataKey: 'tic',
  width: 200
}, {
  key: 'reason',
  title: '原因',
  dataKey: 'reason',
  width: 300
}, {
  key: 'risk',
  title: '威胁度',
  dataKey: 'risk',
  width: 100
}, {
  key: 'credibility',
  title: '可信度',
  dataKey: 'credibility',
  width: 100
}, {
  key: 'action',
  title: '操作',
  dataKey: 'ip',
  width: 100,
  cellRenderer: ({ cellData: ip }: { cellData: string }) => (
    <ElButton type="danger" size="small" onClick={() => deleteFireWallRule(ip)}>删除</ElButton>
  )
}]

const data = computed(() => {
  return threats.value?.map(t => ({
    ...t,
    time: dayjs(t.time).format('YYYY-MM-DD HH:mm:ss'),
    risk: `${riskLevel[t.risk]}(${t.risk})`,
    credibility: `${credibilityLevel[t.credibility]}(${t.credibility})`,
  }))
})
</script>

<template>
  <ElAutoResizer>
    <template #default="{ height, width }">
      <ElTableV2 :columns="columns" :data="data || []" :height="height" :width="width" />
    </template>
  </ElAutoResizer>
</template>