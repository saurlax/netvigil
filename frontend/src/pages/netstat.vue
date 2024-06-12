<script setup lang="ts">
import { ElTable, ElTableColumn, ElTooltip, TableColumnCtx } from 'element-plus'
import { netstats, Netstat } from '../utils'
import { computed } from 'vue'

const getFilter = (records: Netstat[], property: keyof Netstat) => {
  const values = new Set(records.map(record => record[property].toString()))
  return Array.from(values).map(value => ({ text: value, value: value }))
}

const filters = computed(() => {
  return {
    Executable: getFilter(netstats.value, 'Executable'),
    Location: getFilter(netstats.value, 'Location'),
    Reason: getFilter(netstats.value, 'Reason'),
    TIX: getFilter(netstats.value, 'TIX'),
  }
})

const filterHandler = (value: string, row: Netstat, column: TableColumnCtx<Netstat>) => {
  const property = column['property'] as keyof Netstat
  return row[property] === value
}

</script>

<template>
  <ElTable :data="netstats" stripe>
    <ElTableColumn prop="Time" label="时间" sortable>
      <template #default="scope">
        <span>{{ new Date(scope.row.Time).toLocaleString() }}</span>
      </template>
    </ElTableColumn>
    <ElTableColumn prop="LocalIP" label="本地地址" sortable :formatter="row => `${row.LocalIP}:${row.LocalPort}`">
    </ElTableColumn>
    <ElTableColumn prop="RemoteIP" label="远程地址" sortable :formatter="row => `${row.RemoteIP}:${row.RemotePort}`">
    </ElTableColumn>
    <ElTableColumn prop="Executable" label="发起程序" sortable :filters="filters.Executable" :filter-method="filterHandler">
      <template #default="scope">
        <ElTooltip :content="scope.row.Executable">
          <span>{{ scope.row.Executable.match(/([^\\\/]+)$/)[0] }}</span>
        </ElTooltip>
      </template>
    </ElTableColumn>
    <ElTableColumn prop="Location" label="位置" sortable :filters="filters.Location" :filter-method="filterHandler">
    </ElTableColumn>
    <ElTableColumn prop="Reason" label="原因" sortable :filters="filters.Reason" :filter-method="filterHandler">
    </ElTableColumn>
    <ElTableColumn prop="Risk" label="威胁度" sortable></ElTableColumn>
    <ElTableColumn prop="TIX" label="情报中心" sortable :filters="filters.TIX" :filter-method="filterHandler">
    </ElTableColumn>
    <ElTableColumn prop="Confidence" label="可信度" sortable></ElTableColumn>
  </ElTable>
</template>