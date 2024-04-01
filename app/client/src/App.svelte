<script>
	import { Router, Route } from 'svelte-routing';

	import Login from './pages/Login.svelte';
	import Main from './pages/Main.svelte';
</script>

<Router>
	<Route path="/" component={Login}/>
	<Route path="/main" component={Main}/>
</Router>