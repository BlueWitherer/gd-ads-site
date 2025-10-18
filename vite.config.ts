import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  "plugins": [
    react({
      "babel": {
        "plugins": [['babel-plugin-react-compiler']],
      },
    }),
    tailwindcss(),
  ],
  "server": {
    "proxy": {
      '/api': 'http://localhost:8081',
      '/cdn': 'http://localhost:8081',
      '/session': 'http://localhost:8081',
      '/login': 'http://localhost:8081',
      '/logout': 'http://localhost:8081',
      '/callback': 'http://localhost:8081',
      '/ads': 'http://localhost:8081',
      '/proxy': 'http://localhost:8081',
      '/account': 'http://localhost:8081',
    },
  },
  "build": {
    "chunkSizeWarningLimit": 1024,
  },
});