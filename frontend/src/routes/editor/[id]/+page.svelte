<script>
	/** @type {import('./$types').PageData} */
	import { onMount } from 'svelte';
	import toast, { Toaster } from 'svelte-french-toast';
	import Editor from '$lib/components/Editor.svelte';
	import {currentVideoPlayerTime, currentTranscription, audioMode} from '$lib/stores';
	import { CLIENT_API_HOST } from '$lib/utils';


	let video;
	let tolerance = 0.1; // Tolerance level in seconds
	let canPlay = false;
	$: if(canPlay && video && Math.abs(video.currentTime - $currentVideoPlayerTime) > tolerance) {
		console.log(video.currentTime, $currentVideoPlayerTime)
		// When testing in Chrome, it works, just see https://stackoverflow.com/a/67584611
        video.currentTime = $currentVideoPlayerTime;
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
				<Editor />
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
						{#if video && !video.paused}
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
				<Editor />
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
