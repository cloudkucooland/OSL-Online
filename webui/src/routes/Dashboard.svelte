<script lang="ts">
	import { getContext } from 'svelte';
	import { Table, TableBody, TableBodyCell, TableBodyRow } from 'flowbite-svelte';
	import { getDashboard } from '../oo';
	import { push } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Dashboard</title>
</svelte:head>

{#await getDashboard()}
	<h3>... loading ...</h3>
{:then d}
	<Table>
		<TableBody>
			<TableBodyRow>
				<TableBodyCell>Life Vow Members</TableBodyCell>
				<TableBodyCell>{d.LifevowCount}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Life Vow members who gave (Since July 1)</TableBodyCell>
				<TableBodyCell>{d.LifeVowsWhoGave}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Annual Vow Members</TableBodyCell>
				<TableBodyCell>{d.AnnualCount}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Annual Vow members who gave (Since July 1)</TableBodyCell>
				<TableBodyCell>{d.AnnualVowsWhoGave}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Friends</TableBodyCell>
				<TableBodyCell>{d.FriendCount}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Institutional & non-member Doxology Subscribers (account current)</TableBodyCell>
				<TableBodyCell>{d.SubscriberCount}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>This year member giving (since July 1)</TableBodyCell>
				<TableBodyCell>{d.ThisYearGiving}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Last year member giving (July - July)</TableBodyCell>
				<TableBodyCell>{d.LastYearGiving}</TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
{/await}
