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

    const setActiveElection = () => {
        // Assert electionName and # of candidates is valid
        if (electionID == null || electionID === "") {
            alert("Election ID is invalid!");
            return;
        }

        alert(`${electionID} has been set as the active election`);
    };
</script>

<InputSet bind:inputConfigs />
<Button label="Set" onClick={setActiveElection} />
