<script>
	import { deleteTranscription } from '$lib/utils.js';
	export let tr;

	function confirmDelete() {
		if (tr && tr.id) {
            modalConfirmDelete.close();
			deleteTranscription(tr.id)();
		}
	}
</script>

<dialog id="modalConfirmDelete" class="modal">
	<div class="modal-box">
		<h3 class="font-bold text-lg">Confirm Deletion</h3>
		{#if tr}
			<p class="py-4">
				Are you sure you want to delete the transcription for
				<span class="font-bold">{tr.fileName.split('_WHSHPR_')[1]}</span>? This action cannot be
				undone.
			</p>
		{/if}
		<div class="modal-action">
			<button class="btn" onclick="modalConfirmDelete.close()">Cancel</button>
			<button class="btn btn-error" on:click={confirmDelete}>Delete</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop"><button>close</button></form>
</dialog>