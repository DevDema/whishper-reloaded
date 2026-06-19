<script>
	import { STATUS_CONFIG } from './ui/statusConfig';
	import { createEventDispatcher } from 'svelte';

	export let kind = 'done';
	export let interactive = false;

	const dispatch = createEventDispatcher();

	$: cfg = STATUS_CONFIG[kind] ?? STATUS_CONFIG.done;

	function onKeydown(e) {
		if (interactive && (e.key === 'Enter' || e.key === ' ')) {
			e.preventDefault();
			dispatch('activate');
		}
	}
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="group relative flex items-center gap-4 px-5 py-4 rounded-xl bg-card border border-border hover:border-border/80 transition-all duration-200 text-left shadow-[0_1px_3px_rgba(0,0,0,0.12)] {interactive
		? 'cursor-pointer hover:bg-card/80'
		: ''}"
	role={interactive ? 'button' : undefined}
	tabindex={interactive ? 0 : undefined}
	on:click={() => interactive && dispatch('activate')}
	on:keydown={onKeydown}
>
	<!-- Status accent bar -->
	<div class="absolute left-0 top-3 bottom-3 w-0.5 rounded-full {cfg.dot}"></div>

	<!-- Status icon -->
	<div class="flex items-center justify-center w-8 h-8 rounded-lg bg-muted shrink-0">
		<svelte:component
			this={cfg.icon}
			size={15}
			class="{cfg.iconColor} {cfg.spin ? 'animate-spin' : ''}"
		/>
	</div>

	<!-- Info -->
	<div class="flex-1 min-w-0">
		<slot />
	</div>

	<!-- Actions -->
	<div class="flex items-center gap-1.5 shrink-0">
		<slot name="actions" />
	</div>
</div>
