<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { appState } from '../../store.svelte';
	import { goto } from '$app/navigation';
	import { Crepe } from '@milkdown/crepe';
	import { fileToBase64, updateAppState } from '$lib/utils';
	import '@milkdown/crepe/theme/common/style.css';
	import '@milkdown/crepe/theme/frame-dark.css';
	import { Button } from '$lib/components/ui/button';
	import CircleX from '@lucide/svelte/icons/circle-x';
	import {
		ClearCurrentFile,
		GetContentHash,
		GetFileContent,
		SaveFile,
		PickImageFile
	} from '$lib/wailsjs/go/main/App';
	import { main } from '$lib/wailsjs/go/models';
	import Save from '@lucide/svelte/icons/save';
	import { editorViewCtx, parserCtx } from '@milkdown/core';
	import { Slice } from '@milkdown/prose/model';
	import * as Dialog from '$lib/components/ui/dialog/index.js';

	let crepe: Crepe | null = $state(null);
	let editorReady = $state(false);
	let readingNewFile = $state(false);
	let saving = $state(false);
	let contentHash = $state('');
	let isDirty = $state(false);
	let showUnsavedDialog = $state(false);

	$effect(() => {
		if (!appState.currentFile) {
			goto('/');
		}

		readingNewFile = true;

		GetFileContent(appState.currentFile!.path)
			.then(async (data) => {
				contentHash = await GetContentHash(data);
				isDirty = false;
				console.log('Loaded content hash:', contentHash, data);
				setEditorContent(data);
			})
			.catch((error) => {
				console.error('Error loading file:', error);
			})
			.finally(() => {
				readingNewFile = false;
			});
	});

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

	async function saveFile() {
		if (!appState.currentFile || !crepe || !editorReady) return;

		saving = true;
		const currentContent = crepe.getMarkdown();

		try {
			await SaveFile(appState.currentFile.path, currentContent);
			// Update hash with saved content
			contentHash = await GetContentHash(currentContent);
			isDirty = false;
			console.log('File saved successfully');
		} catch (error) {
			console.error('Error saving file:', error);
		} finally {
			saving = false;
		}
	}

	function handleClose() {
		if (isDirty) {
			showUnsavedDialog = true;
		} else {
			closeFile();
		}
	}

	function closeFile() {
		ClearCurrentFile().then((res: main.CurrentFilesState) => {
			updateAppState(appState, res);
		});
	}

	async function saveAndClose() {
		await saveFile();
		closeFile();
	}

	function discardAndClose() {
		showUnsavedDialog = false;
		closeFile();
	}

	onMount(() => {
		crepe = new Crepe({
			root: document.getElementById('editor'),
			defaultValue: '# Hello, Crepe!\n\nStart writing your markdown...',
			featureConfigs: {
				'image-block': {
					async blockOnUpload(file: File) {
						console.log('Image uploaded:', file);

						try {
							// Use Wails Go function to open file dialog
							const filePath = await PickImageFile();
							console.log('Selected image path:', filePath);
							return filePath || '';
						} catch (error) {
							console.error('Error selecting image:', error);
							return '';
						}
					}
				}
			}
		});

		crepe.on((listener) => {
			listener.markdownUpdated((ctx, markdown, prevMarkdown) => {
				GetContentHash(markdown).then((hash) => {
					if (hash !== contentHash && !readingNewFile) {
						isDirty = true;
						console.log('dirty', hash, markdown);
					} else {
						isDirty = false;
					}
				});
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

<div class="flex h-full w-full flex-col overflow-hidden p-4">
	<!-- Topbar -->
	<div
		class="mb-4 flex shrink-0 justify-between rounded-xl bg-white/30 p-4 backdrop-blur-md dark:bg-black/30"
	>
		<div class="flex-1">
			<h1 class=" text-3xl font-extrabold">{appState.currentFile?.name}</h1>
		</div>
		<div class="flex items-center gap-1">
			{#if isDirty}
				<Button variant="outline" class="rounded-full" onclick={saveFile} disabled={saving}>
					<Save />
					{saving ? 'Saving...' : 'Save'}
				</Button>
			{/if}
			<Button
				variant="outline"
				size="icon"
				class="size-8 rounded-full"
				onclick={handleClose}
			>
				<CircleX />
			</Button>
		</div>
	</div>

	<div id="editor" class="min-h-0 w-full flex-1 overflow-hidden rounded-xl"></div>
</div>

<!-- Close without save confirmation -->
<Dialog.Root bind:open={showUnsavedDialog}>
	<Dialog.Content class="bg-card">
		<Dialog.Header>
			<Dialog.Title>Your changes will be lost?</Dialog.Title>
			<Dialog.Description>You've made changes to this file, do you want to discard them?</Dialog.Description>
		</Dialog.Header>
		<div class="flex w-full gap-2 justify-end">
			<Button variant="outline" onclick={() => (showUnsavedDialog = false)}>Cancel</Button>
			<Button onclick={saveAndClose}>Save and Close</Button>
			<Button variant="destructive" onclick={discardAndClose}>Discard & Close</Button>
		</div>
	</Dialog.Content>
</Dialog.Root>

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
		backdrop-filter: blur(12px) !important;
		background: rgba(255, 255, 255, 0.3) !important;
	}

	:global(.dark .milkdown) {
		background: rgba(0, 0, 0, 0.3) !important;
	}
</style>
