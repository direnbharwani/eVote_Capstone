<script>
	import { Router, Route } from 'svelte-routing';

	import Login from './pages/Login.svelte';
	import Main from './pages/Main.svelte';
	import Admin from './pages/Admin.svelte';
</script>

<Router>
	<Route path="/" component={Login}/>
	<Route path="/main" component={Main}/>
	<Route path="/admin" component={Admin}/>
</Router>