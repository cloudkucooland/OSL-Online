<script lang="ts">
import { push } from "svelte-spa-router";
import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
import { getMe, search } from "../oo";
import { toast } from '@zerodevx/svelte-toast';

const me = getMe();
if (me === undefined) {
  push('/Login');
}

let query;
let result;

async function doSearch(e) {
  e.preventDefault();
  e.stopPropagation();

  try {
    result = await search(query);
  } catch(e) {
    console.log(e);
    toast.push(e);
  }
}

async function resetSearch(e) {
  e.preventDefault();
  e.stopPropagation();
  result = '';
}

</script>

<svelte:head>
  <title>OSL Member Manager</title>
</svelte:head>

{#if !result}
<form on:submit={doSearch} >
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
<form on:submit={resetSearch} >
<Table>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell colspan="3">&nbsp;</TableBodyCell>
			<TableBodyCell><Button type="submit">Reset</Button></TableBodyCell>
		</TableBodyRow>
		{#each result as r, i}
		<TableBodyRow>
			<TableBodyCell>{i}</TableBodyCell>
			<TableBodyCell><a href="#/member/{r.ID}">{r.FirstName}</a></TableBodyCell>
			<TableBodyCell><a href="#/member/{r.ID}">{r.LastName}</a></TableBodyCell>
			<TableBodyCell>{r.MemberStatus}</TableBodyCell>
		</TableBodyRow>
		{/each}
	</TableBody>
</Table>
</form>
{/if}

<div>
<p>
<a href="#/subsearch">Subscriber Search</a>
</p>
</div>
