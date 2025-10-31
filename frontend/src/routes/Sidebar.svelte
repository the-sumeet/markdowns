<script>
	import Button from '$lib/components/ui/button/button.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import Plus from '@lucide/svelte/icons/plus';
	import Moon from '@lucide/svelte/icons/moon';
	import Sun from '@lucide/svelte/icons/sun';
</script>

<div
	style="--wails-draggable:drag"
	class="flex h-full flex-col items-center gap-8 overflow-y-auto border-r p-6 pt-20"
>
	<!-- <p class="font-yellowtail">Markdowns</p> -->

	<Button  size="icon" class="size-8 rounded-full">
		<Plus />
	</Button>
</div>
