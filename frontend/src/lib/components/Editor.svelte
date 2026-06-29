<script>
	import { onMount } from 'svelte';
	import toast from 'svelte-french-toast';
	import { editorSettings, currentTranscription, editorHistory, audioMode } from '$lib/stores';
	import EditorSettings from './EditorSettings.svelte';
	import EditorSegment from './EditorSegment.svelte';
	import GoToSegment from './GoToSegment.svelte';
	import IconButton from './ui/IconButton.svelte';
	import Select from './ui/Select.svelte';
	import Spinner from './ui/Spinner.svelte';
	import ThemeToggle from './ui/ThemeToggle.svelte';
	import { ArrowLeft, Save, Volume2, MonitorPlay } from 'lucide-svelte';
	import { _ } from 'svelte-i18n';

	export let language;
	export let segmentsToShow;
	export let setSegmentsToShow;

	const audioExtensions = [
		'mp3', 'm4a', 'wav', 'aac', 'ogg', 'flac', 'aiff', 'wma', 'alac', 'opus', 'amr', 'pcm', 'mka', 'dsd', 'wv', 'ape', 'm3u', 'm3u8', 'pls', 'aif', 'au', 'ra', 'ram', 'ac3', 'dts', 'mid', 'midi', 'mpa', 'mpc', 'oga', 'spx', 'tta', 'voc', 'vox', 'caf', 'snd', 'kar', 'mod', 's3m', 'xm', 'it', 'mtm', 'umx', 'amz', 'mogg', '8svx', 'cda', 'gsm', 'ivs', 'mlp', 'msv', 'nmf', 'shn', 'tak', 'vqf', 'xa'
	];

	function isAudioOnlyFile(fileName) {
		if (!fileName) return false;
		const ext = fileName.split('.').pop().toLowerCase();
		return audioExtensions.includes(ext);
	}

	let isAudioOnly = false;

	$: {
		if ($currentTranscription && $currentTranscription.fileName) {
			isAudioOnly = isAudioOnlyFile($currentTranscription.fileName);
			if (isAudioOnly && !$audioMode) {
				$audioMode = true;
			}
		}
	}

	function loadMore() {
		setSegmentsToShow(segmentsToShow + 10);
	}
	let loadMoreButton;

	async function textFromSegments() {
		let text = '';
		if ($language == 'original') {
			text = $currentTranscription.result.segments
				.map((segment) => segment.text)
				.join(' ')
				.replace(/(\r\n|\n|\r)/gm, ' ');
		} else {
			text = $currentTranscription.translations
				.filter((translation) => translation.targetLanguage == $language)[0]
				.result.segments.map((segment) => segment.text)
				.join(' ')
				.replace(/(\r\n|\n|\r)/gm, ' ');
		}
		return text;
	}

	async function saveChanges() {
		var url = `${CLIENT_API_HOST}/api/transcriptions`;
		if ($language == 'original') {
			$currentTranscription.result.text = await textFromSegments();
		} else {
			$currentTranscription.translations.forEach(async (translation) => {
				if (translation.targetLanguage == $language) translation.result.text = await textFromSegments();
			});
		}

		try {
			const response = await fetch(url, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify($currentTranscription)
			});

			if (!response.ok) {
				if (response.status === 304) {
					if (!$editorSettings.autoSave) {
						toast($_('editor.toasts.noChanges'), { icon: '👐' });
					}
					return;
				} else {
					toast.error($_('editor.toasts.saveError'));
					throw new Error(`HTTP error! status: ${response.status}`);
				}
			}

			if ($editorSettings.autoSave) {
				toast($_('editor.toasts.autosaving'), { icon: 'ℹ️' });
			} else {
				toast.success($_('editor.toasts.saved'));
			}
		} catch (error) {
			toast.error($_('editor.toasts.saveError'));
			console.error('Error:', error);
		}
	}

	import { CLIENT_API_HOST } from '$lib/utils';

	let handleKeyDown;
	let observer;
	onMount(() => {
		observer = new IntersectionObserver((entries) => {
			entries.forEach((entry) => {
				if (entry.isIntersecting) {
					loadMore();
				}
			});
		});
		observer.observe(loadMoreButton);

		editorHistory.set([JSON.parse(JSON.stringify($currentTranscription))]);
		let isUndoing = false;
		handleKeyDown = function (e) {
			if (e.ctrlKey && e.key === 'z' && !isUndoing) {
				isUndoing = true;
				let previousTranscription = null;

				editorHistory.update((history) => {
					if (history.length > 1) {
						history = history.slice(0, -1);
						previousTranscription = { ...history[history.length - 1] };
					}
					return history;
				});

				if (previousTranscription) {
					$currentTranscription = { ...previousTranscription };
				}
				isUndoing = false;
			}

			let isSaving = false;
			if (e.ctrlKey && e.key === 's') {
				e.preventDefault();
				if (!isSaving) {
					isSaving = true;
					saveChanges();
					isSaving = false;
				}
			}
		};

		if (!$editorSettings.autoSave) {
			toast($_('editor.toasts.autosaveDisabled'), { icon: '👋' });
		}

		document.addEventListener('keydown', handleKeyDown);
	});

	let autosaveInterval;
	let autoSaveAux = $editorSettings.autoSave;
	$: if ($editorSettings.autoSave) {
		toast.success($_('editor.toasts.autosaveEnabled'));
		autoSaveAux = true;
		autosaveInterval = setInterval(() => {
			saveChanges();
		}, $editorSettings.autosaveInterval);
	} else {
		if (autoSaveAux == true) {
			toast($_('editor.toasts.autosaveDisabled'), { icon: '👋' });
			autoSaveAux = false;
		}
		clearInterval(autosaveInterval);
	}

	$: displayName = $currentTranscription?.fileName?.split('_WHSHPR_')[1] ?? '';
</script>

{#if $currentTranscription.status != 2}
	<div class="flex items-center justify-center gap-3 py-20 text-muted-foreground">
		<Spinner size={24} />
		<p>
			{$_('editor.waitingTask', {
				values: {
					action:
						$currentTranscription.status == 3
							? $_('editor.actionTranslating')
							: $_('editor.actionTranscribing')
				}
			})}
		</p>
	</div>
{:else}
	<!-- Header -->
	<div class="sticky top-0 z-10 bg-card/60 backdrop-blur-sm border-b border-border px-4 py-3">
		<div class="flex items-center justify-between gap-4 max-w-5xl mx-auto">
			<div class="flex items-center gap-3 min-w-0">
				<a
					href="/"
					title={$_('editor.back')}
					aria-label={$_('editor.back')}
					class="w-8 h-8 shrink-0 rounded-lg flex items-center justify-center border border-border bg-muted text-muted-foreground hover:text-foreground hover:bg-secondary transition-all hover:scale-105 active:scale-95"
				>
					<ArrowLeft size={16} />
				</a>
				<h1 class="text-[1rem] font-semibold truncate min-w-0 text-foreground" title={displayName}>
					{displayName}
				</h1>
			</div>

			<div class="flex items-center gap-2 flex-shrink-0">
				{#if $currentTranscription.translations.length > 0}
					<Select bind:value={$language} class="uppercase !py-2 w-auto min-w-[7rem]">
						<option value="original">✅ {$currentTranscription.result.language}</option>
						{#each $currentTranscription.translations as translation}
							<option value={translation.targetLanguage}>🤖 {translation.targetLanguage}</option>
						{/each}
					</Select>
				{/if}

				<EditorSettings />
				<GoToSegment language={$language} {segmentsToShow} {setSegmentsToShow} />

				<div class="flex items-center gap-1.5">
					<ThemeToggle />
					<IconButton title={$_('editor.save')} on:click={saveChanges}>
						<Save size={16} />
					</IconButton>
					{#if !isAudioOnly}
						<IconButton
							title={$audioMode ? $_('editor.exitAudioMode') : $_('editor.audioMode')}
							on:click={() => ($audioMode = !$audioMode)}
						>
							{#if $audioMode}<MonitorPlay size={16} />{:else}<Volume2 size={16} />{/if}
						</IconButton>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<!-- Segments -->
	<div class="max-w-5xl mx-auto px-4 py-5 space-y-1.5">
		{#if $language == 'original'}
			{#each $currentTranscription.result.segments.slice(0, segmentsToShow) as segment, index (segment.id)}
				<EditorSegment {segment} {index} translationIndex={-1} />
			{/each}
		{:else}
			{#each $currentTranscription.translations as translation, translationIndex}
				{#if translation.targetLanguage == $language}
					{#each translation.result.segments.slice(0, segmentsToShow) as segment, index (segment.id)}
						<EditorSegment {segment} {index} {translationIndex} />
					{/each}
				{/if}
			{/each}
		{/if}

		<button
			bind:this={loadMoreButton}
			class="w-full text-center text-[0.78rem] text-muted-foreground py-4"
		>
			{#if $language == 'original'}
				{#if segmentsToShow >= $currentTranscription.result.segments.length}
					{$_('editor.noMoreSegments')}
				{:else}
					{$_('editor.loadingMore')}
				{/if}
			{:else if segmentsToShow >= $currentTranscription.translations.filter((translation) => translation.targetLanguage == $language)[0].result.segments.length}
				{$_('editor.noMoreSegments')}
			{:else}
				{$_('editor.loadingMore')}
			{/if}
		</button>
	</div>
{/if}
