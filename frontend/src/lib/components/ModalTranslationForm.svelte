<script>
	import toast from 'svelte-french-toast';
	import { CLIENT_API_HOST } from '$lib/utils';
	import { _, locale } from 'svelte-i18n';
	import { getTargetLanguageOptions } from '$lib/languages.js';
	import Modal from './ui/Modal.svelte';
	import Button from './ui/Button.svelte';
	import Select from './ui/Select.svelte';
	import Spinner from './ui/Spinner.svelte';
	import { Languages, X } from 'lucide-svelte';

	export let open = false;
	export let tr;
	export let availableLanguages;
	export let languagesLoading = false;
	let targetLanguage = null;

	$: sourceLanguage = tr?.result?.language;
	$: targetOptions = getTargetLanguageOptions(
		availableLanguages,
		sourceLanguage,
		$locale || 'en'
	);

	$: if (open && tr && targetOptions.length) {
		if (!targetLanguage || !targetOptions.some((o) => o.code === targetLanguage)) {
			targetLanguage =
				targetOptions.find((o) => o.code === 'en')?.code ?? targetOptions[0].code;
		}
	}

	$: if (!open) {
		targetLanguage = null;
	}

	const handleTranslate = (id) => {
		if (!targetLanguage) return;

		open = false;
		toast.success($_('modals.translation.toasts.started'));

		const url = `${CLIENT_API_HOST}/api/translate/${id}/${targetLanguage}`;
		fetch(url).catch((error) => {
			console.error(error);
			toast.error($_('modals.translation.toasts.error'));
		});
	};
</script>

<Modal bind:open size="sm" let:close>
	<div class="p-6">
		<div class="flex items-center justify-between mb-6">
			<div class="flex items-center gap-2.5">
				<div
					class="w-8 h-8 rounded-lg bg-primary/15 border border-primary/25 flex items-center justify-center"
				>
					<Languages size={14} class="text-primary" />
				</div>
				<h2 class="text-foreground text-[0.95rem] font-semibold">
					{$_('modals.translation.title')}
				</h2>
			</div>
			<button
				on:click={close}
				class="w-7 h-7 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
			>
				<X size={14} />
			</button>
		</div>

		{#if tr}
			{#if languagesLoading}
				<div class="flex items-center justify-center py-10 text-muted-foreground">
					<Spinner size={22} />
				</div>
			{:else if targetOptions.length === 0}
				<p class="text-muted-foreground text-[0.85rem] leading-relaxed">
					{$_('modals.translation.noLanguagesAvailable')}
				</p>
			{:else}
				<label
					for="targetLan"
					class="text-muted-foreground mb-2 block text-[0.72rem] uppercase tracking-wider"
				>
					{$_('modals.translation.targetLanguage')}
				</label>
				<Select id="targetLan" bind:value={targetLanguage} class="mb-6">
					{#each targetOptions as option (option.code)}
						<option value={option.code}>{option.label}</option>
					{/each}
				</Select>

				<Button
					class="w-full py-2.5 text-[0.85rem] font-bold tracking-wide"
					disabled={!targetLanguage}
					on:click={() => handleTranslate(tr.id)}
				>
					<Languages size={15} />
					{$_('modals.translation.translate')}
				</Button>
			{/if}
		{/if}
	</div>
</Modal>
