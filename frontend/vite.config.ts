import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    host: 'localhost',
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8095',
      '/ws': { target: 'ws://localhost:8095', ws: true },
      '/health': 'http://localhost:8095',
    },
  },
});
