<script>
	import { getContext } from 'svelte';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		Button,
		Input,
		Label
	} from 'flowbite-svelte';
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
			toast.push(e.message);
		}
	}
</script>

<form on:submit={doLogin}>
	<div class="grid grid-cols-8 gap-4 px-4 py-2">
		<div class="col-span-8">
			Your username is your primary email address as listed in our system.<br />
			<b
				>If this is your first time logging in, you will need to register using the link below to
				have your password emailed to you</b
			>
		</div>
		<div class="col-span-2"><Label for="username" class="block">Primary Email Address</Label></div>
		<div class="col-span-6"><Input type="text" name="username" bind:value={username} /></div>
		<div class="col-span-2"><Label for="password" class="block">Password</Label></div>
		<div class="col-span-6"><Input type="password" name="password" bind:value={password} /></div>
		<div class="col-span-4">&nbsp;</div>
		<div class="col-span-2"><a href="#/register">Register/Lost Password</a></div>
		<div class="col-span-2"><Button type="submit">Login</Button></div>
	</div>
</form>
