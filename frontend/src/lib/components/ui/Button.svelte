<script>
	export let variant = 'primary'; // primary | secondary | ghost | danger
	export let size = 'md'; // sm | md
	export let type = 'button';
	export let disabled = false;
	let className = '';
	export { className as class };

	const base =
		'inline-flex items-center justify-center gap-2 rounded-xl font-semibold transition-all duration-200 disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:scale-100 focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/40';

	const variants = {
		primary:
			'bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg shadow-primary/20 hover:shadow-primary/30 hover:scale-[1.02] active:scale-[0.98]',
		secondary:
			'border border-border bg-muted hover:bg-secondary text-muted-foreground hover:text-foreground',
		ghost: 'text-muted-foreground hover:text-foreground hover:bg-muted',
		danger: 'border border-red-500/20 bg-red-500/10 text-red-400 hover:bg-red-500/20'
	};

	const sizes = {
		sm: 'px-3 py-2 text-[0.8rem]',
		md: 'px-4 py-2.5 text-[0.85rem]'
	};
</script>

<button
	{type}
	{disabled}
	class="{base} {variants[variant]} {sizes[size]} {className}"
	on:click
	{...$$restProps}
>
	<slot />
</button>
