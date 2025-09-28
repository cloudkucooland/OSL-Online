<script lang="ts">
	import { getContext } from 'svelte';
	import { getMeFromServer, getMeGiving } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	async function getBoth() {
		try {
			const m = await getMeFromServer();
			const g = await getMeGiving(m.Id);
			return { m, g };
		} catch (err) {
			toast.push('failed to load: ' + err.message);
			console.log(err);
			throw err;
		}
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Giving History</title>
</svelte:head>

{#await getBoth()}
	<h3>... loading ...</h3>
{:then { m, g }}
	<div>
		Giving History for : <a href="#/me">
			{m.Title}
			{m.FirstName}
			{m.LastName}</a
		>
	</div>

	<section>
		<div class="grid grid-cols-5 gap-4 px-4 py-2">
			<div class="col-span-1">Date</div>
			<div class="col-span-1">Amount</div>
			<div class="col-span-1">Description</div>
			<div class="col-span-1">Check Number</div>
			<div class="col-span-1">Paypal Transaction ID</div>
		</div>
		{#each g as gr}
			<div class="grid grid-cols-5 gap-4 px-4 py-2">
				<div class="col-span-1">{gr.Date}</div>
				<div class="col-span-1">{gr.Amount}</div>
				<div class="col-span-1">{gr.Description}</div>
				<div class="col-span-1">{gr.Check}</div>
				<div class="col-span-1">
					<a href="https://www.paypal.com/unifiedtransactions/details/payment/{gr.Transaction}"
						>{gr.Transaction}</a
					>
				</div>
			</div>
		{/each}
	</section>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}
