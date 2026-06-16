import { browser } from '$app/environment';
import { init, register } from 'svelte-i18n';

const defaultLocale = 'en';
const supportedLocales = ['en', 'es', 'it', 'fr', 'pt'];

register('en', () => import('./locales/en.json'));
register('es', () => import('./locales/es.json'));
register('it', () => import('./locales/it.json'));
register('fr', () => import('./locales/fr.json'));
register('pt', () => import('./locales/pt.json'));

// Picks the locale from the browser language, falling back to the default.
function getInitialLocale() {
	if (!browser) return defaultLocale;

	const navLocale = window.navigator.language?.split('-')[0];
	if (navLocale && supportedLocales.includes(navLocale)) return navLocale;

	return defaultLocale;
}

init({
	fallbackLocale: defaultLocale,
	initialLocale: getInitialLocale()
});
