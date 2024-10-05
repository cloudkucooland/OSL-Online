<script lang="ts">
	import { getContext } from 'svelte';
	import { push, replace } from 'svelte-spa-router';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { getMe, search } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	export let params = {};
	const { me } = getContext('oo');
	// console.log('in HomePage', $me);
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
			result = await search(query);
			if (result == null || result.length == 0) {
				result = [
					{ id: 0, FirstName: 'no results', LastName: '', PreferredName: '', MemberStatus: '' }
				];
			}
			replace(`/search/${query}`);
		} catch (e) {
			console.log(e);
			toast.push(e);
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
	<title>OSL Member Manager : {query}</title>
</svelte:head>

{#if !result}
	<form on:submit={doSearch}>
		<Table>
			<TableBody>
				<TableBodyRow>
					<TableBodyCell>Member Search:</TableBodyCell>
					<TableBodyCell>
						<Input type="text" name="query" bind:value={query} on:change={query} />
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
				<TableBodyRow>
					<TableBodyCell colspan={4}>&nbsp;</TableBodyCell>
					<TableBodyCell><Button type="submit">Reset</Button></TableBodyCell>
				</TableBodyRow>
				{#each result as r, i}
					<TableBodyRow>
						<TableBodyCell>{i}</TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.FirstName}</a></TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.PreferredName}</a></TableBodyCell>
						<TableBodyCell><a href="#/member/{r.ID}">{r.LastName}</a></TableBodyCell>
						<TableBodyCell>{r.MemberStatus}</TableBodyCell>
					</TableBodyRow>
				{/each}
			</TableBody>
		</Table>
	</form>
{/if}

{#if $me && $me.level > 1}
	<div>
		<p>
			<a href="#/subsearch">Subscriber Search</a> | <a href="#/addmember">Add Member</a>
		</p>
	</div>
{/if}
