import {createRouter, createWebHashHistory} from 'vue-router'

const routes = [
  {// 主界面
    path:'/',
    name:'index',
    component:()=>import(/*webpackChunkName:'Index'*/ '@/page/index/index.vue'),
    redirect:'/add_task', //默认显示页面
    //二级路由
    children:[
      {// add_task 添加检测任务
        path: '/add_task',
        name: 'add_task',
        component:()=>import(/*webpackChunkName:'add_task'*/ '@/page/check/add_task.vue'),
        redirect:'quick_check',
        // 三级路由
        children:[
          {// local_scan 本地扫描
            path: '/local_scan',
            name: 'local_scan',
            component:()=>import(/*webpackChunkName:'localScan'*/ '@/page/check/task/localScan.vue'),
          },{// remote_scan 远程扫描
            path: '/remote_scan',
            name: 'remote_scan',
            component:()=>import(/*webpackChunkName:'remoteScan'*/ '@/page/check/task/remoteScan.vue'),
          },{// advanced 高级选项
            path: '/advanced',
            name: 'advanced',
            component:()=>import(/*webpackChunkName:'advanced'*/ '@/page/check/task/advanced.vue'),
          },{// quick_check 快速检测
            path: '/quick_check',
            name: 'quick_check',
            component:()=>import(/*webpackChunkName:'quick_check'*/ '@/page/check/task/quick_check.vue'),
          },
        ]
      },
      {// see_task 查看任务
        path: '/see_task',
        name: 'see_task',
        component:()=>import(/*webpackChunkName:'see_task'*/ '@/page/check/see_task.vue'),
      }
    ],
  }
]


const router = createRouter({
  // 4. 内部提供了 history 模式的实现。为了简单起见，我们在这里使用 hash 模式。
  history: createWebHashHistory(),
  routes, // `routes: routes` 的缩写
})

export default router