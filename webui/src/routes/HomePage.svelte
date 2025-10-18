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
		Select
	} from 'flowbite-svelte';
	import { updateMember, search, searchemail, cleanDateFormat } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';
	import { SearchOutline } from 'flowbite-svelte-icons';

	let { params = {} } = $props();
	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	let query = $state();
	let result = $state();
	let mode = $state('name');

	const modes = [
		{ value: 'name', name: 'Name', selected: true },
		{ value: 'email', name: 'Full/Exact Email Address' }
	];

	if (params.query) {
		// https://svelte.dev/e/state_referenced_locally
		if (() => query != params.query) {
			query = params.query;
			const event = new Event('search', { bubbles: true, cancelable: true });
			doSearch(event);
		}
	}

	async function doSearch(event) {
		event.preventDefault();
		event.stopPropagation();

		if (!query) {
			toast.push('Please enter a search query, mimimum of 3 letters');
			return;
		}

		query = query.trim();

		if (query.length < 3) {
			toast.push('Query too short (A minimum 3 of letters are required)');
			return;
		}

		if (query.indexOf(' ') != -1) {
			const subs = query.split(' ');
			query = '';
			// look for a substring that is long enough
			for (const sub of subs) {
				console.log(sub);
				if (sub.length >= 3) {
					query = sub;
					toast.push('Using: ' + query);
					break;
				}
			}
			if (query == '') {
				toast.push('Invalid query (read the directions, please)');
				resetSearch(event);
				return;
			}
		}

		if (query.indexOf('@') != -1) {
			toast.push('Searching email addresses');
			mode = 'email';
		} else {
			mode = 'name';
		}

		try {
			switch (mode) {
				case 'email':
					result = await searchemail(query);
					break;
				default:
					result = await search(query);
					break;
			}
			if (result == null || result.length == 0) {
				result = [];
			}
			push(`/search/${query}`);
		} catch (err) {
			console.log(err);
			toast.push(err.message);
		}
	}

	async function resetSearch(event) {
		event.preventDefault();
		event.stopPropagation();
		result = [];
		query = '';
		mode = 'name';
		push(`/`);
	}

	async function quickRenew(event, r) {
		event.preventDefault();
		event.stopPropagation();
		try {
			await updateMember(r.ID, 'DateReaffirmation', cleanDateFormat(new Date().toISOString()));
			push(`/member/${r.ID}`);
		} catch (err) {
			console.log(err);
			toast.push(err.message);
		}
	}
</script>

<svelte:head>
	<title>OSL Member Manager : {query}</title>
</svelte:head>

<form onsubmit={doSearch}>
	<Table>
		<TableBody>
			<TableBodyRow>
				<TableBodyCell colspan={3}>
					<Button href="#/me">My Data</Button>
				</TableBodyCell>
			</TableBodyRow>
			{#if !result}
				<TableBodyRow>
					<TableBodyCell colspan={3}>
						{#if mode != 'email'}
							This searches first, last, life-vow, and preferred name as individual fields.<br />
							Do <b>not</b> type full names (e.g. "Bob Smith"), they will not match the individual
							fields.<br />
							Use <b>one</b> name or, better still, a partial name to search.<br />
							A minimum of 3 letters are required.<br />
							Case is ignored.<br />
							<i>
								e.g. <b>"sMi"</b> matches both "<b>Smi</b>thers Boberson" and "Bob
								<b>Smi</b>th".
							</i> <br />
						{/if}
						{#if mode == 'email'}
							This searches by EXACT and FULL email address, not partial matches.
						{/if}
					</TableBodyCell>
				</TableBodyRow>
			{/if}
			<TableBodyRow>
				<TableBodyCell>
					<Select class="mt-2" items={modes} bind:value={mode} />
				</TableBodyCell>
				<TableBodyCell>
					<Input type="text" name="query" bind:value={query} />
				</TableBodyCell>
				<TableBodyCell
					><Button color="green" type="submit"><SearchOutline class="h-6 w-6" /> Search</Button
					></TableBodyCell
				>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
<Table class="w-full">
	<TableBody>
		{#each result as r}
			<TableBodyRow>
				<TableBodyCell><a href="#/member/{r.ID}">{r.FirstName}</a></TableBodyCell>
				<TableBodyCell><a href="#/member/{r.ID}">{r.PreferredName}</a></TableBodyCell>
				<TableBodyCell><a href="#/member/{r.ID}">{r.LastName}</a></TableBodyCell>
				<TableBodyCell>{r.MemberStatus}</TableBodyCell>
				<TableBodyCell>
					{#if $me && $me.level > 1}
						<Button color="purple" onclick={(e) => quickRenew(e, r)}>Quick Renew</Button>
					{/if}
				</TableBodyCell>
			</TableBodyRow>
		{/each}
		<TableBodyRow>
			<TableBodyCell colspan={4}>&nbsp;</TableBodyCell>
			<TableBodyCell>
				{#if result}
					<form onsubmit={resetSearch}>
						<Button color="red" type="submit">Reset</Button>
					</form>
				{/if}
			</TableBodyCell>
		</TableBodyRow>
		{#if result == '' && query.length > 3}
			<TableBodyRow>
				<TableBodyCell colspan={4}>No results for query : {query}</TableBodyCell>
				<TableBodyCell>
					<Button color="purple" onclick={() => push('#/addmember')}>Add Member</Button>
				</TableBodyCell>
			</TableBodyRow>
		{/if}
	</TableBody>
</Table>
