<script>
	import { createEventDispatcher } from 'svelte';
	import { fade, scale } from 'svelte/transition';

	export let open = false;
	export let size = 'md'; // sm | md | lg

	const dispatch = createEventDispatcher();

	const widths = {
		sm: 'max-w-sm',
		md: 'max-w-lg',
		lg: 'max-w-2xl'
	};

	function close() {
		open = false;
		dispatch('close');
	}

	function onKeydown(e) {
		if (e.key === 'Escape' && open) close();
	}
</script>

<svelte:window on:keydown={onKeydown} />

{#if open}
	<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-md"
		on:click|self={close}
		transition:fade={{ duration: 150 }}
		role="presentation"
	>
		<div
			class="relative w-full {widths[size]} rounded-2xl border border-border bg-card text-card-foreground shadow-2xl overflow-hidden max-h-[90vh] overflow-y-auto"
			transition:scale={{ duration: 180, start: 0.96 }}
			role="dialog"
			aria-modal="true"
		>
			<div
				class="absolute top-0 left-8 right-8 h-px bg-gradient-to-r from-transparent via-primary/40 to-transparent"
			></div>
			<slot {close} />
		</div>
	</div>
{/if}
