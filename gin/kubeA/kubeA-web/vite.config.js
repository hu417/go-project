import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

import { join } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  //软链接
  resolve: {
    alias: {
        '@': join(__dirname, './src')
    }
  },
  // 服务配置
  server: {
    // 配置前端服务地址和端口
    host: '0.0.0.0',
    port: 8080,
    // 是否开启 https
    https: false,
    hmr: {
      overlay: false,
    },
    // 代理
    proxy: {
      // 代理所有 /api 的请求
      '/api': {
          // 代理请求之后的请求地址
          target: 'https://www.wuyabala.com/',
          // 跨域
          changeOrigin: true,
          // uri重写
          // rewrite: path => path.replace(/^\/api1/, '')
      }
    }
  }
})
