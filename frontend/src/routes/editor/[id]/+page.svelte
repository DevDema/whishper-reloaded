<script>
	/** @type {import('./$types').PageData} */
	import { onMount } from 'svelte';
	import { Toaster } from 'svelte-french-toast';
	import Editor from '$lib/components/Editor.svelte';
	import { currentVideoPlayerTime, currentTranscription, audioMode } from '$lib/stores';
	import { CLIENT_API_HOST } from '$lib/utils';
	import { _ } from 'svelte-i18n';
	import { writable } from 'svelte/store';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import IconButton from '$lib/components/ui/IconButton.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import { Play, Pause, LocateFixed } from 'lucide-svelte';

	let language = writable('original');
	let segmentsToShow = writable(20);
	let video;
	let tolerance = 0.1; // Tolerance level in seconds
	let canPlay = false;
	$: if (canPlay && video && Math.abs(video.currentTime - $currentVideoPlayerTime) > tolerance) {
		// When testing in Chrome, it works, just see https://stackoverflow.com/a/67584611
		video.currentTime = $currentVideoPlayerTime;
	}

	function getCurrentSegments() {
		if ($language === 'original') {
			return $currentTranscription.result.segments;
		} else {
			return (
				$currentTranscription.translations.filter(
					(translation) => translation.targetLanguage === language
				)[0]?.result.segments || []
			);
		}
	}

	function scrollToSegmentRow(index) {
		const allRows = document.querySelectorAll('[data-segment-row]');
		const targetSegmentRow = allRows[index];
		if (targetSegmentRow) {
			targetSegmentRow.scrollIntoView({ behavior: 'smooth', block: 'center' });
			targetSegmentRow.classList.add('ring-2', 'ring-primary');
			setTimeout(() => {
				targetSegmentRow.classList.remove('ring-2', 'ring-primary');
			}, 2000);
		}
	}

	function navigateToSegment(index) {
		if (index > 0 && index <= getCurrentSegments().length) {
			if (index > $segmentsToShow) {
				segmentsToShow.set(index + 10);
				setTimeout(() => {
					scrollToSegmentRow(index);
				}, 200);
			} else {
				scrollToSegmentRow(index);
			}
		}
	}

	function goToCurrentTimeSegment() {
		const segments = getCurrentSegments();
		if (!segments.length) return;
		let idx = segments.findIndex((seg, i) => {
			const next = segments[i + 1];
			return $currentVideoPlayerTime >= seg.start && (!next || $currentVideoPlayerTime < next.start);
		});

		if (idx === -1) {
			if ($currentVideoPlayerTime < segments[0].start) idx = 0;
			else idx = segments.length - 1;
		}
		navigateToSegment(idx);
	}

	onMount(() => {
		const handleKeyDown = (e) => {
			// F7: Backward 10s
			if (e.code === 'F7' || e.key === 'F7') {
				e.preventDefault();
				$currentVideoPlayerTime = Math.max(0, $currentVideoPlayerTime - 10);
			}

			// F8: Play/Pause toggle
			if (e.code === 'F8' || e.key === 'F8') {
				e.preventDefault();
				if (video) {
					if (video.paused) {
						video.play();
					} else {
						video.pause();
					}
				}
			}

			// F9: Forward 10s
			if (e.code === 'F9' || e.key === 'F9') {
				e.preventDefault();
				if (video && video.duration) {
					$currentVideoPlayerTime = Math.min(video.duration, $currentVideoPlayerTime + 10);
				} else {
					$currentVideoPlayerTime = $currentVideoPlayerTime + 10;
				}
			}
		};

		document.addEventListener('keydown', handleKeyDown);

		return () => {
			document.removeEventListener('keydown', handleKeyDown);
		};
	});

	function fmt(t) {
		if (!t && t !== 0) return '0:00';
		return Math.floor(t / 60) + ':' + String(Math.floor(t % 60)).padStart(2, '0');
	}
</script>

<Toaster
	toastOptions={{
		style: 'background: var(--card); color: var(--foreground); border: 1px solid var(--border);'
	}}
/>
{#if $currentTranscription}
	{#if $audioMode}
		<div class="h-screen overflow-hidden relative bg-background">
			<div class="h-full overflow-auto pb-24">
				<Editor language={language} segmentsToShow={$segmentsToShow} setSegmentsToShow={(v) => segmentsToShow.set(v)} />
			</div>

			<!-- Player bar -->
			<div class="fixed bottom-0 left-0 right-0 bg-card/80 backdrop-blur-md border-t border-border p-4">
				<div class="flex items-center gap-4 max-w-5xl mx-auto">
					<audio
						id="audio"
						bind:this={video}
						on:timeupdate={(e) => ($currentVideoPlayerTime = e.target.currentTime)}
						on:canplay={() => (canPlay = true)}
						on:loadedmetadata={() => (canPlay = true)}
						class="hidden"
					>
						<source src="{CLIENT_API_HOST}/api/video/{$currentTranscription.fileName}" type="video/mp4" />
					</audio>

					<button
						class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center shadow-lg shadow-primary/20 hover:bg-primary/90 transition-all shrink-0"
						on:click={() => video && (video.paused ? video.play() : video.pause())}
					>
						{#if !canPlay}
							<Spinner size={18} />
						{:else if video && !video.paused}
							<Pause size={18} />
						{:else}
							<Play size={18} />
						{/if}
					</button>

					<button
						class="px-3 py-2 rounded-lg border border-border bg-muted text-muted-foreground hover:text-foreground hover:bg-secondary transition-all text-[0.78rem] font-mono shrink-0"
						on:click={() => ($currentVideoPlayerTime = Math.max(0, $currentVideoPlayerTime - 10))}>-10s</button
					>
					<button
						class="px-3 py-2 rounded-lg border border-border bg-muted text-muted-foreground hover:text-foreground hover:bg-secondary transition-all text-[0.78rem] font-mono shrink-0"
						on:click={() => video && video.duration && ($currentVideoPlayerTime = Math.min(video.duration, $currentVideoPlayerTime + 10))}>+10s</button
					>

					<div class="flex-1 min-w-0">
						<input
							type="range"
							class="w-full accent-[var(--primary)]"
							min="0"
							max="100"
							value={video && video.duration ? ($currentVideoPlayerTime / video.duration) * 100 : 0}
							on:input={(e) => video && video.duration && ($currentVideoPlayerTime = (e.target.value / 100) * video.duration)}
						/>
						<div class="flex justify-between text-[0.7rem] font-mono text-muted-foreground mt-1">
							<span>{fmt($currentVideoPlayerTime)}</span>
							<span>{video ? fmt(video.duration) : '0:00'}</span>
						</div>
					</div>

					<IconButton variant="primary" title={$_('editor.goToCurrentSegment')} on:click={goToCurrentTimeSegment}>
						<LocateFixed size={16} />
					</IconButton>

					<Select
						class="w-auto !py-2"
						on:change={(e) => video && (video.playbackRate = e.target.value)}
					>
						<option value={0.5}>0.5x</option>
						<option value={0.75}>0.75x</option>
						<option value={1} selected>1x</option>
						<option value={1.25}>1.25x</option>
						<option value={1.5}>1.5x</option>
						<option value={2}>2x</option>
					</Select>
				</div>
			</div>
		</div>
	{:else}
		<div class="grid h-screen grid-cols-3 overflow-hidden bg-background">
			<div class="col-span-1 overflow-hidden bg-black/60">
				<div class="relative w-full h-full">
					<!-- svelte-ignore a11y-media-has-caption -->
					<video
						id="video"
						controls
						bind:this={video}
						on:timeupdate={(e) => ($currentVideoPlayerTime = e.target.currentTime)}
						on:canplay={() => (canPlay = true)}
						on:loadedmetadata={() => (canPlay = true)}
						class="absolute top-0 left-0 w-full h-full"
					>
						<source src="{CLIENT_API_HOST}/api/video/{$currentTranscription.fileName}" type="video/mp4" />
						<track kind="captions" />
					</video>
				</div>
			</div>
			<div class="col-span-2 overflow-auto bg-background">
				<Editor language={language} segmentsToShow={$segmentsToShow} setSegmentsToShow={(v) => segmentsToShow.set(v)} />
			</div>
		</div>
	{/if}
{:else}
	<div class="flex items-center justify-center w-screen h-screen bg-background text-muted-foreground">
		<Spinner size={40} />
	</div>
{/if}
