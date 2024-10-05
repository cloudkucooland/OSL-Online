<script lang="ts">
	import { getContext } from 'svelte';
	import { getMe, getMember, getGiving, postGiving } from '../oo';
	import { Label, Input, Button, Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	export let params;
	const dd = new Date().toISOString().split('T');
	let postdate = dd[0];
	let amount;
	let description = 'Annual Reaffirmation';
	let check = 0;
	let transaction = 0;

	async function add() {
		console.log(params.id, postdate, amount, description, check, transaction);
		try {
			await postGiving(params.id, postdate, amount, description, check, transaction);
			toast.push(`Posted`);
			push(`#/giving/${params.id}`);
			return true;
		} catch (err) {
			toast.push('failed to change: ' + err);
			console.log(err);
		}
	}

	async function getBoth(id) {
		try {
			const m = await getMember(id);
			const g = await getGiving(id);
			return { m, g };
		} catch (err) {
			toast.push('failed to load: ' + err);
			console.log(err);
			throw err;
		}
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Member</title>
</svelte:head>

{#await getBoth(params.id)}
	<h3>... loading ...</h3>
{:then { m, g }}
	<div>
		Giving History for : <a href="#/member/{m.ID}">
			{m.Title}
			{m.FirstName}
			{m.LastName}</a
		>
		( {m.MemberStatus} )
	</div>

	<section>
		<div class="grid grid-cols-5 gap-4 px-4 py-2">
			<div class="col-span-1">Date</div>
			<div class="col-span-1">Amount</div>
			<div class="col-span-1">Description</div>
			<div class="col-span-1">Check Number</div>
			<div class="col-span-1">Transaction ID</div>
		</div>
		{#each g as gr}
			<div class="grid grid-cols-5 gap-4 px-4 py-2">
				<div class="col-span-1">{gr.Date}</div>
				<div class="col-span-1">{gr.Amount}</div>
				<div class="col-span-1">{gr.Description}</div>
				<div class="col-span-1">{gr.Check}</div>
				<div class="col-span-1">{gr.Transaction}</div>
			</div>
		{/each}
	</section>

	<section>
		<form on:submit={add}>
			<div class="grid grid-cols-5 gap-4 px-4 py-2">
				<div class="col-span-1">
					<Label for="Date" class="block">Date</Label>
					<Input id="Date" bind:value={postdate} />
				</div>
				<div class="col-span-1">
					<Label for="Amount" class="block">Amount</Label>
					<Input id="Amount" bind:value={amount} />
				</div>
				<div class="col-span-1">
					<Label for="Description" class="block">Description</Label>
					<Input id="Description" bind:value={description} />
				</div>
				<div class="col-span-1">
					<Label for="Check" class="block">Check Number (if check)</Label>
					<Input id="Check" bind:value={check} />
				</div>
				<div class="col-span-1">
					<Label for="Transaction" class="block">Transaction (if paypal)</Label>
					<Input id="Transaction" bind:value={transaction} />
				</div>
				<div class="col-span-4">&nbsp;</div>
				<div class="col-span-1">
					<Button type="submit">Add Giving Record</Button>
				</div>
			</div>
		</form>
	</section>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}
