<!-- src/components/Button.svelte -->

<script>
    import { Link } from 'svelte-routing';

    export let label = "Button";
    export let type = "button";
    export let onClick;
    export let linkTo = "/";

    function handleClick() {
        if (onClick) {
            onClick();
        }
    }
</script>

<Link to={linkTo}>
    <button {type} on:click={handleClick}>
        {label}
    </button>
</Link>


<style>
    /* Add your button styles here */
    button {
        width: 50%;
        padding: 20px;
        margin-top: 25%;
        font-family: Josefin Sans;
        font-size: 1.5em;
        border-radius: 10px;
        overflow: hidden;
        cursor: pointer;
        background-color: #ed1a3b;
        color: #ffffff;
        transition: background-color 0.3s ease;
    }

    button:hover {
        background-color: #a01227;
    }
</style>
