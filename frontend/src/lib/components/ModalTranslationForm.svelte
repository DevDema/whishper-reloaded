<script>
    import toast from 'svelte-french-toast';
    import { CLIENT_API_HOST } from '$lib/utils';
    import { _ } from 'svelte-i18n';
    export let tr;
    export let availableLanguages;
    let targetLanguage = null;

    const handleTranslate = (id) => {
        if(targetLanguage) {
            const url = `${CLIENT_API_HOST}/api/translate/${id}/${targetLanguage}`;
            fetch(url)
            .then(() => toast.success($_('modals.translation.toasts.started')))
            .catch(error => {
                console.error(error);
                toast.error($_('modals.translation.toasts.error'))
            });
        }
    }
</script>

<dialog id="modalTranslation" class="modal">
    <form method="dialog" class="flex flex-col items-center justify-center modal-box">
        <button class="absolute btn btn-sm btn-circle btn-ghost right-2 top-2">✕</button>
        {#if tr}
            <h1 class="pb-2 mt-2 font-bold text-center">
                {$_('modals.translation.title')}
            </h1>
            <div>
                <!-- Language picker -->
                <div class="w-full max-w-xs form-control">
                    <label for="target-lan" class="label">
                      <span class="label-text">{$_('modals.translation.targetLanguages', { values: { language: tr.result.language } })}</span>
                    </label>
                    <select bind:value={targetLanguage} name="target-lan" class="select select-bordered">
                      <option disabled selected>{$_('modals.translation.pickOne')}</option>
                      <!-- Iterate all available languages -->
                      {#each availableLanguages as lan}
                        <!-- When we find the source language -->
                        {#if lan.code == tr.result.language}
                            <!-- Iterate all possible target languages -->
                            {#each lan.targets as t}
                                {#if t != tr.result.language}
                                    <option value="{t}">{t}</option>
                                {/if}
                            {/each}
                        {/if}
                      {/each}
                    </select>
                </div>
                <!-- End language picker-->
                <button on:click={handleTranslate(tr.id)} class="mt-5 btn btn-active btn-primary">{$_('modals.translation.translate')}</button>
            </div>
        {/if}
    </form>
    <form method="dialog" class="modal-backdrop">
        <button>{$_('common.close')}</button>
    </form>
</dialog>