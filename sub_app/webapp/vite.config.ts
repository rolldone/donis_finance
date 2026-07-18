import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  base: '/',
  server: {
    host: '0.0.0.0',
    port: 8202,
    proxy: {
      '/api': {
        target: 'http://100.104.55.66:8200',
        changeOrigin: true,
      },
    },
  },
})
