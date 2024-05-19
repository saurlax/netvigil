<template>
  <div ref="chartRef" class="chart"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import * as echarts from 'echarts';

// 定义chartRef和myChart变量
const chartRef = ref<HTMLElement | null>(null);
let myChart: echarts.ECharts | null = null;

const resizeChart = () => {
  if (myChart) {
    myChart.resize();
  }
};

onMounted(() => {
  if (!chartRef.value) {
    console.error("Chart DOM element not found.");
    return;
  }

  // 初始化ECharts实例
  myChart = echarts.init(chartRef.value, 'dark');

  const option: echarts.EChartsOption = {
    legend: {},
    tooltip: {
      trigger: 'axis',
      showContent: false
    },
    dataset: {
      source: [
        ['Date', '2024.5.13/Mon', '2024.5.14/Tue', '2024.5.15/Wed', '2024.5.16/Thu', '2024.5.17/Fri', '2024.5.18/Sat', '2024.5.19/Sun'],
        ['Unknown', 56.5, 82.1, 88.7, 70.1, 53.4, 85.1, 92],
        ['Safe', 51.1, 51.4, 55.1, 53.3, 73.8, 68.7, 93],
        ['Normal', 40.1, 62.2, 69.5, 36.4, 45.2, 32.5, 94],
        ['Suspicious', 25.2, 37.1, 41.2, 18, 33.9, 49.1, 95],
        ['Malicious', 25, 75, 23, 46, 43, 34, 12]
      ]
    },
    xAxis: { type: 'category' },
    yAxis: { gridIndex: 0 },
    grid: { top: '55%' },
    series: [
      {
        type: 'line',
        smooth: true,
        seriesLayoutBy: 'row',
        emphasis: { focus: 'series' }
      },
      {
        type: 'line',
        smooth: true,
        seriesLayoutBy: 'row',
        emphasis: { focus: 'series' }
      },
      {
        type: 'line',
        smooth: true,
        seriesLayoutBy: 'row',
        emphasis: { focus: 'series' }
      },
      {
        type: 'line',
        smooth: true,
        seriesLayoutBy: 'row',
        emphasis: { focus: 'series' }
      },
      {
        type: 'line',
        smooth: true,
        seriesLayoutBy: 'row',
        emphasis: { focus: 'series' }
      },
      {
        type: 'pie',
        id: 'pie',
        radius: '30%',
        center: ['50%', '25%'],
        emphasis: {
          focus: 'self'
        },
        label: {
          formatter: '{b}: {@2024.5.13/Mon} ({d}%)'
        },
        encode: {
          itemName: 'Date',
          value: '2024.5.13/Mon',
          tooltip: '2024.5.13/Mon'
        }
      }
    ]
  };

  myChart.setOption(option);

  myChart.on('updateAxisPointer', function (event: any) {
    const xAxisInfo = event.axesInfo[0];
    if (xAxisInfo) {
      const dimension = xAxisInfo.value + 1;
      myChart?.setOption<echarts.EChartsOption>({
        series: {
          id: 'pie',
          label: {
            formatter: '{b}: {@[' + dimension + ']} ({d}%)'
          },
          encode: {
            value: dimension,
            tooltip: dimension
          }
        }
      });
    }
  });

  // 添加窗口大小变化事件监听器
  window.addEventListener('resize', resizeChart);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart);
  if (myChart) {
    myChart.dispose();
  }
});
</script>

<style scoped>
.chart {
  width: 100%;
  height: 100%;
  min-height: 100vh;
  min-width: 60vw;
}
</style>
