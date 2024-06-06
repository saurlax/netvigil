<template>
  <div class="dashboard">
    <div class="chart-container">
      <v-chart class="chart" :option="option" :theme="'dark'" autoresize />
    </div>
    <div class="last7days">
      <v-chart class="chart" :option="option" :theme="'dark'" autoresize />
      <div>
        <h2 class="title">近七天的监测数据</h2>
        <div class="card-container">
          <div class="card" v-for="(item, index) in last7DaysStats" :key="index">
            <div class="card-title">{{ item.title }}</div>
            <div class="card-value">{{ item.value }} <span>个</span></div>
          </div>
        </div>
        <div class="last7days-chart-container">
          <v-chart class="last7days-chart" :option="last7DaysOption" :theme="'dark'" autoresize />
        </div>
        <v-chart class="chart" :option="last7DaysOption" :theme="'dark'" autoresize />
      </div>
    </div>
  </div>
</template>

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

  const statsChartData = computed(() => {
    const map: Record<string, number[]> = {}
    for (const record of records.value) {
      const day = moment(record.Time).format('YYYY-MM-DD')
      if (!map[day]) map[day] = [0, 0, 0, 0, 0]
      map[day][record.Risk]++
    }
    return Object.entries(map)
      .sort(([a], [b]) => new Date(a).getTime() - new Date(b).getTime())
      .map(([key, val]) => {
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
    for (const stat of statsChartData.value) {
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
      title: {
        text: '总计数据图表',
        left: 'left'
      },
      legend: {
      },
      tooltip: {
        trigger: 'axis'
      },
      dataset: {
        source: statsChartData.value
      },
      xAxis: { type: 'category' },
      yAxis: {},
      grid: [{ top: '50%' }],
      series: [
        {
          type: 'pie',
          radius: '30%',
          data: total.value,
          center: ['50%', '30%'],
          label: {
            formatter: '{b}: {d}%'
          }
        },
        { type: 'line', smooth: true, seriesLayoutBy: 'row' },
        { type: 'line', smooth: true, seriesLayoutBy: 'row' },
        { type: 'line', smooth: true, seriesLayoutBy: 'row' },
        { type: 'line', smooth: true, seriesLayoutBy: 'row' },
        { type: 'line', smooth: true, seriesLayoutBy: 'row' },
      ]
    }
  })

  //七天监测数据
  const last7DaysStats = computed(() => {
    const map = {
      '未知': 0,
      '安全': 0,
      '普通': 0,
      '可疑': 0,
      '恶意': 0
    }

    const last7Days = records.value.filter(record => {
      return moment().diff(moment(record.Time), 'days') <= 7
    })

    for (const record of last7Days) {
      if (record.Risk === 0) map['未知']++
      if (record.Risk === 1) map['安全']++
      if (record.Risk === 2) map['普通']++
      if (record.Risk === 3) map['可疑']++
      if (record.Risk === 4) map['恶意']++
    }

    return [
      { title: '未知', value: map['未知'], trend: 'up' },
      { title: '安全', value: map['安全'], trend: 'up' },
      { title: '普通', value: map['普通'], trend: 'up' },
      { title: '可疑', value: map['可疑'], trend: 'up' },
      { title: '恶意', value: map['恶意'], trend: 'up' }
    ]
  })

  const last7DaysOption = computed(() => {
    return {
      title: {
        text: '近7天数据',
        left: 'left'
      },
      legend: {},
      tooltip: {
        trigger: 'axis'
      },
      dataset: {
        source: last7DaysStats.value
      },
      xAxis: { type: 'category' },
      yAxis: {},
      grid: [{ top: '20%' }],
      series: [
        { type: 'bar', smooth: true, seriesLayoutBy: 'row' }
      ]
    }
  })
</script>

<style scoped>
  .dashboard {
    padding: 20px;
    background-color: #0b082e;
    color: white;
    width: 100%
  }

  .chart-container,
  .last7days,
  .last7days-chart-container {
    margin: 20px 0;
    padding: 20px;
    border-radius: 8px;
    width: 100%;
  }

  .title {
    font-size: 24px;
    text-align: center;
    margin-bottom: 20px;
  }

  .card-container {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-around;
    margin-bottom: 20px;
  }

  .card {
    background-color: #324057;
    border-radius: 8px;
    padding: 20px;
    text-align: center;
    width: 150px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin: 10px;
  }

  .card-title {
    font-size: 16px;
    margin-bottom: 10px;
  }

  .card-value {
    font-size: 24px;
    font-weight: bold;
  }

  .chart {
    height: 800px;
    width: 100%;
  }

  .last7days-chart {
    height: 400px;
    width: 100%;
  }
</style>
