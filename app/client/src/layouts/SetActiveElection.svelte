<!-- src/layouts/SetActiveElection.svelte -->

<script>
    import { onDestroy } from "svelte";
    import { writable } from "svelte/store";
    import { Circle } from "svelte-loading-spinners";
    import axios from "axios";

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

    let activeElection = JSON.parse(localStorage.getItem("activeElection"));

    let loading = false;
    const setActiveElection = async () => {
        // Assert electionID is valid
        if (electionID == null || electionID === "") {
            alert("Election ID is invalid!");
            return;
        }

        loading = true;

        try {
            const response = await axios.get(
                `https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/get-election/${electionID}`,
            );

            activeElection = {
                ElectionID: electionID,
                Name: response.data.Name,
            };

            localStorage.setItem("activeElection", JSON.stringify(activeElection));
        } catch (error) {
            console.error(error);
            alert(error);
        }

        loading = false;
        alert(`${electionID} has been set as the active election`);
    };
</script>

{#if loading}
    <Circle size="60" color="#a01227" unit="px" duration="1s" />
{:else}
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
{/if}

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
