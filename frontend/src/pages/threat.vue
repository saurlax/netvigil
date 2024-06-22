<script setup lang="tsx">
import { ElAutoResizer, ElButton, ElTableV2 } from 'element-plus'
import { credibilityLevel, riskLevel, threats } from '../utils'
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
  width: 100,
  cellRenderer: () => (<ElButton type="danger" size="small">删除</ElButton>)
}]

const data = computed(() => threats.value.map(n => {
  return {
    ...n,
    time: dayjs(n.time).format('YYYY-MM-DD HH:mm:ss'),
    risk: `${riskLevel[n.risk]}(${n.risk})`,
    credibility: `${credibilityLevel[n.credibility]}(${n.credibility})`,
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