<script>
	import { validateURL, CLIENT_API_HOST } from '$lib/utils.js';
	import { env } from '$env/dynamic/public';
	import { uploadProgress } from '$lib/stores';

	import toast from 'svelte-french-toast';

	let errorMessage = '';
	let disableSubmit = true;
	let modelSize = 'small';
	let language = 'auto';
	let sourceUrl = '';
	let fileInput;
	let device = env.PUBLIC_WHISHPER_PROFILE == 'gpu' ? 'cuda' : 'cpu';

	let languages = [
		'auto',
		'ar',
		'be',
		'bg',
		'bn',
		'ca',
		'cs',
		'cy',
		'da',
		'de',
		'el',
		'en',
		'es',
		'fr',
		'it',
		'ja',
		'nl',
		'pl',
		'pt',
		'ru',
		'sk',
		'sl',
		'sv',
		'tk',
		'tr',
		'zh'
	];
	let models = [
		'tiny',
		'tiny.en',
		'base',
		'base.en',
		'small',
		'small.en',
		'medium',
		'medium.en',
		'large-v2',
		'large-v3'
	];
	// Sort the languages
	languages.sort((a, b) => {
		if (a == 'auto') return -1;
		if (b == 'auto') return 1;
		return a.localeCompare(b);
	});

	// New parameters
	let beamSize = 5;
	let initialPrompt = '';
	let hotwords = '';

	// Function that sends the data as a form to the backend
	async function sendForm() {
	   if (sourceUrl && !validateURL(sourceUrl)) {
		   toast.error('You must enter a valid URL.');
		   return;
	   }

	   if (!sourceUrl && !fileInput) {
		   toast.error('No file or URL.');
		   return;
	   }

	   let formData = new FormData();
	   formData.append('language', language);
	   formData.append('modelSize', modelSize);
	   if (device == 'cuda' || device == 'cpu') {
		   formData.append('device', device);
	   } else {
		   formData.append('device', 'cpu');
	   }
	   formData.append('sourceUrl', sourceUrl);
	   if (sourceUrl == '') {
		   formData.append('file', fileInput.files[0]);
	   }
	   // Add new params
	   formData.append('beam_size', beamSize);
	   if (initialPrompt && initialPrompt.trim() !== '') {
		   formData.append('initial_prompt', initialPrompt);
	   }
	   if (hotwords && hotwords.trim() !== '') {
		   // Send as comma-separated string
		   formData.append('hotwords', hotwords);
	   }

	   return new Promise((resolve, reject) => {
		   const xhr = new XMLHttpRequest();

		   // Set up progress event listener
		   xhr.upload.addEventListener('progress', (event) => {
			   if (event.lengthComputable) {
				   const percentCompleted = Math.round((event.loaded * 100) / event.total);
				   uploadProgress.set(percentCompleted);
			   }
		   });

		   // Set up load event listener
		   xhr.addEventListener('load', () => {
			   if (xhr.status === 200) {
				   resolve(xhr.response);
				   toast.success('Success!');
			   } else {
				   reject(xhr.statusText);
				   toast.error('Upload failed');
			   }
			   uploadProgress.set(0); // Reset progress after completion
		   });

		   // Set up error event listener
		   xhr.addEventListener('error', () => {
			   reject(xhr.statusText);
			   toast.error('An error occurred during upload');
			   uploadProgress.set(0); // Reset progress on error
		   });

		   xhr.open('POST', `${CLIENT_API_HOST}/api/transcriptions`);
		   xhr.send(formData);
	   });

	   // Set file and sourceUrl to empty
	   sourceUrl = '';
	   fileInput.value = '';
	   uploadProgress.set(0);

	   toast.success('Success!');
	}

	// Reactive statement
	$: if (sourceUrl && !validateURL(sourceUrl)) {
		errorMessage = 'Enter a valid URL';
		disableSubmit = true;
	} else {
		errorMessage = '';
		disableSubmit = false;
	}
</script>


<style>
	.centered-modal-box {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-width: 350px;
		width: 100%;
		max-width: 480px;
		margin: 0 auto;
	}
	.full-width-btn {
		width: 100%;
		min-width: 200px;
		max-width: 100%;
		display: block;
	}
	.horizontal-selectors {
		display: flex;
		flex-direction: row;
		gap: 1rem;
		width: 100%;
		justify-content: center;
	}
	.horizontal-selectors > div {
		flex: 1 1 0;
		min-width: 0;
	}
</style>

<dialog id="modalNewTranscription" class="modal">
	<form method="dialog" class="modal-box centered-modal-box">
		<button class="absolute btn btn-sm btn-circle btn-ghost right-2 top-2">✕</button>
		{#if errorMessage != ''}
			<div class="alert alert-error">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="w-6 h-6 stroke-current shrink-0"
					fill="none"
					viewBox="0 0 24 24"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
				<span>{errorMessage}</span>
			</div>
		{/if}
			<div class="mt-0 space-y-2 w-full flex flex-col items-center">
				<div class="w-full form-control">
					<label for="file" class="label">
						<span class="label-text">Pick a file</span>
					</label>
					<input
						name="file"
						bind:this={fileInput}
						type="file"
						class="w-full file-input file-input-sm file-input-bordered file-input-primary"
					/>
				</div>

				<div class="w-full form-control">
					<label for="sourceUrl" class="label">
						<span class="label-text">Or a source URL</span>
					</label>
					<input
						name="sourceUrl"
						bind:value={sourceUrl}
						type="text"
						placeholder="https://youtube.com/watch?v=Hd33fCdW"
						class="w-full input input-sm input-bordered input-primary"
					/>
				</div>
			</div>

		<div class="mb-0 divider w-full" />
		<!-- Whisper Configuration -->
		<div class="horizontal-selectors mb-2">
			<div class="form-control">
				<label for="modelSize" class="label">
					<span class="label-text">Whisper model</span>
				</label>
				<select name="modelSize" bind:value={modelSize} class="select select-bordered">
					{#each models as m}
						<option value={m}>{m}</option>
					{/each}
				</select>
			</div>
			<div class="form-control">
				<label for="language" class="label">
					<span class="label-text">Language</span>
				</label>
				<select name="language" bind:value={language} class="select select-bordered">
					{#each languages as l}
						<option value={l}>{l}</option>
					{/each}
				</select>
			</div>
			<div class="form-control">
				<label for="device" class="label">
					<span class="label-text">Device</span>
				</label>
				<select name="device" bind:value={device} class="select select-bordered">
					{#if env.PUBLIC_WHISHPER_PROFILE == 'gpu'}
						<option selected value="cuda">GPU</option>
						<option value="cpu">CPU</option>
					{:else}
						<option selected value="cpu">CPU</option>
						<option disabled value="cuda">GPU</option>
					{/if}
				</select>
			</div>
		</div>

		<div class="mb-0 divider w-full" />

			<div class="flex flex-row flex-wrap gap-4 w-full">
				<div class="w-full form-control">
					<label for="beamSize" class="label">
						<span class="label-text">Beam size</span>
					</label>
					<input
						name="beamSize"
						type="number"
						min="1"
						max="20"
						bind:value={beamSize}
						class="w-full input input-sm input-bordered input-primary"
					/>
				</div>

				<div class="w-full form-control">
					<label for="initialPrompt" class="label">
						<span class="label-text">Initial prompt (optional)</span>
					</label>
					<textarea
						name="initialPrompt"
						bind:value={initialPrompt}
						class="w-full input input-sm input-bordered input-primary"
						placeholder="e.g. context for the model"
						rows="3"
					/>
				</div>

				<div class="w-full form-control">
					<label for="hotwords" class="label">
						<span class="label-text">Hotwords (comma separated, optional)</span>
					</label>
					<input
						name="hotwords"
						type="text"
						bind:value={hotwords}
						class="w-full input input-sm input-bordered input-primary"
						placeholder="e.g. keyword1, keyword2"
					/>
				</div>
			</div>

		<div class="mb-0 divider w-full" />
		<!--Actions-->
		<button class="btn btn-primary full-width-btn mt-2" on:click={sendForm} disabled={disableSubmit}
			>Start</button
		>
	</form>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
