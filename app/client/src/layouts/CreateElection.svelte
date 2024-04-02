<script>
    import { onDestroy } from "svelte";
    import { writable } from "svelte/store";
    import { Circle } from "svelte-loading-spinners";
    import axios from "axios";

    import Button from "../components/Button.svelte";
    import InputSet from "../components/InputSet.svelte";

    let inputConfigs = writable([
        {
            label: "Election Name",
            type: "text",
            value: "",
            placeholder: "testElection",
        },
        {
            label: "Start Time",
            type: "datetime-local",
            value: "",
            placeholder: "testElection",
        },
        {
            label: "End Time",
            type: "datetime-local",
            value: "",
            placeholder: "testElection",
        },
        { label: "# of Candidates", type: "number", value: "" },
    ]);

    let electionName;
    let startTime;
    let endTime;
    let numCandidates;

    const unsubscribe = inputConfigs.subscribe((values) => {
        electionName = values[0].value;
        startTime = values[1].value;
        endTime = values[2].value;
        numCandidates = values[3].value;
    });

    // Don't forget to unsubscribe when the component is destroyed
    onDestroy(() => {
        unsubscribe();
    });

    let loading = false;

    const createElection = async () => {
        // Assert electionName and # of candidates is valid
        if (electionName == null || electionName === "") {
            alert("Election Name is invalid!");
            return;
        }

        if (numCandidates == null || numCandidates < 1 || numCandidates > 10) {
            alert("# of Candidates is invalid! Must be between 1 and 10");
            return;
        }

        let formattedStartTime;
        let formattedEndTime;

        const parsedStartTime = new Date(startTime);
        if (isNaN(parsedStartTime)) {
            alert("Invalid Start Time");
            return;
        } else {
            const year = parsedStartTime.getFullYear();
            const month = String(parsedStartTime.getMonth() + 1).padStart(
                2,
                "0",
            ); // Months are 0-indexed
            const day = String(parsedStartTime.getDate()).padStart(2, "0");
            const hours = String(parsedStartTime.getHours()).padStart(2, "0");
            const minutes = String(parsedStartTime.getMinutes()).padStart(
                2,
                "0",
            );
            const seconds = String(parsedStartTime.getSeconds()).padStart(
                2,
                "0",
            );

            formattedStartTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
        }

        const parsedEndTime = new Date(endTime);
        if (isNaN(parsedEndTime)) {
            alert("Invalid End Time");
            return;
        } else {
            const year = parsedEndTime.getFullYear();
            const month = String(parsedEndTime.getMonth() + 1).padStart(2, "0"); // Months are 0-indexed
            const day = String(parsedEndTime.getDate()).padStart(2, "0");
            const hours = String(parsedEndTime.getHours()).padStart(2, "0");
            const minutes = String(parsedEndTime.getMinutes()).padStart(2, "0");
            const seconds = String(parsedEndTime.getSeconds()).padStart(2, "0");

            formattedEndTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
        }

        let electionID;

        loading = true;

        try {
            const response = await axios.post(
                "https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/create-election",
                {
                    ElectionName: electionName,
                    StartTime: formattedStartTime,
                    EndTime: formattedEndTime,
                    NumCandidates: numCandidates,
                },
            );

            electionID = response.data.ElectionID;
        } catch (error) {
            console.error(error);
            alert(error);
        }

        loading = false;

        alert(`election created with id: ${electionID}`);
    };
</script>

{#if loading}
    <Circle size="60" color="#a01227" unit="px" duration="1s" />
{:else}
    <InputSet bind:inputConfigs />
    <Button label="Create" onClick={createElection} />
{/if}
