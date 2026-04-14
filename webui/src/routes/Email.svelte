<script lang="ts">
	import { getContext } from 'svelte';
	import { Input, Label, Textarea, Button, Select, Card, Heading, Badge } from 'flowbite-svelte';
	import { push } from 'svelte-spa-router';
	import {
		EnvelopeOutline,
		InfoCircleOutline,
		ExclamationCircleOutline
	} from 'flowbite-svelte-icons';
	import { sendemail } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	const oo = getContext('oo');

	// Guard: Redirect if not logged in
	if (!oo.me) {
		push('/Login');
	}

	const towhom = [
		{ value: 'nobody@localhost', name: 'Select Recipient Group...' },
		{ value: 'all', name: 'All Vowed Members' },
		{ value: 'annual', name: 'Annual Vowed Members' },
		{ value: 'life', name: 'Life Vowed Members' },
		{ value: 'friends', name: 'Friends' }
	];

	let whom = $state('nobody@localhost');
	let content = $state('');
	let subject = $state('OSL Announcement');
	let sending = $state(false);

	async function handleSend(event: Event) {
		event.preventDefault();

		if (whom === 'nobody@localhost') {
			toast.push('Please select a recipient group');
			return;
		}

		if (!content.trim()) {
			toast.push('Email body cannot be empty');
			return;
		}

		if (!confirm(`Are you sure you want to send this to ${whom.toUpperCase()}?`)) {
			return;
		}

		sending = true;
		try {
			await sendemail(whom, subject, content);
			toast.push('Broadcast message sent successfully');
			content = '';
			subject = 'OSL Announcement';
			whom = 'nobody@localhost';
		} catch (err: any) {
			console.error(err);
			toast.push('Failed to send: ' + err.message);
		} finally {
			sending = false;
		}
	}
</script>

<svelte:head>
	<title>Mass Mailer | OSL Directory</title>
</svelte:head>

<div class="mx-auto max-w-4xl px-4 py-8">
	<header class="mb-6">
		<Heading tag="h2" class="flex items-center gap-3 text-3xl font-bold text-slate-900">
			<EnvelopeOutline class="h-8 w-8 text-primary-600" />
			Broadcast Email
		</Heading>
		<p class="mt-2 text-slate-500">Send a system-generated email to specific membership groups.</p>
	</header>

	<form onsubmit={handleSend} class="space-y-6">
		<Card size="none" class="border-l-4 border-l-amber-400 bg-amber-50/30 p-6">
			<div class="flex items-start gap-3">
				<InfoCircleOutline class="mt-0.5 h-6 w-6 text-amber-600" />
				<div class="text-sm text-amber-800">
					<strong>Note:</strong> Emails include a standard header and footer. The recipient's name
					will be automatically prepended as <em>"Hi [Name]"</em>.
				</div>
			</div>
		</Card>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-4">
			<div class="space-y-4 md:col-span-4">
				<div>
					<Label for="whom" class="mb-2">Recipients</Label>
					<Select id="whom" items={towhom} bind:value={whom} class="bg-white" />
				</div>
				<div>
					<Label for="subject" class="mb-2">Subject Line</Label>
					<Input id="subject" bind:value={subject} placeholder="Enter email subject..." />
				</div>
			</div>

			<div class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm md:col-span-4">
				<div class="mb-4 border-b border-slate-100 pb-4">
					<span class="font-serif text-lg italic text-slate-400">Hi [sibling name],</span>
				</div>

				<Textarea
					id="content"
					rows={12}
					bind:value={content}
					placeholder="Type your message here..."
					class="border-none p-0 font-sans text-base focus:ring-0"
				/>

				<div class="mt-4 space-y-1 border-t border-slate-100 pt-4 font-serif italic text-slate-500">
					<p>Living the Sacramental Life,</p>
					<p class="font-bold">The Order of Saint Luke</p>
				</div>
			</div>

			<div class="flex items-center justify-between md:col-span-4">
				{#if whom !== 'nobody@localhost'}
					<Badge color="red" class="flex items-center gap-1">
						<ExclamationCircleOutline class="h-4 w-4" />
						Sending to: {towhom.find((t) => t.value === whom)?.name}
					</Badge>
				{:else}
					<div></div>
				{/if}

				<Button
					type="submit"
					color="red"
					size="xl"
					disabled={sending || whom === 'nobody@localhost'}
				>
					{sending ? 'Dispatching...' : 'Send Broadcast'}
				</Button>
			</div>
		</div>
	</form>
</div>
