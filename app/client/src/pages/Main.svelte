<!-- src/components/Button.svelte -->

<script>
  import { onMount } from "svelte";
  import axios from "axios";
  import { Circle } from "svelte-loading-spinners";
  import { navigate } from "svelte-routing";

  import Button from "../components/Button.svelte";
  import Ballot from "../layouts/Ballot.svelte";
  import {
    showAlert,
    alertMessage,
    shouldRedicrect,
    setAlert,
  } from "../stores/alertStore";

  let loading = true;

  let userID = "";
  let activeElection = JSON.parse(localStorage.getItem("activeElection"));
  let registered = false;

  let voterID = "";
  let ballotID = "";

  let candidates = [];
  let selectedCandidate;
  let ballotCast = false;

  let redirect = false;

  onMount(() => {
    userID = sessionStorage.getItem("userID");
    document.title = "eVote Main";

    fetchBallot();
  });

  const fetchBallot = async () => {
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
  };

  const register = async () => {
    loading = true;

    try {
      await axios.post(
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
  };

  const submitVote = async () => {
    try {
      console.log(selectedCandidate);

      await axios.post(
        "https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/submit-vote",
        {
          VoterID: userID,
          BallotID: ballotID,
          CandidateID: "",
        },
      );

      alert("Thank you for submitting your vote!");
    } catch (error) {
      alert("Submission failed, please try again!");
    }
  };
</script>

<main>
  {#if $showAlert}
    <div ckass="alert">
      <p>{$alertMessage}</p>
      <button on:click={() => showAlert.set(false)}>Close</button>
    </div>
  {/if}

  <body>
    {#if loading}
      <Circle size="60" color="#a01227" unit="px" duration="1s" />
    {:else if registered}
      <Ballot
        bind:voterID
        bind:ballotID
        bind:candidates
        bind:ballotCast
        bind:selectedCandidate
      />

      {#if !ballotCast && candidates.length > 0}
        <Button label="Submit Vote" onClick={submitVote} />
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
