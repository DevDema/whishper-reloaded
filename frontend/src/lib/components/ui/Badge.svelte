<script>
	export let mono = false;
	let className = '';
	export { className as class };
</script>

<span
	class="shrink-0 inline-flex items-center px-2 py-0.5 rounded-md border text-[0.65rem] leading-none {mono
		? 'font-mono'
		: ''} {className}"
>
	<slot />
</span>
