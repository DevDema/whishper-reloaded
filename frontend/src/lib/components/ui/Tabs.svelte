<script>
	import { createEventDispatcher } from 'svelte';

	export let tabs = []; // [{ id, label }]
	export let active = tabs.length ? tabs[0].id : undefined;

	const dispatch = createEventDispatcher();

	function select(id) {
		active = id;
		dispatch('change', id);
	}
</script>

<div class="flex items-center gap-1 border-b border-border/0">
	{#each tabs as tab}
		<button
			on:click={() => select(tab.id)}
			class="px-4 py-2 transition-all border-b-2 -mb-px text-[0.82rem] {active === tab.id
				? 'border-primary text-foreground font-medium'
				: 'border-transparent text-muted-foreground hover:text-foreground'}"
		>
			{tab.label}
		</button>
	{/each}
</div>
