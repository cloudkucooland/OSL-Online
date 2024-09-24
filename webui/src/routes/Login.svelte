<script>
import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
import { toast } from '@zerodevx/svelte-toast';
import { push } from "svelte-spa-router";
import { getJWT } from "../oo";

const jwt = localStorage.getItem('jwt');
if (jwt) {
  toast.push('Logged out');
  localStorage.removeItem('jwt');
}

// export let data;
let username;
let password;

async function doLogin(event) {
  event.preventDefault();
  event.stopPropagation();
  try {
    await getJWT(username, password);
    push('/');
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
		</TableBody>
	</Table>
</form>
