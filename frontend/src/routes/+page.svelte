<script>
	import { Toaster } from 'svelte-french-toast';
	import { transcriptions, uploadProgress, currentPage, loadingTranscriptions } from '$lib/stores';
	import { browser, dev } from '$app/environment';
	import { CLIENT_WS_HOST } from '$lib/utils';
	import { onMount, onDestroy } from 'svelte';
	import ModalTranscriptionForm from '$lib/components/ModalTranscriptionForm.svelte';
	import ModalDownloadOptions from '$lib/components/ModalDownloadOptions.svelte';
	import ModalTranslationForm from '$lib/components/ModalTranslationForm.svelte';
	import ModalRenameFile from '$lib/components/ModalRenameFile.svelte';
	import SuccessTranscription from '$lib/components/SuccessTranscription.svelte';
	import PendingTranscription from '$lib/components/PendingTranscription.svelte';
	import PendingTranslation from '$lib/components/PendingTranslation.svelte';
	import ErrorTranscription from '$lib/components/ErrorTranscription.svelte';
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
	const UPDATE_INTERVAL = 1_000; // ms

    $: totalPages = Math.ceil($transcriptions.length / itemsPerPage);

    $: paginatedTranscriptions = $transcriptions.slice().reverse()
        .slice(($currentPage - 1) * itemsPerPage, $currentPage * itemsPerPage);

		import { tick } from "svelte";

	let searchText = "";
	let searchTimeout;
	let searching = false; // active state
	let filteredTranscriptions = [];

	// --- Debounced search Logic ----

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
			// Switch back to paginated
		}
	}

	function doSearch() {
		const text = searchText.toLowerCase();
		filteredTranscriptions = $transcriptions.filter(tr => {
			return tr.fileName.toLowerCase().includes(text);
		}).reverse();
		searching = true;
	}

	function clearSearch() {
		searchText = "";
		filteredTranscriptions = [];
		searching = false;
	}

    function prevPage() {
        if ($currentPage > 1) currentPage.set($currentPage - 1);
    }

    function nextPage() {
        if ($currentPage < totalPages) currentPage.set($currentPage + 1);
    }

    const getAvailableLangs = async () => {
		return fetch(`${env.PUBLIC_TRANSLATION_API_HOST}/languages`)
		.then(res => res.json())
		.then(data => {
			if (data && data.length > 0) {
					availableLanguages = data;
					languagesAvailable = true;
					languagesError = null;
					languagesLoading = false;
			} else {
				throw new Error("No languages returned");
			}
		})
		.catch(err => {
			languagesError = "Could not fetch available languages.";
			languagesLoading = false;
		});
    };

	const checkTranscriptionsAvailability = async () => {
		return fetch(`${env.PUBLIC_API_HOST}/api/status`)
			.then(res => {
				if (!res.ok) throw new Error("Transcription service unavailable");
				return res.json();
			})
			.then(data => {
				if (data.status === "ok") {
					transcriptionServiceAvailable = true;
					transcriptionServiceStatus = data.service_message || '';
					transcriptionServiceError = '';
				} else {
					transcriptionServiceAvailable = false;
					transcriptionServiceStatus = data.service_message || '';
					transcriptionServiceError = 'Transcription service is reporting an error. New transcriptions are unavailable.';
				}
			})
			.catch(error => {
				transcriptionServiceAvailable = false;
				transcriptionServiceStatus = '';
				transcriptionServiceError = 'Could not reach transcription service. Cannot create new transcriptions.';
			});
	};

	const fetchData = async () => {
		loadingTranscriptions.set(true);
		const endpoint = browser ? `${env.PUBLIC_API_HOST}/api/transcriptions` : `${env.PUBLIC_INTERNAL_API_HOST}/api/transcriptions`;

		return fetch(endpoint).then(response => {
			return response.json();
		}).then(ts => {
			if (ts) {
				transcriptions.update(_ => ts.length > 0 ? ts : []);
			} else {
				transcriptions.update(_ => []);
			} 
		}).finally(ts => {
			loadingTranscriptions.set(false);
		});
	}

    onMount(async () => {
		if ($transcriptions.length === 0) {
			await Promise.all([
				getAvailableLangs(),
				checkTranscriptionsAvailability(),
				fetchData()
			]);
		} else {
			await Promise.all([
				getAvailableLangs(),
				checkTranscriptionsAvailability()
			]);
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
            // use update to update the store
            transcriptions.update(transcriptions => {
                let index = transcriptions.findIndex(tr => tr.id === update.id);
                if (index >= 0) {
                    // replace the item at index
                    transcriptions[index] = update;
                } else {
                    // add the new item
                    transcriptions.push(update);
                }
                return transcriptions; // return a new object to trigger reactivity
            });
        };
	}

	let downloadTranscription = null;
	let handleDownload = (event) => {
		downloadTranscription = event.detail; // this will be the transcription to download
		modalDownloadOptions.showModal(); // show the modal
	};
	let translateTranscription = null;
	let handleTranslate = (event) => {
		translateTranscription = event.detail; // this will be the transcription to translate
		modalTranslation.showModal(); // show the modal
	};
	let renameTranscription = null;
	let handleRename = (event) => {
		renameTranscription = event.detail; // this will be the transcription to rename
		modalRename.showModal(); // show the modal
	};
</script>

<Toaster />
<ModalDownloadOptions tr={downloadTranscription} />
<ModalRenameFile tr={renameTranscription} />

{#if !languagesError}
	<ModalTranslationForm tr={translateTranscription} availableLanguages={availableLanguages} />
{/if}
{#if transcriptionServiceAvailable}
	<ModalTranscriptionForm />
{/if}

<header>
	<h1 class="flex items-center justify-center mt-8 space-x-4 text-4xl font-bold">
		<span>
			<img class="w-20 h-20" src="/logo.svg" alt="Logo: a cloud whispering" />
		</span>
		<span> Whishper Reloaded</span>
	</h1>
	<h2 class="font-mono text-center text-md opacity-70">{data.randomSentence}</h2>
</header>

<main class="w-4/6 mx-auto mt-4 mb-8 card bg-neutral text-neutral-content">
	{#if !transcriptionServiceAvailable}
		<div class="flex items-center justify-between bg-red-200 border-l-4 border-red-400 p-4 mb-2 text-red-900 font-semibold">
			<span>🛑 Transcription service is unavailable: {transcriptionServiceError}</span>
			<button 
				on:click={checkTranscriptionsAvailability}
				class="ml-4 px-2 py-1 bg-red-400 hover:bg-red-500 text-white text-sm rounded"
				title="Refresh"
			>
				Refresh
			</button>
		</div>
	{/if}
	{#if languagesError}
        <div class="flex items-center justify-between bg-yellow-200 border-l-4 border-yellow-400 p-4 mb-4 text-yellow-900 font-semibold">
            <span>⚠️ Language features are unavailable: {languagesError}</span>
			<button 
				on:click={getAvailableLangs}
				class="ml-4 px-2 py-1 bg-yellow-400 hover:bg-red-500 text-white text-sm rounded"
				title="Refresh"
			>
				Refresh
			</button>
        </div>
    {/if}
	{#if $uploadProgress > 0}
		<div class="flex flex-col items-center justify-center px-4 pt-4 my-4">
			<progress class="w-full mx-2 progress progress-success" value="{$uploadProgress}" max="100"></progress>
			<span>Uploading...</span>
		</div>
	{:else }
		<button
			class="max-w-md mx-auto mt-8 btn btn-primary btn-md"
			onclick="modalNewTranscription.showModal()"
			disabled={!transcriptionServiceAvailable}>✨ new transcription</button
		>
	{/if}
	{#if $loadingTranscriptions}
		<div class="flex justify-center items-center py-20">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else}
		<div class="max-w-md mx-auto mt-4 mb-8 w-full">
 		<div class="relative flex items-center w-full">
			<input
			type="text"
			bind:value={searchText}
			placeholder="Search transcriptions..."
			class="input input-bordered w-full pr-12"
			on:input={onSearchInput}
			/>
			{#if searchText.length >= 4 && searching}
			<button
				class="absolute right-3 flex items-center justify-center btn btn-sm btn-ghost p-0 min-h-0 h-8 w-8"
				on:click={clearSearch}
				title="Clear"
				tabindex="0"
				type="button"
				style="font-size: 1.45rem; line-height:1"
			>
				✖
			</button>
			{/if}
		</div>
		{#if searchText.length >= 4 && searching}
			<p class="text-xs opacity-60 mt-2 mb-0">{filteredTranscriptions.length} result{filteredTranscriptions.length === 1 ? '' : 's'} found.</p>
		{:else if searchText.length > 0 && searchText.length < 4}
			<p class="text-xs text-warning mt-2 mb-0">Type at least 4 characters to search.</p>
		{/if}
		</div>
		<div class="items-center mb-0 text-center card-body">
			{#if searching && searchText.length >= 4}
				{#if filteredTranscriptions.length > 0}
					{#each filteredTranscriptions as tr (tr.id)}
						{#if tr.status == 2}
							<SuccessTranscription {tr} on:rename={handleRename} on:download={handleDownload} on:translate={handleTranslate} languagesAvailable={languagesAvailable} />
						{/if}
						{#if tr.status < 2 && tr.status >= 0}
							<PendingTranscription {tr} />
						{/if}
						{#if tr.status == 3}
							<PendingTranslation {tr} />
						{/if}
						{#if tr.status < 0}
							<ErrorTranscription {tr} />
						{/if}
					{/each}
				{:else}
					<p class="text-2xl font-bold text-center">🔮 No transcriptions found 🔮</p>
				{/if}
			{:else}
				{#if $transcriptions.length > 0}
					{#each paginatedTranscriptions as tr (tr.id)}
						{#if tr.status == 2}
							<SuccessTranscription {tr} on:rename={handleRename} on:download={handleDownload} on:translate={handleTranslate} languagesAvailable={languagesAvailable} />
						{/if}
						{#if tr.status < 2 && tr.status >= 0}
							<PendingTranscription {tr} />
						{/if}
						{#if tr.status == 3}
							<PendingTranslation {tr} />
						{/if}
						{#if tr.status < 0}
							<ErrorTranscription {tr} />
						{/if}
					{/each}

					<!-- Pagination only if not searching -->
					<div class="flex justify-center space-x-4 my-4">
						<button on:click={prevPage} disabled={$currentPage === 1} class="btn btn-sm btn-ghost">Previous</button>
						<span>Page {$currentPage} of {totalPages}</span>
						<button on:click={nextPage} disabled={$currentPage === totalPages} class="btn btn-sm btn-ghost">Next</button>
					</div>
				{:else}
					<p class="text-2xl font-bold text-center">🔮 No transcriptions yet 🔮</p>
				{/if}
			{/if}
		</div>
	{/if}
</main>

<footer class="text-center py-4 text-sm opacity-70">
	<p>Whishper version: {data.version}</p>
	<p>
		<a href="https://github.com/DevDema/whishper" class="link">Whishper-Reloaded</a> is a fork of
		<a href="https://github.com/pluja/whishper" class="link">Whishper</a>.
		Praise the idea of the original creator.
	</p>
</footer>
