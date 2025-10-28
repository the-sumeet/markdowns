<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import MoreHorizontalIcon from '@lucide/svelte/icons/more-horizontal';
</script>

<div class="relative flex h-full flex-col gap-8 p-2 md:px-10 xl:px-14">
	<!-- <InputGroup.Root class="border-0 bg-transparent focus-visible:ring-0 dark:bg-transparent">
		<InputGroup.Input placeholder="Search..." />
		<InputGroup.Addon>
			<SearchIcon />
		</InputGroup.Addon>
	</InputGroup.Root> -->

	<div class="flex flex-1 flex-col gap-12 overflow-y-auto pt-12 w-full">
		<div class="flex flex-col">
			<h1 class="text-6xl font-extrabold">Foo</h1>
			<p class="text-muted-foreground">/Foo/bar/baz/...</p>
		</div>

		<div class="flex w-full justify-center">
			<div class="max-w-8xl flex flex-1 flex-wrap">
				{#each Array(100) as i}
					<div class="xl:1/4 w-full p-2 md:w-1/2 lg:w-1/3 xl:w-1/4">
						<Card.Root class="">
							<Card.Header>
								<Card.Title>Card Title</Card.Title>
								<Card.Description>Card Description</Card.Description>
							</Card.Header>
							<Card.Content>
								<p>Card Content</p>
							</Card.Content>
							<Card.Footer>
								<p>Card Footer</p>
							</Card.Footer>
						</Card.Root>
					</div>
				{/each}
			</div>
		</div>

		<InputGroup.Root class="absolute bottom-4 inset-x-2 md:inset-x-10 xl:inset-x-14 max-w-8xl">
			<InputGroup.Input placeholder="Enter file name" />
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
