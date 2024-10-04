<script>
	import { getContext } from 'svelte';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	// import { push } from 'svelte-spa-router';
	import { postRegister } from '../oo';

	// const { me } = getContext('oo');
	// if ($me !== undefined) { }

	let username;
	let submitted = false;

	async function doRegister(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			submitted = await postRegister(username);
		} catch (e) {
			console.log(e);
			toast.push(e);
		}
	}
</script>

<form on:submit={doRegister}>
	<Table>
		<TableBody>
			{#if !submitted}
				<TableBodyRow>
					<TableBodyCell>Email Address:</TableBodyCell>
					<TableBodyCell><Input type="text" name="username" bind:value={username} /></TableBodyCell>
				</TableBodyRow>
				<TableBodyRow>
					<TableBodyCell>&nbsp;</TableBodyCell>
					<TableBodyCell><Button type="submit">Register/Recover</Button></TableBodyCell>
				</TableBodyRow>
			{:else}
				<TableBodyRow>
					<TableBodyCell
						>Request submitted. Please check your email for your password.</TableBodyCell
					>
				</TableBodyRow>
			{/if}
		</TableBody>
	</Table>
</form>
