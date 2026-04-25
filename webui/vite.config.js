import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [svelte(), tailwindcss()],
	base: '/oo/',
	resolve: {
		alias: {
			$lib: path.resolve('./src/lib')
		}
	},
	build: {
		outDir: 'dist',
		emptyOutDir: true
	}
});
