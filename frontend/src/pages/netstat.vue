<script setup lang="tsx">
import { ElAutoResizer, ElTableV2, ElTooltip, ElPagination, ElMessage } from 'element-plus'
import { useRouter } from 'vue-router';
import { computed, ref, watchEffect } from 'vue'
import dayjs from 'dayjs'
import { Netstat, user } from '../utils';
import axios from 'axios';

const router = useRouter()
const netstats = ref<Netstat[]>([])
const total = ref(0)
const page = ref(1)
const limit = ref(100)

watchEffect(() => {
  router.push({ query: { page: page.value, limit: limit.value } })
  axios.get(`/api/netstats?limit=${limit.value}&page=${page.value}`, {
    headers: {
      Authorization: `Bearer ${user.value?.token}`
    }
  }).then(res => {
    netstats.value = res.data.netstats
    total.value = res.data.total
  }).catch(e => {
    if (e.response.status === 401) {
      router.push('/login')
    } else {
      ElMessage.error(e.response.data.error ?? e.message)
    }
  })
})


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
  key: 'src',
  title: '本地地址',
  dataKey: 'src',
  width: 200
}, {
  key: 'dst',
  title: '远程地址',
  dataKey: 'dst',
  width: 200
}, {
  key: 'location',
  title: '位置',
  dataKey: 'location',
  width: 300
}, {
  key: 'executable',
  title: '发起程序',
  dataKey: 'executable',
  width: 300,
  cellRenderer: ({ cellData: executable }: { cellData: string }) => {
    const match = executable.match(/([^\\\/]+)$/)
    return <ElTooltip content={executable}>{match ? match[0] : executable}</ElTooltip>
  }
}]

const data = computed(() => netstats.value.map(n => {
  return {
    ...n,
    time: n.time ? dayjs(n.time).format('YYYY-MM-DD HH:mm:ss') : '未知时间',
    src: `${n.srcIP}:${n.srcPort}`,
    dst: `${n.dstIP}:${n.dstPort}`,
    location: n.location ?? '未知位置',
    executable: n.executable ?? '未知程序'
  }
}))
</script>

<template>
  <ElAutoResizer>
    <template #default="{ height, width }">
      <ElTableV2 :columns="columns" :data="data || []" :height="height - 50" :width="width" />
      <ElPagination v-model:current-page="page" v-model:page-size="limit" :page-sizes="[100, 200, 500]" :total="total"
        layout="total, sizes, prev, pager, next, jumper" />
    </template>
  </ElAutoResizer>
</template>
