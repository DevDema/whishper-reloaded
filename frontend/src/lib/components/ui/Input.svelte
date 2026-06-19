<script>
	export let value = '';
	let className = '';
	export { className as class };

	const base =
		'w-full px-3.5 py-2.5 rounded-xl bg-muted border border-border text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary/40 transition-all text-[0.85rem]';
</script>

<input
	class="{base} {className}"
	bind:value
	on:input
	on:change
	on:keydown
	on:focus
	on:blur
	{...$$restProps}
/>
