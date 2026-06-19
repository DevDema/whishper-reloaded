<script>
	export let tr;
	import { deleteTranscription } from '$lib/utils.js';
	import { _ } from 'svelte-i18n';
	import { Trash2 } from 'lucide-svelte';
	import TranscriptionRow from './TranscriptionRow.svelte';
	import IconButton from './ui/IconButton.svelte';

	// Progress comes from the backend as a fraction (0.0 - 1.0).
	$: progress = tr.progress ?? 0;
	$: percent = Math.round(progress * 100);
	$: hasProgress = progress > 0;
	$: downloadingModel = tr.downloadingModel ?? false;
	$: displayName = tr.fileName != '' ? tr.fileName.split('_WHSHPR_')[1] : tr.id;
</script>

<TranscriptionRow kind="processing">
	<p class="text-foreground text-[0.88rem] font-medium truncate" title={displayName}>{displayName}</p>
	{#if downloadingModel}
		<p class="text-muted-foreground text-[0.75rem] font-mono mt-0.5">
			{$_('transcription.status.downloadingModel', { values: { model: tr.modelSize } })}
		</p>
	{:else if hasProgress}
		<p class="text-muted-foreground text-[0.75rem] font-mono mt-0.5">
			{$_('transcription.status.transcribing', { values: { percent } })}
		</p>
	{:else}
		<p class="text-muted-foreground text-[0.75rem] font-mono mt-0.5">
			{$_('transcription.status.waitingTranscription')}
		</p>
	{/if}
	{#if hasProgress}
		<div class="mt-2 h-1.5 w-full rounded-full bg-muted overflow-hidden">
			<div class="h-full rounded-full bg-blue-400 transition-all duration-300" style="width:{percent}%"></div>
		</div>
	{/if}

	<svelte:fragment slot="actions">
		<IconButton
			variant="danger"
			title={$_('transcription.actions.delete')}
			on:click={() => deleteTranscription(tr.id)}
		>
			<Trash2 size={14} />
		</IconButton>
	</svelte:fragment>
</TranscriptionRow>
