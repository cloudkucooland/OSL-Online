<script>
	import { setContext, getContext } from 'svelte';
	import { writable } from 'svelte/store';
	import Router from 'svelte-spa-router';
	import {
		Footer,
		FooterCopyright,
		Navbar,
		NavBrand,
		NavLi,
		NavUl,
		NavHamburger,
		Dropdown,
		DropdownItem
	} from 'flowbite-svelte';
	import { SvelteToast } from '@zerodevx/svelte-toast';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';

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
	import MeGiving from './routes/MeGiving.svelte';
	import Giving from './routes/Giving.svelte';
	import Changelog from './routes/Changelog.svelte';
	import ChapterBrowser from './routes/ChapterBrowser.svelte';
	import ChapterEditor from './routes/ChapterEditor.svelte';
	import Leadership from './routes/Leadership.svelte';
	import { getMe } from './oo';
	import LocalityBrowser from './routes/LocalityBrowser.svelte';
	import Email from './routes/Email.svelte';
	import Dashboard from './routes/Dashboard.svelte';
	import Notes from './routes/Notes.svelte';
	import Necrology from './routes/Necrology.svelte';

	const _init = getMe();
	setContext('oo', { me: writable(_init) });
	const { me } = getContext('oo');

	const routes = {
		'/': HomePage,
		'/search/:query': HomePage,

		'/login': Login,
		'/register': Register,
		'/reports': Reports,

		'/me': Me,
		'/mygiving': MeGiving,
		'/member/:id': Member,
		'/addmember/': AddMember,
		'/giving/:id': Giving,
		'/notes/:id': Notes,
		'/changelog/:id': Changelog,

		'/subsearch': SubSearch,
		'/subscriber/:id': Subscriber,
		'/addsubscriber/': AddSubscriber,
		'/chapterbrowser/': ChapterBrowser,
		'/chapterbrowser/:id': ChapterBrowser,
		'/chaptereditor/': ChapterEditor,
		'/chaptereditor/:id': ChapterEditor,
		'/leadership/': Leadership,
		'/leadership/:id': Leadership,
		'/localitybrowser/': LocalityBrowser,
		'/localitybrowser/:loc': LocalityBrowser,
		'/email/': Email,
		'/dashboard/': Dashboard,
		'/necrology/': Necrology,
		'*': HomePage
	};
</script>

<svelte:window />
<header class="w-full flex-none bg-white">
	<Navbar>
		<NavBrand href="#/">
			<span class="self-center text-xl whitespace-nowrap">OSL Member Directory</span>
		</NavBrand>
		<NavHamburger />
		<NavUl>
			{#if $me && $me.sub}
				<NavLi class="cursor-pointer">
					Me <ChevronDownOutline class="text-primary-800 ms-2 inline h-6 w-6" />
				</NavLi>
				<Dropdown class="z-20 w-44">
					<DropdownItem href="#/me">My Data</DropdownItem>
					<DropdownItem href="#/mygiving">My Giving</DropdownItem>
					<DropdownItem href="#/Login">Log out</DropdownItem>
				</Dropdown>

				<NavLi class="cursor-pointer">
					Lists <ChevronDownOutline class="text-primary-800 ms-2 inline h-6 w-6" />
				</NavLi>
				<Dropdown class="z-20 w-44">
					<DropdownItem href="#/chapterbrowser">Chapters</DropdownItem>
					<DropdownItem href="#/localitybrowser">Localities</DropdownItem>
					<DropdownItem href="#/leadership">Leadership</DropdownItem>
					<DropdownItem href="#/necrology">Necrology</DropdownItem>
				</Dropdown>

				{#if $me.level >= 1}
					<NavLi href="#/reports">Reports</NavLi>
				{/if}

				{#if $me.level >= 2}
					<NavLi class="cursor-pointer">
						Admin Tools<ChevronDownOutline class="text-primary-800 ms-2 inline h-6 w-6" />
					</NavLi>
					<Dropdown>
						<DropdownItem href="#/email">Email membership</DropdownItem>
						<DropdownItem href="#/subsearch">Subscribers</DropdownItem>
						<DropdownItem href="#/addmember">Add Member/Friend</DropdownItem>
						<DropdownItem href="#/Dashboard">Dashboard</DropdownItem>
						<DropdownItem href="#/chaptereditor">Chapter Editor</DropdownItem>
					</Dropdown>
				{/if}
			{/if}
		</NavUl>
	</Navbar>
	<SvelteToast />
</header>

<main>
	<Router {routes}></Router>
</main>

<Footer class="start-0 bottom-0 border-t py-2.5 sm:px-4">
	<FooterCopyright href="https://saint-luke.net/" by="The Order of St. Luke ®" year={2025} />
</Footer>
