<script lang="ts">
	import { getContext } from 'svelte';
	import { Label, Input, Button, Card, Heading } from 'flowbite-svelte';
	import { push } from 'svelte-spa-router';
	import { UserAddOutline } from 'flowbite-svelte-icons';
	import { createMember } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	const oo = getContext('oo');
	if (!oo.me) {
		push('/Login');
	}

	let firstname = $state('');
	let lastname = $state('');
	let processing = $state(false);

	async function handleCreate(event: Event) {
		event.preventDefault();

		if (!firstname || !lastname) {
			toast.push('Both names are required');
			return;
		}

		processing = true;
		try {
			const id = await createMember(firstname, lastname);
			toast.push(`Member created successfully`);
			push(`/member/${id}`);
		} catch (err: any) {
			console.error(err);
			toast.push('Failed to create: ' + err.message);
		} finally {
			processing = false;
		}
	}
</script>

<svelte:head>
	<title>Add New Member | OSL Directory</title>
</svelte:head>

<div class="mx-auto max-w-2xl px-4 py-8">
	<Card size="md" padding="xl" class="shadow-sm">
		<div class="mb-6 flex items-center gap-3">
			<div class="bg-primary-100 rounded-lg p-2">
				<UserAddOutline class="text-primary-600 h-6 w-6" />
			</div>
			<div>
				<Heading tag="h2" class="text-2xl font-bold text-slate-900">Add New Member</Heading>
				<p class="text-sm text-slate-500">Create a new base record for the directory.</p>
			</div>
		</div>

		<form onsubmit={handleCreate} class="space-y-6">
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="FirstName">First Name</Label>
					<Input id="FirstName" placeholder="e.g. Robert" bind:value={firstname} required />
				</div>
				<div class="space-y-2">
					<Label for="LastName">Last Name</Label>
					<Input id="LastName" placeholder="e.g. Smith" bind:value={lastname} required />
				</div>
			</div>

			<div class="flex justify-end gap-3 border-t border-slate-100 pt-4">
				<Button color="alternative" onclick={() => push('/')}>Cancel</Button>
				<Button type="submit" color="primary" disabled={processing}>
					{processing ? 'Creating...' : 'Create Member'}
				</Button>
			</div>
		</form>
	</Card>
</div>
