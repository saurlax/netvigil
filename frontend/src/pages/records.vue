<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { ElTable, ElTableColumn } from 'element-plus'

const data = ref([])

const update = async () => {
  const res = await axios.get('/api/records')
  data.value = res.data
}

onMounted(update)
</script>

<template>
  <ElTable :data="data" stripe>
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