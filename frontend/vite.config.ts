import {defineConfig} from 'vite'
import dotenv from 'dotenv'
import path from 'node:path';
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite';

dotenv.config({ path: "../.env"})

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
      react(), 
      tailwindcss()
  ],
  resolve: {
    alias: {
        '@': path.resolve(import.meta.dirname, './src'),
        '@wailsjs': path.resolve(import.meta.dirname, './wailsjs'),
        '@components': path.resolve(import.meta.dirname, './src/components'),
        '@pages': path.resolve(import.meta.dirname, './src/pages'),
        '@util': path.resolve(import.meta.dirname, './src/util'),
        '@api': path.resolve(import.meta.dirname, './src/api'),
        '@contexts': path.resolve(import.meta.dirname, './src/contexts'),
        '@assets': path.resolve(import.meta.dirname, './src/assets'),
    }
  },
  server: {
      port: 3000,
  },
})
