import { createApp } from 'vue'
// import './style.css'
import App from './App.vue'
import './index.css'
import 'vue-sonner/style.css' // vue-sonner v2 requires this import


import router from "./router/index.ts";
import pinia from "./stores/index.ts";
// createApp(App).use(router).mount('#app')

// @ts-ignore 忽略 vue-i18n 模块的类型检查
import { createI18n } from 'vue-i18n'

const i18n = createI18n({
  locale: "en-US",
  legacy: false,
});

const app = createApp(App);
app.use(router);
app.use(pinia);
// Set i18n instance on app
app.use(i18n);

app.mount("#app");
