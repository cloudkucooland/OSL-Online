<script lang="ts">
	import { getMember, getChangelog } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	let { params } = $props();

	async function getBoth(id) {
		try {
			const m = await getMember(id);
			const c = await getChangelog(id);
			return { m, c };
		} catch (err) {
			toast.push('failed to load: ' + err.message);
			console.log(err);
			throw err;
		}
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Change Log</title>
</svelte:head>

{#await getBoth(params.id)}
	<h3>... loading ...</h3>
{:then { m, c }}
	<div>
		Change Log for : <a href="#/member/{m.ID}">
			{m.Title}
			{m.FirstName}
			{m.LastName}</a
		>
		( {m.MemberStatus} )
	</div>

	<section>
		<div class="grid grid-cols-5 gap-4 px-4 py-2">
			<div class="col-span-1">Date</div>
			<div class="col-span-1">Changer</div>
			<div class="col-span-1">Field</div>
			<div class="col-span-1">Value</div>
		</div>
		{#each c as cr}
			<div class="grid grid-cols-5 gap-4 px-4 py-2">
				<div class="col-span-1">{cr.Date}</div>
				<div class="col-span-1">{cr.Changer}</div>
				<div class="col-span-1">{cr.Field}</div>
				<div class="col-span-1">{cr.Value}</div>
			</div>
		{/each}
	</section>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}
