<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { ElButton, ElTable, ElTableColumn } from 'element-plus'

const data = ref([])

const update = async () => {
  const res = await axios.get('/api/iprecords')
  data.value = res.data
  console.log(data)

}

onMounted(update)
</script>

<template>
  <ElTable :data="data" stripe>
    <ElTableColumn prop="ip" label="IP" sortable></ElTableColumn>
    <ElTableColumn prop="risk" label="Risk" sortable></ElTableColumn>
    <ElTableColumn prop="Reason" label="Reason" sortable></ElTableColumn>
    <ElTableColumn prop="description" label="Description"></ElTableColumn>
    <ElTableColumn prop="confirmed_by" label="Confirmed By" sortable></ElTableColumn>
    <ElTableColumn label="Operation">
      <template>
        <ElButton type="text" size="small">Edit</ElButton>
        <ElButton type="text" size="small">Delete</ElButton>
      </template>
    </ElTableColumn>
  </ElTable>
</template>