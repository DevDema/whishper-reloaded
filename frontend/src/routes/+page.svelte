<script>
	import { Toaster } from 'svelte-french-toast';
	import { transcriptions, uploadProgress, currentPage, loadingTranscriptions } from '$lib/stores';
	import { browser } from '$app/environment';
	import { CLIENT_WS_HOST } from '$lib/utils';
	import { onMount, onDestroy, tick } from 'svelte';
	import ModalTranscriptionForm from '$lib/components/ModalTranscriptionForm.svelte';
	import ModalDownloadOptions from '$lib/components/ModalDownloadOptions.svelte';
	import ModalTranslationForm from '$lib/components/ModalTranslationForm.svelte';
	import ModalRenameFile from '$lib/components/ModalRenameFile.svelte';
	import ModalUploadJSON from '$lib/components/ModalUploadJSON.svelte';
	import SuccessTranscription from '$lib/components/SuccessTranscription.svelte';
	import RunningTranscription from '$lib/components/RunningTranscription.svelte';
	import PendingTranscription from '$lib/components/PendingTranscription.svelte';
	import PendingTranslation from '$lib/components/PendingTranslation.svelte';
	import ErrorTranscription from '$lib/components/ErrorTranscription.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import SearchField from '$lib/components/ui/SearchField.svelte';
	import Banner from '$lib/components/ui/Banner.svelte';
	import Pagination from '$lib/components/ui/Pagination.svelte';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import ThemeToggle from '$lib/components/ui/ThemeToggle.svelte';
	import { Plus, Layers, Search } from 'lucide-svelte';
	import { _ } from 'svelte-i18n';
	import { env } from '$env/dynamic/public';

	let availableLanguages = [];
	let languagesError = null;
	let languagesLoading = true;
	let languagesAvailable = false;
	let transcriptionServiceAvailable = false;
	let transcriptionServiceStatus = '';
	let transcriptionServiceError = '';
	let intervalId;
	const itemsPerPage = 10;

	$: totalPages = Math.ceil($transcriptions.length / itemsPerPage);

	$: paginatedTranscriptions = $transcriptions
		.slice()
		.reverse()
		.slice(($currentPage - 1) * itemsPerPage, $currentPage * itemsPerPage);

	let searchText = '';
	let searchTimeout;
	let searching = false; // active state
	let filteredTranscriptions = [];

	// --- Debounced search logic ----

	function onSearchInput(event) {
		searchText = event.target.value;
		searching = false;

		if (searchTimeout) clearTimeout(searchTimeout);

		if (searchText.length >= 4) {
			searchTimeout = setTimeout(async () => {
				await tick();
				doSearch();
			}, 1000);
		} else {
			filteredTranscriptions = []; // Clear if <4 chars
		}
	}

	function doSearch() {
		const text = searchText.toLowerCase();
		filteredTranscriptions = $transcriptions
			.filter((tr) => tr.fileName.toLowerCase().includes(text))
			.reverse();
		searching = true;
	}

	function clearSearch() {
		searchText = '';
		filteredTranscriptions = [];
		searching = false;
	}

	const getAvailableLangs = async () => {
		return fetch(`${env.PUBLIC_TRANSLATION_API_HOST}/languages`)
			.then((res) => res.json())
			.then((data) => {
				if (data && data.length > 0) {
					availableLanguages = data;
					languagesAvailable = true;
					languagesError = null;
					languagesLoading = false;
				} else {
					throw new Error('No languages returned');
				}
			})
			.catch(() => {
				languagesError = $_('layout.serviceErrors.languagesFetch');
				languagesLoading = false;
			});
	};

	const checkTranscriptionsAvailability = async () => {
		return fetch(`${env.PUBLIC_API_HOST}/api/status`)
			.then((res) => {
				if (!res.ok) throw new Error('Transcription service unavailable');
				return res.json();
			})
			.then((data) => {
				if (data.status === 'ok') {
					transcriptionServiceAvailable = true;
					transcriptionServiceStatus = data.service_message || '';
					transcriptionServiceError = '';
				} else {
					transcriptionServiceAvailable = false;
					transcriptionServiceStatus = data.service_message || '';
					transcriptionServiceError = $_('layout.serviceErrors.serviceReporting');
				}
			})
			.catch(() => {
				transcriptionServiceAvailable = false;
				transcriptionServiceStatus = '';
				transcriptionServiceError = $_('layout.serviceErrors.serviceUnreachable');
			});
	};

	const fetchData = async () => {
		loadingTranscriptions.set(true);
		const endpoint = browser
			? `${env.PUBLIC_API_HOST}/api/list-transcriptions`
			: `${env.PUBLIC_INTERNAL_API_HOST}/api/list-transcriptions`;

		return fetch(endpoint)
			.then((response) => response.json())
			.then((ts) => {
				if (ts) {
					transcriptions.update((_) => (ts.length > 0 ? ts : []));
				} else {
					transcriptions.update((_) => []);
				}
			})
			.finally(() => {
				loadingTranscriptions.set(false);
			});
	};

	onMount(async () => {
		if ($transcriptions.length === 0) {
			await Promise.all([getAvailableLangs(), checkTranscriptionsAvailability(), fetchData()]);
		} else {
			await Promise.all([getAvailableLangs(), checkTranscriptionsAvailability()]);
		}

		intervalId = setInterval(checkTranscriptionsAvailability, 30_000);
		connect();
	});

	onDestroy(() => {
		clearInterval(intervalId);
		if (socket) {
			socket.close(1000);
		}
	});

	let socket;
	export let data;

	function connect() {
		if (!browser) {
			console.log('Server, not connecting');
			return;
		}

		let new_uri = '';
		var loc = window.location;
		if (loc.protocol === 'https:') {
			new_uri = 'wss:';
		} else {
			new_uri = 'ws:';
		}
		new_uri += '//' + (CLIENT_WS_HOST == '' ? loc.host : CLIENT_WS_HOST);
		new_uri += '/ws/transcriptions';
		console.log('Connecting to: ', new_uri);
		socket = new WebSocket(new_uri);

		socket.onopen = () => console.log('WebSocket is connected...');
		socket.onerror = (error) => console.log('WebSocket Error: ', error);
		socket.onclose = (event) => {
			console.log('WebSocket is closed with code: ', event.code, ' and reason: ', event.reason);
			setTimeout(() => {
				console.log('Reconnecting...');
				connect();
			}, 1000);
		};

		socket.onmessage = (event) => {
			let update = JSON.parse(event.data);
			transcriptions.update((transcriptions) => {
				let index = transcriptions.findIndex((tr) => tr.id === update.id);
				if (index >= 0) {
					transcriptions[index] = update;
				} else {
					transcriptions.push(update);
				}
				return transcriptions;
			});
		};
	}

	// --- Modal state ---
	let newTranscriptionOpen = false;
	let downloadOpen = false;
	let translateOpen = false;
	let renameOpen = false;
	let uploadOpen = false;

	let downloadTranscription = null;
	let translateTranscription = null;
	let renameTranscription = null;
	let uploadTranscription = null;

	const handleDownload = (event) => {
		downloadTranscription = event.detail;
		downloadOpen = true;
	};
	const handleTranslate = (event) => {
		translateTranscription = event.detail;
		translateOpen = true;
	};
	const handleRename = (event) => {
		renameTranscription = event.detail;
		renameOpen = true;
	};
	const handleUpload = (event) => {
		uploadTranscription = event.detail;
		uploadOpen = true;
	};
</script>

<Toaster
	toastOptions={{
		style: 'background: var(--card); color: var(--foreground); border: 1px solid var(--border);'
	}}
/>
<ModalDownloadOptions bind:open={downloadOpen} tr={downloadTranscription} />
<ModalRenameFile bind:open={renameOpen} tr={renameTranscription} />
<ModalUploadJSON bind:open={uploadOpen} tr={uploadTranscription} />

{#if !languagesError}
	<ModalTranslationForm bind:open={translateOpen} tr={translateTranscription} {availableLanguages} />
{/if}
{#if transcriptionServiceAvailable}
	<ModalTranscriptionForm bind:open={newTranscriptionOpen} />
{/if}

<div class="min-h-screen bg-background">
	<!-- Header -->
	<header class="border-b border-border/50 bg-card/40 backdrop-blur-sm sticky top-0 z-10">
		<div class="max-w-5xl mx-auto px-6 py-4 flex items-center justify-between">
			<div class="flex items-center gap-3">
				<div
					class="w-9 h-9 rounded-xl bg-primary/15 border border-primary/30 flex items-center justify-center overflow-hidden"
				>
					<img class="w-6 h-6" src="/logo.svg" alt={$_('header.logoAlt')} />
				</div>
				<div>
					<h1 class="text-foreground text-[1.05rem] font-semibold leading-tight">
						{$_('header.title')}
					</h1>
					<p class="text-muted-foreground text-[0.72rem]">{$_(`taglines.${data.taglineIndex}`)}</p>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<ThemeToggle />
				<span
					class="px-2 py-0.5 rounded-md bg-muted border border-border text-muted-foreground text-[0.72rem] font-mono"
				>
					v{data.version}
				</span>
			</div>
		</div>
	</header>

	<div class="max-w-5xl mx-auto px-6 py-8 space-y-6">
		<!-- Service errors -->
		{#if !transcriptionServiceAvailable}
			<Banner variant="error" retryLabel={$_('common.refresh')} on:click={checkTranscriptionsAvailability}>
				{$_('home.serviceUnavailable', { values: { error: transcriptionServiceError } })}
			</Banner>
		{/if}
		{#if languagesError}
			<Banner variant="warning" retryLabel={$_('common.refresh')} on:click={getAvailableLangs}>
				{$_('home.languagesUnavailable', { values: { error: languagesError } })}
			</Banner>
		{/if}

		<!-- Upload progress -->
		{#if $uploadProgress > 0}
			<div class="flex flex-col items-center justify-center gap-2 rounded-xl border border-border bg-card px-4 py-4">
				<div class="h-2 w-full rounded-full bg-muted overflow-hidden">
					<div class="h-full rounded-full bg-primary transition-all duration-300" style="width:{$uploadProgress}%"></div>
				</div>
				<span class="text-muted-foreground text-[0.8rem]">{$_('common.uploading')}</span>
			</div>
		{:else}
			<!-- Actions row -->
			<div class="flex items-center gap-3">
				<Button on:click={() => (newTranscriptionOpen = true)} disabled={!transcriptionServiceAvailable}>
					<Plus size={16} />
					{$_('home.newTranscription')}
				</Button>

				<SearchField
					class="flex-1"
					bind:value={searchText}
					placeholder={$_('home.search.placeholder')}
					on:input={onSearchInput}
					on:clear={clearSearch}
				/>

				<div class="flex items-center gap-1.5 text-muted-foreground text-[0.78rem] shrink-0">
					<Layers size={13} />
					<span>{$_('home.fileCount', { values: { count: $transcriptions.length } })}</span>
				</div>
			</div>
		{/if}

		<!-- Search hints -->
		{#if searchText.length >= 4 && searching}
			<p class="text-xs text-muted-foreground -mt-3">
				{$_('home.search.results', { values: { count: filteredTranscriptions.length } })}
			</p>
		{:else if searchText.length > 0 && searchText.length < 4}
			<p class="text-xs text-amber-400 -mt-3">{$_('home.search.minChars')}</p>
		{/if}

		<!-- List -->
		{#if $loadingTranscriptions}
			<div class="flex justify-center items-center py-20 text-muted-foreground">
				<Spinner size={32} />
			</div>
		{:else if searching && searchText.length >= 4}
			{#if filteredTranscriptions.length > 0}
				<div class="space-y-2">
					{#each filteredTranscriptions as tr (tr.id)}
						{#if tr.status == 2}
							<SuccessTranscription {tr} on:rename={handleRename} on:download={handleDownload} on:translate={handleTranslate} on:upload={handleUpload} {languagesAvailable} />
						{/if}
						{#if tr.status == 1}<RunningTranscription {tr} />{/if}
						{#if tr.status == 0}<PendingTranscription {tr} />{/if}
						{#if tr.status == 3}<PendingTranslation {tr} />{/if}
						{#if tr.status < 0}<ErrorTranscription {tr} />{/if}
					{/each}
				</div>
			{:else}
				<div class="text-center py-16 text-muted-foreground">
					<Search size={32} class="mx-auto mb-3 opacity-30" />
					<p class="text-[0.9rem]">{$_('home.empty.noResults')}</p>
				</div>
			{/if}
		{:else if $transcriptions.length > 0}
			<div class="space-y-2">
				{#each paginatedTranscriptions as tr (tr.id)}
					{#if tr.status == 2}
						<SuccessTranscription {tr} on:rename={handleRename} on:download={handleDownload} on:translate={handleTranslate} on:upload={handleUpload} {languagesAvailable} />
					{/if}
					{#if tr.status == 1}<RunningTranscription {tr} />{/if}
					{#if tr.status == 0}<PendingTranscription {tr} />{/if}
					{#if tr.status == 3}<PendingTranslation {tr} />{/if}
					{#if tr.status < 0}<ErrorTranscription {tr} />{/if}
				{/each}
			</div>

			<Pagination
				bind:page={$currentPage}
				{totalPages}
				prevLabel={$_('home.pagination.previous')}
				nextLabel={$_('home.pagination.next')}
			/>
		{:else}
			<div class="text-center py-16 text-muted-foreground">
				<Layers size={32} class="mx-auto mb-3 opacity-30" />
				<p class="text-[0.9rem]">{$_('home.empty.none')}</p>
			</div>
		{/if}
	</div>

	<!-- Footer -->
	<footer class="border-t border-border/30 py-4 mt-8">
		<div class="max-w-5xl mx-auto px-6 text-center text-muted-foreground text-[0.72rem]">
			{$_('footer.version', { values: { version: data.version } })} ·
			<a href="https://github.com/DevDema/whishper" class="hover:text-foreground transition-colors underline">Whishper-Reloaded</a>
			{$_('footer.forkOf')}
			<a href="https://github.com/pluja/whishper" class="hover:text-foreground transition-colors underline">Whishper</a>.
			{$_('footer.praise')}
		</div>
	</footer>
</div>
