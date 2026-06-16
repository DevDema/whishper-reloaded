<script>
    import toast from 'svelte-french-toast';
    import { renameTranscription } from '$lib/utils.js';
    import { _ } from 'svelte-i18n';

    export let tr;  // transcription object passed from parent
    let newFileName = '';
    let fileExtension = '';
    let isSubmitting = false;
    let initialized = false;
    
    // Extract filename and extension when tr changes
    $: if (tr && !initialized) {
        const parts = tr.fileName.split('WHSHPR_');
        if (parts.length > 1) {
            // Extract the extension from the second part
            const lastPart = parts[1];
            const extIndex = lastPart.lastIndexOf('.');
            if (extIndex !== -1) {
                fileExtension = lastPart.slice(extIndex);
                newFileName = lastPart.slice(0, extIndex);
            } else {
                newFileName = lastPart;
                fileExtension = '';
            }
        } else {
            // Handle files without WHSHPR_ prefix
            const extIndex = tr.fileName.lastIndexOf('.');
            if (extIndex !== -1) {
                fileExtension = tr.fileName.slice(extIndex);
                newFileName = tr.fileName.slice(0, extIndex);
            } else {
                newFileName = tr.fileName;
                fileExtension = '';
            }
        }
        initialized = true;
    }

    // Get the original name for comparison (without extension)
    $: originalName = tr ? (() => {
        const parts = tr.fileName.split('WHSHPR_');
        const nameWithExt = parts.length > 1 ? parts[1] : tr.fileName;
        const extIndex = nameWithExt.lastIndexOf('.');
        return extIndex !== -1 ? nameWithExt.slice(0, extIndex) : nameWithExt;
    })() : '';

    $: canSubmit = tr && newFileName && newFileName !== originalName;

    async function closeModal() {
        tr = null;
        initialized = false;
        newFileName = '';
        document.getElementById('modalRename').close();
    }

    async function handleSubmit() {
        if (!canSubmit) return;
        
        isSubmitting = true;
        try {
            const fullNewFileName = newFileName + fileExtension;
            await renameTranscription(tr.id, fullNewFileName);
            toast.success($_('modals.rename.toasts.success'));
            closeModal();
        } catch (error) {
            console.error('Error renaming file:', error);
            toast.error($_('modals.rename.toasts.error'));
        } finally {
            isSubmitting = false;
        }
    }
</script>

<dialog class="modal" id="modalRename">
    <form method="dialog" class="modal-box w-11/12 max-w-xl">
        <button 
            type="submit"
            class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
            on:click={closeModal} 
        >✕</button>
        
        <h3 class="font-bold text-lg mb-4">{$_('modals.rename.title')}</h3>
        
        {#if tr}
            <div class="form-control w-full">
                <label class="label" for="newFileName">
                    <span class="label-text">{$_('modals.rename.newFileName')}</span>
                </label>
                <div class="join w-full">
                    <input
                        type="text"
                        id="newFileName"
                        bind:value={newFileName}
                        class="input input-bordered join-item w-full"
                        placeholder={$_('modals.rename.placeholder')}
                    />
                    {#if fileExtension}
                        <div class="join-item flex items-center px-3 bg-base-200 border border-base-300">
                            <span class="text-base-content/70">{fileExtension}</span>
                        </div>
                    {/if}
                </div>
            </div>

            <div class="modal-action">
                <button
                    type="submit"
                    class="btn btn-primary"
                    disabled={!canSubmit || isSubmitting}
                    on:click={handleSubmit}
                >
                    {#if isSubmitting}
                        <span class="loading loading-spinner"></span>
                    {:else}
                        {$_('modals.rename.rename')}
                    {/if}
                </button>
            </div>
        {/if}
    </form>
</dialog>
