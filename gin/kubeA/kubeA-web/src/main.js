import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

// 引入路由
import router from './router'
// 引入antd
import Antd from 'ant-design-vue'
// 引入暗黑风格主题
import 'ant-design-vue/dist/antd.dark.css'
// 引入icon图标
import * as Icons from '@ant-design/icons-vue'

// createApp(App).mount('#app')
const app = createApp(App)

// 图标组件全局注册
for (const i in Icons) {
    app.component(i, Icons[i])
}

// 注册
app.use(Antd)
app.use(router)

app.mount('#app')
