<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Button,
		Input,
		Heading,
		Card,
		Spinner,
		Badge
	} from 'flowbite-svelte';
	import { SearchOutline, CloseCircleOutline, BuildingOutline } from 'flowbite-svelte-icons';
	import { subsearch } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	const oo = getContext('oo');

	// Guard: Fix the $me reference
	if (!oo.me) {
		push('/Login');
	}

	let query = $state('');
	let result = $state(null);
	let searching = $state(false);

	async function doSearch(event?: Event) {
		if (event) event.preventDefault();
		if (!query.trim()) return;

		searching = true;
		try {
			result = await subsearch(query);
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			searching = false;
		}
	}

	function clearSearch() {
		query = '';
		result = null;
	}
</script>

<svelte:head>
	<title>Subscriber Search | OSL</title>
</svelte:head>

<div class="mx-auto max-w-4xl space-y-6 px-4 py-8">
	<header class="flex items-center gap-3">
		<div class="rounded-lg bg-blue-100 p-2">
			<BuildingOutline class="h-6 w-6 text-blue-600" />
		</div>
		<div>
			<Heading tag="h2" class="text-2xl font-bold text-slate-900">Institutional Subscribers</Heading
			>
			<p class="text-sm italic text-slate-500">
				Search Doxology subscriptions for libraries and institutions.
			</p>
		</div>
	</header>

	<Card size="none" class="border-slate-200 bg-slate-50/50 p-6 shadow-sm">
		<form onsubmit={doSearch} class="flex gap-2">
			<div class="relative flex-grow">
				<div class="pointer-events-none absolute inset-y-0 start-0 flex items-center ps-3">
					<SearchOutline class="h-4 w-4 text-gray-500" />
				</div>
				<Input
					type="text"
					placeholder="Search by institution name or attention field..."
					class="bg-white ps-10"
					bind:value={query}
				/>
			</div>
			<Button type="submit" color="blue" disabled={searching || !query.trim()}>
				{#if searching}<Spinner size="4" class="mr-2" />{/if}
				Search
			</Button>
			{#if result}
				<Button color="alternative" onclick={clearSearch} title="Clear Results">
					<CloseCircleOutline class="h-5 w-5" />
				</Button>
			{/if}
		</form>
	</Card>

	{#if result}
		<Card size="none" class="overflow-hidden border-slate-200 shadow-md">
			<Table hoverable={true}>
				<TableHead class="border-b bg-slate-100">
					<TableHeadCell class="w-12 text-center">#</TableHeadCell>
					<TableHeadCell>Institution Name</TableHeadCell>
					<TableHeadCell>Attention</TableHeadCell>
					<TableHeadCell class="text-right">Action</TableHeadCell>
				</TableHead>
				<TableBody>
					{#each result as r, i}
						<TableBodyRow class="group cursor-pointer" onclick={() => push(`/subscriber/${r.ID}`)}>
							<TableBodyCell class="text-center font-mono text-xs text-slate-400"
								>{i + 1}</TableBodyCell
							>
							<TableBodyCell class="font-bold text-slate-900 group-hover:text-blue-700">
								{r.Name}
							</TableBodyCell>
							<TableBodyCell>{r.Attn || '—'}</TableBodyCell>
							<TableBodyCell class="text-right">
								<Badge
									color="blue"
									outline
									class="transition-colors group-hover:bg-blue-600 group-hover:text-white"
								>
									View Record
								</Badge>
							</TableBodyCell>
						</TableBodyRow>
					{:else}
						<TableBodyRow>
							<TableBodyCell colspan={4} class="text-center py-16 bg-slate-50">
								<p class="text-slate-500 font-medium">
									No subscribers found matching "<span class="text-slate-900">{query}</span>"
								</p>
								<Button
									color="alternative"
									size="xs"
									class="mt-4"
									onclick={() => push('/addsubscriber')}
								>
									Add New Subscriber
								</Button>
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		</Card>
	{/if}
</div>
