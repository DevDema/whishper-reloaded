<script>
	import { editorSettings } from '$lib/stores';
	import { currentVideoPlayerTime, currentTranscription, editorHistory } from '$lib/stores';
	import { afterUpdate } from 'svelte';
	import { _ } from 'svelte-i18n';
	import IconButton from './ui/IconButton.svelte';
	import { Clock, ArrowUpToLine, ArrowDownToLine, Split, Trash2 } from 'lucide-svelte';
	export let index;
	export let segment;
	export let translationIndex;

	let isActive = false;

	function getCps(s) {
		let duration = s.end - s.start;
		let charCount = s.text.length;
		let cps = charCount / duration;
		cps = Math.round(cps * 100) / 100;
		return Math.trunc(cps);
	}

	function deleteSegment(index, callback) {
		const source =
			translationIndex == -1
				? $currentTranscription.result.segments
				: $currentTranscription.translations[translationIndex].result.segments;
		source.splice(index, 1);
		$currentTranscription = { ...$currentTranscription };
		callback();
	}

	function splitSegment(index, callback) {
		const source =
			translationIndex == -1
				? $currentTranscription.result.segments
				: $currentTranscription.translations[translationIndex].result.segments;

		const segment = source[index];
		const words = segment.text.split(' ');
		const half = Math.ceil(words.length / 2);
		const firstHalf = words.slice(0, half).join(' ');
		const secondHalf = words.slice(half).join(' ');

		const duration = segment.end - segment.start;
		const midTime = segment.start + duration / 2;

		segment.text = firstHalf;
		segment.end = midTime;

		const newSegment = {
			id: JSON.stringify(Date.now()),
			start: midTime,
			end: segment.end + duration / 2,
			text: secondHalf,
			words: []
		};
		source.splice(index + 1, 0, newSegment);

		$currentTranscription = { ...$currentTranscription };
		callback();
	}

	// Text changes only save after 6 keystrokes
	let keystrokes = 0;
	function handleKeystrokes() {
		keystrokes++;
		if (keystrokes > 6) {
			handleHistory();
			keystrokes = 0;
		}
	}

	function handleHistory() {
		let currentT = JSON.parse(JSON.stringify($currentTranscription));
		editorHistory.update((value) => {
			return [...value, currentT];
		});
	}

	function insertSegmentAbove(index, callback) {
		const source =
			translationIndex == -1
				? $currentTranscription.result.segments
				: $currentTranscription.translations[translationIndex].result.segments;
		source.splice(index, 0, {
			id: JSON.stringify(Date.now()),
			start: 0,
			end: 0,
			text: '',
			words: []
		});
		$currentTranscription = { ...$currentTranscription };
		callback();
	}

	function insertSegmentBelow(index, callback) {
		const source =
			translationIndex == -1
				? $currentTranscription.result.segments
				: $currentTranscription.translations[translationIndex].result.segments;
		source.splice(index + 1, 0, {
			id: JSON.stringify(Date.now()),
			start: 0,
			end: 0,
			text: '',
			words: []
		});
		$currentTranscription = { ...$currentTranscription };
		callback();
	}

	// For custom tab navigation
	let textInputEl;
	function handleTab(e) {
		if (e.key === 'Tab') {
			e.preventDefault();
			const allInputs = Array.from(document.querySelectorAll('[data-segment-text-input]'));
			const idx = allInputs.indexOf(textInputEl);
			if (e.shiftKey) {
				if (idx > 0) {
					allInputs[idx - 1].focus();
				}
			} else {
				if (idx !== -1 && idx < allInputs.length - 1) {
					allInputs[idx + 1].focus();
				}
			}
		}
	}
	afterUpdate(() => {
		if (textInputEl) {
			textInputEl.setAttribute('data-segment-text-input', '');
		}
	});

	$: if (segment.start <= $currentVideoPlayerTime && $currentVideoPlayerTime <= segment.end) {
		isActive = true;
	} else {
		isActive = false;
	}

	$: cpsValue = getCps(segment);
	$: cpsHigh = cpsValue > 16;
</script>

<div
	data-segment-row
	data-start={segment.start}
	class="group rounded-xl border bg-card transition-all px-4 py-3 {isActive
		? 'border-amber-400/50 ring-1 ring-amber-400/30 bg-amber-500/5'
		: 'border-border hover:border-border/80'}"
>
	<div class="flex items-start gap-3">
		<!-- Index -->
		<span class="text-muted-foreground text-[0.7rem] font-mono mt-2.5 w-6 text-right shrink-0 tabular-nums">
			{index}
		</span>

		<!-- Timestamps -->
		<div class="flex flex-col gap-1.5 shrink-0">
			<div class="flex items-center gap-1">
				<input
					class="w-20 px-2 py-1.5 rounded-lg bg-muted border border-border text-foreground text-[0.78rem] font-mono focus:outline-none focus:ring-2 focus:ring-primary/40"
					type="number"
					step="0.01"
					bind:value={segment.start}
					on:input={(e) => ($currentVideoPlayerTime = e.target.value)}
					on:input={handleHistory}
					on:click={(e) => {
						if ($editorSettings.seekOnClick) $currentVideoPlayerTime = e.target.value;
					}}
				/>
				<IconButton
					variant="primary"
					class="w-7 h-7"
					title={$_('editor.segment.setToCurrentTime')}
					on:click={() => {
						segment.start = $currentVideoPlayerTime;
						handleHistory();
					}}
				>
					<Clock size={13} />
				</IconButton>
			</div>
			<div class="flex items-center gap-1">
				<input
					class="w-20 px-2 py-1.5 rounded-lg bg-muted border border-border text-foreground text-[0.78rem] font-mono focus:outline-none focus:ring-2 focus:ring-primary/40"
					type="number"
					step="0.01"
					bind:value={segment.end}
					on:input={(e) => ($currentVideoPlayerTime = e.target.value)}
					on:input={handleHistory}
					on:click={(e) => {
						if ($editorSettings.seekOnClick) $currentVideoPlayerTime = e.target.value;
					}}
				/>
				<IconButton
					variant="primary"
					class="w-7 h-7"
					title={$_('editor.segment.setToCurrentTime')}
					on:click={() => {
						segment.end = $currentVideoPlayerTime;
						handleHistory();
					}}
				>
					<Clock size={13} />
				</IconButton>
			</div>
		</div>

		<!-- Text -->
		<div class="flex-1 min-w-0">
			<!-- svelte-ignore a11y-no-static-element-interactions -->
			<div
				bind:this={textInputEl}
				bind:textContent={segment.text}
				on:input={handleKeystrokes}
				on:keydown={handleTab}
				class="w-full min-h-[2.5rem] p-2.5 rounded-lg bg-muted border text-[0.85rem] font-mono leading-relaxed text-foreground focus:outline-none focus:ring-2 focus:ring-primary/40 {cpsHigh
					? 'border-red-500/60'
					: 'border-border'}"
				contenteditable="true"
				tabindex="0"
			/>
			<div class="flex items-center gap-3 mt-1.5 text-[0.68rem] font-mono text-muted-foreground">
				<span class={cpsHigh ? 'text-red-400' : ''}>{$_('editor.segment.cps', { values: { value: cpsValue } })}</span>
				<span class="text-border">·</span>
				<span>{$_('editor.segment.duration', { values: { value: Math.round((segment.end - segment.start) * 100) / 100 } })}</span>
			</div>
		</div>

		<!-- Actions -->
		<div class="flex flex-col gap-1 shrink-0 opacity-60 group-hover:opacity-100 transition-opacity">
			<IconButton variant="primary" size="sm" title={$_('editor.segment.insertAbove')} on:click={() => insertSegmentAbove(index, handleHistory)}>
				<ArrowUpToLine size={13} />
			</IconButton>
			<IconButton variant="primary" size="sm" title={$_('editor.segment.insertBelow')} on:click={() => insertSegmentBelow(index, handleHistory)}>
				<ArrowDownToLine size={13} />
			</IconButton>
			<IconButton variant="primary" size="sm" title={$_('editor.segment.split')} on:click={() => splitSegment(index, handleHistory)}>
				<Split size={13} />
			</IconButton>
			<IconButton variant="danger" size="sm" title={$_('transcription.actions.delete')} on:click={() => deleteSegment(index, handleHistory)}>
				<Trash2 size={13} />
			</IconButton>
		</div>
	</div>
</div>
