<script>
	import { uploadJSON } from '$lib/utils';
	import toast from 'svelte-french-toast';
	import { _ } from 'svelte-i18n';
	import Modal from './ui/Modal.svelte';
	import Button from './ui/Button.svelte';
	import Spinner from './ui/Spinner.svelte';
	import { Upload, X, FileJson } from 'lucide-svelte';

	export let open = false;
	export let tr;

	let fileInput;
	let selectedFile = null;
	let uploading = false;

	function handleFileSelect(event) {
		const file = event.target.files[0];
		if (file) {
			if (file.type !== 'application/json' && !file.name.endsWith('.json')) {
				toast.error($_('modals.upload.toasts.invalidFile'));
				return;
			}
			selectedFile = file;
		}
	}

	async function handleUpload() {
		if (!selectedFile || !tr) {
			toast.error($_('modals.upload.toasts.noFile'));
			return;
		}

		uploading = true;

		try {
			const fileContent = await selectedFile.text();
			const jsonData = JSON.parse(fileContent);
			const result = await uploadJSON(tr.id, jsonData);

			if (result.noChanges) {
				toast($_('modals.upload.toasts.noChanges'), { icon: '👐' });
			} else {
				toast.success($_('modals.upload.toasts.success'));
				clearFile();
				open = false;
			}
		} catch (error) {
			console.error('Upload error:', error);
			if (error instanceof SyntaxError) {
				toast.error($_('modals.upload.toasts.invalidFormat'));
			} else {
				toast.error(error.message || $_('modals.upload.toasts.failed'));
			}
		} finally {
			uploading = false;
		}
	}

	function clearFile() {
		selectedFile = null;
		if (fileInput) fileInput.value = '';
	}
</script>

<Modal bind:open size="md" let:close>
	<div class="p-6">
		<div class="flex items-center justify-between mb-5">
			<div class="flex items-center gap-2.5">
				<div class="w-8 h-8 rounded-lg bg-primary/15 border border-primary/25 flex items-center justify-center">
					<Upload size={14} class="text-primary" />
				</div>
				<h2 class="text-foreground text-[0.95rem] font-semibold">{$_('modals.upload.title')}</h2>
			</div>
			<button
				on:click={close}
				class="w-7 h-7 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
			>
				<X size={14} />
			</button>
		</div>

		{#if tr}
			<p class="text-muted-foreground text-[0.8rem] mb-4">
				{$_('modals.upload.description')}<br />
				<span class="font-mono text-foreground">{tr.fileName.split('_WHSHPR_')[1]}</span>
			</p>

			<label for="jsonFile" class="text-muted-foreground mb-2 block text-[0.72rem] uppercase tracking-wider">
				{$_('modals.upload.selectFile')}
			</label>

			<button
				on:click={() => fileInput?.click()}
				class="w-full flex items-center gap-3 rounded-xl border-2 border-dashed border-border hover:border-primary/40 hover:bg-muted/50 px-4 py-5 transition-all text-left {uploading
					? 'opacity-60 pointer-events-none'
					: ''}"
			>
				<div class="w-10 h-10 rounded-lg bg-muted border border-border flex items-center justify-center shrink-0">
					<FileJson size={18} class="text-muted-foreground" />
				</div>
				<div class="min-w-0">
					{#if selectedFile}
						<p class="text-foreground text-[0.85rem] font-mono truncate">{selectedFile.name}</p>
					{:else}
						<p class="text-foreground text-[0.85rem] font-medium">{$_('modals.upload.selectFile')}</p>
						<p class="text-muted-foreground text-[0.72rem]">.json</p>
					{/if}
				</div>
			</button>
			<input
				bind:this={fileInput}
				type="file"
				id="jsonFile"
				accept=".json,application/json"
				on:change={handleFileSelect}
				class="hidden"
				disabled={uploading}
			/>

			<div class="flex items-center justify-end gap-2 mt-6">
				<Button variant="secondary" size="sm" on:click={close}>{$_('common.cancel')}</Button>
				<Button size="sm" on:click={handleUpload} disabled={!selectedFile || uploading}>
					{#if uploading}
						<Spinner size={14} />
						{$_('common.uploading')}
					{:else}
						<Upload size={14} />
						{$_('modals.upload.upload')}
					{/if}
				</Button>
			</div>
		{:else}
			<div class="flex items-center justify-center py-12 text-muted-foreground">
				<Spinner size={28} />
			</div>
		{/if}
	</div>
</Modal>
