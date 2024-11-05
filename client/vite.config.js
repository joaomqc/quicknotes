import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const devApiUrl = "http://192.168.1.211:8000";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    port: 8080,
    proxy: {
      "/api/": {
        target: devApiUrl,
        changeOrigin: true,
      },
    },
  },
})
