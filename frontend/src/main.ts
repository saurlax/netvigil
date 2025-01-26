import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./app.vue";
import ElementPlus from "element-plus";
import zhCn from "element-plus/es/locale/lang/zh-cn";
import Login from "./pages/login.vue";
import Index from "./pages/index.vue";
import Threat from "./pages/threat.vue";
import Netstat from "./pages/netstat.vue";
import Config from "./pages/config.vue";

import "element-plus/dist/index.css";
import "./style.css";
import Client from "./pages/client.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      component: Login,
      name: "login",
    },
    {
      path: "/",
      component: Index,
      name: "home",
    },
    {
      path: "/threat",
      component: Threat,
      name: "threat",
    },
    {
      path: "/netstat",
      component: Netstat,
      name: "netstat",
    },
    {
      path: "/config",
      component: Config,
      name: "config",
    },
    {
      path: "/client",
      component: Client,
      name: "client",
    },
  ],
});

const app = createApp(App);
app.use(router);
app.use(ElementPlus, {
  locale: zhCn,
});
app.mount("#app");
