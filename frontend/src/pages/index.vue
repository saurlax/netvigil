<template>
  <div class="dashboard">
    <div class="chart-container">
      <v-chart class="chart" :option="option" :theme="'dark'" autoresize />
    </div>
    <div class="last7days">
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
      </div>
    </div>
    <div id="map" class="map-container"></div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { DatasetComponent, GridComponent, LegendComponent, TitleComponent, TooltipComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import moment from 'moment'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import markerIcon from 'leaflet/dist/images/marker-icon-2x.png'
import { fetchGeoLocation, records } from '../utils'
import { icon } from 'leaflet'

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
    未知: 0,
    安全: 0,
    普通: 0,
    可疑: 0,
    恶意: 0
  }
  for (const stat of statsChartData.value) {
    map.未知 += stat.Unknown
    map.安全 += stat.Safe
    map.普通 += stat.Normal
    map.可疑 += stat.Suspicious
    map.恶意 += stat.Malicious
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
    tooltip: {
      trigger: 'axis'
    },
    dataset: {
      source: statsChartData.value
    },
    xAxis: { type: 'category' },
    yAxis: {},
    grid: [{ top: '50%' }],
    color: ['#63dbe8', '#91cd77', '#fc8251', '#fffd55', '#ed1c24'],
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
      { type: 'line', smooth: true, seriesLayoutBy: 'row', color: '#63dbe8' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row', color: '#91cd77' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row', color: '#fc8251' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row', color: '#fffd55' },
      { type: 'line', smooth: true, seriesLayoutBy: 'row', color: '#ed1c24' },
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

onMounted(async () => {
  const initialZoomLevel = 2; // 初始缩放级别
  const southWest = L.latLng(-85, -180);
  const northEast = L.latLng(85, 180);
  const bounds = L.latLngBounds(southWest, northEast);

  const map = L.map('map', {
    attributionControl: false,
    zoom: initialZoomLevel, // 设置初始缩放级别
    minZoom: initialZoomLevel, // 设置最小缩放级别为初始缩放级别
    maxBounds: bounds, // 设置最大拖动边界
    maxBoundsViscosity: 1.0 // 设置最大拖动边界的粘滞度，值越大边界效果越明显
  }).setView([0, 0], initialZoomLevel);

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors'
  }).addTo(map);

  const defaultIcon = L.icon({
    iconUrl: markerIcon,
    iconAnchor: [10, 41]
  })

  // Fetch geolocations and add markers
  for (const record of records.value) {
    const geoLocation = await fetchGeoLocation(record.RemoteIP);
    const marker = L.marker([geoLocation.lat, geoLocation.lon], { icon: defaultIcon }).addTo(map);
    marker.bindPopup(`IP:${record.RemoteIP},风险等级: ${record.Risk}`).openPopup();
  }
  // const ip = '46.27.127.255';
  // const geoLocation = await fetchGeoLocation(ip);
  // const marker = L.marker([geoLocation.lat, geoLocation.lon]).addTo(map);
  // marker.bindPopup(`Ip: ${ip},风险等级: 2`).openPopup(); 
});
</script>

<style scoped>
.dashboard {
  padding: 20px;
  background-color: #0b082e;
  color: white;
  width: 100%;
  box-sizing: border-box;
  /* 添加这行以确保padding不会导致超出宽度 */
}

.chart-container,
.last7days,
.last7days-chart-container {
  margin: 20px 0;
  padding: 20px;
  border-radius: 8px;
  width: 100%;
  max-width: 100%;
  /* 确保宽度不会超过容器 */
  border: 2px solid #3a3f5c;
  /* 增加边框 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  /* 增加阴影 */
  box-sizing: border-box;
  /* 确保padding不会导致超出宽度 */
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
  border: 1px solid #3a3f5c;
  /* 增加边框 */
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
  height: 80vh;
  /* 使用视口高度，确保图表根据屏幕大小自适应 */
  width: 100%;
  max-width: 100%;
}

.last7days-chart {
  height: 40vh;
  /* 使用视口高度，确保图表根据屏幕大小自适应 */
  width: 100%;
  max-width: 100%;
}

.map-container {
  height: 80vh;
  /* 使用视口高度，确保地图根据屏幕大小自适应 */
  width: 100%;
  max-width: 100%;
  margin-top: 20px;
  border: 2px solid #3a3f5c;
  /* 增加边框 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  /* 增加阴影 */
}
</style>
