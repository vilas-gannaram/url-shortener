import { defineConfig, loadEnv } from 'vite';
import { cwd } from 'node:process';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, cwd());

	return {
		plugins: [tailwindcss(), react()],
		server: {
			proxy: {
				'/api': {
					target: env.VITE_API_URL,
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/api/, ''),
				},
			},
			port: 3000,
		},
	};
});
