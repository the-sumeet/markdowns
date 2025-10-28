<script>
	import Button from '$lib/components/ui/button/button.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import Plus from '@lucide/svelte/icons/plus';
</script>

<div style="--wails-draggable:drag" class="flex h-full flex-col items-center border-r pt-20 gap-8 overflow-y-auto p-6">
	<!-- <p class="font-yellowtail">Markdowns</p> -->

	<Button size="icon" class="size-8 rounded-full">
		<Plus  />
	</Button>
</div>
