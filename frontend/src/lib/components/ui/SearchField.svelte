<script>
	import { Search, X } from 'lucide-svelte';
	import { createEventDispatcher } from 'svelte';

	export let value = '';
	export let placeholder = '';
	let className = '';
	export { className as class };

	const dispatch = createEventDispatcher();

	function clear() {
		value = '';
		dispatch('clear');
	}
</script>

<div class="relative {className}">
	<Search
		size={15}
		class="absolute left-3.5 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none"
	/>
	<input
		type="text"
		bind:value
		{placeholder}
		on:input
		class="w-full pl-10 pr-10 py-2.5 rounded-xl bg-muted border border-border text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary/40 transition-all text-[0.85rem]"
	/>
	{#if value}
		<button
			type="button"
			on:click={clear}
			title="Clear"
			class="absolute right-2.5 top-1/2 -translate-y-1/2 w-6 h-6 rounded-md flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-secondary transition-all"
		>
			<X size={14} />
		</button>
	{/if}
</div>
