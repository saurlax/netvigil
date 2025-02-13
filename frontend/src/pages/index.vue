<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import 'echarts'
import 'echarts-gl'
import VChart from 'vue-echarts'

const stats = ref<any[]>([])
const sevenDayThreatPieChart = ref<{ name: string, value: number }[]>([])
const geoLocationFrequency = ref<Record<string, number>>({})
const ticFrequency = ref<Record<string, number>>({})

const mainOption = computed(() => {
  return {
    dataset: {
      source: stats.value,
    },
    color: ['#5470c6', '#91cc75', '#73c0de', '#fac858', '#ee6666'],
    legend: {},
    tooltip: {},
    series: [],
    globe: {
      baseTexture: '/assets/earth.webp',
      heightTexture: '/assets/height.webp',
      // environment: '/assets/starfield.webp',
      top: '10%',
      globeRadius: 60,
      displacementScale: 0.1,
      shading: 'lambert',
      layers: [
        {
          type: 'blend',
          blendTo: 'emission',
          texture: '/assets/night.webp'
        },
        {
          type: 'overlay',
          texture: '/assets/clouds.webp',
          shading: 'lambert',
          distance: 5
        }
      ]
    }
  }
})

onMounted(async () => {
  try {
    // 发送请求获取数据
    const response = await fetch('/api/stats')
    if (!response.ok) {
      throw new Error('Failed to fetch stats data')
    }
    const data = await response.json()

    sevenDayThreatPieChart.value = data.sevenDayThreatPieChart
    geoLocationFrequency.value = data.geoLocationFrequency
    ticFrequency.value = data.ticFrequency
  } catch (error) {
    console.error('Error fetching stats data:', error)
  }
})

const sevenDaysPieView = computed(() => {
  return {
    color: ['#5470c6', '#91cc75', '#73c0de', '#fac858', '#ee6666'],
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
    },
    series: [
      {
        name: '近七日威胁度',
        type: 'pie',
        radius: '55%',
        center: ['50%', '50%'],
        data: sevenDayThreatPieChart.value, // 绑定后端数据
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowOffsetY: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
      },
    ],
  }
})

const ticFrequencyPieView = computed(() => {
  const data = Object.keys(ticFrequency.value).map(key => {
    const value = ticFrequency.value[key];
    return {
      name: key,
      value: value,
    };
  });

  return {
    color: ['#5470c6', '#91cc75', '#73c0de', '#fac858', '#ee6666'],
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
    },
    series: [
      {
        name: '情报来源占比',
        type: 'pie',
        radius: '55%',
        center: ['50%', '50%'],
        data: data,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowOffsetY: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
      },
    ],
  };
});

const pieViewOption3 = computed(() => {
  return {
    dataset: {
      source: stats.value,
    }, backgroundColor: '',
    color: ['#5470c6', '#91cc75', '#73c0de', '#fac858', '#ee6666'],
    tooltip: {},
    series: [
      {
        type: 'pie',
        seriesLayoutBy: 'row',
        encode: { itemName: 0, value: 8 },
      },
    ],
  }
})

const barFrequencyOption = computed(() => {
  return {
    dataset: {
      source: stats.value,
    },
    color: ['#73c0de', '#fac858', '#ee6666'],
    tooltip: {},
    xAxis: { type: 'category' },
    yAxis: { type: 'value' },
    series: [
      {
        type: 'bar',
        encode: { x: 0, y: [3, 4, 5] }, // 使用“可疑、恶意”的统计数据
      },
    ],
  }
})

const barGeoRankingOption = computed(() => {
  const geoData = Object.entries(geoLocationFrequency.value)
    .filter(([location, _]) => location !== '' && location != null)
    .sort((a, b) => a[1] - b[1])

  return {
    dataset: {
      source: [
        ['地理位置', '频率'],
        ...geoData,
      ],
    },
    color: ['#91cc75'],
    tooltip: {},
    xAxis: { type: 'category' }, 
    yAxis: { type: 'value' },
    series: [
      {
        type: 'bar',
        encode: { x: 0, y: 1 },
      },
    ],
  }
})

const lineTrendOption = computed(() => {
  return {
    dataset: {
      source: stats.value,
    },
    color: ['#5470c6'],
    tooltip: {},
    xAxis: { type: 'category' },
    yAxis: { type: 'value' },
    series: [
      {
        type: 'line',
        encode: { x: 0, y: 8 }, // 使用趋势数据列
      },
    ],
  }
})
</script>

<template>
  <div class="panel">
    <VChart class="main-chart" :option="mainOption" theme="dark" autoresize />
    <div class="panel-title">
      <span>qwq</span>
      <h2>全球威胁态势</h2>
      <span>awa</span>
    </div>
    <div class="panel-stats">
      
    </div>
    <div class="subviews">
      <div class="subview-chart">
        <h3>近七日威胁度</h3>
        <VChart :option="sevenDaysPieView" theme="dark" autoresize />
      </div>
      <div class="subview-chart">
        <h3>可疑及以上威胁度的频率</h3>
        <VChart :option="barFrequencyOption" theme="dark" autoresize />
      </div>
      <div class="subview-chart">
        <h3>地理位置排名</h3>
        <VChart :option="barGeoRankingOption" theme="dark" autoresize />
      </div>
      <div class="subview-chart">
        <h3>威胁度走势</h3>
        <VChart :option="lineTrendOption" theme="dark" autoresize />
      </div>
      <div class="subview-chart">
        <h3>情报来源占比</h3>
        <VChart :option="ticFrequencyPieView" theme="dark" autoresize />
      </div>
      <div class="subview-chart">
        <h3>To be continue</h3>
        <VChart :option="pieViewOption3" theme="dark" autoresize />
      </div>
    </div>
  </div>
</template>

<style scoped>
.panel {
  position: relative;
  height: 100vh;
}

.panel-title {
  position: absolute;
  top: 0;
  width: 100%;
  height: 10%;
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-around;
  user-select: none;
}

.panel-title h2 {
  font-size: 1.5rem;
}

.panel-stats {
  position: absolute;
  top: 10%;
  left: 28%;
  right: 28%;
  height: 100px;
  margin: 0 10px;
  padding: 10px;
  color: gold;
  border: #4992f14d 1px solid;
  overflow-y: auto;
}

.main-chart {
  position: absolute;
  top: 0;
  height: 100%;
}

.subviews {
  position: absolute;
  top: 10%;
  width: 100%;
  height: 90%;
  display: grid;
  grid-template-columns: repeat(2, 28%);
  grid-template-rows: repeat(3, 30%);
  justify-content: space-between;
  gap: 20px;
  pointer-events: none;
}

.subview-chart {
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  background-image: radial-gradient(#4482c99c, transparent);
  border: #4992f14d 1px solid;
  pointer-events: auto;
}

.subview-chart h3 {
  margin: 0;
  padding: 10px 0;
  color: white;
  text-shadow: 0 1px 2px black;
  font-size: 16px;
  text-align: center;
  font-weight: normal;
  user-select: none;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin: 10px 0;
}

table th, table td {
  padding: 8px 12px;
  text-align: center;
  border: 1px solid #ddd;
}

table th {
  background-color: #2a2a2a;
  color: white;
}

table td {
  background-color: #333;
  color: white;
}

table tr:nth-child(even) {
  background-color: #444;
}

</style>
