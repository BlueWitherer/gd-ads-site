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
      '/api': 'http://localhost:3000',
      '/cdn': 'http://localhost:3000',
      '/session': 'http://localhost:3000',
      '/login': 'http://localhost:3000',
      '/logout': 'http://localhost:3000',
      '/callback': 'http://localhost:3000',
      '/ads': 'http://localhost:3000',
      '/proxy': 'http://localhost:3000',
      '/account': 'http://localhost:3000',
    },
  },
  "build": {
    "chunkSizeWarningLimit": 1024,
  },
});