<script>
	import { editorSettings } from '$lib/stores';
	import { _ } from 'svelte-i18n';
	import Checkbox from './ui/Checkbox.svelte';
</script>

<div class="flex items-center gap-4 px-3 py-1.5 rounded-lg border border-border bg-muted/50">
	<label
		class="flex items-center gap-2 cursor-pointer select-none"
		title={$_('editor.settings.seekOnClickTooltip')}
	>
		<Checkbox bind:checked={$editorSettings.seekOnClick} />
		<span class="text-muted-foreground text-[0.78rem]">{$_('editor.settings.seekOnClick')}</span>
	</label>

	<label
		class="flex items-center gap-2 cursor-pointer select-none"
		title={$_('editor.settings.autosaveTooltip')}
	>
		<Checkbox bind:checked={$editorSettings.autoSave} />
		<span class="text-muted-foreground text-[0.78rem]">{$_('editor.settings.autosave')}</span>
	</label>
</div>
