<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
	import FolderIcon from '@lucide/svelte/icons/folder';
	import ShareIcon from '@lucide/svelte/icons/share';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import Folder from '@lucide/svelte/icons/folder';
	import File from '@lucide/svelte/icons/file';
	import SquareArrowDown from '@lucide/svelte/icons/square-arrow-down';
	import { OpenFile, ListFiles, DeleteFile, RenameFile, UpdateWindowTitleWithCurrentDir } from '$lib/wailsjs/go/main/App';
	import { onMount } from 'svelte';
	import type { main } from '$lib/wailsjs/go/models';
	import { appState } from '../../store.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';

	let files: main.FileEntry[] = $state([]);

	// Dialog state
	let showRenameDialog = $state(false);
	let showDeleteDialog = $state(false);
	let selectedFile: main.FileEntry | null = $state(null);
	let newName = $state('');
	let errorMessage = $state('');

	// Handle rename
	async function handleRename() {
		errorMessage = '';
		if (!selectedFile) return;

		try {
			await RenameFile(selectedFile.path, newName);
			showRenameDialog = false;
			newName = '';
			selectedFile = null;
			// Trigger refresh
			appState.refreshTrigger = (appState.refreshTrigger || 0) + 1;
		} catch (err) {
			errorMessage = String(err);
		}
	}

	// Handle delete
	async function handleDelete() {
		errorMessage = '';
		if (!selectedFile) return;

		try {
			await DeleteFile(selectedFile.path);
			showDeleteDialog = false;
			selectedFile = null;
			// Trigger refresh
			appState.refreshTrigger = (appState.refreshTrigger || 0) + 1;
		} catch (err) {
			errorMessage = String(err);
		}
	}

	// Open rename dialog
	function openRenameDialog(file: main.FileEntry) {
		selectedFile = file;
		newName = file.name;
		errorMessage = '';
		showRenameDialog = true;
	}

	// Open delete dialog
	function openDeleteDialog(file: main.FileEntry) {
		selectedFile = file;
		errorMessage = '';
		showDeleteDialog = true;
	}

	$effect(() => {
		// Watch both currentDir and refreshTrigger to reload files
		appState.refreshTrigger;
		if (appState.currentDir) {
			ListFiles(appState.currentDir).then((res: main.FileEntry[]) => {
				files = res;
			});
		}
	});

	const sidebar = useSidebar();

	onMount(() => {});
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<!-- <Sidebar.GroupLabel>Projects</Sidebar.GroupLabel> -->
	<Sidebar.Menu>
		{#if appState.currentDir !== '/'}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton
					onclick={() => {
						if (appState.currentDir) {
							const parentDir = appState.currentDir.split('/').slice(0, -1).join('/') || '/';
							ListFiles(parentDir).then((res: main.FileEntry[]) => {
								files = res;
								appState.currentDir = parentDir;
								// Update window title
								UpdateWindowTitleWithCurrentDir();
							});
						}
					}}
				>
					<span>..</span>
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		{/if}
		{#each files as file (file.path)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton
					disabled={!file.isDirectory &&
						!file.name.endsWith('.md') &&
						!file.name.endsWith('.markdown')}
					onclick={() => {
						OpenFile(file.path).then((res) => {
							appState.currentDir = res.currentDir;
							appState.currentFile = res.currentFile;
						});
					}}
				>
					<!-- <item.icon /> -->
					{#if file.isDirectory}
						<Folder />
					{:else if file.name.endsWith('.md') || file.name.endsWith('.markdown')}
						<SquareArrowDown />
					{:else}
						<File />
					{/if}
					<span>{file.name}</span>
				</Sidebar.MenuButton>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuAction showOnHover {...props}>
								<EllipsisIcon />
								<span class="sr-only">More</span>
							</Sidebar.MenuAction>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content
						class="w-48"
						side={sidebar.isMobile ? 'bottom' : 'right'}
						align={sidebar.isMobile ? 'end' : 'start'}
					>
						<DropdownMenu.Item
							onclick={(e) => {
								e.preventDefault();
								openRenameDialog(file);
							}}
						>
							<ShareIcon class="text-muted-foreground" />
							<span>Rename</span>
						</DropdownMenu.Item>
						<!-- <DropdownMenu.Separator /> -->
						<DropdownMenu.Item
							onclick={(e) => {
								e.preventDefault();
								openDeleteDialog(file);
							}}
						>
							<Trash2Icon class="text-muted-foreground" />
							<span>Delete</span>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>

<!-- Rename Dialog -->
<Dialog.Root bind:open={showRenameDialog}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Rename {selectedFile?.isDirectory ? 'Directory' : 'File'}</Dialog.Title>
			<Dialog.Description>
				Enter a new name for "{selectedFile?.name}".
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="newName">New name</Label>
				<Input
					id="newName"
					bind:value={newName}
					placeholder="Enter new name"
					onkeydown={(e) => {
						if (e.key === 'Enter') {
							handleRename();
						}
					}}
				/>
			</div>
			{#if errorMessage}
				<p class="text-sm text-red-500">{errorMessage}</p>
			{/if}
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showRenameDialog = false)}>Cancel</Button>
			<Button onclick={handleRename} disabled={!newName.trim()}>Rename</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Delete Confirmation Dialog -->
<Dialog.Root bind:open={showDeleteDialog}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Delete {selectedFile?.isDirectory ? 'Directory' : 'File'}</Dialog.Title>
			<Dialog.Description>
				Are you sure you want to delete "{selectedFile?.name}"?
				{#if selectedFile?.isDirectory}
					<span class="block mt-2 text-red-500 font-medium">
						Warning: This will delete the directory and all its contents.
					</span>
				{/if}
				This action cannot be undone.
			</Dialog.Description>
		</Dialog.Header>
		{#if errorMessage}
			<p class="text-sm text-red-500">{errorMessage}</p>
		{/if}
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showDeleteDialog = false)}>Cancel</Button>
			<Button variant="destructive" onclick={handleDelete}>Delete</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
