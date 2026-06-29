<script>
	export let tr;
	import { _ } from 'svelte-i18n';
	import TranscriptionRow from './TranscriptionRow.svelte';

	$: displayName = tr.fileName != '' ? tr.fileName.split('_WHSHPR_')[1] : tr.id;
</script>

<TranscriptionRow kind="translating">
	<p class="text-foreground text-[0.88rem] font-medium truncate" title={displayName}>{displayName}</p>
	<p class="text-muted-foreground text-[0.75rem] font-mono mt-0.5">
		{$_('transcription.status.waitingTranslation')}
	</p>
	<p class="text-muted-foreground text-[0.7rem] font-mono mt-0.5 flex items-center gap-2 flex-wrap">
		{#if tr.modelSize}
			<span>{$_('transcription.meta.model', { values: { model: tr.modelSize } })}</span>
			<span class="text-border">·</span>
		{/if}
		<span>{$_('transcription.meta.language', { values: { language: tr.result?.language ?? tr.language ?? 'auto' } })}</span>
	</p>
</TranscriptionRow>
