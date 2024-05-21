<script setup lang="ts">
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { DatasetComponent, GridComponent, LegendComponent, TitleComponent, TooltipComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { computed } from 'vue'
import { records } from '../utils'
import moment from 'moment'

use([
  CanvasRenderer,
  PieChart,
  LineChart,
  BarChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  DatasetComponent
])

const stats = computed(() => {
  const map: Record<string, number[]> = {}
  for (const record of records.value) {
    const day = moment(record.Time).format('YYYY-MM-DD')
    if (!map[day]) map[day] = [0, 0, 0, 0, 0]
    map[day][record.Risk]++
  }
  return Object.entries(map).map(([key, val]) => {
    return {
      Date: key,
      Unknown: val[0],
      Safe: val[1],
      Normal: val[2],
      Suspicious: val[3],
      Malicious: val[4]
    }
  })
})

const total = computed(() => {
  const map = {
    Unknown: 0,
    Safe: 0,
    Normal: 0,
    Suspicious: 0,
    Malicious: 0
  }
  for (const stat of stats.value) {
    map.Unknown += stat.Unknown
    map.Safe += stat.Safe
    map.Normal += stat.Normal
    map.Suspicious += stat.Suspicious
    map.Malicious += stat.Malicious
  }
  return Object.entries(map).map(([key, val]) => {
    return {
      name: key,
      value: val
    }
  })
})

const option = computed(() => {
  return {
    legend: {},
    tooltip: {
      trigger: 'axis'
    },
    dataset: {
      source: stats.value
    },
    xAxis: { type: 'category' },
    yAxis: {},
    grid: [{ top: '50%' }],
    series: [
      {
        type: 'pie',
        radius: '30%',
        data: total.value,
        center: ['50%', '25%']
      },
      { type: 'line', smooth: true, seriesLayoutBy: 'row' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row' },
    ]
  }
})
</script>

<template>
  <v-chart class="chart" :option="option" autoresize />
</template>

<style scoped>
.chart {
  height: 800px;
}
</style>