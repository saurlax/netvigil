import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './app.vue'
import ElementPlus from 'element-plus'
import Index from './pages/index.vue'
import Iprecords from './pages/iprecords.vue'

import 'element-plus/dist/index.css'
import './style.css'
import Tix from './pages/tix.vue'
import Login from './pages/login.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [{
    path: '/login',
    component: Login,
    name: 'login'
  }, {
    path: '/',
    component: Index,
    name: 'home'
  }, {
    path: '/iprecords',
    component: Iprecords,
    name: 'iprecords'
  }, {
    path: '/tix',
    component: Tix,
    name: 'tix'
  }]
})

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.mount('#app')
