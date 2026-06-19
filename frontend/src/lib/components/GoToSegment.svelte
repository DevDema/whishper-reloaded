<script>
	import { editorSettings, currentTranscription, currentVideoPlayerTime } from '$lib/stores';
	import { _ } from 'svelte-i18n';
	import IconButton from './ui/IconButton.svelte';
	import { ArrowRight, LocateFixed } from 'lucide-svelte';
	export let language;
	export let segmentsToShow;
	export let setSegmentsToShow;

	let segmentIndex = '';

	function navigateToSegment() {
		const index = parseInt(segmentIndex);
		if (index > 0 && index <= getCurrentSegments().length) {
			if (index > segmentsToShow) {
				setSegmentsToShow(index + 10);
				setTimeout(() => {
					scrollToSegmentRow(index);
				}, 200);
			} else {
				scrollToSegmentRow(index);
			}
		}
	}

	function scrollToSegmentRow(index) {
		let allRows = document.querySelectorAll('[data-segment-row]');
		const targetSegmentRow = allRows[index - 1]; // Convert to 0-based index

		if (targetSegmentRow) {
			targetSegmentRow.scrollIntoView({ behavior: 'smooth', block: 'center' });
			targetSegmentRow.classList.add('ring-2', 'ring-primary');
			setTimeout(() => {
				targetSegmentRow.classList.remove('ring-2', 'ring-primary');
			}, 2000);
		} else {
			setTimeout(() => {
				const retryRows = document.querySelectorAll('[data-segment-row]');
				const retryTargetRow = retryRows[index - 1];
				if (retryTargetRow) {
					retryTargetRow.scrollIntoView({ behavior: 'smooth', block: 'center' });
					retryTargetRow.classList.add('ring-2', 'ring-primary');
					setTimeout(() => {
						retryTargetRow.classList.remove('ring-2', 'ring-primary');
					}, 2000);
				}
			}, 300);
		}
	}

	function getCurrentSegments() {
		if (language === 'original') {
			return $currentTranscription.result.segments;
		} else {
			return (
				$currentTranscription.translations.filter(
					(translation) => translation.targetLanguage === language
				)[0]?.result.segments || []
			);
		}
	}
</script>

<div class="flex items-center gap-1.5 px-2 py-1.5 rounded-lg border border-border bg-muted/50">
	<input
		type="number"
		min="0"
		max={getCurrentSegments().length}
		placeholder={$_('editor.goTo.placeholder', { values: { max: getCurrentSegments().length } })}
		bind:value={segmentIndex}
		on:keydown={(e) => e.key === 'Enter' && navigateToSegment()}
		class="w-24 px-2 py-1.5 rounded-md bg-muted border border-border text-foreground text-[0.78rem] focus:outline-none focus:ring-2 focus:ring-primary/40"
	/>
	<IconButton
		variant="primary"
		title={$_('editor.goTo.goToSegment')}
		on:click={navigateToSegment}
		disabled={!segmentIndex || parseInt(segmentIndex) < 0 || parseInt(segmentIndex) > getCurrentSegments().length}
	>
		<ArrowRight size={15} />
	</IconButton>
	<IconButton
		variant="primary"
		title={$_('editor.goTo.goAndSeek')}
		on:click={() => {
			navigateToSegment();
			if ($editorSettings.seekOnClick && segmentIndex) {
				$currentVideoPlayerTime = getCurrentSegments()[parseInt(segmentIndex)]?.start;
			}
		}}
		disabled={!segmentIndex || parseInt(segmentIndex) < 1 || parseInt(segmentIndex) > getCurrentSegments().length}
	>
		<LocateFixed size={15} />
	</IconButton>
</div>
