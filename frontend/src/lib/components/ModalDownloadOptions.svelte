<script>
	import { downloadSRT, downloadTXT, downloadJSON, downloadVTT, CLIENT_API_HOST } from '$lib/utils';
	import toast from 'svelte-french-toast';
	import { _ } from 'svelte-i18n';
	import Modal from './ui/Modal.svelte';
	import Checkbox from './ui/Checkbox.svelte';
	import Spinner from './ui/Spinner.svelte';
	import { X, Download, Copy, Film, Check } from 'lucide-svelte';

	export let open = false;
	export let tr;

	const FORMATS = ['srt', 'vtt', 'json', 'txt'];

	let subtitleFormat = 'srt';
	let language = 'original';
	let includeWords = false;
	let copied = false;

	function downloadSubtitle() {
		let segments = [];
		let text = '';
		let title = 'subtitles';

		if (language == 'original') {
			segments = tr.result.segments;
			text = tr.result.text;
			title = tr.fileName.split('_WHSHPR_')[1];
		} else {
			for (const translation of tr.translations) {
				if (translation.targetLanguage == language) {
					segments = translation.result.segments;
					text = translation.result.text;
					title = tr.fileName.split('_WHSHPR_')[1];
					break;
				}
			}
		}

		if (segments.length == 0 || text == '') {
			toast.error($_('modals.download.toasts.noData'));
			return;
		}

		if (subtitleFormat == 'srt') {
			downloadSRT(segments, title);
		} else if (subtitleFormat == 'vtt') {
			downloadVTT(segments, title);
		} else if (subtitleFormat == 'json') {
			downloadJSON(tr.result, title, includeWords);
		} else if (subtitleFormat == 'txt') {
			downloadTXT(text, title);
		}
	}

	function downloadMedia() {
		let link = document.createElement('a');
		link.href = `${CLIENT_API_HOST}/api/video/${tr.fileName}`;
		link.target = '_blank';
		link.download = tr.fileName.split('_WHSHPR_')[1] + '.mp4';
		link.click();
	}

	async function copyText() {
		let text = '';
		if (language == 'original') {
			text = tr.result.text;
		} else {
			for (const translation of tr.translations) {
				if (translation.targetLanguage == language) {
					text = translation.result.text;
					break;
				}
			}
		}
		try {
			await navigator.clipboard.writeText(text);
			toast.success($_('modals.download.toasts.copied'));
			copied = true;
			setTimeout(() => (copied = false), 1500);
		} catch (err) {
			console.error('Error in copying text: ', err);
		}
	}
</script>

<Modal bind:open size="md" let:close>
	<div class="p-6">
		<div class="flex items-start justify-between mb-6">
			<div>
				<h2 class="text-foreground text-[0.95rem] font-semibold">{$_('modals.download.title')}</h2>
				{#if tr}
					<p class="text-muted-foreground text-[0.72rem] font-mono mt-0.5 truncate max-w-[16rem]">
						{tr.fileName.split('_WHSHPR_')[1]}
					</p>
				{/if}
			</div>
			<button
				on:click={close}
				class="w-7 h-7 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all shrink-0"
			>
				<X size={14} />
			</button>
		</div>

		{#if tr}
			<div class="space-y-4 mb-7">
				<!-- Format -->
				<div>
					<p class="text-muted-foreground mb-2 text-[0.72rem] uppercase tracking-wider">{$_('modals.download.fileFormat')}</p>
					<div class="flex items-center gap-1.5 flex-wrap">
						{#each FORMATS as f}
							<button
								on:click={() => (subtitleFormat = f)}
								class="px-3 py-1.5 rounded-lg border transition-all text-[0.78rem] font-medium uppercase {subtitleFormat ===
								f
									? 'bg-primary border-primary/50 text-primary-foreground shadow-lg shadow-primary/20'
									: 'border-border bg-muted text-muted-foreground hover:text-foreground hover:border-border/80'}"
							>
								{f}
							</button>
						{/each}
					</div>
				</div>

				<!-- Language -->
				<div>
					<p class="text-muted-foreground mb-2 text-[0.72rem] uppercase tracking-wider">{$_('modals.download.textLanguage')}</p>
					<div class="flex items-center gap-1.5 flex-wrap">
						<button
							on:click={() => (language = 'original')}
							class="px-3 py-1.5 rounded-lg border transition-all text-[0.78rem] font-medium uppercase {language ===
							'original'
								? 'bg-primary border-primary/50 text-primary-foreground shadow-lg shadow-primary/20'
								: 'border-border bg-muted text-muted-foreground hover:text-foreground hover:border-border/80'}"
						>
							{tr.result.language}
						</button>
						{#each tr.translations as translation}
							<button
								on:click={() => (language = translation.targetLanguage)}
								class="px-3 py-1.5 rounded-lg border transition-all text-[0.78rem] font-medium uppercase {language ===
								translation.targetLanguage
									? 'bg-primary border-primary/50 text-primary-foreground shadow-lg shadow-primary/20'
									: 'border-border bg-muted text-muted-foreground hover:text-foreground hover:border-border/80'}"
							>
								{translation.targetLanguage}
							</button>
						{/each}
					</div>
				</div>

				{#if subtitleFormat === 'json'}
					<label class="flex items-center gap-2 cursor-pointer">
						<Checkbox bind:checked={includeWords} />
						<span class="text-muted-foreground text-[0.78rem]">{$_('modals.download.includeWords')}</span>
					</label>
				{/if}
			</div>

			<!-- Actions -->
			<div class="grid grid-cols-3 gap-2">
				<button
					on:click={downloadSubtitle}
					title={$_('modals.download.tooltipFile', { values: { format: subtitleFormat, language } })}
					class="flex flex-col items-center gap-2 py-3.5 rounded-xl border border-emerald-500/20 bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 hover:border-emerald-500/40 transition-all hover:scale-[1.02] active:scale-[0.98]"
				>
					<Download size={17} />
					<span class="text-[0.7rem] font-semibold tracking-wider">{$_('modals.download.file')}</span>
				</button>
				<button
					on:click={copyText}
					title={$_('modals.download.tooltipCopy', { values: { language } })}
					class="flex flex-col items-center gap-2 py-3.5 rounded-xl border transition-all hover:scale-[1.02] active:scale-[0.98] {copied
						? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-400'
						: 'border-blue-500/20 bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 hover:border-blue-500/40'}"
				>
					{#if copied}<Check size={17} />{:else}<Copy size={17} />{/if}
					<span class="text-[0.7rem] font-semibold tracking-wider">{$_('modals.download.copy')}</span>
				</button>
				<button
					on:click={downloadMedia}
					title={$_('modals.download.tooltipMedia')}
					class="flex flex-col items-center gap-2 py-3.5 rounded-xl border border-violet-500/20 bg-violet-500/10 text-violet-400 hover:bg-violet-500/20 hover:border-violet-500/40 transition-all hover:scale-[1.02] active:scale-[0.98]"
				>
					<Film size={17} />
					<span class="text-[0.7rem] font-semibold tracking-wider">{$_('modals.download.media')}</span>
				</button>
			</div>
		{:else}
			<div class="flex items-center justify-center py-12 text-muted-foreground">
				<Spinner size={28} />
			</div>
		{/if}
	</div>
</Modal>
