<script>
import Router from "svelte-spa-router";
import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Footer } from 'flowbite-svelte';
import { SvelteToast } from '@zerodevx/svelte-toast';
import { DarkMode } from 'flowbite-svelte';

import HomePage from "./routes/HomePage.svelte";
import Login from "./routes/Login.svelte";
import Downloads from "./routes/Downloads.svelte";
import Reports from "./routes/Reports.svelte";
import Member from "./routes/Member.svelte";
import AddMember from "./routes/AddMember.svelte";
import SubSearch from "./routes/SubSearch.svelte";
import Subscriber from "./routes/Subscriber.svelte";
import AddSubscriber from "./routes/AddSubscriber.svelte";

/* import "./app.css"; */

const routes = {
  "/": HomePage,
  "/Search/:query": HomePage,

  "/login": Login,
  "/download": Downloads,
  "/reports": Reports,

  "/member/:id": Member,
  "/addmember/": AddMember,

  "/subsearch": SubSearch,
  "/subscriber/:id": Subscriber,
  "/addsubscriber/": AddSubscriber,
  "*": HomePage,
};</script>

<Navbar >
  <NavBrand href="#/">
    <span class="dark:text-white">OSL Member Manager</span>
  </NavBrand>
  <NavHamburger></NavHamburger>
  <NavUl>
    <NavLi href="#/download">Mail Merge Lists</NavLi>
    <NavLi href="#/reports">Reports</NavLi>
  </NavUl>
</Navbar>
<SvelteToast></SvelteToast>

<main>
  <Router {routes}></Router>
</main>

<Footer class="sm:px-4 py-2.5 bottom-0 start-0 border-t">
  <p class="text-muted small text-end">
  <DarkMode />
    Copyright © <a href="https://saint-luke.net/">
      The Order of St. Luke ®
    </a> 2024. All Rights Reserved. &nbsp;|&nbsp;
    <a href="https://github.com/cloudkucooland/OSL-Online"> GitHub </a>
  </p>
</Footer>
