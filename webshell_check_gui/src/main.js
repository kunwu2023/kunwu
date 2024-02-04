import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import elTableInfiniteScroll from 'el-table-infinite-scroll'

//全局引入css
// import 'element-plus/theme-chalk/dark/css-vars.css'
// import '../style/headtap.css'
import '../style/normalize.css'

import store from './store'

const app = createApp(App).use(store)
app.use(router)
app.use(ElementPlus, {locale: zhCn,})
app.use(elTableInfiniteScroll)
app.mount('#app')
