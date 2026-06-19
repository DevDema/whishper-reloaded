<script>
	import { deleteTranscription } from '$lib/utils.js';
	import { _ } from 'svelte-i18n';
	import { Trash2 } from 'lucide-svelte';
	import TranscriptionRow from './TranscriptionRow.svelte';
	import Button from './ui/Button.svelte';
	export let tr;
</script>

<TranscriptionRow kind="error">
	<p class="text-foreground text-[0.88rem] font-medium truncate">id-{tr.id}</p>
	<p class="text-muted-foreground text-[0.75rem] font-mono mt-0.5">
		{$_('transcription.status.failed')}
	</p>

	<svelte:fragment slot="actions">
		<Button variant="danger" size="sm" on:click={() => deleteTranscription(tr.id)}>
			<Trash2 size={14} />
			{$_('transcription.actions.delete')}
		</Button>
	</svelte:fragment>
</TranscriptionRow>
