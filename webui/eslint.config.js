import js from '@eslint/js';
import svelteParser from 'svelte-eslint-parser';
import globals from 'globals';
import { globalIgnores } from 'eslint/config';

export default [
	js.configs.recommended,
	globalIgnores(['build/**/*', 'dist/**/*']),
	{
		files: ['src/**/*.svelte'],
		languageOptions: {
			parser: svelteParser,
			globals: { ...globals.browser, browser: true, nodeBuiltin: true, node: true, es6: true }
		}
	},
	{
		files: ['src/*.js'],
		languageOptions: {
			globals: { ...globals.browser, browser: true, nodeBuiltin: true, node: true, es6: true }
		}
	}
];
