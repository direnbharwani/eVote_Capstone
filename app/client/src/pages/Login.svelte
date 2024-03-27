<!-- src/pages/Login.svelte -->
<script>
    import Button from '../components/Button.svelte'

    let userID;
    let password;   // Used for validation since html does not enforce min and max lengthss

    function handleLogin() {
        if (userID == null || userID.length < 6) {
            alert("Login ID must be more than 6 characters long");
            event.preventDefault();
            return;
        }

        if (password == null || password.length < 8 || password.length > 20) {
            alert("Password must be between 8 and 20 characters long");
            event.preventDefault();
            return;
        }

        sessionStorage.setItem("userID", userID);
        console.log("Logging in user: ", userID);
    }
</script>


<main>
    <body>
        <div id="evote-header">
            <h1>e-Vote</h1>
        </div>

        <div id="login">
            <div class="input-container">
                <label class="login-label" for="login-id-input-text-box" id="login-id-label">
                    Login ID
                </label>
                <input type="username" bind:value={userID} class="login-input-box" id="login-id-input-text-box" minlength="6">
            </div>
    
            <div class="input-container">
                <label class="login-label" for="password-input-text-box" id="password-label">
                    Password
                </label>
                <input type="password" bind:value={password} class="login-input-box" id="password-input-text-box" minlength="8"
                    maxlength="20">
            </div>
    
            <Button label="Continue" linkTo={"/ballot"} onClick={handleLogin} />
    
            <p id="out"></p>
        </div>
    </body>
</main>

<style>
    @font-face {
        font-family: Josefin Sans;
        src: url('/fonts/JosefinSans-VariableFont_wght.ttf');
        font-weight: normal;
        font-style: normal;
    }

    body {
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        height: 100vh;
        margin: auto;
        background-color: #f4f4f4;
    }

    #evote-header {
        width: 100%;
        padding: 20px;
        margin-bottom: 10%;
        text-align: center;
        font-family: Josefin Sans, sans-serif;
        font-size: 3em;
    }

    #login {
        width: 50%;
        padding: 20px;
        font-family: Josefin Sans, sans-serif;
        font-size: 2em;
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

    .login-label {
        flex: 1;
        text-align: left;
    }

    .login-input-box {
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