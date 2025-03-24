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

	async function doSearch(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			result = await subsearch(query);
		} catch (err) {
			console.log(err);
			toast.push(err.message);
		}
	}

	async function resetSearch(event) {
		event.preventDefault();
		event.stopPropagation();
		result = [];
	}
</script>

<svelte:head>
	<title>OSL Subscriber Manager</title>
</svelte:head>

<form onsubmit={doSearch}>
	<Table>
		<TableBody>
			<TableBodyRow>
				<TableBodyCell>Subscriber Search:</TableBodyCell>
				<TableBodyCell>
					<Input type="text" name="query" bind:value={query} />
				</TableBodyCell>
				<TableBodyCell>
					<Button type="submit" color="green">Search</Button>
				</TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
{#if result}
	<Table>
		<TableBody>
			{#each result as r, i}
				<TableBodyRow>
					<TableBodyCell>{i}</TableBodyCell>
					<TableBodyCell><a href="#/subscriber/{r.ID}">{r.Name}</a></TableBodyCell>
					<TableBodyCell><a href="#/subscriber/{r.ID}">{r.Attn}</a></TableBodyCell>
				</TableBodyRow>
			{/each}
			<TableBodyRow>
				<TableBodyCell colspan={3}>
					<form onsubmit={resetSearch}>
						<Button type="submit" color="red">Reset</Button>
					</form>
				</TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
{/if}
