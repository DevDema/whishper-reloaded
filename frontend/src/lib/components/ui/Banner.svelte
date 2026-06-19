<script>
	import { XCircle, AlertTriangle, RefreshCw } from 'lucide-svelte';
	import { slide } from 'svelte/transition';

	export let variant = 'error'; // error | warning
	export let retryLabel = '';

	const variants = {
		error: {
			wrap: 'border-red-500/20 bg-red-500/5 text-red-400',
			btn: 'border-red-500/30 bg-red-500/10 hover:bg-red-500/20 text-red-400',
			icon: XCircle
		},
		warning: {
			wrap: 'border-amber-500/20 bg-amber-500/5 text-amber-400',
			btn: 'border-amber-500/30 bg-amber-500/10 hover:bg-amber-500/20 text-amber-400',
			icon: AlertTriangle
		}
	};

	$: cfg = variants[variant];
</script>

<div
	transition:slide={{ duration: 200 }}
	class="flex items-center gap-3 px-4 py-3 rounded-xl border text-[0.82rem] {cfg.wrap}"
>
	<svelte:component this={cfg.icon} size={15} class="shrink-0" />
	<span class="flex-1"><slot /></span>
	{#if retryLabel}
		<button
			on:click
			class="ml-auto flex items-center gap-1.5 px-3 py-1 rounded-lg border transition-colors text-[0.75rem] {cfg.btn}"
		>
			<RefreshCw size={12} />
			{retryLabel}
		</button>
	{/if}
</div>
