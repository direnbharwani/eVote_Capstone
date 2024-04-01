<!-- src/components/Button.svelte -->

<script>
  import { onMount } from "svelte";
  import axios from "axios";
  import { Circle } from "svelte-loading-spinners";

  import Button from "../components/Button.svelte";
  import Ballot from "../layouts/Ballot.svelte";

  let loading = true;

  let userID = "";
  let activeElection = JSON.parse(localStorage.getItem("activeElection"));
  let registered = false;

  let voterID = "";
  let ballotID = "";

  let candidates = [];
  let ballotCast = false;

  onMount(() => {
    userID = sessionStorage.getItem("userID");
    document.title = "eVote Main";

    fetchBallot();
  });

  async function fetchBallot() {
    loading = true;

    try {
      const response = await axios.post(
        "https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/read-vote",
        {
          NRIC: userID,
          ElectionID: activeElection.ElectionID,
        },
      );

      registered = true;

      voterID = response.data.VoterID;
      ballotID = response.data.BallotID;
      candidates = response.data.Candidates;
    } catch (error) {
      console.error("No record found: ", error);
    }

    loading = false;
  }

  async function register() {
    loading = true;

    try {
      const response = await axios.post(
        "https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/register",
        {
          NRIC: userID,
          ElectionID: activeElection.ElectionID,
        },
      );

      await fetchBallot();
    } catch (error) {
      console.error("Failed to register: ", error);
    }

    loading = false;
  }
</script>

<main>
  <body>
    {#if loading}
      <Circle size="60" color="#a01227" unit="px" duration="1s" />
    {:else if registered}
      <Ballot bind:voterID bind:ballotID bind:candidates bind:ballotCast />

      {#if !ballotCast && candidates.length > 0}
          <Button label="Submit Vote" />
      {/if}
    {:else}
      <h2 align="center">
        Register to cast your vote for {activeElection.Name}!
      </h2>
      <Button label="Register" onClick={register} />
    {/if}
  </body>
</main>

<style>
  body {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    height: 100vh;
    margin: auto;
    background-color: #f4f4f4;
  }
</style>
