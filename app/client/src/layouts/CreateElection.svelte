<script>
    import { onDestroy } from "svelte";
    import { writable } from "svelte/store";
    import Button from "../components/Button.svelte";
    import InputSet from "../components/InputSet.svelte";

    let inputConfigs = writable([
        {
            label: "Election Name",
            type: "text",
            value: "",
            placeholder: "testElection",
        },
        { label: "# of Candidates", type: "number", value: "" },
    ]);
    let inputValues = writable([]);

    let electionName;
    let numCandidates;

    const unsubscribe = inputConfigs.subscribe((values) => {
        electionName = values[0].value;
        numCandidates = values[1].value;
    });

    // Don't forget to unsubscribe when the component is destroyed
    onDestroy(() => {
        unsubscribe();
    });

    const createElection = () => {
        // Assert electionName and # of candidates is valid
        if (electionName == null || electionName === "") {
            alert("Election Name is invalid!");
            return;
        }

        if (numCandidates == null || numCandidates < 1 || numCandidates > 10) {
            alert("# of Candidates is invalid! Must be between 1 and 10");
            return;
        }

        alert(
            `election ${electionName} with ${numCandidates} candidates created with id: e-e2u8b1e-djnu3i19jr-dmi391sa`,
        );
    };
</script>

<InputSet bind:inputConfigs />
<Button label="Create" onClick={createElection} />
