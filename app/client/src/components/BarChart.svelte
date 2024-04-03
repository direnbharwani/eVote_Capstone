<!-- src/components/BarChart.svelte -->

<script>
  export let candidates = [];

  // Reactive statement to recalculate maxNumVotes and sortedCandidates when candidates changes
  $: maxNumVotes = Math.max(...candidates.map((candidate) => candidate.NumVotes));
  $: sortedCandidates = [...candidates].sort((a, b) => b.NumVotes - a.NumVotes);

  // Calculate the percentage width for each candidate's bar
  function calculateWidth(numVotes) {
    return (numVotes / maxNumVotes) * 100 + "%";
  }
</script>

<section class="bar-graph bar-graph-horizontal">
  {#each sortedCandidates as candidate, index}
    <div class={`bar-${index + 1}`}>
      <span class="candidate-name">{candidate.Name}</span>
      <div
        class="bar"
        style={`width: ${calculateWidth(candidate.NumVotes)}`}
        data={`${candidate.NumVotes}`}
      ></div>
    </div>
  {/each}
</section>

<style>
  @font-face {
    font-family: Tauri;
    src: url("/fonts/Tauri-Regular.ttf");
    font-weight: normal;
    font-style: normal;
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

  .bar-graph-horizontal .bar {
    border-radius: 3px;
    height: 55px;
    float: left;
    overflow: hidden;
    position: relative;
    width: 0;
  }

  .bar-graph-horizontal .candidate-name {
    font-family: Tauri;
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
      width: var(width);
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
</style>
