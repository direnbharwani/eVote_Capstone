<!-- src/components/Candidates.svelte -->

<script>
  export let candidates = [];
  export let selectedCandidate;

  export let ballotCast = false;

  // Verify that candidates have the expected structure
  if (
    candidates.some(
      (obj) =>
        !obj.hasOwnProperty("CandidateID") ||
        !obj.hasOwnProperty("Name") ||
        typeof obj.Voted !== "boolean",
    )
  ) {
    console.error("Invalid candidate object structure!");
  }

  // Find the candidate with Voted=true and set it as the selected option
  candidates.forEach((candidate) => {
    if (candidate.Voted) {
      ballotCast = true;
      selectedCandidate = candidate.CandidateID;
    }
  });
</script>

<div id="radio-buttons">
  {#each candidates as candidate, index}
    <div class="radio-container">
      <input
        type="radio"
        id={`option${index}`}
        name="options"
        value={candidate.CandidateID}
        bind:group={selectedCandidate}
        class="styled-radio"
        disabled={ballotCast && selectedCandidate}
      />
      <label for={`option${index}`} class="radio-label">{candidate.Name}</label>
    </div>
  {/each}
</div>

<style>
  #radio-buttons {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 5%;
  }

  .radio-container {
    display: flex;
    align-items: center;
    margin-bottom: 3%; /* Space between radio containers */
  }

  .styled-radio {
    display: none; /* Hide the default radio buttons */
    cursor: pointer;
  }

  .radio-label {
    cursor: pointer;
    margin-right: 10px; /* Space between the label and the radio button */
    position: relative;
    font-family: Tauri, sans-serif;
    font-size: 1.5em;
    text-transform: capitalize;
  }

  .radio-label::before {
    content: "";
    display: inline-block;
    width: 50px; /* Adjust width of the tick box */
    height: 50px; /* Adjust height of the tick box */
    border: 1px solid #000;
    margin-right: 50px; /* Space between the label and the tick box */
    vertical-align: middle; /* Align vertically with the text */
    position: absolute;
    left: -20%;
    top: 50%;
    transform: translateY(-50%);
    box-sizing: border-box; /* Include border in the size */
    padding: 2px;
  }

  input[type="radio"] + .radio-label::after {
    content: "✔";
    font-family: "lucida grande";
    font-size: 1.5em; /* Adjust size of the tick mark */
    line-height: 20px; /* Adjust line-height to vertically center the tick mark */
    text-align: center;
    display: none; /* Initially hidden */
    position: absolute;
    left: -15%; /* Position inside the tick box */
    top: 50%; /* Vertically center inside the tick box */
    transform: translateY(-50%);
  }

  input[type="radio"]:checked + .radio-label::after {
    display: inline-block; /* Show the tick mark when checked */
  }

  /* Adjust the padding for the label to account for the tick box */
  input[type="radio"] + .radio-label {
    padding-left: 30px; /* Space for the radio button and the tick box */
  }
</style>
