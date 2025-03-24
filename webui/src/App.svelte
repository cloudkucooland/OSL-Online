<script>
	import { setContext, getContext, onMount } from 'svelte';
	import { writable, readable } from 'svelte/store';
	import Router from 'svelte-spa-router';
	import {
		DarkMode,
		Footer,
		FooterCopyright,
		FooterLink,
		Navbar,
		NavBrand,
		NavLi,
		NavUl,
		NavHamburger
	} from 'flowbite-svelte';
	import { SvelteToast } from '@zerodevx/svelte-toast';

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
	import Giving from './routes/Giving.svelte';
	import Changelog from './routes/Changelog.svelte';
	import ChapterBrowser from './routes/ChapterBrowser.svelte';
	import Leadership from './routes/Leadership.svelte';
	import { getMe, getChapters } from './oo';
	import LocalityBrowser from './routes/LocalityBrowser.svelte';
	import Email from './routes/Email.svelte';

	const _init = getMe();
	setContext('oo', { me: writable(_init), chapters: readable(getChapters()) });
	const { me, chapters } = getContext('oo');

	const routes = {
		'/': HomePage,
		'/search/:query': HomePage,

		'/login': Login,
		'/register': Register,
		'/reports': Reports,

		'/me': Me,
		'/member/:id': Member,
		'/addmember/': AddMember,
		'/giving/:id': Giving,
		'/changelog/:id': Changelog,

		'/subsearch': SubSearch,
		'/subscriber/:id': Subscriber,
		'/addsubscriber/': AddSubscriber,
		'/chapterbrowser/': ChapterBrowser,
		'/chapterbrowser/:id': ChapterBrowser,
		'/leadership/': Leadership,
		'/leadership/:id': Leadership,
		'/localitybrowser/': LocalityBrowser,
		'/localitybrowser/:loc': LocalityBrowser,
		'/email/': Email,
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
			<NavLi href="#/chapterbrowser">Chapters</NavLi>
			<NavLi href="#/localitybrowser">Localities</NavLi>
			<NavLi href="#/leadership">Leadership</NavLi>
			{#if $me && $me.level >= 1}
				<NavLi href="#/reports">Reports</NavLi>{/if}
			{#if $me && $me.level == 2}
				<NavLi href="#/email">Email</NavLi>
			{/if}
			{#if $me}<NavLi href="#/Login">Log out</NavLi>{/if}
		</NavUl>
	</Navbar>
	<div class="flex gap-10"></div>
	<SvelteToast />
</header>

<main>
	<Router {routes}></Router>
</main>

<Footer class="start-0 bottom-0 border-t py-2.5 sm:px-4">
	<FooterCopyright href="/" by="The Order of St. Luke ®" year={2025} />
	<!-- <DarkMode /> -->
</Footer>
