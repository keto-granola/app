import react from '@vitejs/plugin-react'
import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react()],
  build: {
    emptyOutDir: true,
    manifest: true,
    rollupOptions: {
      input: {
        'add-to-cart': resolve(__dirname, 'src/public/islands/entries/add-to-cart.ts'),
        'admin-dashboard': resolve(__dirname, 'src/admin/islands/entries/admin.tsx'),
      },
      output: {
        dir: '../backend/internal/webassets/dist',
      },
    },
  },
})
