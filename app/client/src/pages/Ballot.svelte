<!-- src/pages/Login.svelte -->

<script>
    import { onMount } from 'svelte';

    let userID;

    onMount(() => {
        userID = sessionStorage.getItem("userID");
    })

</script>

<main>
    <body>
        <h1 align="center">Hello {userID}!</h1>
    </body>
</main>