<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { Select, Label, Card, Badge, Heading, Spinner } from 'flowbite-svelte';
	import { BookOpenOutline, EnvelopeOutline, PhoneOutline } from 'flowbite-svelte-icons';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getChapters, getChapterMembers } from '../oo';

	let { params = {} } = $props();
	let chaps = $state([]);
	let members = $state([]);
	let selected = $state();
	let prior = $state();
	let loading = $state(false);

	const oo = getContext('oo');

	async function updateData() {
		if (!selected) return;
		loading = true;
		try {
			members = await getChapterMembers(selected);
			const chap = oo.chaptercache.find((c: any) => c.ID == selected);
			prior = chap?.Prior;
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	}

	async function handleSelectChange() {
		if (selected) {
			push(`/chapterbrowser/${selected}`);
			await updateData();
		}
	}

	onMount(async () => {
		if (!oo.me) {
			push('/Login');
			return;
		}
		
		if (!oo.chaptercache) {
			try {
				oo.chaptercache = await getChapters();
			} catch (err: any) {
				toast.push(err.message);
			}
		}
		chaps = oo.chaptercache || [];

		if (params.id) {
			selected = params.id;
			await updateData();
		}
	});
</script>

<svelte:head>
	<title>Chapter Browser | OSL Directory</title>
</svelte:head>

<div class="mx-auto max-w-6xl space-y-6 p-4">
	<Card size="none" class="border-slate-200 bg-white p-6 shadow-sm">
		<div class="mb-6 flex items-center gap-3 text-primary-900">
			<BookOpenOutline size="lg" />
			<Heading tag="h2" class="text-2xl font-bold">Chapter Browser</Heading>
		</div>

		<div class="space-y-2">
			<Label for="chapter-select" class="text-sm font-semibold text-slate-700">Select Chapter</Label>
			<Select
				id="chapter-select"
				items={chaps}
				bind:value={selected}
				onchange={handleSelectChange}
				placeholder="Choose a chapter to view members..."
			/>
		</div>
	</Card>

	{#if loading}
		<div class="flex justify-center py-20">
			<Spinner color="purple" size="12" />
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each members as m}
				<Card
					size="none"
					class="group border-slate-100 p-5 shadow-sm transition-all hover:border-primary-200 hover:bg-slate-50"
				>
					<div class="flex h-full flex-col">
						<div class="mb-3 flex items-start justify-between">
							<a
								href="#/member/{m.ID}"
								class="text-lg font-bold text-slate-900 transition-colors group-hover:text-primary-700"
							>
								{oslname(m)}
							</a>
							{#if prior == m.ID}
								<Badge color="purple" class="text-[10px] uppercase">Prior</Badge>
							{:else}
								<Badge color="indigo" class="text-[10px] uppercase">{m.MemberStatus}</Badge>
							{/if}
						</div>

						<div class="mt-auto space-y-2 border-t border-slate-100 pt-3 text-sm">
							{#if m.PrimaryEmail}
								<div class="flex items-center gap-2 text-slate-600">
									<EnvelopeOutline size="xs" class="text-slate-400" />
									<a
										href="mailto:{m.PrimaryEmail}"
										class="truncate hover:text-primary-600 hover:underline"
									>
										{m.PrimaryEmail}
									</a>
								</div>
							{/if}
							{#if m.PrimaryPhone}
								<div class="flex items-center gap-2 text-slate-600">
									<PhoneOutline size="xs" class="text-slate-400" />
									<a href="tel:{m.PrimaryPhone}" class="hover:text-primary-600">
										{m.PrimaryPhone}
									</a>
								</div>
							{/if}
						</div>
					</div>
				</Card>
			{:else}
				{#if selected}
					<div class="col-span-full py-16 text-center text-slate-400 italic">
						No members found in this chapter.
					</div>
				{/if}
			{/each}
		</div>
	{/if}
</div>