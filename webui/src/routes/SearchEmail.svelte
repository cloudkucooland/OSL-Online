<script lang="ts">
	import { getContext } from 'svelte';
	import { push, replace } from 'svelte-spa-router';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { getMe, searchemail } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	export let params = {};
	const { me } = getContext('oo');
	if ($me === undefined) {
		replace('/Login');
	}

	let query;
	let result;

	if (params.query) {
		query = params.query;
		const e = new Event('search', { bubbles: true, cancelable: true });
		doSearch(e);
	}

	async function doSearch(e) {
		e.preventDefault();
		e.stopPropagation();

		try {
			result = await searchemail(query);
			if (result == null || result.length == 0) {
				result = [
					{ id: 0, FirstName: 'no results', LastName: '', PreferredName: '', MemberStatus: '' }
				];
			}
			replace(`/emailsearch/${query}`);
		} catch (e) {
			console.log(e);
			toast.push(e.message);
		}
	}

	async function resetSearch(e) {
		e.preventDefault();
		e.stopPropagation();
		push(`/`);
		result = '';
	}
</script>

<svelte:head>
	<title>OSL Member Manager : Email Search</title>
</svelte:head>

{#if !result}
	<form on:submit={doSearch}>
		<Table>
			<TableBody>
				<TableBodyRow>
					<TableBodyCell colspan={2}>
						Seach by exact-match email address<br />
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
	<form on:submit={resetSearch}>
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
