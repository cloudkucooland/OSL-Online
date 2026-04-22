<script lang="ts">
	import { getContext } from 'svelte';
	import { Label, Input, Button, Card, Heading } from 'flowbite-svelte';
	import { push } from 'svelte-spa-router';
	import { UserAddOutline } from 'flowbite-svelte-icons';
	import { createMember } from '../oo'; // Replace with createSubscriber if added to oo.js
	import { toast } from '@zerodevx/svelte-toast';

	const oo = getContext('oo');

	// Guard: Redirect if not logged in
	if (!oo.me) {
		push('/Login');
	}

	let firstname = $state('');
	let lastname = $state('');
	let processing = $state(false);

	async function handleCreate(event: Event) {
		event.preventDefault();

		if (!firstname || !lastname) {
			toast.push('First and Last name are required');
			return;
		}

		processing = true;
		try {
			// Using createMember as the placeholder since subscribers
			// usually start as a basic person record
			const id = await createMember(firstname, lastname);
			toast.push(`Subscriber record created`);
			push(`/subscriber/${id}`);
		} catch (err: any) {
			console.error(err);
			toast.push('Failed to create: ' + err.message);
		} finally {
			processing = false;
		}
	}
</script>

<svelte:head>
	<title>Add Subscriber | OSL Directory</title>
</svelte:head>

<div class="mx-auto max-w-2xl px-4 py-8">
	<Card size="md" padding="xl" class="border-t-4 border-t-indigo-500 shadow-sm">
		<div class="mb-6 flex items-center gap-3">
			<div class="rounded-lg bg-indigo-100 p-2">
				<UserAddOutline class="h-6 w-6 text-indigo-600" />
			</div>
			<div>
				<Heading tag="h2" class="text-2xl font-bold text-slate-900">Add Doxology Subscriber</Heading
				>
				<p class="text-sm text-slate-500">Create a new record for a non-member subscriber.</p>
			</div>
		</div>

		<form onsubmit={handleCreate} class="space-y-6">
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="FirstName">First Name</Label>
					<Input id="FirstName" placeholder="First Name" bind:value={firstname} required />
				</div>
				<div class="space-y-2">
					<Label for="LastName">Last Name</Label>
					<Input id="LastName" placeholder="Last Name" bind:value={lastname} required />
				</div>
			</div>

			<div class="flex justify-end gap-3 border-t border-slate-100 pt-4">
				<Button color="alternative" onclick={() => push('/subsearch')}>Cancel</Button>
				<Button type="submit" color="primary" disabled={processing}>
					{processing ? 'Creating...' : 'Add Subscriber'}
				</Button>
			</div>
		</form>
	</Card>
</div>
