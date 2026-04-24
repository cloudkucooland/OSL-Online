<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { Select, Card, Badge, Label, Heading, Spinner } from 'flowbite-svelte';
	import { UsersGroupOutline, EnvelopeOutline, PhoneOutline } from 'flowbite-svelte-icons';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getLeaders } from '../oo';

	const oo = getContext('oo');

	// Guard: Ensure user is logged in
	if (!oo.me) {
		push('/Login');
	}

	let { params = {} } = $props();
	let categories = [
		{ name: 'Priors', value: 'prior' },
		{ name: 'Abbatial Appointment', value: 'council' },
		{ name: 'Canon', value: 'canon' },
		{ name: 'General Officer', value: 'elected' }
	];
	let leaders = $state([]);
	let selected = $state();
	let loading = $state(false);

	async function chooseType() {
		if (!selected) return;
		loading = true;
		try {
			leaders = await getLeaders(selected);
			push(`#/leadership/${selected}`);
		} catch (e: any) {
			console.error(e);
			toast.push(e.message);
		} finally {
			loading = false;
		}
	}

	onMount(async () => {
		if (params.id) {
			selected = params.id;
			await chooseType();
		}
	});
</script>

<svelte:head>
	<title>Leadership Browser | OSL Directory</title>
</svelte:head>

<div class="mx-auto max-w-6xl space-y-6 p-4">
	<Card size="none" class="border-slate-200 bg-white p-6 shadow-sm">
		<div class="mb-6 flex items-center gap-3 text-primary-900">
			<UsersGroupOutline size="lg" />
			<Heading tag="h2" class="text-2xl font-bold">Leadership Browser</Heading>
		</div>

		<div class="space-y-2">
			<Label for="category-select" class="text-sm font-semibold text-slate-700">Select Category</Label>
			<Select
				id="category-select"
				items={categories}
				bind:value={selected}
				onchange={chooseType}
				placeholder="Choose a leadership group..."
			/>
		</div>
	</Card>

	{#if loading}
		<div class="flex justify-center py-20">
			<Spinner color="purple" size="12" />
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each leaders as m}
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
							<Badge color="purple" class="text-[10px] uppercase">
								{m.MemberStatus}
							</Badge>
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
						No members found in this leadership category.
					</div>
				{/if}
			{/each}
		</div>
	{/if}
</div>