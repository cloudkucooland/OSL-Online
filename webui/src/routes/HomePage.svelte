<script lang="ts">
	import { getContext } from 'svelte';
	import { push, replace } from 'svelte-spa-router';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { search } from '../oo';
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
	<title>OSL Member Manager : {query}</title>
</svelte:head>

{#if !result}
	<form on:submit={doSearch}>
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

{#if $me && $me.level > 1}
	<div>
		<p>
			<a href="#/subsearch">Subscriber Search</a> | <a href="#/searchemail">Email Search</a> |
			<a href="#/addmember">Add Member</a>
		</p>
	</div>
{/if}
