<script>
    export let inputConfigs = [];
</script>

<div id="input-set">
    {#each $inputConfigs as config, index}
        <div class="input-container">
            <label class="input-box-label" for={`input-text-box-${index}`}
                >{config.label}</label
            >
            {#if config.type === "text"}
                <input
                    type="text"
                    bind:value={config.value}
                    class="input-box"
                    id={`input-text-box-${index}`}
                    placeholder={config.placeholder}
                />
            {:else if config.type === "number"}
                <input
                    type="number"
                    bind:value={config.value}
                    class="input-box"
                    id={`input-text-box-${index}`}
                />
            {:else}
                <input
                    type="password"
                    bind:value={config.value}
                    class="input-box"
                    id={`input-text-box-${index}`}
                />
            {/if}
        </div>
    {/each}
</div>

<style>
    @font-face {
        font-family: Josefin Sans;
        src: url("/fonts/JosefinSans-VariableFont_wght.ttf");
        font-weight: normal;
        font-style: normal;
    }

    #input-set {
        width: 50%;
        padding: 20px;
        font-family:
            Josefin Sans,
            sans-serif;
        font-size: 1.5em;
        text-align: center;
        /* box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); */
        /* border: 1px solid; */
    }

    .input-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px; /* Add margin below each input container */
        /* border: 1px solid; */
    }

    .input-box-label {
        flex: 1;
        text-align: left;
    }

    .input-box {
        flex: 2;
        padding: 10px;
        border: 1px solid #7e7b7b70;
        border-radius: 5px;
        box-sizing: border-box;
        position: relative;
        margin-top: 10px;
        margin-bottom: 10px;
    }
</style>
