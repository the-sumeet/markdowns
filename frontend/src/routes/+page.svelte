<script lang="ts">
	import { Crepe } from '@milkdown/crepe';
	import '@milkdown/crepe/theme/common/style.css';
	import '@milkdown/crepe/theme/frame-dark.css';
	import { onMount, onDestroy, untrack } from 'svelte';
	import { appState } from '../store.svelte';
	import { GetFileContent, SaveFile } from '$lib/wailsjs/go/main/App';
	import { editorViewCtx, parserCtx } from '@milkdown/core';
	import { Slice } from '@milkdown/prose/model';

	let crepe: Crepe | null = $state(null);
	let editorReady = $state(false);
	let isDirty = $state(false);
	let currentFileBeingEdited: string | null = $state(null);
	let saveTimeout: NodeJS.Timeout | null = null;
	let lastSavedContent: string = '';
	let changeCheckInterval: NodeJS.Timeout | null = null;
	let previousFile: string | undefined = undefined;
	let isSaving = $state(false);
	let isLoading = $state(false);

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

	// Function to save file
	async function saveCurrentFile() {
		// Prevent concurrent saves
		if (isSaving || !currentFileBeingEdited || !isDirty) return;

		isSaving = true;
		try {
			const currentContent = crepe?.getMarkdown();

			// CRITICAL: Never save if content is undefined or editor not ready
			if (currentContent === undefined || !editorReady) {
				console.error('Cannot save: editor not ready or content undefined');
				return;
			}

			// CRITICAL: Validate content before saving
			// Only allow empty save if lastSavedContent was also empty (new file scenario)
			if (currentContent.trim() === '' && lastSavedContent.trim() !== '') {
				console.error('PREVENTED DATA LOSS: Attempted to save empty content when file had data');
				return;
			}

			await SaveFile(currentFileBeingEdited, currentContent);
			console.log('File saved:', currentFileBeingEdited);
			lastSavedContent = currentContent;
			isDirty = false;
		} catch (error) {
			console.error('Error saving file:', error);
		} finally {
			isSaving = false;
		}
	}

	// Check for content changes periodically
	function checkForChanges() {
		if (!crepe || !editorReady || !currentFileBeingEdited) return;

		const currentContent = crepe.getMarkdown();
		// Only trigger if content changed AND not already dirty (avoid repeated triggers)
		if (currentContent !== lastSavedContent && !isDirty) {
			console.log('Content changed, marking as dirty');
			isDirty = true;
			scheduleAutoSave();
		}
	}

	// Debounced auto-save (3 seconds)
	function scheduleAutoSave() {
		// Clear existing timeout
		if (saveTimeout) {
			clearTimeout(saveTimeout);
		}

		// Schedule new save after 3 seconds
		saveTimeout = setTimeout(() => {
			saveCurrentFile();
		}, 3000);
	}

	// Watch for file changes and load content
	$effect(() => {
		const newFile = appState.currentFile;

		// Only load if the file actually changed and not currently loading
		if (newFile && newFile !== previousFile && !isLoading) {
			// Update previousFile immediately to prevent race conditions
			previousFile = newFile;

			// Use untrack to prevent isDirty from becoming a dependency
			untrack(() => {
				// Save current file before switching if it's dirty
				if (isDirty && currentFileBeingEdited && currentFileBeingEdited !== newFile) {
					saveCurrentFile().then(() => {
						loadNewFile(newFile);
					});
				} else {
					loadNewFile(newFile);
				}
			});
		}
	});

	// Load new file content
	async function loadNewFile(filePath: string | undefined) {
		console.log('Loading file:', filePath);
		if (!filePath || isLoading) return;

		isLoading = true;
		try {
			const data = await GetFileContent(filePath);

			// CRITICAL: Only update state AFTER successfully loading AND setting content
			if (!editorReady) {
				console.error('Cannot load file: editor not ready');
				return;
			}

			// Set editor content first
			setEditorContent(data);

			// Only update state after successful content load
			lastSavedContent = data;
			currentFileBeingEdited = filePath;
			isDirty = false;
			console.log('File content loaded successfully');
		} catch (error) {
			console.error('Error loading file:', error);
		} finally {
			isLoading = false;
		}
	}

	// Convert File to base64 data URI
	function fileToBase64(file: File): Promise<string> {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.readAsDataURL(file);
			reader.onload = () => resolve(reader.result?.toString() || '');
			reader.onerror = (error) => reject(error);
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

		crepe.create().then(() => {
			editorReady = true;

			// Start polling for changes every 500ms
			changeCheckInterval = setInterval(checkForChanges, 500);
		});
	});

	onDestroy(() => {
		// Save before destroying if dirty
		if (isDirty && currentFileBeingEdited) {
			saveCurrentFile();
		}

		// Clear intervals and timeouts
		if (saveTimeout) {
			clearTimeout(saveTimeout);
		}
		if (changeCheckInterval) {
			clearInterval(changeCheckInterval);
		}

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
