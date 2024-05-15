<script setup lang="ts">
import { ElTable, ElTableColumn, ElTooltip } from 'element-plus'
import { records } from '../store'

</script>

<template>
  <ElTable :data="records" stripe>
    <ElTableColumn prop="Time" label="Time" sortable></ElTableColumn>
    <ElTableColumn prop="LocalAddr" label="Local Address" sortable></ElTableColumn>
    <ElTableColumn prop="RemoteAddr" label="Remote Address" sortable></ElTableColumn>
    <ElTableColumn prop="Executable" label="Executable" sortable>
      <template #default="scope">
        <ElTooltip :content="scope.row.Executable">
          <span>{{ scope.row.Executable.match(/([^\\\/]+)$/)[0] }}</span>
        </ElTooltip>
      </template>
    </ElTableColumn>
    <ElTableColumn prop="Location" label="Location" sortable></ElTableColumn>
    <ElTableColumn prop="Reason" label="Reason" sortable></ElTableColumn>
    <ElTableColumn prop="Risk" label="Risk" sortable></ElTableColumn>
    <ElTableColumn prop="TIX" label="TIX" sortable></ElTableColumn>
    <ElTableColumn prop="Confidence" label="Confidence" sortable></ElTableColumn>
  </ElTable>
</template>