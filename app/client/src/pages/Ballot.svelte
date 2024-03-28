<!-- src/components/Button.svelte -->

<script>
    import { onMount } from "svelte";
    import axios from 'axios';

    import Button from "../components/Button.svelte";
    import CandidateList from "../components/CandidateList.svelte";

    let userID = "";
    let ballotID = "";

    let candidates = [];
    let ballotCast = false;

    onMount(async () => {
        userID = sessionStorage.getItem("userID");
        ballotID = sessionStorage.getItem("ballotID");

        
        await fetchBallot();
    });

    async function fetchBallot() {
        try {
            const response = await axios.post('https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/read-vote', {
                VoterID: userID,
	            BallotID: ballotID
            });

            candidates = response.data.Candidates;

        } catch (error) {
            console.error('No ballot found: ', error);
        }
    }

</script>

<main>
    <body>
        <div id="ballot-paper">
            <!-- Voter and Ballot IDs at the top of the ballot paper -->
            <div id="ids">
                <div class="id-item" id="voterID">
                    <b>VoterID:</b> <span>{userID}</span>
                </div>
                <div class="id-item" id="ballotID">
                    <b>BallotID:</b> <span>{ballotID}</span>
                </div>
            </div>

            {#if candidates.length > 0}
                <CandidateList {candidates} bind:ballotCast />
            {:else}
                <h4 align="center">No Candidates Found!</h4>
            {/if}
        </div>

        {#if !ballotCast && candidates.length > 0}
            <Button label="Submit Vote" />
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

    #ballot-paper {
        width: 75%;
        font-family: Consolas;
        border: 2px solid;
    }

    #ids {
        display: grid;
        grid-template-columns: auto 1fr; /* Label takes auto width, value takes the rest */
        grid-gap: 0; /* Space between rows */
        margin-bottom: 5%;
    }

    .id-item {
        display: flex;
        align-items: center;
    }

    .id-item b {
        width: 100px; /* Fixed width for the label */
        text-align: left; /* Align text to the right */
    }

    .id-item span {
        flex-grow: 1; /* Value takes the remaining space */
        padding-left: 10px; /* Add padding to the left of the value */
    }

    #voterID {
        grid-row: 1;
    }

    #ballotID {
        grid-row: 2;
    }
</style>
