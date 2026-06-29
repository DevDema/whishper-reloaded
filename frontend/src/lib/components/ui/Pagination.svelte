<script>
	import { ChevronsLeft, ChevronsRight } from 'lucide-svelte';

	export let page = 1;
	export let totalPages = 1;
	export let prevLabel = 'First';
	export let nextLabel = 'Last';

	$: pages = Array.from({ length: Math.min(7, totalPages) }, (_, i) =>
		totalPages <= 7
			? i + 1
			: page <= 4
			? i + 1
			: page >= totalPages - 3
			? totalPages - 6 + i
			: page - 3 + i
	);

	function go(p) {
		if (p >= 1 && p <= totalPages) page = p;
	}
</script>

{#if totalPages > 1}
	<div class="flex items-center justify-center gap-3 pt-2">
		<button
			on:click={() => go(1)}
			disabled={page === 1}
			class="flex items-center gap-1.5 px-3 py-2 rounded-lg border border-border bg-card hover:bg-muted disabled:opacity-30 disabled:cursor-not-allowed transition-colors text-muted-foreground hover:text-foreground text-[0.8rem]"
		>
			<ChevronsLeft size={14} />
			{prevLabel}
		</button>

		<div class="flex items-center gap-1">
			{#each pages as p}
				<button
					on:click={() => go(p)}
					class="w-8 h-8 rounded-lg transition-colors text-[0.8rem] {p === page
						? 'bg-primary text-primary-foreground'
						: 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
				>
					{p}
				</button>
			{/each}
		</div>

		<button
			on:click={() => go(totalPages)}
			disabled={page === totalPages}
			class="flex items-center gap-1.5 px-3 py-2 rounded-lg border border-border bg-card hover:bg-muted disabled:opacity-30 disabled:cursor-not-allowed transition-colors text-muted-foreground hover:text-foreground text-[0.8rem]"
		>
			{nextLabel}
			<ChevronsRight size={14} />
		</button>
	</div>
{/if}
