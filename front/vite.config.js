import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    rollupOptions: {
      input: 'main.js',
      output: {
        dir: '../backend/router/static',
        entryFileNames: 'main.js',
        format: 'es',
      },
    },
  },
});