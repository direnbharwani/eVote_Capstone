<script>
    import Button from "../components/Button.svelte";
    import CandidateList from "../components/CandidateList.svelte";

    export let voterID;
    export let ballotID;

    export let candidates = [];
    export let ballotCast;
</script>

<div id="ballot-paper">
    <!-- Voter and Ballot IDs at the top of the ballot paper -->
    <div id="ids">
        <div class="id-item" id="voterID">
            <b>VoterID:</b> <span>{voterID}</span>
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

<style>
    #ballot-paper {
        width: 75%;
        border: 2px solid;
    }

    #ids {
        display: grid;
        grid-template-columns: auto 1fr; /* Label takes auto width, value takes the rest */
        grid-gap: 0; /* Space between rows */
        margin-bottom: 5%;
        font-family: JetBrains Mono;
    }

    .id-item {
        display: flex;
        align-items: center;
    }

    .id-item b {
        width: 100px; /* Fixed width for the label */
        text-align: left; /* Align text to the right */
        padding-left: 10px;
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
