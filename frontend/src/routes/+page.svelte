<script lang="ts">
	import { Crepe } from '@milkdown/crepe';
	import '@milkdown/crepe/theme/common/style.css';
	import '@milkdown/crepe/theme/frame-dark.css';
	import { onMount, onDestroy } from 'svelte';
	import { appState } from '../store.svelte';
	import {
		GetContentHash,
		GetCurrentFilesState,
		GetFileContent,
		SaveFile
	} from '$lib/wailsjs/go/main/App';
	import { editorViewCtx, parserCtx } from '@milkdown/core';
	import { Slice } from '@milkdown/prose/model';
	import { fileToBase64 } from '$lib/utils';
	import { main } from '$lib/wailsjs/go/models';

	let crepe: Crepe | null = $state(null);
	let editorReady = $state(false);
	let readingNewFile = $state(false);
	let saving = $state(false);

	$inspect(appState);

	$effect(() => {
		if (appState.currentFile && editorReady) {
			// New file should not be selected unless we're done with the content of the
			// current markdown content.
			readingNewFile = true;
			GetFileContent(appState.currentFile)
				.then((data) => {
					setEditorContent(data);
					appState.staleContent = false;
				})
				.catch((error) => {
					console.error('Error loading file:', error);
				})
				.finally(() => {
					readingNewFile = false;
				});
		}
	});

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
				view.dispatch(state.tr.replace(0, state.doc.content.size, new Slice(doc.content, 0, 0)));
			});
		} catch (error) {
			console.error('Error setting editor content:', error);
		}
	}

	function saveContentToFile() {
		if (!appState.currentFile || !crepe || !editorReady) return;

		const currentContent = crepe?.getMarkdown();

		saving = true;

		SaveFile(appState.currentFile, currentContent)
			.then(async () => {
				// Update content hash after saving
				const newState = await GetCurrentFilesState();
				appState.currentDir = newState.currentDir;
				appState.currentFile = newState.currentFile;
				appState.contentHash = newState.contentHash;
			})
			.catch((error) => {
				console.error('Error saving file:', error);
			})
			.finally(() => {
				saving = false;
			});
	}

	onMount(() => {
		// Create editor instance with base64 uploader
		crepe = new Crepe({
			root: document.getElementById('editor'),
			defaultValue: '# Hello, Crepe!\n\nStart writing your markdown...',
			featureConfigs: {
				'image-block': {
					async blockOnUpload(file: File) {
						console.log('Converting image to base64...');

						try {
							// Convert file to base64 data URI
							const base64DataUri = await fileToBase64(file);
							console.log('Image converted to base64');
							return base64DataUri;
						} catch (error) {
							console.error('Error converting image to base64:', error);
							return '';
						}
					}
				}
			}
		});

		// Register content change listener
		crepe.on((listener) => {
			listener.markdownUpdated((ctx, markdown, prevMarkdown) => {
				if (markdown !== prevMarkdown) {
					GetContentHash(markdown).then((hash) => {
						console.log('Content hash updated:', hash);
						if (appState.contentHash !== hash) {
							appState.staleContent = true;
						} else {
							appState.staleContent = false;
						}
					});
				}
			});
		});

		crepe.create().then(() => {
			editorReady = true;
		});
	});

	onDestroy(() => {
		// Clean up when component is destroyed
		editorReady = false;
		crepe?.destroy();
	});
</script>

<div class="h-full">
	<div id="editor" class="h-full rounded-xl"></div>
</div>

<style>
	:global(.milkdown) {
		border-radius: 12px;
		overflow-y: hidden;
		height: 100% !important;
		margin: 0px !important;
		width: 100% !important;
		max-width: 100% !important;
		max-height: 100% !important;
		border: 0px !important;
		overflow-y: auto !important;
	}
</style>
