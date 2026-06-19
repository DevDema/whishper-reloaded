<script>
	export let checked = false;
	export let disabled = false;
	export let label = '';
</script>

<label class="flex items-center gap-2 cursor-pointer group select-none">
	<button
		type="button"
		role="switch"
		aria-checked={checked}
		aria-label={label}
		{disabled}
		on:click={() => (checked = !checked)}
		class="relative rounded-full transition-all shrink-0 disabled:opacity-40 disabled:cursor-not-allowed {checked
			? 'bg-primary'
			: 'bg-muted border border-border'}"
		style="height:1.15rem;width:2.1rem;"
	>
		<span
			class="absolute top-1/2 -translate-y-1/2 rounded-full bg-white shadow-sm transition-all"
			style="width:0.85rem;height:0.85rem;left:{checked ? 'calc(100% - 0.85rem - 2px)' : '2px'};"
		></span>
	</button>
	{#if label}
		<span
			class="text-muted-foreground group-hover:text-foreground transition-colors text-[0.78rem]"
		>
			{label}
		</span>
	{/if}
</label>
