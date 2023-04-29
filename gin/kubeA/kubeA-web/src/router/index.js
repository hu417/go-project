// 初始化路由

// 导入router路由模式
import { createRouter, createWebHistory } from 'vue-router'


// 导入进度条
import nProgress from "nprogress";
import "nprogress/nprogress.css"

// 引入layout布局
import Layout from "@/layout/Layout.vue"

// 定义路由规则
const routes = [
    {
        path: '/',
        redirect: "/home",
    },
    {
        path: '/login',
        component: () => import('@/views/common/Login.vue'),  //视图组件
        meta: {  // meta元信息
            title: "登录", 
            requireAuth: false
        },
    },
    // 引入布局组件
    {   
        path: "/home",
        component: Layout,
        children: [
            {
                path: "/home",
                name: "概览",
                icon: "fund-outlined",
                meta: {title: "概览", requireAuth: true},
                component: () => import('@/views/home/Home.vue'),
            }
        ]
    },
    {
        path: "/cluster",
        name: "集群",
        component: Layout,
        icon: "cloud-server-outlined",
        children: [
            {
                path: "/cluster/node",
                name: "Node",
                meta: {title: "Node", requireAuth: true},
                component: () => import('@/views/cluster/Node.vue'),
            },
            {
                path: "/cluster/namespace",
                name: "Namespace",
                meta: {title: "Namespace", requireAuth: true},
                component: () => import('@/views/cluster/Namespace.vue'),
            },
            {
                path: "/cluster/pv",
                name: "PV",
                meta: {title: "PV", requireAuth: true},
                component: () => import('@/views/cluster/Pv.vue'),
            }
        ]
    },
    {
        path: "/workload",
        name: "工作负载",
        component: Layout,
        icon: "block-outlined",
        children: [
            {
                path: "/workload/pod",
                name: "Pod",
                meta: {title: "Pod", requireAuth: true},
                component: () => import('@/views/workload/Pod.vue'),
            },
            {
                path: "/workload/deployment",
                name: "Deployment",
                meta: {title: "Deployment", requireAuth: true},
                component: () => import('@/views/workload/Deployment.vue'),
            },
            {
                path: "/workload/daemonset",
                name: "DaemonSet",
                meta: {title: "DaemonSet", requireAuth: true},
                component: () => import('@/views/workload/DaemonSet.vue'),
            },
            {
                path: "/workload/statefulset",
                name: "StatefulSet",
                meta: {title: "StatefulSet", requireAuth: true},
                component: () => import('@/views/workload/StatefulSet.vue'),
            },
   
        ]
    },
]

// 创建路由实例
const router = createRouter({
    /**
     * hash模式：createWebHashHistory，路由后面+#
     * history模式：createWebHistory
     */
    history: createWebHistory(),
    routes,
})

// 进度条 //
// 定义进度条
nProgress.inc(100) 
// 进度条配置,easing: 动画字符串，speed: 动画速度，showSpinner: 进度环显示/隐藏
nProgress.configure({easing: 'ease',speed: 600, showSpinner: false})

// 结合路由守卫，开启和关闭进度条
router.beforeEach((to,form,next) =>{
    // 路由跳转之前，开启进度条
    nProgress.start()

    // 设置to页面的title
    if (to.meta.title) {
        document.title = to.meta.title
    } else {
        document.title = "kubeA"
    }
    
    // 放行
    next()
})

// 路由跳转完成后
router.afterEach(() => {
    // 关闭进度条
    nProgress.done()
})




// 暴露路由实例
export default router

