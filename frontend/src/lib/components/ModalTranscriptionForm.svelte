<script>
	import { validateURL, CLIENT_API_HOST } from '$lib/utils.js';
	import { env } from '$env/dynamic/public';
	import { uploadProgress } from '$lib/stores';
	import { _ } from 'svelte-i18n';
	import toast from 'svelte-french-toast';
	import Modal from './ui/Modal.svelte';
	import Button from './ui/Button.svelte';
	import Input from './ui/Input.svelte';
	import Textarea from './ui/Textarea.svelte';
	import Select from './ui/Select.svelte';
	import Switch from './ui/Switch.svelte';
	import { Sparkles, X, Upload, FileAudio, Mic } from 'lucide-svelte';

	export let open = false;

	let errorMessage = '';
	let disableSubmit = true;
	let modelSize = 'small';
	let language = 'auto';
	let sourceUrl = '';
	let fileInput;
	let selectedFile = null;
	let dragging = false;
	let submitting = false;
	let device = env.PUBLIC_WHISHPER_PROFILE == 'gpu' ? 'cuda' : 'cpu';

	let languages = [
		'auto', 'ar', 'be', 'bg', 'bn', 'ca', 'cs', 'cy', 'da', 'de', 'el', 'en', 'es', 'fr', 'it',
		'ja', 'nl', 'pl', 'pt', 'ru', 'sk', 'sl', 'sv', 'tk', 'tr', 'zh'
	];
	let models = [
		'tiny', 'tiny.en', 'base', 'base.en', 'small', 'small.en', 'medium', 'medium.en', 'large-v2', 'large-v3'
	];
	languages.sort((a, b) => {
		if (a == 'auto') return -1;
		if (b == 'auto') return 1;
		return a.localeCompare(b);
	});

	// New parameters
	let beamSize = 5;
	let initialPrompt = '';
	let hotwords = '';

	// VAD settings
	let enableVad = false;
	let vadThreshold = '';
	let vadMinSpeech = '';
	let vadMinSilence = '';

	function handleDrop(e) {
		e.preventDefault();
		dragging = false;
		const f = e.dataTransfer.files[0];
		if (f) selectedFile = f;
	}

	function handleFileSelect(e) {
		const f = e.target.files?.[0];
		if (f) selectedFile = f;
	}

	function clearFile() {
		selectedFile = null;
		if (fileInput) fileInput.value = '';
	}

	async function sendForm() {
		if (sourceUrl && !validateURL(sourceUrl)) {
			toast.error($_('modals.transcription.toasts.invalidUrl'));
			return;
		}

		if (!sourceUrl && !selectedFile) {
			toast.error($_('modals.transcription.toasts.noFileOrUrl'));
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
			formData.append('file', selectedFile);
		}
		formData.append('beam_size', beamSize);
		if (initialPrompt && initialPrompt.trim() !== '') {
			formData.append('initial_prompt', initialPrompt);
		}
		if (hotwords && hotwords.trim() !== '') {
			formData.append('hotwords', hotwords);
		}
		if (enableVad) {
			formData.append('vad_filter', 'true');
			if (vadThreshold !== '') formData.append('vad_threshold', vadThreshold);
			if (vadMinSpeech !== '') formData.append('vad_min_speech_duration_ms', vadMinSpeech);
			if (vadMinSilence !== '') formData.append('vad_min_silence_duration_ms', vadMinSilence);
		}

		return new Promise((resolve, reject) => {
			const xhr = new XMLHttpRequest();

			xhr.upload.addEventListener('progress', (event) => {
				if (event.lengthComputable) {
					const percentCompleted = Math.round((event.loaded * 100) / event.total);
					uploadProgress.set(percentCompleted);
				}
			});

			xhr.addEventListener('load', () => {
				if (xhr.status === 200) {
					resolve(xhr.response);
					toast.success($_('modals.transcription.toasts.success'));
				} else {
					reject(xhr.statusText);
					toast.error($_('modals.transcription.toasts.uploadFailed'));
				}
				uploadProgress.set(0);
			});

			xhr.addEventListener('error', () => {
				reject(xhr.statusText);
				toast.error($_('modals.transcription.toasts.uploadError'));
				uploadProgress.set(0);
			});

			xhr.open('POST', `${CLIENT_API_HOST}/api/transcriptions`);
			xhr.send(formData);
		});
	}

	async function handleStart() {
		if (submitting) return;
		submitting = true;
		try {
			await sendForm();
			clearFile();
			sourceUrl = '';
			open = false;
		} catch (e) {
			console.error(e);
		} finally {
			submitting = false;
		}
	}

	$: if (sourceUrl && !validateURL(sourceUrl)) {
		errorMessage = $_('modals.transcription.errorValidUrl');
		disableSubmit = true;
	} else {
		errorMessage = '';
		disableSubmit = false;
	}
</script>

<Modal bind:open size="lg" let:close>
	<!-- Header -->
	<div class="flex items-center justify-between px-6 pt-5 pb-4 border-b border-border">
		<div class="flex items-center gap-2.5">
			<div class="w-8 h-8 rounded-lg bg-primary/20 border border-primary/30 flex items-center justify-center">
				<Sparkles size={15} class="text-primary" />
			</div>
			<h3 class="text-foreground text-[0.95rem] font-semibold">{$_('home.newTranscription')}</h3>
		</div>
		<button
			on:click={close}
			class="w-7 h-7 rounded-lg flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
		>
			<X size={15} />
		</button>
	</div>

	<!-- Body -->
	<div class="px-6 py-5 space-y-5">
		{#if errorMessage}
			<div class="flex items-center gap-2 px-3 py-2 rounded-lg border border-red-500/20 bg-red-500/5 text-red-400 text-[0.8rem]">
				{errorMessage}
			</div>
		{/if}

		<!-- Drop zone -->
		<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
		<div
			on:dragover={(e) => { e.preventDefault(); dragging = true; }}
			on:dragleave={() => (dragging = false)}
			on:drop={handleDrop}
			on:click={() => fileInput?.click()}
			class="relative flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-8 cursor-pointer transition-all {dragging
				? 'border-primary/60 bg-primary/5'
				: selectedFile
				? 'border-emerald-500/40 bg-emerald-500/5'
				: 'border-border hover:border-primary/40 hover:bg-muted/50'}"
			role="button"
			tabindex="0"
		>
			<input bind:this={fileInput} type="file" class="hidden" accept="audio/*,video/*" on:change={handleFileSelect} />
			{#if selectedFile}
				<div class="w-12 h-12 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center">
					<FileAudio size={22} class="text-emerald-400" />
				</div>
				<div class="text-center">
					<p class="text-foreground text-[0.85rem] font-medium">{selectedFile.name}</p>
					<p class="text-muted-foreground text-[0.72rem] mt-0.5">
						{(selectedFile.size / 1024 / 1024).toFixed(1)} MB
					</p>
				</div>
				<button
					on:click|stopPropagation={clearFile}
					class="text-muted-foreground hover:text-foreground transition-colors text-[0.72rem]"
				>
					{$_('common.cancel')}
				</button>
			{:else}
				<div class="w-12 h-12 rounded-xl bg-muted border border-border flex items-center justify-center">
					<Upload size={20} class="text-muted-foreground" />
				</div>
				<div class="text-center">
					<p class="text-foreground text-[0.85rem] font-medium">{$_('modals.transcription.pickFile')}</p>
					<p class="text-muted-foreground text-[0.72rem] mt-0.5">{$_('modals.transcription.dropHint')}</p>
				</div>
			{/if}
		</div>

		<!-- Source URL -->
		<div class="space-y-1.5">
			<label for="sourceUrl" class="text-muted-foreground text-[0.75rem]">{$_('modals.transcription.sourceUrl')}</label>
			<Input id="sourceUrl" bind:value={sourceUrl} placeholder={$_('modals.transcription.sourceUrlPlaceholder')} />
		</div>

		<!-- Model / Language / Device -->
		<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
			<div class="space-y-1.5">
				<label for="modelSize" class="text-muted-foreground text-[0.75rem]">{$_('modals.transcription.whisperModel')}</label>
				<Select id="modelSize" bind:value={modelSize}>
					{#each models as m}<option value={m}>{m}</option>{/each}
				</Select>
			</div>
			<div class="space-y-1.5">
				<label for="language" class="text-muted-foreground text-[0.75rem]">{$_('modals.transcription.language')}</label>
				<Select id="language" bind:value={language}>
					{#each languages as l}<option value={l}>{l}</option>{/each}
				</Select>
			</div>
			<div class="space-y-1.5">
				<label for="device" class="text-muted-foreground text-[0.75rem]">{$_('modals.transcription.device')}</label>
				<Select id="device" bind:value={device}>
					{#if env.PUBLIC_WHISHPER_PROFILE == 'gpu'}
						<option value="cuda">{$_('modals.transcription.deviceGpu')}</option>
						<option value="cpu">{$_('modals.transcription.deviceCpu')}</option>
					{:else}
						<option value="cpu">{$_('modals.transcription.deviceCpu')}</option>
						<option disabled value="cuda">{$_('modals.transcription.deviceGpu')}</option>
					{/if}
				</Select>
			</div>
		</div>

		<!-- Advanced -->
		<div class="grid grid-cols-1 gap-3">
			<div class="space-y-1.5">
				<label for="beamSize" class="text-muted-foreground text-[0.75rem]">{$_('modals.transcription.beamSize')}</label>
				<Input id="beamSize" type="number" min="1" max="20" bind:value={beamSize} />
			</div>
			<div class="space-y-1.5">
				<label for="initialPrompt" class="text-muted-foreground text-[0.75rem]">{$_('modals.transcription.initialPrompt')}</label>
				<Textarea id="initialPrompt" rows="3" bind:value={initialPrompt} placeholder={$_('modals.transcription.initialPromptPlaceholder')} />
			</div>
			<div class="space-y-1.5">
				<label for="hotwords" class="text-muted-foreground text-[0.75rem]">{$_('modals.transcription.hotwords')}</label>
				<Input id="hotwords" bind:value={hotwords} placeholder={$_('modals.transcription.hotwordsPlaceholder')} />
			</div>
		</div>

		<!-- VAD -->
		<div class="rounded-xl border border-border bg-muted/40 p-4 space-y-3">
			<Switch bind:checked={enableVad} label={$_('modals.transcription.enableVad')} />
			{#if enableVad}
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
					<div class="space-y-1.5">
						<label for="vadThreshold" class="text-muted-foreground text-[0.72rem]">{$_('modals.transcription.vadThreshold')}</label>
						<Input id="vadThreshold" type="number" min="0" max="1" step="0.01" bind:value={vadThreshold} placeholder={$_('modals.transcription.vadThresholdPlaceholder')} />
					</div>
					<div class="space-y-1.5">
						<label for="vadMinSpeech" class="text-muted-foreground text-[0.72rem]">{$_('modals.transcription.vadMinSpeech')}</label>
						<Input id="vadMinSpeech" type="number" min="0" bind:value={vadMinSpeech} placeholder={$_('modals.transcription.vadMinSpeechPlaceholder')} />
					</div>
					<div class="space-y-1.5">
						<label for="vadMinSilence" class="text-muted-foreground text-[0.72rem]">{$_('modals.transcription.vadMinSilence')}</label>
						<Input id="vadMinSilence" type="number" min="0" bind:value={vadMinSilence} placeholder={$_('modals.transcription.vadMinSilencePlaceholder')} />
					</div>
				</div>
			{/if}
		</div>

		{#if $uploadProgress > 0}
			<div class="space-y-1.5">
				<div class="h-2 w-full rounded-full bg-muted overflow-hidden">
					<div class="h-full rounded-full bg-primary transition-all duration-300" style="width:{$uploadProgress}%"></div>
				</div>
				<p class="text-center text-muted-foreground text-[0.75rem]">{$_('common.uploading')}</p>
			</div>
		{/if}
	</div>

	<!-- Footer -->
	<div class="flex items-center justify-end gap-2 px-6 py-4 border-t border-border">
		<Button variant="secondary" on:click={close}>{$_('common.cancel')}</Button>
		<Button on:click={handleStart} disabled={disableSubmit || submitting || (!selectedFile && !sourceUrl)}>
			<Mic size={14} />
			{$_('modals.transcription.start')}
		</Button>
	</div>
</Modal>
