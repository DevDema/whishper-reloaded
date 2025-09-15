<script>
	import { editorSettings, currentTranscription, currentVideoPlayerTime } from '$lib/stores';
		export let language;
		export let segmentsToShow;
		export let setSegmentsToShow;

    // Segment navigation
	let segmentIndex = '';
	
	   function navigateToSegment() {
		   const index = parseInt(segmentIndex);
		   if (index > 0 && index <= getCurrentSegments().length) {
			   // Check if we need to load more segments to reach the target
			   if (index > segmentsToShow) {
				   setSegmentsToShow(index + 10);
				   setTimeout(() => {
					   scrollToSegmentRow(index);
				   }, 200)
			   } else {
				   // Target segment should already be loaded
				   scrollToSegmentRow(index);
			   }
		   }
	   }
	
	function scrollToSegmentRow(index) {
		// Try multiple selectors to find the segment rows
		let allRows = document.querySelectorAll('table.table tbody tr[data-start]');
		if (allRows.length === 0) {
			// Fallback: try without table class
			allRows = document.querySelectorAll('tbody tr[data-start]');
		}
		if (allRows.length === 0) {
			// Fallback: try any tr with data-start
			allRows = document.querySelectorAll('tr[data-start]');
		}
		
		console.log(`Found ${allRows.length} rows, looking for index ${index - 1}`);
		const targetSegmentRow = allRows[index - 1]; // Convert to 0-based index
		
		if (targetSegmentRow) {
			console.log(`Navigating to segment ${index}`, targetSegmentRow);
			targetSegmentRow.scrollIntoView({ 
				behavior: 'smooth', 
				block: 'center' 
			});
			// Briefly highlight the target row
			targetSegmentRow.classList.add('ring-2', 'ring-primary');
			setTimeout(() => {
				targetSegmentRow.classList.remove('ring-2', 'ring-primary');
			}, 2000);
		} else {
			console.log(`Target segment row not found for index ${index}`);
			// If still not found, try again after a longer delay
			setTimeout(() => {
				let retryRows = document.querySelectorAll('table.table tbody tr[data-start]');
				if (retryRows.length === 0) {
					retryRows = document.querySelectorAll('tbody tr[data-start]');
				}
				if (retryRows.length === 0) {
					retryRows = document.querySelectorAll('tr[data-start]');
				}
				
				console.log(`Retry: Found ${retryRows.length} rows`);
				const retryTargetRow = retryRows[index - 1];
				if (retryTargetRow) {
					console.log(`Retry success: Navigating to segment ${index}`);
					retryTargetRow.scrollIntoView({ 
						behavior: 'smooth', 
						block: 'center' 
					});
					retryTargetRow.classList.add('ring-2', 'ring-primary');
					setTimeout(() => {
						retryTargetRow.classList.remove('ring-2', 'ring-primary');
					}, 2000);
				} else {
					console.log(`Retry failed: Still can't find segment ${index}`);
				}
			}, 300);
		}
	}
	
	   function getCurrentSegments() {
		   if (language === 'original') {
			   return $currentTranscription.result.segments;
		   } else {
			   return $currentTranscription.translations
				   .filter(translation => translation.targetLanguage === language)[0]
				   ?.result.segments || [];
		   }
	   }
</script>

<div class="menu menu-horizontal bg-base-200 gap-2 p-2 rounded-box">
    <input
        type="number"
        min="0"
        max={getCurrentSegments().length}
        placeholder="0-{getCurrentSegments().length}"
        bind:value={segmentIndex}
        on:keydown={(e) => e.key === 'Enter' && navigateToSegment()}
        class="input input-sm w-30 input-bordered"
    />
    <!-- Go to segment button (purple style) -->
    <button
        on:click={navigateToSegment}
        class="btn btn-primary btn-xs md:btn-sm btn-square"
        title="Go to segment"
        disabled={!segmentIndex || parseInt(segmentIndex) < 0 || parseInt(segmentIndex) > getCurrentSegments().length}
    >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
        </svg>
    </button>
    <!-- Go & Seek button (purple style) -->
    <button
        on:click={() => { navigateToSegment(); if ($editorSettings.seekOnClick && segmentIndex) { $currentVideoPlayerTime = getCurrentSegments()[parseInt(segmentIndex)]?.start; } }}
        class="btn btn-primary btn-xs md:btn-sm btn-square"
        title="Go & Seek to segment start"
        disabled={!segmentIndex || parseInt(segmentIndex) < 1 || parseInt(segmentIndex) > getCurrentSegments().length}
    >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            <circle cx="6" cy="12" r="2" fill="currentColor" />
        </svg>
    </button>
</div>