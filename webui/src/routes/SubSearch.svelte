<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { subsearch } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	let query = $state();
	let result = $state();

	async function doSearch(e) {
		e.preventDefault();
		e.stopPropagation();

		try {
			result = await subsearch(query);
		} catch (e) {
			console.log(e);
			toast.push(e.message);
		}
	}

	async function resetSearch(e) {
		e.preventDefault();
		e.stopPropagation();
		result = '';
	}
</script>

<svelte:head>
	<title>OSL Subscriber Manager</title>
</svelte:head>

{#if !result}
	<form onsubmit={doSearch}>
		<Table>
			<TableBody>
				<TableBodyRow>
					<TableBodyCell>Subscriber Search:</TableBodyCell>
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
	<form onsubmit={resetSearch}>
		<Table>
			<TableBody>
				<TableBodyRow>
					<TableBodyCell colspan="2">&nbsp;</TableBodyCell>
					<TableBodyCell><Button type="submit">Reset</Button></TableBodyCell>
				</TableBodyRow>
				{#each result as r, i}
					<TableBodyRow>
						<TableBodyCell>{i}</TableBodyCell>
						<TableBodyCell><a href="#/subscriber/{r.ID}">{r.Name}</a></TableBodyCell>
						<TableBodyCell><a href="#/subscriber/{r.ID}">{r.Attn}</a></TableBodyCell>
					</TableBodyRow>
				{/each}
			</TableBody>
		</Table>
	</form>
{/if}

<div>
	<p>
		<a href="#/">Member Search</a>
	</p>
</div>
