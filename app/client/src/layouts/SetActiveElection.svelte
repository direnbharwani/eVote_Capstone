<!-- src/layouts/SetActiveElection.svelte -->

<script>
    import { onDestroy } from "svelte";
    import { writable } from "svelte/store";
    import Button from "../components/Button.svelte";
    import InputSet from "../components/InputSet.svelte";

    let inputConfigs = writable([
        { label: "Election ID", type: "text", value: "" },
    ]);

    let electionID;
    const unsubscribe = inputConfigs.subscribe((values) => {
        electionID = values[0].value;
    });

    // Don't forget to unsubscribe when the component is destroyed
    onDestroy(() => {
        unsubscribe();
    });

    let activeElection;

    const setActiveElection = async() => {
        // Assert electionID and electionNAme is valid
        if (electionID == null || electionID === "") {
            alert("Election ID is invalid!");
            return;
        }

        // localStorage.setItem(
        //     "activeElection",
        //     JSON.stringify({
        //         ElectionID: "e-test5",
        //         Name: "TestElection5",
        //     }),
        // );

        alert(`${electionID} has been set as the active election`);
    };
</script>

<h4 align="center">Active Election:</h4>
{#if activeElection != null}
    <div id="active-election">
        {#each Object.keys(activeElection) as key}
            <div class="item">
                <code>{key}: {activeElection[key]}</code>
            </div>
        {/each}
    </div>
{:else}
    <code>null</code>
{/if}

<InputSet bind:inputConfigs />
<Button label="Set" onClick={setActiveElection} />

<style>
    #active-election {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .item {
        text-align: center;
        margin-bottom: 10px;
    }
</style>
