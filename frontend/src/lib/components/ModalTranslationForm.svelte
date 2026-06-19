<script>
	import toast from 'svelte-french-toast';
	import { CLIENT_API_HOST } from '$lib/utils';
	import { _ } from 'svelte-i18n';
	import Modal from './ui/Modal.svelte';
	import Button from './ui/Button.svelte';
	import Select from './ui/Select.svelte';
	import { Languages, X } from 'lucide-svelte';

	export let open = false;
	export let tr;
	export let availableLanguages;
	let targetLanguage = null;

	const handleTranslate = (id) => {
		if (targetLanguage) {
			const url = `${CLIENT_API_HOST}/api/translate/${id}/${targetLanguage}`;
			fetch(url)
				.then(() => {
					toast.success($_('modals.translation.toasts.started'));
					open = false;
				})
				.catch((error) => {
					console.error(error);
					toast.error($_('modals.translation.toasts.error'));
				});
		}
	};
</script>

<Modal bind:open size="md" let:close>
	<div class="p-6">
		<div class="flex items-center justify-between mb-5">
			<div class="flex items-center gap-2.5">
				<div class="w-8 h-8 rounded-lg bg-primary/15 border border-primary/25 flex items-center justify-center">
					<Languages size={14} class="text-primary" />
				</div>
				<h2 class="text-foreground text-[0.95rem] font-semibold">{$_('modals.translation.title')}</h2>
			</div>
			<button
				on:click={close}
				class="w-7 h-7 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
			>
				<X size={14} />
			</button>
		</div>

		{#if tr}
			<label for="targetLan" class="text-muted-foreground mb-2 block text-[0.72rem] uppercase tracking-wider">
				{$_('modals.translation.targetLanguages', { values: { language: tr.result.language } })}
			</label>
			<Select id="targetLan" bind:value={targetLanguage}>
				<option disabled selected value={null}>{$_('modals.translation.pickOne')}</option>
				{#each availableLanguages as lan}
					{#if lan.code == tr.result.language}
						{#each lan.targets as t}
							{#if t != tr.result.language}
								<option value={t}>{t}</option>
							{/if}
						{/each}
					{/if}
				{/each}
			</Select>

			<div class="flex items-center justify-end gap-2 mt-6">
				<Button variant="secondary" size="sm" on:click={close}>{$_('common.cancel')}</Button>
				<Button size="sm" on:click={() => handleTranslate(tr.id)} disabled={!targetLanguage}>
					<Languages size={14} />
					{$_('modals.translation.translate')}
				</Button>
			</div>
		{/if}
	</div>
</Modal>
