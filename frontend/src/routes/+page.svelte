<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import { GetFileContentPreview, ListFiles } from '$lib/wailsjs/go/main/App';
	import type { main } from '$lib/wailsjs/go/models';
	import MoreHorizontalIcon from '@lucide/svelte/icons/more-horizontal';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { isMarkdownFile } from '$lib/utils';
	import Folder from '@lucide/svelte/icons/folder';

	let hoveredCard = $state<number | null>(null);
	let files: main.FileEntry[] = $state([]);

	$effect(() => {
		if (files) {
			files.forEach((file) => {
				if (file.isDirectory) return;
				if (!isMarkdownFile(file.name)) return;
				console.log('Fetching preview for file:', file.path);
				if (!fileContentPreview[file.path]) {
					GetFileContentPreview(file.path).then((res) => {
						console.log('Received preview for file:', file.path, res);
						fileContentPreview[file.path] = res;
					});
				}
			});
		}
	});
	let fileContentPreview: {
		[fileID: string]: string;
	} = $state({});

	onMount(() => {
		ListFiles('').then((res) => {
			files = res;
		});
	});
</script>

<div class="relative flex h-full flex-col gap-8 p-2 md:px-10 xl:px-14">
	<!-- <InputGroup.Root class="border-0 bg-transparent focus-visible:ring-0 dark:bg-transparent">
		<InputGroup.Input placeholder="Search..." />
		<InputGroup.Addon>
			<SearchIcon />
		</InputGroup.Addon>
	</InputGroup.Root> -->

	<div class="flex w-full flex-1 flex-col gap-12 overflow-y-auto pb-32 pt-12">
		<div class="flex flex-col">
			<h1 class="text-6xl font-extrabold">Foo</h1>
			<p class="text-muted-foreground">/Foo/bar/baz/...</p>
		</div>

		<div class="flex w-full justify-center">
			<div class="max-w-8xl flex flex-1 flex-wrap items-stretch">
				{#each files as file, i}
					<div
						class="w-full p-2 md:w-1/2 lg:w-1/3 xl:w-1/4 flex"
						onmouseenter={() => (hoveredCard = i)}
						onmouseleave={() => (hoveredCard = null)}
					>
						<Card.Root class="w-full flex flex-col hover:shadow-md transition-shadow duration-200">
							<Card.Header>
								<Card.Title class="truncate">
									{#if file.isDirectory}
										<Folder class="inline-block mr-2 size-5 text-muted-foreground" />
									{/if}
									{file.name}
								</Card.Title>
								<!-- <Card.Description>{file.description}</Card.Description> -->
							</Card.Header>
							<Card.Content class="min-w-0 flex-1">
								{#if file.isDirectory || !isMarkdownFile(file.name)}
									
								{:else if fileContentPreview[file.path]}
									<p class="text-muted-foreground text-sm whitespace-pre-wrap truncate overflow-hidden">
										{fileContentPreview[file.path]}
									</p>
								{:else}
									<div class="flex flex-col gap-1">
										<Skeleton class="h-[20px] max-w-[128px] rounded-full" />
										<Skeleton class="h-[20px] max-w-[256px] rounded-full" />
										<Skeleton class="h-[20px] max-w-[256px] rounded-full" />
										<Skeleton class="h-[20px] max-w-[256px] rounded-full" />
									</div>
								{/if}
							</Card.Content>
							<Card.Footer class="flex justify-between">
								<p class="text-muted-foreground text-sm">
									{new Date(file.modTime).toLocaleString()}
								</p>

								{#if hoveredCard === i}
									<div transition:fade={{ duration: 200 }}>
										<Button size="icon" class="size-8 rounded-full">
											<Trash2 />
										</Button>
									</div>
								{:else}
									<div class="invisible">
										<Button size="icon" class="size-8 rounded-full">
											<Trash2 />
										</Button>
									</div>
								{/if}
							</Card.Footer>
						</Card.Root>
					</div>
				{/each}
			</div>
		</div>

		<div class="absolute inset-x-2 bottom-10 md:px-8">
			<InputGroup.Root class="!h-auto rounded-xl backdrop-blur  ">
				<InputGroup.Input class="p-8 !text-xl" placeholder="Search in folder..." />
				<InputGroup.Addon align="inline-end">
					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}
								<InputGroup.Button {...props} variant="ghost" aria-label="More" size="icon-xs">
									<MoreHorizontalIcon />
								</InputGroup.Button>
							{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end">
							<DropdownMenu.Item>Settings</DropdownMenu.Item>
							<DropdownMenu.Item>Copy path</DropdownMenu.Item>
							<DropdownMenu.Item>Open location</DropdownMenu.Item>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</InputGroup.Addon>
			</InputGroup.Root>
		</div>
	</div>
</div>
