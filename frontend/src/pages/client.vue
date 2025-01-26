<script setup lang="ts">
import { ElForm, ElFormItem, ElInput, ElSkeleton, ElSwitch, ElTag } from 'element-plus';
import { user } from '../utils';
import { ref } from 'vue';

const data = ref(null);

const columns = [];

</script>

<template>
  <ElAutoResizer>
    <template #default="{ height, width }">
      <ElTableV2 :columns="columns" :data="data || []" :height="height" :width="width" />
    </template>
  </ElAutoResizer>
</template>