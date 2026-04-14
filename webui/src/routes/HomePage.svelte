<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		Button,
		Input,
		Select,
		Card,
		Badge,
		Heading
	} from 'flowbite-svelte';
	import { updateMember, search, searchemail, cleanDateFormat } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';
	import {
		SearchOutline,
		RefreshOutline,
		UserAddOutline,
		InfoCircleOutline
	} from 'flowbite-svelte-icons';

	let { params = {} } = $props();
	const oo = getContext('oo');

	// If no user is logged in, redirect
	if (!oo.me) {
		push('/Login');
	}

	let query = $state('');
	let result = $state(null);
	let mode = $state('name');

	$effect(() => {
		if (params.query) {
			query = params.query;
			// optional: trigger search here if you want it automatic
		}
	});

	const modes = [
		{ value: 'name', name: 'Name' },
		{ value: 'email', name: 'Full Email' }
	];

	// Run search automatically if URL has a query
	$effect(() => {
		if (params.query && !result) {
			query = params.query;
			handleSearch();
		}
	});

	async function handleSearch(event?: Event) {
		if (event) {
			event.preventDefault();
		}

		if (!query || query.trim().length < 3) {
			toast.push('A minimum 3 letters are required');
			return;
		}

		const cleanedQuery = query.trim();

		// Handle multi-word strings by grabbing the first valid chunk
		if (cleanedQuery.includes(' ')) {
			const subs = cleanedQuery.split(' ');
			const firstValid = subs.find((s) => s.length >= 3);
			if (firstValid) {
				query = firstValid;
				toast.push(`Using: ${query}`);
			} else {
				toast.push('Invalid query format');
				return;
			}
		}

		// Auto-detect email mode
		if (query.includes('@')) {
			mode = 'email';
		}

		try {
			result = mode === 'email' ? await searchemail(query) : await search(query);
			if (!result) result = [];
			push(`/search/${query}`);
		} catch (err: any) {
			console.error(err);
			toast.push(err.message);
		}
	}

	function resetSearch() {
		result = null;
		query = '';
		mode = 'name';
		push('/');
	}

	async function quickRenew(r: any) {
		try {
			await updateMember(r.ID, 'DateReaffirmation', cleanDateFormat(new Date().toISOString()));
			push(`/member/${r.ID}`);
			toast.push(`Renewed ${r.LastName}`);
		} catch (err: any) {
			toast.push(err.message);
		}
	}
</script>

<svelte:head>
	<title>OSL Directory : {query || 'Search'}</title>
</svelte:head>

<div class="space-y-6">
	<section
		class="flex flex-col gap-4 rounded-lg border border-slate-200 bg-slate-50 p-6 md:flex-row md:items-end"
	>
		<div class="flex-grow space-y-2">
			<label for="search-input" class="text-sm font-medium text-slate-700">Find Members</label>
			<div class="flex gap-2">
				<Select items={modes} bind:value={mode} class="w-32 md:w-48" />
				<Input
					id="search-input"
					type="text"
					placeholder={mode === 'name' ? 'Search names...' : 'Exact email...'}
					bind:value={query}
					onkeydown={(e) => e.key === 'Enter' && handleSearch()}
				>
					<SearchOutline slot="left" class="h-5 w-5 text-slate-400" />
				</Input>
			</div>
		</div>
		<div class="flex gap-2">
			<Button color="alternative" onclick={resetSearch} disabled={!query}>Reset</Button>
			<Button color="primary" onclick={() => handleSearch()}>
				<SearchOutline class="mr-2 h-5 w-5" /> Search
			</Button>
		</div>
	</section>

	{#if !result}
		<Card size="xl" padding="lg" class="border-l-4 border-l-primary-600 shadow-sm">
			<div class="flex items-start gap-3">
				<InfoCircleOutline class="mt-1 h-6 w-6 text-primary-600" />
				<div class="space-y-2 text-slate-600">
					<Heading tag="h4" class="text-lg font-semibold text-slate-800">Search Tips</Heading>
					<ul class="list-inside list-disc space-y-1 text-sm">
						<li>Search by individual fields (First, Last, Preferred, or Life-vow names).</li>
						<li>
							<strong>Avoid full names</strong> (e.g. "Bob Smith") as they won't match individual fields.
						</li>
						<li>Minimum of 3 letters required; case is ignored.</li>
						<li>
							Email search requires the <strong>exact</strong> and <strong>full</strong> address.
						</li>
						<li>
							Don't search for yourself; use the "Me" menu above. Those who most need to read this
							message are the least likely to read it....
						</li>
					</ul>
				</div>
			</div>
		</Card>
	{:else if result.length > 0}
		<div class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
			<Table hoverable={true}>
				<thead>
					<tr class="border-b border-slate-200 bg-slate-50">
						<TableBodyCell class="font-bold">First Name</TableBodyCell>
						<TableBodyCell class="font-bold">Preferred</TableBodyCell>
						<TableBodyCell class="font-bold">Last Name</TableBodyCell>
						<TableBodyCell class="font-bold">Status</TableBodyCell>
						<TableBodyCell class="text-right font-bold">Actions</TableBodyCell>
					</tr>
				</thead>
				<TableBody>
					{#each result as r}
						<TableBodyRow class="cursor-pointer" onclick={() => push(`/member/${r.ID}`)}>
							<TableBodyCell class="font-medium text-primary-700">{r.FirstName}</TableBodyCell>
							<TableBodyCell>{r.LifeVowedName || r.PreferredName || '-'}</TableBodyCell>
							<TableBodyCell class="font-semibold text-slate-900">{r.LastName}</TableBodyCell>
							<TableBodyCell>
								<Badge
									color={r.MemberStatus.includes('Vows')
										? 'green'
										: r.MemberStatus === 'Deceased'
											? 'purple'
											: r.MemberStatus === 'Removed'
												? 'red'
												: 'indigo'}
								>
									{r.MemberStatus}
								</Badge>
							</TableBodyCell>
							<TableBodyCell class="text-right" onclick={(e) => e.stopPropagation()}>
								{#if oo.me.level > 1}
									<Button size="xs" color="purple" outline onclick={() => quickRenew(r)}>
										<RefreshOutline class="mr-1 h-4 w-4" /> Renew
									</Button>
								{/if}
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		</div>
	{:else if result.length === 0}
		<div class="rounded-lg border-2 border-dashed border-slate-300 bg-slate-50 py-12 text-center">
			<div class="mb-4">No results found for "<span class="font-bold">{query}</span>"</div>
			<Button color="alternative" onclick={() => push('/addmember')}>
				<UserAddOutline class="mr-2 h-5 w-5" /> Add New Member
			</Button>
		</div>
	{/if}
</div>
