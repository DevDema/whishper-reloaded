<script>
	import { Loader2 } from 'lucide-svelte';
	export let size = 18;
	let className = '';
	export { className as class };
</script>

<Loader2 {size} class="animate-spin {className}" />
