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
		<a href="#/reports" on:click={() => report('expired')}
			>"Active" but not renewed for 2 years or longer</a
		>
	</li>
	<li><a href="#/reports" on:click={() => report('email')}>All Member Email List</a></li>
	<li><a href="#/reports" on:click={() => report('annual')}>Annual Vows Mail Merge</a></li>
	<li><a href="#/reports" on:click={() => report('life')}>Life Vows Mail Merge</a></li>
	<li>
		<a href="#/reports" on:click={() => report('doxprint')}
			>Doxology Mailing List (member + subscriber)</a
		>
	</li>
	<li>
		<a href="#/reports" on:click={() => report('doxemail')}
			>Doxology Email List (member + subscriber)</a
		>
	</li>
	<li>
		<a href="#/reports" on:click={() => report('fontemail')}
			>Font Email List (for google groups update)</a
		>
	</li>
	<li>
		<a href="#/reports" on:click={() => report('allsubscribers')}
			>All Subscribers</a
		>
	</li>
</ul>
