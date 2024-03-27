<script>
	import { Router, Route } from 'svelte-routing';

	import Login from './pages/Login.svelte';
	import Ballot from './pages/Ballot.svelte';
</script>

<Router>
	<Route path="/" component={Login}/>
	<Route path="/ballot" component={Ballot}/>
</Router>