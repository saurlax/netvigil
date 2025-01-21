<script setup lang="tsx">
import { ElAutoResizer, ElButton, ElTableV2, ElMessageBox } from 'element-plus'
import { credibilityLevel, riskLevel, threats } from '../utils'
import { computed } from 'vue'
import { user } from '../utils'
import dayjs from 'dayjs'

const DelFireWall = async (id: number, ip: string) => {
  try {
    await ElMessageBox.confirm(
      `确定要解除对 ${ip} 的封锁吗？`,
      '删除威胁ip',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )

    const res = await fetch("/api/threats", {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${user.value?.token}` },
      body: JSON.stringify({ id: Number(id), action: "remove" })
    })

    const result = await res.json()
    if (res.status === 400) {
      ElMessageBox.error(`请求参数错误: ${result.error}`)
      return
    }

    if (result.success) {
      threats.value = threats.value.filter(n => n.id !== id)
      ElMessageBox.success(`成功删除 ${ip}`)
    } else {
      ElMessageBox.error(`删除失败: ${result.error}`)
    }
  } catch (err) {
    console.log(err)
  }
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
  width: 100,
  title: '操作',
  cellRenderer: ({ rowData }) => (<ElButton type="danger" size="small" onClick={() => DelFireWall(rowData.id, rowData.ip)}>删除</ElButton>)
}]

const data = computed(() => {
  if (!threats.value) return []
  return threats.value.map(n => ({
    ...n,
    time: dayjs(n.time).format('YYYY-MM-DD HH:mm:ss'),
    risk: `${riskLevel[n.risk]}(${n.risk})`,
    credibility: `${credibilityLevel[n.credibility]}(${n.credibility})`,
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