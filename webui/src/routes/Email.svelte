<script lang="ts">
	import { getContext } from 'svelte';
	import { Input, Label, Textarea, Button, Select } from 'flowbite-svelte';
	import { push } from 'svelte-spa-router';
	import { sendemail } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	const towhom = [
		{ value: 'nobody@localhost', name: 'None' },
		{ value: 'all', name: 'All Vowed Members' },
		{ value: 'annual', name: 'Annual Vowed Members' },
		{ value: 'life', name: 'Life Vowed Members' },
		{ value: 'friends', name: 'Friends' }
	];

	let whom = $state('nobody@localhost');
	let content = $state('');
	let subject = $state('OSL Announcement');

	async function send(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			const id = await sendemail(whom, subject, content);
			toast.push(`Message Sent`);
			content = '';
			subject = 'OSL Announcment';
		} catch (err) {
			console.log(err);
			toast.push('failed to create: ' + err.message);
		}
		return true;
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Email</title>
</svelte:head>

<form onsubmit={send}>
	<section>
		<div class="grid grid-cols-4 gap-4 px-4 py-2">
			<div class="col-span-4">
				<Select id="whom" items={towhom} bind:value={whom} />
			</div>
			<div class="col-span-4">
				<Label>Subject</Label>
				<Input id="subject" bind:value={subject} />
			</div>
			<div class="col-span-4">
				<h3>Hi [sibling name]</h3>
				<Textarea id="content" rows={8} bind:value={content} />
				<h3>Living the Sacramental Life</h3>
				<h3>Yours truly,</h3>
				<h3>The Order of Saint Luke</h3>
			</div>
			<div class="col-span-3"></div>
			<div class="col-span-1">
				<Button type="submit" color="red">Send</Button>
			</div>
		</div>
	</section>
</form>
