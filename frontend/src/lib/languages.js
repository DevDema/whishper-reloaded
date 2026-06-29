/**
 * Maps ISO 639-1 language codes to human-readable labels for the UI.
 * The backend always receives the raw `code`; only the label is shown to the user.
 */

const FALLBACK_LABELS = {
	auto: 'Automatic',
	ar: 'Arabic',
	be: 'Belarusian',
	bg: 'Bulgarian',
	bn: 'Bengali',
	ca: 'Catalan',
	cs: 'Czech',
	cy: 'Welsh',
	da: 'Danish',
	de: 'German',
	el: 'Greek',
	en: 'English',
	es: 'Spanish',
	fr: 'French',
	it: 'Italian',
	ja: 'Japanese',
	nl: 'Dutch',
	pl: 'Polish',
	pt: 'Portuguese',
	ru: 'Russian',
	sk: 'Slovak',
	sl: 'Slovenian',
	sv: 'Swedish',
	tk: 'Turkmen',
	tr: 'Turkish',
	zh: 'Chinese'
};

/**
 * @param {string} code
 * @param {string} [locale]
 */
export function getLanguageLabel(code, locale = 'en') {
	if (!code) return '';
	const normalized = code.toLowerCase();

	try {
		const display = new Intl.DisplayNames([locale], { type: 'language' });
		const name = display.of(normalized);
		if (name) {
			return name.charAt(0).toUpperCase() + name.slice(1);
		}
	} catch {
		// Intl not available or locale unsupported
	}

	return FALLBACK_LABELS[normalized] ?? normalized.toUpperCase();
}

/**
 * @param {Array<{ code: string, targets: string[] }>} availableLanguages
 * @param {string} sourceLanguage
 * @param {string} [locale]
 * @returns {Array<{ code: string, label: string }>}
 */
export function getTargetLanguageOptions(availableLanguages, sourceLanguage, locale = 'en') {
	if (!availableLanguages?.length || !sourceLanguage) return [];

	const entry = availableLanguages.find((l) => l.code === sourceLanguage);
	if (!entry?.targets?.length) return [];

	return entry.targets
		.filter((code) => code !== sourceLanguage)
		.map((code) => ({
			code,
			label: getLanguageLabel(code, locale)
		}))
		.sort((a, b) => a.label.localeCompare(b.label, locale));
}
