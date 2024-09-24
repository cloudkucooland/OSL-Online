<script lang="ts">
import { push } from "svelte-spa-router";
import { getMe } from "../oo";

const me = getMe();
if (me === undefined) {
  push('/Login');
}
</script>

<svelte:head>
  <title>OSL Member Manager: Subscriber</title>
</svelte:head>

Subscriber

