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
	<div class="grid grid-cols-8 gap-4 px-4 py-2">
		{#if !submitted}
			<div class="col-span-8">
				Use the primary email address as listed in our system. Other email addresses will be
				ignored.
			</div>
			<div class="col-span-2">
				<Label for="username" class="block">Primary Email Address</Label>
			</div>
			<div class="col-span-6"><Input type="text" name="username" bind:value={username} /></div>
			<div class="col-span-6">&nbsp;</div>
			<div class="col-span-2"><Button type="submit">Register/Recover</Button></div>
		{:else}
			<div class="col-span-8">Request submitted. Please check your email for your password.</div>
		{/if}
	</div>
</form>
