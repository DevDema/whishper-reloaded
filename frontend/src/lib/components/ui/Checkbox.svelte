<script>
	export let checked = false;
	export let disabled = false;
	let className = '';
	export { className as class };
</script>

<input
	type="checkbox"
	bind:checked
	{disabled}
	on:change
	class="h-4 w-4 rounded-md border border-border bg-muted accent-[var(--primary)] focus:outline-none focus:ring-2 focus:ring-primary/40 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed {className}"
	{...$$restProps}
/>
