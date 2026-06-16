<script>
    import { uploadJSON } from '$lib/utils';
    import toast from 'svelte-french-toast';
    import { _ } from 'svelte-i18n';
    export let tr;

    let fileInput;
    let selectedFile = null;
    let uploading = false;

    function handleFileSelect(event) {
        const file = event.target.files[0];
        if (file) {
            if (file.type !== 'application/json' && !file.name.endsWith('.json')) {
                toast.error($_('modals.upload.toasts.invalidFile'));
                return;
            }
            selectedFile = file;
        }
    }

    async function handleUpload() {
        if (!selectedFile || !tr) {
            toast.error($_('modals.upload.toasts.noFile'));
            return;
        }

        uploading = true;

        try {
            const fileContent = await selectedFile.text();
            const jsonData = JSON.parse(fileContent);

            const result = await uploadJSON(tr.id, jsonData);

            if (result.noChanges) {
                toast($_('modals.upload.toasts.noChanges'), {
                    icon: '👐'
                });
            } else {
                toast.success($_('modals.upload.toasts.success'));

                // Reset state and close modal
                clearFile();
                document.getElementById('modalUploadJSON').close();
            }  
        } catch (error) {
            console.error('Upload error:', error);
            if (error instanceof SyntaxError) {
                toast.error($_('modals.upload.toasts.invalidFormat'));
            } else {
                toast.error(error.message || $_('modals.upload.toasts.failed'));
            }
        } finally {
            uploading = false;
        }
    }

    function clearFile() {
        selectedFile = null;
        fileInput.value = '';
    }
</script>

<dialog id="modalUploadJSON" class="modal">
    <form method="dialog" class="modal-box flex flex-col items-center justify-center">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
        {#if tr}
        <h1 class="text-center font-bold mt-2 pb-2">{$_('modals.upload.title')}</h1>
        <p class="text-sm opacity-70 text-center mb-4">
            {$_('modals.upload.description')}<br>
            <span class="font-mono">{tr.fileName.split("_WHSHPR_")[1]}</span>
        </p>

        <div class="form-control w-full max-w-md">
            <label for="jsonFile" class="label">
                <span class="label-text font-bold">{$_('modals.upload.selectFile')}</span>
            </label>
            <input
                bind:this={fileInput}
                type="file"
                id="jsonFile"
                accept=".json,application/json"
                on:change={handleFileSelect}
                class="file-input file-input-bordered w-full"
                disabled={uploading}
            />
            {#if selectedFile}
                <div class="mt-2 p-2 bg-base-200 rounded flex justify-between items-center">
                    <span class="text-sm font-mono">{selectedFile.name}</span>
                    <button
                        type="button"
                        on:click={clearFile}
                        class="btn btn-xs btn-ghost"
                        disabled={uploading}
                    >✕</button>
                </div>
            {/if}
        </div>

        <div class="modal-action">
            <button
                type="button"
                on:click={handleUpload}
                class="btn btn-primary"
                disabled={!selectedFile || uploading}
            >
                {#if uploading}
                    <span class="loading loading-spinner loading-sm"></span>
                    {$_('common.uploading')}
                {:else}
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                        <path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2"></path>
                        <path d="M7 9l5 -5l5 5"></path>
                        <path d="M12 4l0 12"></path>
                    </svg>
                    {$_('modals.upload.upload')}
                {/if}
            </button>
        </div>
        {:else}
            <div class="flex items-center justify-center w-full h-32">
                <span class="loading loading-bars loading-lg" />
            </div>
        {/if}
    </form>
    <form method="dialog" class="modal-backdrop">
        <button>{$_('common.close')}</button>
    </form>
</dialog>