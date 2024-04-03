<!-- src/pages/Main.svelte -->

<script>
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import axios from "axios";
  import { Circle } from "svelte-loading-spinners";
  import { toast, SvelteToast } from "@zerodevx/svelte-toast";

  import Button from "../components/Button.svelte";
  import Ballot from "../layouts/Ballot.svelte";

  let loading = true;

  let userID = "";
  let activeElection = null;
  let registered = false;

  let voterID = "";
  let ballotID = "";

  let candidates = [];
  let selectedCandidate;
  let ballotCast = false;

  let activeElectionExists = false;

  onMount(() => {
    userID = sessionStorage.getItem("userID");
    document.title = "eVote Main";

    let activeElectionData = localStorage.getItem("activeElection");
    if (activeElectionData === null) {
      return;
    }

    activeElection = JSON.parse(activeElectionData);
    activeElectionExists = true;

    if (activeElection.ElectionID == null || activeElection.ElectionID === "") {
      activeElectionExists = false;
      return;
    }

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
    loading = true;

    try {
      console.log(selectedCandidate);

      await axios.post(
        "https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/submit-vote",
        {
          VoterID: voterID,
          BallotID: ballotID,
          CandidateID: selectedCandidate,
        },
      );

      toast.push("Success!", {
        theme: {
          "--toastColor": "mintcream",
          "--toastBackground": "rgba(72,187,120,0.9)",
          "--toastBarBackground": "#2F855A",
        },
      });

      await fetchBallot();
    } catch (error) {
      toast.push("Failed to submit vote", {
        theme: {
          "--toastColor": "white",
          "--toastBackground": "#ed1a3b",
          "--toastBarBackground": "#A01227",
        },
      });
    }

    loading = false;
  };
</script>

<main>
  <SvelteToast />

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
    {:else if activeElectionExists}
      <h2 align="center">
        Register to cast your vote for {activeElection.Name}!
      </h2>
      <Button label="Register" onClick={register} />
    {:else}
      <h2 align="center">No active election set!</h2>
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
