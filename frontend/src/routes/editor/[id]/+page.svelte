<script>
	/** @type {import('./$types').PageData} */
	import { onMount } from 'svelte';
	import toast, { Toaster } from 'svelte-french-toast';
	import Editor from '$lib/components/Editor.svelte';
	import {currentVideoPlayerTime, currentTranscription, audioMode} from '$lib/stores';
	import { CLIENT_API_HOST } from '$lib/utils';
	import { writable } from 'svelte/store';

	let language = writable('original');
	let segmentsToShow = writable(20);
	let video;
	let tolerance = 0.1; // Tolerance level in seconds
	let canPlay = false;
	$: if(canPlay && video && Math.abs(video.currentTime - $currentVideoPlayerTime) > tolerance) {
		console.log(video.currentTime, $currentVideoPlayerTime)
		// When testing in Chrome, it works, just see https://stackoverflow.com/a/67584611
        video.currentTime = $currentVideoPlayerTime;
    }

	// Helper to get current segments (copied from GoToSegment)
	function getCurrentSegments() {
		if ($language === 'original') {
			return $currentTranscription.result.segments;
		} else {
			return $currentTranscription.translations
				.filter(translation => translation.targetLanguage === language)[0]
				?.result.segments || [];
		}
	}

	// Helper to scroll to a segment row (copied from GoToSegment)
	function scrollToSegmentRow(index) {
		let allRows = document.querySelectorAll('table.table tbody tr[data-start]');
		if (allRows.length === 0) {
			allRows = document.querySelectorAll('tbody tr[data-start]');
		}
		if (allRows.length === 0) {
			allRows = document.querySelectorAll('tr[data-start]');
		}
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
            // Check if we need to load more segments to reach the target
            if (index > $segmentsToShow) {
                segmentsToShow.set(index + 10);
                setTimeout(() => {
                    scrollToSegmentRow(index);
                }, 200)
            } else {
                // Target segment should already be loaded
                scrollToSegmentRow(index);
            }
        }
    }

	// Find the segment index for the current time and scroll to it
	function goToCurrentTimeSegment() {
		const segments = getCurrentSegments();
		if (!segments.length) return;
		// Find the last segment whose start is <= current time
		let idx = segments.findIndex((seg, i) => {
			const next = segments[i + 1];
			return $currentVideoPlayerTime >= seg.start && (!next || $currentVideoPlayerTime < next.start);
		});
		console.log('Current index:', idx);

		if (idx === -1) {
			// If before first segment, go to first
			if ($currentVideoPlayerTime < segments[0].start) idx = 0;
			// If after last segment, go to last
			else idx = segments.length - 1;
		}
		navigateToSegment(idx);
	}
	// Keyboard shortcuts for media control
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
</script>

<Toaster />
{#if $currentTranscription}
	{#if $audioMode}
		<div class="h-screen overflow-hidden relative">
			<div class="h-full overflow-auto bg-content pb-20">
				<Editor language={language} segmentsToShow={$segmentsToShow} setSegmentsToShow={(v) => segmentsToShow.set(v)} />
			</div>
			<div class="fixed bottom-0 left-0 right-0 bg-base-200 p-4 border-t">
				<div class="flex items-center gap-4 max-w-full">
					<audio id="audio" 
						   bind:this={video}
						   on:timeupdate={(e) => $currentVideoPlayerTime = e.target.currentTime}
						   on:canplay={() => canPlay = true}
						   on:loadedmetadata={() => canPlay = true}
						   class="hidden">
						<source src="{CLIENT_API_HOST}/api/video/{$currentTranscription.fileName}" type="video/mp4" />
					</audio>
					
					<button class="btn btn-circle btn-primary" on:click={() => video && (video.paused ? video.play() : video.pause())}> 
						{#if !canPlay}
							<span class="loading loading-spinner loading-sm"></span>
						{:else if video && !video.paused}
							<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
								<path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z"/>
							</svg>
						{:else}
							<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24">
								<path d="M8 5v14l11-7z"/>
							</svg>
						{/if}
					</button>
					
					<button class="btn btn-sm" on:click={() => $currentVideoPlayerTime = Math.max(0, $currentVideoPlayerTime - 10)}>-10s</button>
					<button class="btn btn-sm" on:click={() => video && video.duration && ($currentVideoPlayerTime = Math.min(video.duration, $currentVideoPlayerTime + 10))}>+10s</button>
					
					<div class="flex-1 mx-4">
						<input 
							type="range" 
							class="range range-primary w-full" 
							min="0" 
							max="100" 
							value={video && video.duration ? ($currentVideoPlayerTime / video.duration) * 100 : 0}
							on:input={(e) => video && video.duration && ($currentVideoPlayerTime = (e.target.value / 100) * video.duration)}
						/>
						<div class="flex justify-between text-xs text-base-content/70 mt-1">
							<span>{video ? Math.floor($currentVideoPlayerTime / 60) + ':' + String(Math.floor($currentVideoPlayerTime % 60)).padStart(2, '0') : '0:00'}</span>
							<span>{video ? Math.floor(video.duration / 60) + ':' + String(Math.floor(video.duration % 60)).padStart(2, '0') : '0:00'}</span>
						</div>
					</div>
					
					   <!-- Go to current time segment button -->
					   <button
						   class="btn btn-primary btn-xs md:btn-sm btn-square"
						   title="Go to segment at current time"
						   on:click={() => goToCurrentTimeSegment()}
					   >
						   <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							   <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
							   <circle cx="6" cy="12" r="2" fill="currentColor" />
						   </svg>
					   </button>
					   <select class="select select-sm" on:change={(e) => video && (video.playbackRate = e.target.value)}>
						   <option value={0.5}>0.5x</option>
						   <option value={0.75}>0.75x</option>
						   <option value={1} selected>1x</option>
						   <option value={1.25}>1.25x</option>
						   <option value={1.5}>1.5x</option>
						   <option value={2}>2x</option>
					   </select>
					</div>
			</div>
		</div>
	{:else}
		<div class="grid h-screen grid-cols-3 overflow-hidden">
			<div class="col-span-1 overflow-hidden bg-transparent">
				<div class="relative w-full h-full">
					<video id="video" 
						   controls
						   bind:this={video}
						   on:timeupdate={(e) => $currentVideoPlayerTime = e.target.currentTime}
						   on:canplay={() => canPlay = true}
						   on:loadedmetadata={() => canPlay = true}
						   class="absolute top-0 left-0 w-full h-full">
						<source src="{CLIENT_API_HOST}/api/video/{$currentTranscription.fileName}" type="video/mp4" />
						<track kind="captions" />
					</video>
				</div>
			</div>
			<div class="col-span-2 overflow-auto bg-content">
				<Editor language={language} segmentsToShow={$segmentsToShow} setSegmentsToShow={(v) => segmentsToShow.set(v)} />
			</div>
		</div>
	{/if}
{:else}
	<div class="flex items-center justify-center w-screen h-screen">
		<h1>
			<span class="loading loading-bars loading-lg" />
		</h1>
	</div>
{/if}
