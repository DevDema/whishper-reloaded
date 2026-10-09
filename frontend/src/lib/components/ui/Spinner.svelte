<script>
	import { LoaderCircle } from 'lucide-svelte';
	export let size = 18;
	let className = '';
	export { className as class };
</script>

<LoaderCircle {size} class="animate-spin {className}" />
