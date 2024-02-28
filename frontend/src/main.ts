import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './app.vue'
import ElementPlus from 'element-plus'
import Index from './pages/index.vue'
import Iprecords from './pages/iprecords.vue'

import 'element-plus/dist/index.css'
import './style.css'

const router = createRouter({
  history: createWebHistory(),
  routes: [{
    path: '/',
    component: Index,
    name: 'home'
  }, {
    path: '/iprecords',
    component: Iprecords,
    name: 'iprecords'
  }]
})

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.mount('#app')
