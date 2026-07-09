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
        'add-to-cart': resolve(__dirname, 'src/islands/entries/add-to-cart.ts'),
      },
      output: {
        dir: '../backend/internal/webassets/dist',
      },
    },
  },
})
