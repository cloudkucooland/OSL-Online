<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { Select, Label, Card, Badge } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getChapters, getChapterMembers } from '../oo';

	let { params = {} } = $props();
	let chaps = $state([]);
	let members = $state([]);
	let selected = $state();
	let prior = $state();

	const oo = getContext('oo');

	async function updateData() {
		try {
			members = await getChapterMembers(selected);
		} catch (err: any) {
			toast.push(err.message);
		}
		const chap = oo.chaptercache.find((c) => c.ID == selected);
		prior = chap.Prior;
	}

	async function handleSelectChange() {
		if (selected) {
			push(`#/chapterbrowser/${selected}`);
			await updateData();
		}
	}

	async function load() {
		if (oo.chaptercache) return oo.chaptercache;
		return await getChapters();
	}

	onMount(async () => {
		if (!oo.chaptercache) {
			oo.chaptercache = await getChapters();
		}
		chaps = oo.chaptercache;

		if (params.id) {
			selected = params.id;
			await updateData();
		}
	});
</script>

{#await load()}
	<div class="flex justify-center p-12">
		<h3 class="animate-pulse text-slate-400">Loading Chapters...</h3>
	</div>
{:then}
	<div class="mx-auto max-w-4xl space-y-6 p-4">
		<Card size="none" class="border-slate-200 p-6 shadow-sm">
			<Label class="mb-2 text-lg font-bold">Chapter Browser</Label>
			<Select
				items={chaps}
				bind:value={selected}
				onchange={handleSelectChange}
				placeholder="Select a chapter..."
			/>
		</Card>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			{#each members as m}
				<Card size="none" class="border-slate-100 p-4 transition-colors hover:bg-slate-50">
					<div class="flex flex-col">
						<div class="flex items-start justify-between">
							<a href="#/member/{m.ID}" class="text-lg font-bold text-primary-700 hover:underline">
								{oslname(m)}
							</a>
							{#if prior == m.ID}
								<Badge color="purple" class="text-[10px] uppercase">Prior</Badge>
							{/if}
						</div>

						<div class="mt-2 space-y-1 text-sm text-slate-600">
							{#if m.PrimaryEmail}
								<a href="mailto:{m.PrimaryEmail}" class="block hover:text-primary-600"
									>{m.PrimaryEmail}</a
								>
							{/if}
							{#if m.PrimaryPhone}
								<a href="tel:{m.PrimaryPhone}" class="block hover:text-primary-600"
									>{m.PrimaryPhone}</a
								>
							{/if}
						</div>
					</div>
				</Card>
			{:else}
				{#if selected}
					<div class="col-span-full py-12 text-center text-slate-400 italic">
						No members found in this chapter.
					</div>
				{/if}
			{/each}
		</div>
	</div>
{:catch error}
	<Card color="red" class="mx-auto mt-8">
		<p class="font-bold">Error loading chapters:</p>
		<p>{error.message}</p>
	</Card>
{/await}
