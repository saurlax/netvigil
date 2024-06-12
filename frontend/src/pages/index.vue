<script setup lang="ts">
import { computed } from 'vue'
import 'echarts'
import 'echarts-gl'
import VChart from 'vue-echarts'

const option = computed(() => {
  return {
    globe: {
      baseTexture: '/assets/earth.webp',
      heightTexture: '/assets/height.webp',
      environment: '/assets/starfield.webp',
      globeRadius: 80,
      displacementScale: 0.1,
      shading: 'lambert',
      light: {
        main: {
          intensity: 1
        },
        ambient: {
          intensity: 0.2
        },
      },
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
</script>

<template>
  <VChart class="chart" :option="option" theme="dark" autoresize />
</template>

<style scoped>
.chart {
  width: 100%;
  height: 100vh;
}
</style>
