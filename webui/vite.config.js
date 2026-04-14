import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

export default defineConfig({
	plugins: [svelte()],
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
