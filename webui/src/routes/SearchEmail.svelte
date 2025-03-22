<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { searchemail } from '../oo';
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

	async function doSearch(event) {
		event.preventDefault();
		event.stopPropagation();

		if (!query || query == '' || query == 'undefined') {
			toast.push('please enter a search query, a full email address');
			return;
		}

		try {
			result = await searchemail(query);
			if (result == null || result.length == 0) {
				result = [
					{ id: 0, FirstName: 'no results', LastName: '', PreferredName: '', MemberStatus: '' }
				];
			}
			push(`#/searchemail/${query}`);
		} catch (err) {
			console.log(err);
			toast.push(err.message);
		}
	}

	async function resetSearch(event) {
		event.preventDefault();
		event.stopPropagation();
		result = '';
		query = '';
		push(`#/searchemail`);
	}
</script>

<svelte:head>
	<title>OSL Member Manager : Email Search</title>
</svelte:head>

{#if !result}
	<form onsubmit={doSearch}>
		<Table>
			<TableBody>
				<TableBodyRow>
					<TableBodyCell colspan={2}>
						Seach by exact-match email address<br />
					</TableBodyCell>
				</TableBodyRow>
				<TableBodyRow>
					<TableBodyCell>Email Search:</TableBodyCell>
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
	<form onsubmit={resetSearch}>
		<Table>
			<TableBody>
				{#each result as r, i}
					<TableBodyRow>
						<TableBodyCell>{i}</TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.FirstName}</a></TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.PreferredName}</a></TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.LastName}</a></TableBodyCell>
						<TableBodyCell>{r.MemberStatus}</TableBodyCell>
					</TableBodyRow>
				{/each}
				<TableBodyRow>
					<TableBodyCell colspan={4}>&nbsp;</TableBodyCell>
					<TableBodyCell><Button type="submit">Reset</Button></TableBodyCell>
				</TableBodyRow>
			</TableBody>
		</Table>
	</form>
{/if}
