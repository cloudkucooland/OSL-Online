<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { Select, Input, Label, Card, Heading, Button, Spinner, Hr } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { getChapters } from '../oo';

	const oo = getContext('oo');

	// Guard: Ensure admin access
	if (!oo.me || oo.me.level < 3) {
		push('/Login');
	}

	let { params = {} } = $props();
	let selected = $state(0);
	let chaps = $state([]);
	let name = $state('');
	let prior = $state(0);
	let loading = $state(true);

	onMount(async () => {
		try {
			if (oo.chaptercache) {
				chaps = oo.chaptercache;
			} else {
				chaps = await getChapters();
			}
			if (params.id) {
				selected = parseInt(params.id);
				const c = chaps.find((i) => i.ID == selected);
				if (c) {
					name = c.Name;
					prior = c.Prior;
				}
			}
		} catch (err) {
			toast.push('Error loading chapters');
		} finally {
			loading = false;
		}
	});

	async function load() {
		loading = true;
		if (oo.chaptercache) {
			chaps = oo.chaptercache;
		} else {
			chaps = await getChapters();
		}
		return chaps;
	}

	async function chooseChapter() {
		const selectedChap = chaps.find((c) => c.ID == selected);
		if (selectedChap) {
			name = selectedChap.Name;
			prior = selectedChap.Prior;
			push(`/chaptereditor/${selected}`);
		} else {
			name = '';
			prior = 0;
		}
	}

	async function updateChapter(event) {
		console.log(event);
		// Placeholder for your actual update API call
		console.log('Saving Chapter:', { id: selected, name, prior });
		toast.push(`Saving changes for ${name}...`);
		// await saveChapterFromServer({ ID: selected, Name: name, Prior: prior });
	}
</script>

<svelte:head>
	<title>OSL Directory: Chapter Editor</title>
</svelte:head>

<div class="mx-auto max-w-4xl space-y-6 p-6">
	<Card size="none" class="border-slate-200 p-6 shadow-sm">
		<Heading tag="h2" class="mb-6 text-2xl font-bold text-slate-900">Chapter Editor</Heading>

		<div class="space-y-4">
			<div>
				<Label for="chapter-select" class="mb-2">Select Chapter</Label>
				<Select
					id="chapter-select"
					items={chaps.map((c) => ({ value: c.ID, name: c.Name }))}
					bind:value={selected}
					onchange={chooseChapter}
					placeholder="Choose a chapter to edit..."
				/>
			</div>
		</div>

		{#if loading}
			<div class="flex justify-center p-12">
				<Spinner color="purple" size="12" />
			</div>
		{:else if selected > 0}
			<Hr class="my-8" />

			<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
				<div>
					<Label class="mb-2">Chapter Name</Label>
					<Input bind:value={name} placeholder="e.g. St. Luke Chapter" />
				</div>
				<div>
					<Label class="mb-2">Prior / Leader ID</Label>
					<Input type="number" bind:value={prior} placeholder="Member ID" />
				</div>
			</div>

			<div class="mt-8 flex justify-end">
				<Button color="primary" onclick={updateChapter}>Save Chapter Changes</Button>
			</div>
		{/if}
	</Card>
</div>
