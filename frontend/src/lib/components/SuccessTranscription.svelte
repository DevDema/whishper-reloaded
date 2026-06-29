<!-- SuccessTranscription.svelte -->
<script>
	import { createEventDispatcher } from 'svelte';
	import { goto } from '$app/navigation';
	import { deleteTranscription, getFullTranscription } from '$lib/utils.js';
	import toast from 'svelte-french-toast';
	import { _ } from 'svelte-i18n';
	import { Pencil, Download, Upload, Trash2 } from 'lucide-svelte';
	import TranscriptionRow from './TranscriptionRow.svelte';
	import IconButton from './ui/IconButton.svelte';
	import Badge from './ui/Badge.svelte';
	import Spinner from './ui/Spinner.svelte';
	import { STATUS_CONFIG } from './ui/statusConfig';

	export let tr;
	// Kept for compatibility with the list page; translation is handled inside the editor.
	export let languagesAvailable = false;

	const dispatch = createEventDispatcher();

	// Tracks which action is currently loading the full transcription, so we can show a
	// spinner and disable the buttons. The list view only holds a lightweight transcription
	// (no result/segments), so the full record must be fetched before opening any modal.
	let loadingAction = null;

	const runAction = async (action) => {
		if (loadingAction) return;
		loadingAction = action;
		try {
			const full = await getFullTranscription(tr.id);
			dispatch(action, full); // parent opens the modal once the full data is ready
		} catch (error) {
			console.error(error);
			toast.error($_('transcription.toasts.loadFailed'));
		} finally {
			loadingAction = null;
		}
	};

	// Action handlers stop propagation so clicking a button never triggers the card navigation.
	const onAction = (action) => (e) => {
		e.stopPropagation();
		runAction(action);
	};

	const onDelete = (e) => {
		e.stopPropagation();
		deleteTranscription(tr.id);
	};

	const openEditor = () => goto(`/editor/${tr.id}`);

	$: displayName = tr.fileName.split('_WHSHPR_')[1];
	$: duration = $_('transcription.meta.duration', {
		values: {
			duration: new Date(Math.round(tr.duration ?? tr.result?.duration ?? 0) * 1000)
				.toISOString()
				.substr(11, 8)
		}
	});
	$: translations = $_('transcription.meta.translations', {
		values: { count: tr.translations.length }
	});
	$: words = $_('transcription.meta.words', {
		values: { count: tr.words_count ?? tr.result?.text.split(' ').length ?? 0 }
	});
</script>

<TranscriptionRow kind="done" interactive on:activate={openEditor}>
	<div class="flex items-center gap-2 mb-1 flex-wrap">
		<span class="text-foreground truncate text-[0.88rem] font-medium" title={displayName}>
			{displayName}
		</span>
		<Badge mono class={STATUS_CONFIG.done.badge}>{$_('transcription.statusLabel.done')}</Badge>
		{#if tr.vad_filter}
			<span
				title={$_('transcription.meta.vadTooltip', {
					values: {
						threshold: tr.vad_threshold ?? 0.5,
						minSpeech: tr.vad_min_speech_duration_ms ?? 250,
						minSilence: tr.vad_min_silence_duration_ms ?? 500
					}
				})}
			>
				<Badge class="border-primary/20 bg-primary/10 text-primary">VAD</Badge>
			</span>
		{/if}
	</div>
	<div class="flex items-center gap-3 text-muted-foreground text-[0.75rem] font-mono flex-wrap">
		<span>{duration}</span>
		<span class="text-border">·</span>
		<span>{words}</span>
		<span class="text-border">·</span>
		<span>{translations}</span>
		{#if tr.modelSize}
			<span class="text-border">·</span>
			<span>{$_('transcription.meta.model', { values: { model: tr.modelSize } })}</span>
		{/if}
		<span class="text-border">·</span>
		<span>{$_('transcription.meta.language', { values: { language: tr.result?.language ?? tr.language ?? 'auto' } })}</span>
	</div>

	<svelte:fragment slot="actions">
		<div class="flex items-center gap-1.5 opacity-0 group-hover:opacity-100 focus-within:opacity-100 transition-opacity">
			<IconButton title={$_('transcription.actions.rename')} on:click={onAction('rename')} disabled={!!loadingAction}>
				{#if loadingAction === 'rename'}<Spinner size={13} />{:else}<Pencil size={13} />{/if}
			</IconButton>
			<IconButton title={$_('transcription.actions.download')} on:click={onAction('download')} disabled={!!loadingAction}>
				{#if loadingAction === 'download'}<Spinner size={13} />{:else}<Download size={13} />{/if}
			</IconButton>
			<IconButton title={$_('transcription.actions.uploadJson')} on:click={onAction('upload')} disabled={!!loadingAction}>
				{#if loadingAction === 'upload'}<Spinner size={13} />{:else}<Upload size={13} />{/if}
			</IconButton>
			<IconButton variant="danger" title={$_('transcription.actions.delete')} on:click={onDelete}>
				<Trash2 size={13} />
			</IconButton>
		</div>
	</svelte:fragment>
</TranscriptionRow>
