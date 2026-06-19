<script>
	export let variant = 'default'; // default | primary | danger
	export let size = 'md'; // md (32px) | sm (28px)
	export let type = 'button';
	export let disabled = false;
	export let title = '';
	let className = '';
	export { className as class };

	const base =
		'shrink-0 rounded-lg flex items-center justify-center border transition-all hover:scale-105 active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:scale-100 focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/40';

	const sizes = {
		md: 'w-8 h-8',
		sm: 'w-7 h-7'
	};

	const variants = {
		default:
			'border-border bg-muted text-muted-foreground hover:text-foreground hover:bg-secondary',
		primary:
			'border-primary/50 bg-primary text-primary-foreground hover:bg-primary/90 shadow-lg shadow-primary/20',
		danger: 'border-red-500/20 bg-red-500/10 text-red-400 hover:bg-red-500/20'
	};
</script>

<button
	{type}
	{disabled}
	{title}
	class="{base} {sizes[size]} {variants[variant]} {className}"
	on:click
	{...$$restProps}
>
	<slot />
</button>
