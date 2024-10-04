<script>
	import { getContext } from 'svelte';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { getJWT } from '../oo';

	const { me } = getContext('oo');

	const jwt = localStorage.getItem('jwt');
	if (jwt) {
		localStorage.removeItem('jwt');
		$me = undefined;
		toast.push('Logged out');
	}

	let username;
	let password;

	async function doLogin(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			await getJWT(username, password);
			// push('/');
			window.location.href = '';
		} catch (e) {
			console.log(e);
			toast.push(e);
		}
	}
</script>

<form on:submit={doLogin}>
	<Table>
		<TableBody>
			<TableBodyRow>
				<TableBodyCell>Username:</TableBodyCell>
				<TableBodyCell><Input type="text" name="username" bind:value={username} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Password:</TableBodyCell>
				<TableBodyCell
					><Input type="password" name="password" bind:value={password} /></TableBodyCell
				>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>&nbsp;</TableBodyCell>
				<TableBodyCell><Button type="submit">Login</Button></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell colspan="2"><a href="#/register">Register/Lost Password</a></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
