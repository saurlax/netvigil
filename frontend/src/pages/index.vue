<script setup lang="ts">
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { computed, reactive } from 'vue'
import { records } from '../store'

use([
  CanvasRenderer,
  PieChart,
  TitleComponent,
  TooltipComponent
])

const sum = computed(() => {
  const sum = [0, 0, 0, 0, 0]
  records.value.forEach((record) => {
    sum[record.Risk]++
  })
  return sum
})

const option = computed(() => {
  return {
    title: {
      text: 'Network Traffic',
      left: 'center',
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b} : {c} ({d}%)',
    },
    series: [
      {
        name: 'Traffic Sources',
        type: 'pie',
        data: [
          { value: sum.value[0], name: 'Unknown' },
          { value: sum.value[1], name: 'Safe' },
          { value: sum.value[2], name: 'Normal' },
          { value: sum.value[3], name: 'Suspicious' },
          { value: sum.value[4], name: 'Malicious' }
        ],
      },
    ],
  }
})
</script>

<template>
  <v-chart class="chart" :option="option" autoresize />
  <div class="info">
    <div class="news">
      <div><b>安全新闻</b></div>
      <div>JSOutProx新版本针对APAC和MENA地区的金融服务和组织</div>
      <div>Tag：JSOutProx, Resecurity</div>
      <div>事件概述：</div>
      <div>
        安全公司Resecurity检测到JSOutProx的新版本，该版本针对亚太地区和中东北非地区的金融服务和组织。JSOutProx是一个复杂的攻击框架，利用JavaScript和.NET进行攻击。一旦执行，恶意软件使框架能够加载各种插件，对目标进行额外的恶意活动。该恶意软件首次在2019年被识别，最初被归因于SOLAR
        SPIDER的网络钓鱼活动，这些活动将JSOutProx RAT传递给非洲、中东、南亚和东南亚的金融机构。在2024年2月8日左右，该活动的高峰期，沙特阿拉伯的一家主要系统集成商报告了针对其一家主要银行的客户的事件。
      </div>
    </div>
    <div class="stats">
      <div><b>统计信息</b></div>
      <div>未知：{{ sum[0] }}</div>
      <div>安全：{{ sum[1] }}</div>
      <div>正常：{{ sum[2] }}</div>
      <div>可疑：{{ sum[3] }}</div>
      <div>恶意：{{ sum[4] }}</div>
    </div>
  </div>
</template>

<style scoped>
.chart {
  height: 400px;
}

.info {
  display: flex;
  justify-content: space-between;
  font-size: 0.9rem;
}

.news {
  max-width: 30%;
  font-size: 0.8rem;
}
</style>
