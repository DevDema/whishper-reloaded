<script>
	import toast from 'svelte-french-toast';
	import { renameTranscription } from '$lib/utils.js';
	import { _ } from 'svelte-i18n';
	import Modal from './ui/Modal.svelte';
	import Button from './ui/Button.svelte';
	import { Pencil, X } from 'lucide-svelte';

	export let open = false;
	export let tr; // transcription object passed from parent

	let newFileName = '';
	let fileExtension = '';
	let isSubmitting = false;
	let lastId;

	// Parse filename + extension whenever a different transcription is provided.
	$: if (tr && tr.id !== lastId) {
		const parts = tr.fileName.split('WHSHPR_');
		const nameWithExt = parts.length > 1 ? parts[1] : tr.fileName;
		const extIndex = nameWithExt.lastIndexOf('.');
		if (extIndex !== -1) {
			fileExtension = nameWithExt.slice(extIndex);
			newFileName = nameWithExt.slice(0, extIndex);
		} else {
			newFileName = nameWithExt;
			fileExtension = '';
		}
		lastId = tr.id;
	}

	$: originalName = tr
		? (() => {
				const parts = tr.fileName.split('WHSHPR_');
				const nameWithExt = parts.length > 1 ? parts[1] : tr.fileName;
				const extIndex = nameWithExt.lastIndexOf('.');
				return extIndex !== -1 ? nameWithExt.slice(0, extIndex) : nameWithExt;
		  })()
		: '';

	$: canSubmit = tr && newFileName && newFileName !== originalName;

	async function handleSubmit() {
		if (!canSubmit) return;
		isSubmitting = true;
		try {
			const fullNewFileName = newFileName + fileExtension;
			await renameTranscription(tr.id, fullNewFileName);
			toast.success($_('modals.rename.toasts.success'));
			open = false;
		} catch (error) {
			console.error('Error renaming file:', error);
			toast.error($_('modals.rename.toasts.error'));
		} finally {
			isSubmitting = false;
		}
	}
</script>

<Modal bind:open size="md" let:close>
	<div class="p-6">
		<div class="flex items-center justify-between mb-5">
			<div class="flex items-center gap-2.5">
				<div class="w-8 h-8 rounded-lg bg-primary/15 border border-primary/25 flex items-center justify-center">
					<Pencil size={14} class="text-primary" />
				</div>
				<h2 class="text-foreground text-[0.95rem] font-semibold">{$_('modals.rename.title')}</h2>
			</div>
			<button
				on:click={close}
				class="w-7 h-7 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
			>
				<X size={14} />
			</button>
		</div>

		{#if tr}
			<label for="renameInput" class="text-muted-foreground mb-2 block text-[0.72rem] uppercase tracking-wider">
				{$_('modals.rename.newFileName')}
			</label>

			<div
				class="flex items-stretch rounded-xl border border-border bg-muted overflow-hidden mb-6 transition-all focus-within:ring-2 focus-within:ring-primary/35 focus-within:border-primary/35"
			>
				<!-- svelte-ignore a11y-autofocus -->
				<input
					id="renameInput"
					type="text"
					bind:value={newFileName}
					on:keydown={(e) => e.key === 'Enter' && handleSubmit()}
					placeholder={$_('modals.rename.placeholder')}
					autofocus
					class="flex-1 px-3.5 py-2.5 bg-transparent text-foreground placeholder:text-muted-foreground focus:outline-none text-[0.875rem]"
				/>
				{#if fileExtension}
					<div class="flex items-center px-3 border-l border-border text-muted-foreground bg-secondary/40 shrink-0 text-[0.78rem] font-mono">
						{fileExtension}
					</div>
				{/if}
			</div>

			<div class="flex items-center justify-end gap-2">
				<Button variant="secondary" size="sm" on:click={close}>{$_('common.cancel')}</Button>
				<Button size="sm" on:click={handleSubmit} disabled={!canSubmit || isSubmitting}>
					{$_('modals.rename.rename')}
				</Button>
			</div>
		{/if}
	</div>
</Modal>
