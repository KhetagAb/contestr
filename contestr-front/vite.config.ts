import path from "node:path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const apiTarget = process.env.VITE_API_PROXY_TARGET ?? "http://127.0.0.1:8080";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  build: {
    sourcemap: "inline",
  },
  server: {
    host: true,
    port: 5173,
    strictPort: true,
    proxy: {
      "/api": { target: apiTarget, changeOrigin: true },
    },
  },
});
