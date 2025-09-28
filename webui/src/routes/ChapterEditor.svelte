<script lang="ts">
	import { Select, Input } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { getChapters } from '../oo';

	let { params = {} } = $props();
	let selected = $state(0);
	let chaps = $state([]);
	let name = $state('');
	let prior = $state(0);

	if (params.id) {
		selected = params.id;
		const e = new Event('search', { bubbles: true, cancelable: true });
		chooseChapter(e);
	}

	async function chooseChapter(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			const selectedChap = chaps.find((c) => {
				return c.ID == selected;
			});
			if (typeof selectedChap == 'undefined') {
				name = 'Unknown';
				prior = 0;
			} else {
				name = selectedChap.Name;
				prior = selectedChap.Prior;
			}
			await push(`#/chaptereditor/${selected}`);
		} catch (err) {
			console.log(err);
			toast.push(err.message);
		}
	}

	async function load() {
		chaps = await getChapters();
		return chaps;
	}

	async function updateChapter(event) {
		console.log(event);
	}
</script>

{#await load()}
	<h3>... loading ...</h3>
{:then}
	<div class="grid grid-cols-8 gap-4 px-4 py-2">
		<div class="col-span-8">Chapter Editor</div>
		<div class="col-span-8">
			<Select class="mt-2" items={chaps} bind:value={selected} onchange={chooseChapter} />
		</div>
	</div>
	{#if selected != 0}
		<div class="grid grid-cols-8 gap-4 px-4 py-2">
			<div class="col-span-2">Chapter Name</div>
			<div class="col-span-6">
				<Input type="text" name="name" bind:value={name} onchange={updateChapter} />
			</div>
			<div class="col-span-2">Chapter Prior (ID number for now)</div>
			<div class="col-span-6">
				<Input type="text" name="prior" bind:value={prior} onchange={updateChapter} />
			</div>
		</div>
	{/if}
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}
