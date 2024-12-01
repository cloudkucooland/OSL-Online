<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { updateMember, search, vcard } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	let { params = {} } = $props();
	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	let query = $state();
	let result = $state();

	if (params.query) {
		query = params.query;
		const e = new Event('search', { bubbles: true, cancelable: true });
		doSearch(e);
	}

	async function doSearch(e) {
		e.preventDefault();
		e.stopPropagation();

		if (!query) {
			toast.push('Please enter a search query, mimimum of 3 letters');
			return;
		}

		query = query.trim();

		if (query.length < 3) {
			toast.push('Query too short (A minimum 3 of letters are required)');
			return;
		}

		try {
			result = await search(query);
			if (result == null || result.length == 0) {
				result = [
					{ id: 0, FirstName: 'no results', LastName: '', PreferredName: '', MemberStatus: '' }
				];
			}
			push(`/search/${query}`);
		} catch (err) {
			console.log(err);
			toast.push(err.message);
		}
	}

	async function resetSearch(e) {
		e.preventDefault();
		e.stopPropagation();
		result = '';
		query = '';
		push(`/`);
	}

	async function quickRenew(e, r) {
		e.preventDefault();
		e.stopPropagation();
		try {
			const dd = new Date().toISOString().split('T');
			await updateMember(r.ID, "DateReaffirmation", dd[0]);
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

{#if !result}
	<form onsubmit={doSearch}>
		<Table>
			<TableBody>
				<TableBodyRow>
					<TableBodyCell colspan={2}>
						This searches first, last, life-vow, and preferred name as individual fields.<br />
						Do <b>not</b> type full names (e.g. "Bob Smith"), they will not match the individual
						fields.<br />
						Use <b>one</b> name or, better still, a partial name to search.<br />
						A minimum of 3 letters are required.<br />
						Case is ignored.<br />
						<i>e.g. <b>"sMi"</b> matches both "<b>Smi</b>thers Boberson" and "Bob <b>Smi</b>th".</i
						><br />
					</TableBodyCell>
				</TableBodyRow>
				<TableBodyRow>
					<TableBodyCell>Member Search:</TableBodyCell>
					<TableBodyCell>
						<Input type="text" name="query" bind:value={query} />
					</TableBodyCell>
				</TableBodyRow>
				<TableBodyRow>
					<TableBodyCell>&nbsp;</TableBodyCell>
					<TableBodyCell><Button type="submit">Search</Button></TableBodyCell>
				</TableBodyRow>
			</TableBody>
		</Table>
	</form>
{:else}
	<form>
		<Table>
			<TableBody>
				{#each result as r, i}
					<TableBodyRow>
						<TableBodyCell><a href="#/member/{r.ID}">{r.FirstName}</a></TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.PreferredName}</a></TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.LastName}</a></TableBodyCell>
						<TableBodyCell>{r.MemberStatus}</TableBodyCell>
						<TableBodyCell><Button onclick={() => vcard(r.ID)}>Add to Contacts</Button></TableBodyCell>
						{#if $me && $me.level > 1}
						<TableBodyCell><Button onclick={(e) => quickRenew(e, r)}>Quick Renew</Button></TableBodyCell>
						{/if}
					</TableBodyRow>
				{/each}
				<TableBodyRow>
					<TableBodyCell colspan={4}>&nbsp;</TableBodyCell>
					<TableBodyCell><Button onclick={resetSearch}>Reset</Button></TableBodyCell>
				</TableBodyRow>
			</TableBody>
		</Table>
	</form>
{/if}

{#if $me && $me.level > 1}
	<div>
		<p>
			<a href="#/subsearch">Subscriber Search</a> | <a href="#/searchemail">Email Search</a> |
			<a href="#/addmember">Add Member</a>
		</p>
	</div>
{/if}
