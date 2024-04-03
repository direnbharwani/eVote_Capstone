<!-- src/pages/Admin.svelte -->

<script>
  import CreateElection from "../layouts/CreateElection.svelte";
  import SetElection from "../layouts/SetActiveElection.svelte";
  import CountVotes from "../layouts/CountVotes.svelte";

  export let layout = 0;
</script>

<main>
  <body>
    <nav>
      <ul>
        <li>
          <a href="/admin" on:click|preventDefault={() => (layout = 0)}
            >Create Sample Election</a
          >
        </li>
        <li>
          <a href="/admin" on:click|preventDefault={() => (layout = 1)}
            >Set Active Election</a
          >
        </li>
        <li>
          <a href="/admin" on:click|preventDefault={() => (layout = 2)}>Count Votes</a>
        </li>
      </ul>
    </nav>

    <div id="content-body">
      {#if layout === 0}
        <CreateElection />
      {:else if layout === 1}
        <SetElection />
      {:else if layout === 2}
        <CountVotes />
      {:else}
        <h1>Page Not Found</h1>
      {/if}
    </div>
  </body>
</main>

<style>
  * {
    box-sizing: border-box;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  body {
    padding: 0;
    margin: 0;
    font-family: "Poppins", sans-serif;
  }

  nav {
    padding: 5px 5%;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    box-shadow:
      rgba(50, 50, 93, 0.25) 0px 2px 5px -1px,
      rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
    z-index: 1;
    background-color: #a01227;
  }

  nav ul {
    list-style: none;
    display: flex;
  }

  nav ul li {
    margin-left: 1.5rem;
  }
  nav ul li a {
    text-decoration: none;
    color: #fff;
    font-size: 95%;
    font-weight: 400;
    padding: 4px 8px;
    border-radius: 5px;
  }

  nav ul li a:hover {
    background-color: #76212e;
  }

  #content-body {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    height: 100vh;
    margin: auto;
  }
</style>
