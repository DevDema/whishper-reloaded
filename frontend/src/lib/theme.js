import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function getInitial() {
	if (!browser) return 'dark';
	const stored = localStorage.getItem('theme');
	if (stored === 'dark' || stored === 'light') return stored;
	return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function applyTheme(value) {
	if (!browser) return;
	document.documentElement.classList.toggle('dark', value === 'dark');
	try {
		localStorage.setItem('theme', value);
	} catch (e) {
		/* ignore */
	}
}

function createTheme() {
	const initial = getInitial();
	const { subscribe, set, update } = writable(initial);

	// Keep the DOM/localStorage in sync with the initial value on the client.
	applyTheme(initial);

	return {
		subscribe,
		set: (value) => {
			applyTheme(value);
			set(value);
		},
		toggle: () =>
			update((value) => {
				const next = value === 'dark' ? 'light' : 'dark';
				applyTheme(next);
				return next;
			})
	};
}

export const theme = createTheme();
