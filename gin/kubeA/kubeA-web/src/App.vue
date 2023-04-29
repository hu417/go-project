<template>
  <div>
    <!-- 路由占位符 -->
    <router-view></router-view>
  </div>
</template>

<script setup>

</script>

<style>
  html,body {
    /* 屏幕浏览器最大的宽和高 */
    width: 100vw;
    height: 100vh;
    /* 最小的宽 */
    min-width: 1425px;
  }

  #nprogress .bar {
    /* 自定义进度条颜色和宽度 */
    background-color: aqua;
    width: 3px;
  }
</style>






<!-- <script setup>
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div>
    <a href="https://vitejs.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo" />
    </a>
    <a href="https://vuejs.org/" target="_blank">
      <img src="./assets/vue.svg" class="logo vue" alt="Vue logo" />
    </a>
  </div>
  <HelloWorld msg="Vite + Vue" />
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style> -->
