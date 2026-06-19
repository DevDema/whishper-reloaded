<script>
	import { Sun, Moon } from 'lucide-svelte';
	import { theme } from '$lib/theme';
</script>

<button
	on:click={theme.toggle}
	title={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
	aria-label="Toggle theme"
	class="relative w-8 h-8 rounded-lg border border-border bg-muted hover:bg-secondary text-muted-foreground hover:text-foreground transition-all flex items-center justify-center overflow-hidden"
>
	{#if $theme === 'dark'}
		<Sun size={14} />
	{:else}
		<Moon size={14} />
	{/if}
</button>
