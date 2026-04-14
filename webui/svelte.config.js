import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default {
	preprocess: vitePreprocess(),
	compilerOptions: {
		// 'undefined' (the default) allows Svelte 5 to auto-detect.
		// Setting it to 'true' forces runes everywhere and breaks old libraries.
		runes: undefined
	}
};
