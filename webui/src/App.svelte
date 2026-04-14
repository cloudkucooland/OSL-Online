<script>
	import { setContext, onMount } from 'svelte';
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

	// Route Imports
	import HomePage from './routes/HomePage.svelte';
	import Login from './routes/Login.svelte';
	import Logout from './routes/Logout.svelte';
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
	import LocalityBrowser from './routes/LocalityBrowser.svelte';
	import Email from './routes/Email.svelte';
	import Dashboard from './routes/Dashboard.svelte';
	import Notes from './routes/Notes.svelte';
	import Necrology from './routes/Necrology.svelte';

	import { getMe } from './oo';

	let oo = $state({
		// oo.me is the JSON token
		me: getMe(),
		chaptercache: null, // only load the chapter list once per session
		namecache: null // resolve ID => names
	});
	setContext('oo', oo);

	/* onMount(async () => {
		// prime the caches here, not needed for chaptercache
	}); */

	const routes = {
		'/': HomePage,
		'/search/:query': HomePage,
		'/login': Login,
		'/logout': Logout,
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

<div class="flex min-h-screen flex-col bg-gray-50 text-slate-900">
	<header
		class="sticky top-0 z-40 w-full flex-none border-b border-slate-200 bg-white/80 backdrop-blur-md"
	>
		<Navbar containerClass="mx-auto flex flex-wrap justify-between items-center px-4 py-2">
			<NavBrand href="#/">
				<span class="self-center text-2xl font-semibold tracking-tight text-primary-900">
					OSL Directory
				</span>
			</NavBrand>

			<NavHamburger />

			<NavUl class="transition-all">
				{#if oo.me && oo.me.level}
					<NavLi class="group flex cursor-pointer items-center gap-1">
						Me <ChevronDownOutline size="sm" class="group-hover:text-primary-600" />
					</NavLi>
					<Dropdown class="z-50 w-44">
						<DropdownItem href="#/me">My Data</DropdownItem>
						<DropdownItem href="#/mygiving">My Giving</DropdownItem>
						<DropdownItem href="#/Logout" class="mt-1 border-t text-red-600">Log out</DropdownItem>
					</Dropdown>

					<NavLi class="group flex cursor-pointer items-center gap-1">
						Lists <ChevronDownOutline size="sm" class="group-hover:text-primary-600" />
					</NavLi>
					<Dropdown class="z-50 w-44">
						<DropdownItem href="#/chapterbrowser">Chapters</DropdownItem>
						<DropdownItem href="#/localitybrowser">Localities</DropdownItem>
						<DropdownItem href="#/leadership">Leadership</DropdownItem>
						<DropdownItem href="#/necrology">Necrology</DropdownItem>
						<DropdownItem href="#/subsearch">Doxology Subscribers</DropdownItem>
					</Dropdown>

					{#if oo.me.level > 0}
						<NavLi href="#/reports">Reports</NavLi>
					{/if}

					{#if oo.me.level > 2}
						<NavLi class="group flex cursor-pointer items-center gap-1">
							Admin <ChevronDownOutline size="sm" class="group-hover:text-primary-600" />
						</NavLi>
						<Dropdown class="z-50 w-48">
							<DropdownItem href="#/email">Email Membership</DropdownItem>
							<DropdownItem href="#/addmember">Add Member/Friend</DropdownItem>
							<DropdownItem href="#/Dashboard">System Dashboard</DropdownItem>
							<DropdownItem href="#/chaptereditor">Chapter Editor</DropdownItem>
						</Dropdown>
					{/if}
				{/if}
			</NavUl>
		</Navbar>
	</header>

	<main class="container mx-auto flex-grow px-4 py-8">
		<Router {routes} basepath="/oo"></Router>
	</main>

	<Footer class="mt-auto border-t border-slate-200 bg-white py-6">
		<div
			class="container mx-auto flex flex-col items-center justify-between gap-4 px-4 md:flex-row"
		>
			<FooterCopyright
				href="https://saint-luke.net/"
				by="The Order of St. Luke ®"
				year={new Date().getFullYear()}
			/>
			<div class="text-sm text-slate-500">Living the Sacramental Life</div>
		</div>
	</Footer>
</div>

<SvelteToast />

<style>
	:global(body) {
		font-family:
			'Inter',
			system-ui,
			-apple-system,
			sans-serif;
	}
</style>
