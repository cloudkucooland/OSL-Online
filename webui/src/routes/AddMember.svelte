<script lang="ts">
	import { getContext } from 'svelte';
	import { Label, Input, Button } from 'flowbite-svelte';
	import { push } from 'svelte-spa-router';
	import { createMember } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	let firstname = $state('first');
	let lastname = $state('last');

	async function create(event) {
		event.preventDefault();
		event.stopPropagation();

		try {
			const id = await createMember(firstname, lastname);
			toast.push(`created`);
			push(`/member/${id}`);
		} catch (err) {
			console.log(err);
			toast.push('failed to create: ' + err.message);
		}
		return true;
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Create member</title>
</svelte:head>

<form onsubmit={create}>
	<section>
		<div class="grid grid-cols-2 gap-4 px-4 py-2">
			<div class="col-span-1">
				<Label for="FirstName" class="block">First Name</Label>
				<Input id="FirstName" bind:value={firstname} />
			</div>
			<div class="col-span-1">
				<Label for="LastName" class="block">Last Name</Label>
				<Input id="LastName" bind:value={lastname} />
			</div>
			<div class="col-span-2">
				<Button type="submit" color="green">Add</Button>
			</div>
		</div>
	</section>
</form>
