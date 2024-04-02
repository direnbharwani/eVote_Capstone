<!-- src/components/BarChart.svelte -->

<script>
    export let candidates = [];

    $: sortedCandidates = candidates.sort((a, b) => b.NumVotes - a.NumVotes);
</script>

<body>
    <section class="bar-graph bar-graph-horizontal">
        {#each sortedCandidates as candidate, index}
            <div class={`bar-${index + 1}`}>
                <span class="candidate-name">{candidate.Name}</span>
                <div
                    class="bar"
                    style={`width: ${candidate.NumVotes}px`}
                    data={`${candidate.NumVotes}`}
                ></div>
            </div>
        {/each}
    </section>
</body>

<style>
    body {
        display: flex;
        justify-content: space-evenly;
        align-items: center;
        flex-direction: column;
        height: 100vh;
        margin: auto;
        font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
        padding: 20px;
    }

    /* Adjust the bar width based on NumVotes */
    .bar-graph .bar {
        background-color: #a01227;
        -webkit-animation: show-bar 1.2s forwards;
        -moz-animation: show-bar 1.2s forwards;
        animation: show-bar 1.2s forwards;
    }

    .bar-graph .bar::after {
        -webkit-animation: fade-in-text 2.2s 0.1s forwards;
        -moz-animation: fade-in-text 2.2s 0.1s forwards;
        animation: fade-in-text 2.2s 0.1s forwards;
        color: #fff;
        content: attr(data);
        font-weight: 700;
        position: absolute;
        right: 16px;
        top: 17px;
    }

    /* Bar Graph Horizontal */
    .bar-graph-horizontal {
        max-width: 100%;
    }

    .bar-graph-horizontal > div {
        float: left;
        margin-bottom: 8px;
        width: 100%;
    }

    .bar-graph-horizontal .year {
        float: left;
        margin-top: 18px;
        width: 100%;
    }

    .bar-graph-horizontal .bar {
        border-radius: 3px;
        height: 55px;
        float: left;
        overflow: hidden;
        position: relative;
        width: 0;
    }

    .bar-graph-horizontal .candidate-name {
        display: block;
        -webkit-animation: fade-in-text 2.2s 0.1s forwards;
        -moz-animation: fade-in-text 2.2s 0.1s forwards;
        animation: fade-in-text 2.2s 0.1s forwards;
        opacity: 0;
        float: left;
        margin-top: 18px;
        margin-bottom: 5px;
        width: 100%;
    }

    /* Bar Graph Horizontal Animations */
    @keyframes show-bar {
        0% {
            width: 0;
        }
        100% {
            width: 100%;
        }
    }

    @keyframes fade-in-text {
        0% {
            opacity: 0;
        }
        100% {
            opacity: 1;
        }
    }

    @-webkit-keyframes show-bar {
        0% {
            width: 0;
        }
        100% {
            width: 100%;
        }
    }

    @-webkit-keyframes fade-in-text {
        0% {
            opacity: 0;
        }
        100% {
            opacity: 1;
        }
    }
</style>
