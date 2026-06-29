<script>
	import { ChevronDown } from 'lucide-svelte';
	export let value = undefined;
	let className = '';
	export { className as class };

	const base =
		'w-full appearance-none px-3 py-2.5 pr-9 rounded-xl bg-muted border border-border text-foreground focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary/40 transition-all text-[0.82rem]';
</script>

<div class="relative">
	<select class="{base} {className}" bind:value on:change {...$$restProps}>
		<slot />
	</select>
	<ChevronDown
		size={14}
		class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none"
	/>
</div>
