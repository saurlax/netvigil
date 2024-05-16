<script setup lang="ts">
import { ElTable, ElTableColumn, ElTooltip, TableColumnCtx } from 'element-plus'
import { records, IPRecord } from '../store'
import { computed } from 'vue'

const getFilter = (records: IPRecord[], property: keyof IPRecord) => {
  const values = new Set(records.map(record => record[property].toString()))
  return Array.from(values).map(value => ({ text: value, value: value }))
}

const filters = computed(() => {
  return {
    Executable: getFilter(records.value, 'Executable'),
    Location: getFilter(records.value, 'Location'),
    Reason: getFilter(records.value, 'Reason'),
    TIX: getFilter(records.value, 'TIX'),
  }
})

const filterHandler = (value: string, row: IPRecord, column: TableColumnCtx<IPRecord>) => {
  const property = column['property'] as keyof IPRecord
  return row[property] === value
}

</script>

<template>
  <ElTable :data="records" stripe>
    <ElTableColumn prop="Time" label="Time" sortable>
      <template #default="scope">
        <span>{{ new Date(scope.row.Time).toLocaleString() }}</span>
      </template>
    </ElTableColumn>
    <ElTableColumn prop="LocalAddr" label="Local Address" sortable></ElTableColumn>
    <ElTableColumn prop="RemoteAddr" label="Remote Address" sortable></ElTableColumn>
    <ElTableColumn prop="Executable" label="Executable" sortable :filters="filters.Executable"
      :filter-method="filterHandler">
      <template #default="scope">
        <ElTooltip :content="scope.row.Executable">
          <span>{{ scope.row.Executable.match(/([^\\\/]+)$/)[0] }}</span>
        </ElTooltip>
      </template>
    </ElTableColumn>
    <ElTableColumn prop="Location" label="Location" sortable :filters="filters.Location" :filter-method="filterHandler">
    </ElTableColumn>
    <ElTableColumn prop="Reason" label="Reason" sortable :filters="filters.Reason" :filter-method="filterHandler">
    </ElTableColumn>
    <ElTableColumn prop="Risk" label="Risk" sortable></ElTableColumn>
    <ElTableColumn prop="TIX" label="TIX" sortable :filters="filters.TIX" :filter-method="filterHandler">
    </ElTableColumn>
    <ElTableColumn prop="Confidence" label="Confidence" sortable></ElTableColumn>
  </ElTable>
</template>