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
	import { ClearCurrentFile } from '$lib/wailsjs/go/main/App';
	import { main } from '$lib/wailsjs/go/models';

	let crepe: Crepe | null = $state(null);
	let editorReady = $state(false);
	let readingNewFile = $state(false);
	let saving = $state(false);

	console.log('Current file in note page:', appState.currentFile);

	$effect(() => {
		if (!appState.currentFile) {
			goto('/');
		}
	});

	onMount(() => {
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
		});
	});

	onDestroy(() => {
		// Clean up when component is destroyed
		editorReady = false;
		crepe?.destroy();
		appState.saveFile = undefined;
	});
</script>

<div class="flex h-full w-full flex-col overflow-hidden p-4">
	<!-- Topbar -->
	<div
		class="mb-4 flex flex-shrink-0 justify-between rounded-xl bg-white/30 p-4 backdrop-blur-md dark:bg-black/30"
	>
		<div class="flex-1">
            <h1 class=" text-3xl font-extrabold">{appState.currentFile?.name}</h1>
        </div>
		<div class="flex">
			<Button
            variant="outline"
				size="icon"
				class="size-8 rounded-full"
				onclick={() => {
					ClearCurrentFile().then((res: main.CurrentFilesState) => {
						updateAppState(appState, res);
					});
				}}
			>
				<CircleX />
			</Button>
		</div>
	</div>

	<div id="editor" class="min-h-0 flex-1 w-full rounded-xl overflow-hidden"></div>
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
		backdrop-filter: blur(12px) !important;
		background: rgba(255, 255, 255, 0.3) !important;
	}

	:global(.dark .milkdown) {
		background: rgba(0, 0, 0, 0.3) !important;
	}
</style>
