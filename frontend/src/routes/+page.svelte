<script>
	import { Toaster } from 'svelte-french-toast';
	import { transcriptions, uploadProgress } from '$lib/stores';
	import { browser, dev } from '$app/environment';
	import { CLIENT_WS_HOST } from '$lib/utils';
	import { onMount, onDestroy } from 'svelte';
	import ModalTranscriptionForm from '$lib/components/ModalTranscriptionForm.svelte';
	import ModalDownloadOptions from '$lib/components/ModalDownloadOptions.svelte';
	import ModalTranslationForm from '$lib/components/ModalTranslationForm.svelte';
	import SuccessTranscription from '$lib/components/SuccessTranscription.svelte';
	import PendingTranscription from '$lib/components/PendingTranscription.svelte';
	import PendingTranslation from '$lib/components/PendingTranslation.svelte';
	import ErrorTranscription from '$lib/components/ErrorTranscription.svelte';
	import { env } from '$env/dynamic/public';

	let availableLanguages = [];
    let languagesError = null; 
    let languagesLoading = true;
    let languagesAvailable = false; 
	let transcriptionServiceAvailable = true;
	let transcriptionServiceStatus = '';
	let transcriptionServiceError = '';
	let intervalId;
	let page = 1;
    const itemsPerPage = 10;
	let updateQueue = [];
	let updateTimeout = null;
	const UPDATE_INTERVAL = 1_000; // ms

    $: totalPages = Math.ceil($transcriptions.length / itemsPerPage);

    $: paginatedTranscriptions = $transcriptions.slice().reverse()
        .slice((page - 1) * itemsPerPage, page * itemsPerPage);

    function prevPage() {
        if (page > 1) page -= 1;
    }

    function nextPage() {
        if (page < totalPages) page += 1;
    }

    const getAvailableLangs = () => {
        const fetchLanguages = () => {
            fetch(`${env.PUBLIC_TRANSLATION_API_HOST}/languages`)
            .then(res => res.json())
            .then(data => {
                if (data && data.length > 0) {
                        availableLanguages = data;
                        languagesAvailable = true;
                        languagesError = null;
                        languagesLoading = false;
                        clearInterval(fetchLanguagesInterval);
				} else {
					throw new Error("No languages returned");
				}
            })
			.catch(err => {
				languagesError = "Could not fetch available languages.";
				languagesLoading = false;
			});
        };
        fetchLanguages();
    };



	const checkTranscriptionsAvailability = () => {
		fetch(`${env.PUBLIC_API_HOST}/api/status`)
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



    onMount(async () => {
        getAvailableLangs();
		checkTranscriptionsAvailability();
		intervalId = setInterval(checkTranscriptionsAvailability, 30_000);
		connect();
    });

	onDestroy(() => {
        clearInterval(intervalId);

		if (socket) {
			socket.close(1000);
		}
    });
	
	//export let data;
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
			updateQueue.push(update);
			if (!updateTimeout) {
				updateTimeout = setTimeout(flushUpdates, UPDATE_INTERVAL);
			}
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

	function flushUpdates() {
		transcriptions.update(transcripts => {
			let updated = [...transcripts];
			// Take up to 10 items
			const chunk = updateQueue.splice(0, 10); 
			for (let update of chunk) {
				const index = updated.findIndex(tr => tr.id === update.id);
				if (index >= 0) {
					console.log("update", index, ":", update);
					updated[index] = update;
				} else {
					console.log("push update:", update);
					updated.push(update);
				}
			}
			// If there are more items left, schedule another flush soon
			if (updateQueue.length > 0 && !updateTimeout) {
				updateTimeout = setTimeout(flushUpdates, UPDATE_INTERVAL); // e.g., UPDATE_INTERVAL = 100
			} else {
				updateTimeout = null;
			}
			return updated;
		});
	}
</script>

<Toaster />
<ModalDownloadOptions tr={downloadTranscription} />
<ModalTranslationForm tr={translateTranscription} availableLanguages={availableLanguages} />
<ModalTranscriptionForm />

<header>
	<h1 class="flex items-center justify-center mt-8 space-x-4 text-4xl font-bold">
		<span>
			<img class="w-20 h-20" src="/logo.svg" alt="Logo: a cloud whispering" />
		</span>
		<span> Whishper </span>
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
        <div class="flex items-center bg-yellow-200 border-l-4 border-yellow-400 p-4 mb-4 text-yellow-900 font-semibold">
            <span>⚠️ Language features are unavailable: {languagesError}</span>
			<button 
				on:click={getAvailableLangs}
				class="ml-4 px-2 py-1 bg-red-400 hover:bg-red-500 text-white text-sm rounded"
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
	<div class="items-center mb-0 text-center card-body">
		{#if $transcriptions.length > 0}
			 {#each paginatedTranscriptions as tr (tr.id)}
				{#if tr.status == 2}
					<SuccessTranscription {tr} on:download={handleDownload} on:translate={handleTranslate} languagesAvailable={languagesAvailable} />
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

			<div class="flex justify-center space-x-4 my-4">
				<button on:click={prevPage} disabled={page === 1} class="btn btn-sm btn-ghost">Previous</button>
				<span>Page {page} of {totalPages}</span>
				<button on:click={nextPage} disabled={page === totalPages} class="btn btn-sm btn-ghost">Next</button>
			</div>
		{:else}
			<p class="text-2xl font-bold text-center">🔮 No transcriptions yet 🔮</p>
		{/if}
	</div>
</main>
