<template>
    <!-- 使用layout布局 -->
    <a-layout> 
        <!-- 使用affix固钉: 固定顶部div -->
        <a-affix>
            <a-layout-header>
                <!-- 这是header -->
                <!-- 平台信息 -->
                <!-- 使用float让元素同一行 -->
                <div style="float:left">
                    <img :src="kubeA" alt="" style="height: 40px;margin-bottom:10px">
                    <span style="font-size: 25px;padding: 0 25px">kubeA</span>
                </div>
                <!-- 集群信息 -->
                <!-- theme主题颜色; horizontal: 横向展开 -->
                <a-menu
                    style="float:left; width:250px; line-height: 64px"
                    v-model:selectedKey="selectedKey1"
                    theme="dark"
                    mode="horizontal"
                >
                    <a-menu-item v-for="item in clusterList" :key="item">
                        {{ item }}
                    </a-menu-item>
                </a-menu>

                <!-- 用户信息 -->
                <div style="float:right">
                    <img :src="avator" alt="" style="height:40px;border-radius:50%;margin-right:20px">
                    <!-- dropdown下拉菜单 -->
                    <!-- overlayStyle与admin的间距 -->
                    <a-dropdown :overlayStyle="{paddingTop: '10px'}">
                        <a href="#">
                            Admin
                            <down-outlined />
                        </a>
                        <template #overlay>
                            <a-menu>
                              <a-menu-item>
                                <a href="#">用户信息</a>
                              </a-menu-item>
                              <a-menu-item>
                                <a href="#" @click="chagePass">修改密码</a>
                              </a-menu-item>
                              <a-menu-item>
                                <a href="#" @click="logout">退出登录</a>
                              </a-menu-item>
                            </a-menu>
                          </template>
                    </a-dropdown>
                </div>
            </a-layout-header>
        </a-affix>
        <!-- 中部布局 -->
        <a-layout style="height:calc(100vh - 68px)">
            <!-- sider侧边栏 -->
            <!-- collapsed处理展开和收缩  -->
            <a-layout-sider style="width: 240px; background: rgb(31, 30, 30)" v-model:collapsed="collapsed" collapsible>   
                这是侧边栏
                <!-- selectedKeys表示点击选中的栏目,用于a-menu-item -->
                <!-- openKeys表示展开的栏目，用于a-sub-menu -->
                <!-- openChange事件监听 SubMenu 展开/关闭的回调 -->
                <!-- inline纵向 -->
                <a-menu 
                    :selectedKeys="selectedKey2"
                    :openKeys="openKeys"
                    @openChange="onOpenChange"
                    mode="inline"
                    :style="{
                        height:'100%',
                        boderRight: 0.5
                    }"
                    >
                    <template v-for="menu in routers" :key="menu">
                        <!-- 处理无子路由的情况 -->
                        <!-- routeChange用于跳转和选中 -->
                        <a-menu-item
                            v-if="menu.children && menu.children.lenght == 1"
                            :index="menu.children[0].path"
                            :key="menu.children[0].path"
                            @click="routeChange('item',menu.children[0].path)"
                            >
                            <template #icon>
                                <component :is="menu.children[0].icon"/>
                            </template>
                            <span> {{ menu.children[0].name }} </span>
                        </a-menu-item>
                        <!-- 处理有子路由的情况 -->
                        <a-sub-menu
                            v-else-if="menu.children && menu.children.lenght > 1"
                            :index="menu.path"
                            :key="menu.path"
                            @click="routeChange(sub,children.path)"
                            >
                            <template #icon>
                                <component :is="menu.icon"/>
                            </template>
                            <!-- 父级路由 -->
                            <template #title>
                                <div>
                                    <span :class="[collapsed ? 'is-callapse':'']"> {{ menu.name }} </span>
                                </div>
                            </template>
                            <!-- 子级路由 -->
                            <a-menu-item
                                v-for="child in menu.children"
                                :key="child.path"
                                :index="child.path"
                                @click="routeChange('sub',child.path)"
                                >
                                <span> {{ child.name }} </span>
                            </a-menu-item>
                        </a-sub-menu>
                    </template>
                </a-menu>
                  
            </a-layout-sider>
            <!-- main部分 -->
            <a-layout style="padding: 0 20px;">
                <!-- home主页 -->
                <a-layout-content
                :style="{
                    background: 'rgb(31, 30, 30)',
                    padding: '24px', 
                    margin: 0, 
                    minHeight: '280px', 
                    overflowY: 'auto'}">  
                    这是home主页
                    <router-view></router-view>
                </a-layout-content>
                <!-- 底部 -->
                <a-layout-footer style="text-align: center">
                    Vue3 ©2018 Created by Ant UED
                  </a-layout-footer>
            </a-layout>
        </a-layout>
    </a-layout>
</template>

<script>
import { ref,onMounted, initCustomFormatter } from "vue"
import kubeA from "@/assets/k8s-metrics.png"
import avator from "@/assets/avator.png"
import { useRouter } from "vue-router" // useRouter获取路由信息
export default ({
    setup() {
        
        const collapsed = ref(false)
        // 选择
        const selectedKey1 = ref([])
        // 集群
        const clusterList = ref([
            "Test1",
            "Test2",
        ])

        // 侧边栏 //
        // 路由信息
        const routers = ref([])  // 路由信息
        const selectedKey2 = ref([]) // 侧边栏信息
        const openKeys = ref([]) // 选择的信息
        const router = useRouter() // 使用useRouter()方法获取路由配置和当前页面的路由信息
        console.log(router)



//////// 方法区 ////////   

        // 导航栏点击切换页面，以及处理选择的情况
        function routeChange(type,path){
            // 判断点击是否为sub父栏目，如果不是，则关闭其他父栏目
            if(type != "sub"){
                openKeys.value = []
            }
            // 选择当前path对应的栏目，单独的item或者子item
            selectedKeys2.value = [path]
            // 页面跳转
            if(router.currentRoute.value.path != path){
                router.push(path) // 跳转
            }
        }

        // 
        function getRouter(val) {
            selectedKey2.value = [val[1].path]
            openKeys.value = [val[0].path]
        }
        // 专门用于sub的打开和关闭
        function onOpenChange(val) {
            // 匹配这个val(path)是否已经打开，如果没有，则打开，并关闭其他
            const latestOpenKey = val.find(key => openKeys.value.indexOf(key) == -1)
            openKeys.value = latestOpenKey ? [latestOpenKey]:[]

        }

        // 生命周期钩子
        onMounted(()=>{
            routers.value = router.options.routes
            getRouter(router.currentRoute.value.matched)
        })
        

        // 密码修改
        function chagePass(){
            console.log("修改密码")
        }
        // 退出登录
        function logout(){
            console.log("退出登录")
            alert("退出成功")
            // 移除用户名
            // 移除token
            // 跳转到登录页
        }
        return {
            collapsed,
            kubeA,
            avator,
            selectedKey1,
            clusterList,
            routers,
            selectedKey2,
            openKeys,
            router,
            routeChange,
            onOpenChange,
            getRouter,
            chagePass,
            logout,
        }
    }
})

</script>

<style>
    .ant-layout-header {
        padding: 0px 30px !important;
        color: rgb(239, 239, 239); 
    }

    .ant-layout-footer {
        padding: 5px 50px !important;
        color: rgb(239, 239, 239);
        
    }

    .is-callapse {
        display: none;
    }

</style>