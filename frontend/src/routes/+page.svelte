<script lang="ts">
	import { Crepe } from '@milkdown/crepe';
	import '@milkdown/crepe/theme/common/style.css';
	import '@milkdown/crepe/theme/frame-dark.css';
	import { onMount, onDestroy } from 'svelte';
	import { appState } from '../store.svelte';
	import { GetFileContent } from '$lib/wailsjs/go/main/App';
	import { editorViewCtx, parserCtx } from '@milkdown/core';
	import { Slice } from '@milkdown/prose/model';

	let content: string = $state('');
	let crepe: Crepe | null = $state(null);
	let editorReady = $state(false);

	// Function to set markdown content in the editor
	function setEditorContent(markdown: string) {
		if (!crepe || !editorReady) return;

		try {
			crepe.editor.action((ctx) => {
				const view = ctx.get(editorViewCtx);
				const parser = ctx.get(parserCtx);
				const doc = parser(markdown);

				if (!doc) return;

				const state = view.state;
				view.dispatch(
					state.tr.replace(0, state.doc.content.size, new Slice(doc.content, 0, 0))
				);
			});
		} catch (error) {
			console.error('Error setting editor content:', error);
		}
	}

	// Watch for file changes and load content
	$effect(() => {
		if (appState.currentFile) {
			GetFileContent(appState.currentFile).then((data) => {
				content = data;
				console.log('File content loaded:', content);
				// Update editor content if editor is ready
				setEditorContent(content);
			});
		}
	});

	onMount(() => {
		// Create editor instance
		crepe = new Crepe({
			root: document.getElementById('editor'),
			defaultValue: '# Hello, Crepe!\n\nStart writing your markdown...'
		});

		crepe.create().then(() => {
			editorReady = true;
			// Set initial content if available
			if (content) {
				setEditorContent(content);
			}
		});
	});

	onDestroy(() => {
		// Clean up when component is destroyed
		editorReady = false;
		crepe?.destroy();
	});
</script>

<div class="h-full ">
  <div id="editor" class="h-full rounded-xl"></div>
</div>

<style>
	:global(.milkdown) {
    border-radius: 12px;
    overflow-y: hidden;
		height: 100% !important;
		margin: 0px !important;
		width: 100% !important;
    max-width: 100%!important;
    max-height: 100%!important;
    border: 0px!important;
    overflow-y: auto !important;
	}
</style>
