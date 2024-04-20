<script setup lang="ts">
import { ElTable, ElTableColumn } from 'element-plus'
import { records } from '../store'

</script>

<template>
  <ElTable :data="records" stripe>
    <ElTableColumn prop="Time" label="Time"></ElTableColumn>
    <ElTableColumn prop="LocalAddr" label="Local Address"></ElTableColumn>
    <ElTableColumn prop="RemoteAddr" label="Remote Address" sortable></ElTableColumn>
    <ElTableColumn prop="TIX" label="TIX" sortable></ElTableColumn>
    <ElTableColumn prop="Location" label="Location"></ElTableColumn>
    <ElTableColumn prop="Risk" label="Risk" sortable></ElTableColumn>
    <ElTableColumn prop="Reason" label="Reason"></ElTableColumn>
    <ElTableColumn prop="Executable" label="Executable"></ElTableColumn>
    <ElTableColumn prop="Confidence" label="Confidence"></ElTableColumn>
  </ElTable>
</template>