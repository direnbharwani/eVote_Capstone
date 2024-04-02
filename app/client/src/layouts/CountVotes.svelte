<!-- src/layouts/CountVotes.svelte -->

<script>
    import { onDestroy } from "svelte";
    import { writable } from "svelte/store";
    import { Circle } from "svelte-loading-spinners";
    import axios from "axios";

    import BarChart from "../components/BarChart.svelte";
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

    let loading = false;
    let counted = false;

    let candidates = [];

    const countVotes = async () => {
        // Assert electionName and # of candidates is valid
        if (electionID == null || electionID === "") {
            alert("Election ID is invalid!");
            return;
        }

        loading = true;

        try {
            const response = await axios.post(
                "https://dt1nck5gqd.execute-api.ap-southeast-1.amazonaws.com/dev/count-votes",
                {
                    SignerID: "testVoter0",
                    ElectionID: electionID,
                },
            );

            candidates = response.data.Results;

            counted = true;
        } catch (error) {
            console.error(error);
            alert(error);
        }

        loading = false;
    };

    const reset = () => {
        counted = false;
    };
</script>

{#if loading}
    <Circle size="60" color="#a01227" unit="px" duration="1s" />
    <h1 id="counting" align="center">
        Counting<span class="ellipsis-anim"><span>.</span><span>.</span><span>.</span></span>
    </h1>
{:else if !counted}
    <InputSet bind:inputConfigs />
    <Button label="Count" onClick={countVotes} />
{:else}
    <BarChart bind:candidates />
    <Button label="Back" onClick={reset} />
{/if}

<style>
    @font-face {
        font-family: Josefin Sans;
        src: url("/fonts/JosefinSans-VariableFont_wght.ttf");
        font-weight: normal;
        font-style: normal;
    }

    #counting {
        font-family: Josefin Sans;
    }

    .ellipsis-anim span {
	    opacity: 0;
	    -webkit-animation: ellipsis-dot 1s infinite;
	    animation: ellipsis-dot 1s infinite;
	}
	
	.ellipsis-anim span:nth-child(1) {
	    -webkit-animation-delay: 0.0s;
	    animation-delay: 0.0s;
	}
	.ellipsis-anim span:nth-child(2) {
	    -webkit-animation-delay: 0.1s;
	    animation-delay: 0.1s;
	}
	.ellipsis-anim span:nth-child(3) {
	    -webkit-animation-delay: 0.2s;
	    animation-delay: 0.2s;
	}
	
	@-webkit-keyframes ellipsis-dot {
	      0% { opacity: 0; }
	     50% { opacity: 1; }
	    100% { opacity: 0; }
	}
	
	@keyframes ellipsis-dot {
	      0% { opacity: 0; }
	     50% { opacity: 1; }
	    100% { opacity: 0; }
	}
</style>
