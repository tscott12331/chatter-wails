import { defineConfig } from "vite";
import dotenv from 'dotenv'
import react from "@vitejs/plugin-react";
import wails from "@wailsio/runtime/plugins/vite";
import tailwindcss from '@tailwindcss/vite';
import path from "node:path";

dotenv.config({ path: "../.env"})

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  plugins: [
      react(), 
      wails("./bindings"),
      tailwindcss(),
  ],
  resolve: {
    alias: {
        '@': path.resolve(import.meta.dirname, './src'),
        '@wailsjs': path.resolve(import.meta.dirname, './bindings'),
        '@components': path.resolve(import.meta.dirname, './src/components'),
        '@pages': path.resolve(import.meta.dirname, './src/pages'),
        '@util': path.resolve(import.meta.dirname, './src/util'),
        '@api': path.resolve(import.meta.dirname, './src/api'),
        '@contexts': path.resolve(import.meta.dirname, './src/contexts'),
        '@assets': path.resolve(import.meta.dirname, './src/assets'),
    }
  },
});
