<script>
    import { editorSettings } from '$lib/stores';
    import { _ } from 'svelte-i18n';
</script>

<div class="flex flex-row items-center justify-center p-3 space-x-4 m-2">
    <div
        class="form-control tooltip"
        data-tip={$_('editor.settings.seekOnClickTooltip')}
    >
        <label class="label cursor-pointer space-x-2">
            <span class="label-text">{$_('editor.settings.seekOnClick')}</span>
            <input bind:checked={$editorSettings.seekOnClick} type="checkbox" class="checkbox checkbox-sm" />
        </label>
    </div>

    <div class="form-control tooltip" data-tip={$_('editor.settings.autosaveTooltip')}>
        <label class="label cursor-pointer space-x-2">
            <span class="label-text">{$_('editor.settings.autosave')}</span>
            <input bind:checked={$editorSettings.autoSave} type="checkbox" class="checkbox checkbox-sm" />
        </label>
    </div>
</div>