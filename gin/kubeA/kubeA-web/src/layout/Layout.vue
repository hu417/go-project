<template>
    <a-layout>
        <a-affix>
            <!-- 头部，固钉-->
            <a-layout-header>
                <!-- 平台名 -->
                <div style="float: left;">
                    <img style="height: 40px;margin-bottom: 10px;" :src="kubeLogo" />
                    <span
                        style="padding: 0 50px 0 20px;font-size: 25px;font-weight: bold;color: cornflowerblue">KubeCk</span>
                </div>
                <!-- 集群选择 -->
                <!-- theme主题颜色; horizontal: 横向展开 -->
                <a-menu style="float: left;width: 250px;" 
                    v-model:selectedKeys="selectedKeys1" 
                    theme="dark"
                    mode="horizontal" 
                    :style="{ lineHeight: '64px' }"
                >
                    <a-menu-item v-for="(item) in clusterList" :key="item">
                        {{ item }}
                    </a-menu-item>
                </a-menu>
                <!-- 用户信息 -->
                <div style="float: right;">
                    <img style="height: 40px;border-radius: 50%;margin-right: 14px;" :src="avator" />
                    <!-- :overlayStyle与admin的间距 -->
                    <!-- <a-dropdown :overlayStyle="{paddingTop: '10px'}"></a-dropdown> -->
                    <a-dropdown style="margin-top: 10px;">
                        <a class="ant-dropdown-link" @click.prevent>Admin <down-outlined /></a>
                        <template #overlay>
                            <a-menu>
                                <a-menu-item>
                                    <a @click="logout()">退出登录</a>
                                </a-menu-item>
                                <a-menu-item>
                                    <a href="javascript:;">修改密码</a>
                                </a-menu-item>
                            </a-menu>
                        </template>
                    </a-dropdown>
                </div>
            </a-layout-header>
        </a-affix>
        <a-layout style="height:calc(100vh - 68px)">
            <!-- 侧边栏 -->
            <!--  collapsed处理展开和收缩 -->
            <a-layout-sider style="width='240px';background: gainsboro" 
                v-model:collapsed="collapsed" 
                collapsible
            >
                <!--  selectedKeys表示点击选中的栏目,用于a-menu-item -->
                <!--  openKeys表示展开的栏目，用于a-sub-menu -->
                <!--  openChange事件监听 SubMenu 展开/关闭的回调 -->
                <a-menu :selectedKeys="selectedKey2" 
                    :openKeys="openKeys" 
                    @openChange="onOpenChange" 
                    mode="inline"
                    :style="{ height: '100%', borderRight: 0 }">
                    <template v-for="menu in routers" :key="menu">
                        <a-menu-item v-if="menu.children && menu.children.length == 1" 
                            :index="menu.children[0].path"
                            :key="menu.children[0].path" 
                            @click="routerChange('item', menu.children[0].path)"
                        >
                            <!-- routerChange()自定义方法 -->
                            <component :is="menu.children[0].icon" />
                            <span>{{ menu.children[0].name }}</span>
                        </a-menu-item>
                        <!-- 处理有子路由的情况，也就是有父栏目 -->
                        <a-sub-menu v-else-if="menu.children && menu.children.length > 1" 
                            :index="menu.path"
                            :key="menu.path"
                        >
                            <template #title>
                                <span>
                                    <component :is="menu.icon" />
                                    <span :class="[collapsed ? 'is-collapse' : '']">
                                        {{ menu.name }}
                                    </span>
                                </span>
                            </template>
                            <!-- 子路由的处理 -->
                            <a-menu-item v-for="child in menu.children" :key="child.path" :index="child.path"
                                @click="routerChange('sub', child.path)">
                                <span>{{ child.name }}</span>
                            </a-menu-item>
                        </a-sub-menu>
                    </template>
                </a-menu>
            </a-layout-sider>
            <a-layout style="padding: 0 24px">
                <!-- main部分-->
                <!-- 面包屑部分-->
                <a-breadcrumb style="margin: 12px 0">
                    <a-breadcrumb-item>工作台</a-breadcrumb-item>
                    <template v-for="(matched, index) in router.currentRoute.value.matched" :key="index">
                        <a-breadcrumb-item v-if="matched.name" href="{{router}}">
                            {{ matched.name }}
                        </a-breadcrumb-item>
                    </template>
                </a-breadcrumb>
                <a-layout-content :style="{
                        background: 'rgb(31,30,30)',
                        padding: '24px',
                        margin: 0,
                        minHeight: '280px',
                        overflowY: 'auto'
                    }">
                    <router-view></router-view>
                </a-layout-content>
                <!-- footer部分-->
                <a-layout-footer style="text-align: center">
                    ©2023 Created By Yi Devops
                </a-layout-footer>
            </a-layout>
        </a-layout>
    </a-layout>
</template>

<script>
import { ref, onMounted } from 'vue'
import kubeLogo from '@/assets/k8s-metrics.png'
import avator from '@/assets/avator.png'
import { useRouter } from "vue-router";


export default ({
    setup() {
        const collapsed = ref(false)
        const selectedKeys1 = ref([])
        const clusterList = ref([
            'TST-1',
            'TST-2'
        ])
        // 退出登录
        const logout = () => {
            localStorage.removeItem('username');
            localStorage.removeItem('token');
            router.push('/login')
        }
        // 侧边栏
        const routers = ref([])
        const selectedKey2 = ref([])
        const openKeys = ref([])
        const router = useRouter()
        const routerChange = (type, path) => {
            if (type != 'sub') {
                openKeys.value = []
            }
            selectedKey2.value = [path]
            // 判断当前路由地址是否与点击的path一致
            if (router.currentRoute.value.path != path) {
                // console.log(path)
                // 获取当前路由地址，并切换
                router.push(path)
            }
        }
        const getRouter = (val) => {
            // console.log(val)
            // 将请求的路由赋值到选择的节点上
            selectedKey2.value = [val[1].path]
            // openkeys是请求的父级
            openKeys.value = [val[0].path]
        }
        // 用于只显示选中的父栏目内容，关闭其他父栏目
        const onOpenChange = (val) => {
            const latestOpenKey = val.find(key => openKeys.value.indexOf(key) == -1);
            openKeys.value = latestOpenKey ? [latestOpenKey] : []
            //console.log(val)
        }
        // 生命周期钩子
        onMounted(() => {
            // 获取所有路由
            routers.value = router.options.routes
            // console.log(routers.value)

            // 获取当前路由
            // console.log(router.currentRoute.value.matched)
            getRouter(router.currentRoute.value.matched)
        })
        return {
            collapsed,
            kubeLogo,
            selectedKeys1,
            clusterList,
            avator,
            logout,
            routers,
            selectedKey2,
            openKeys,
            router,
            getRouter,
            routerChange,
            onOpenChange
        }

    }
})

</script>

<style>
.ant-layout-header {
    padding: 0 40px !important;
}

.ant-layout-content::-webkit-scrollbar {
    width: 6px;
}

.ant-layout-content::-webkit-scrollbar-track {
    background-color: rgb(164, 162, 162);
}

.ant-layout-content::-webkit-scrollbar-thumb {
    background-color: #666;
}

.ant-layout-footer {
    padding: 5px 50px !important;
    color: rgb(31, 30, 30);
}

.is-collapse {
    display: none;
}

.ant-layout-sider {
    background: #141414 !important;
    overflow-y: auto;
}

.ant-layout-sider::-webkit-scrollbar {
    display: none;
}

/* 下拉框与admin的距离 */
.ant-dropdown-content {
    margin-top: 10px;
}

.ant-menu-item {
    margin-top: 0 !important;
}</style>


