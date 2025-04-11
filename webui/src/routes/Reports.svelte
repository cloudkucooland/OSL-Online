<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { report } from '../oo';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Reports</title>
</svelte:head>

<ul>
	<li>
		<a href="#/reports" onclick={() => report('expired')}
			>"Active" but not renewed for 2 years or longer</a
		>
	</li>
	<li><a href="#/reports" onclick={() => report('email')}>All Member Email List</a></li>
	<li><a href="#/reports" onclick={() => report('avery')}>Avery Labels</a></li>
	<li><a href="#/reports" onclick={() => report('annual')}>Annual Vows Mail Merge</a></li>
	<li><a href="#/reports" onclick={() => report('life')}>Life Vows Mail Merge</a></li>
	<li>
		<a href="#/reports" onclick={() => report('doxprint')}
			>Doxology Mailing List (member + subscriber)</a
		>
	</li>
	<li>
		<a href="#/reports" onclick={() => report('allsubscribers')}>All Subscribers</a>
	</li>
	<li><a href="#/reports" onclick={() => report('barb')}>Sr. Barb's Monthly Prayer list</a></li>
</ul>
