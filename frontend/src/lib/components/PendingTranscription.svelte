<script>
	export let tr;
	import { deleteTranscription } from '$lib/utils.js';
	import { _ } from 'svelte-i18n';
	import { Trash2 } from 'lucide-svelte';
	import TranscriptionRow from './TranscriptionRow.svelte';
	import IconButton from './ui/IconButton.svelte';

	$: displayName = tr.fileName != '' ? tr.fileName.split('_WHSHPR_')[1] : tr.id;
</script>

<TranscriptionRow kind="queued">
	<p class="text-foreground text-[0.88rem] font-medium truncate" title={displayName}>{displayName}</p>
	<p class="text-muted-foreground text-[0.75rem] font-mono mt-0.5">
		{$_('transcription.status.pending', { values: { language: tr.language || 'auto' } })}
	</p>

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
