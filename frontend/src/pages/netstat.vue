<script setup lang="tsx">
import { ElAutoResizer, ElTableV2, ElTooltip } from 'element-plus'
import { netstats } from '../utils'
import { computed } from 'vue'
import dayjs from 'dayjs'

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
    time: dayjs(n.time).format('YYYY-MM-DD HH:mm:ss'),
    src: `${n.srcIP}:${n.srcPort}`,
    dst: `${n.dstIP}:${n.dstPort}`,
  }
}))
</script>

<template>
  <ElAutoResizer>
    <template #default="{ height, width }">
      <ElTableV2 :columns :data :height :width />
    </template>
  </ElAutoResizer>
</template>