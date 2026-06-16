import { waitLocale } from 'svelte-i18n';
import '$lib/i18n';

/** @type {import('./$types').LayoutLoad} */
export async function load() {
	// Ensure the active locale's messages are loaded before the first render.
	await waitLocale();
}
