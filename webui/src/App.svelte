<script>
	import { setContext, getContext } from 'svelte';
	import { writable } from 'svelte/store';
	import Router from 'svelte-spa-router';
	import {
		Navbar,
		NavBrand,
		NavLi,
		NavUl,
		NavHamburger,
		Footer,
		FooterCopyright,
		FooterLink
	} from 'flowbite-svelte';
	import { SvelteToast } from '@zerodevx/svelte-toast';
	import { DarkMode } from 'flowbite-svelte';

	import HomePage from './routes/HomePage.svelte';
	import Login from './routes/Login.svelte';
	import Reports from './routes/Reports.svelte';
	import Member from './routes/Member.svelte';
	import AddMember from './routes/AddMember.svelte';
	import SubSearch from './routes/SubSearch.svelte';
	import Subscriber from './routes/Subscriber.svelte';
	import AddSubscriber from './routes/AddSubscriber.svelte';
	import Register from './routes/Register.svelte';
	import Me from './routes/Me.svelte';
	import { getMe } from './oo';

	setContext('oo', { me: writable(getMe()) });
	const { me } = getContext('oo');
	// console.log("from context", me);

	const routes = {
		'/': HomePage,
		'/search/:query': HomePage,

		'/login': Login,
		'/register': Register,
		'/reports': Reports,

		'/me': Me,
		'/member/:id': Member,
		'/addmember/': AddMember,

		'/subsearch': SubSearch,
		'/subscriber/:id': Subscriber,
		'/addsubscriber/': AddSubscriber,
		'*': HomePage
	};
</script>

<svelte:window />
<header class="w-full flex-none bg-white dark:bg-slate-950">
	<Navbar>
		<NavBrand href="#/">
			<span class="whitspace-nowrap self-center text-xl dark:text-white">OSL Member Directory</span>
		</NavBrand>
		<NavHamburger></NavHamburger>
		<NavUl>
			<NavLi href="#/me">Me</NavLi>
			{#if $me && $me.level > 1}
				<NavLi href="#/reports">Reports</NavLi>{/if}
		</NavUl>
	</Navbar>
	<div class="flex gap-10"></div>
	<SvelteToast />
</header>

<main>
	<Router {routes}></Router>
</main>

<Footer class="bottom-0 start-0 border-t py-2.5 sm:px-4">
	<FooterCopyright href="/" by="The Order of St. Luke ®" year={2024} />
	{#if $me}<FooterLink href="#/Login">Log out</FooterLink>{/if}
	<DarkMode />
</Footer>
