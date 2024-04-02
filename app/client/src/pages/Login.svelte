<!-- src/pages/Login.svelte -->
<script>
  import { onMount, onDestroy } from "svelte";
  import { writable } from "svelte/store";
  import { navigate } from 'svelte-routing';

  import Button from "../components/Button.svelte";
  import InputSet from "../components/InputSet.svelte";

  let inputConfigs = writable([
    { label: "Login ID", type: "text", value: "", placeholder: "S9876543A" },
    { label: "Password", type: "password", value: "" },
  ]);
  let inputValues = writable([]);

  let userID;
  let password; // Used for validation since html does not enforce min and max lengths

  const unsubscribe = inputConfigs.subscribe((values) => {
    userID = values[0].value;
    password = values[1].value;
  });

  // Don't forget to unsubscribe when the component is destroyed
  onDestroy(() => {
    unsubscribe();
  });

  onMount(() => {
    userID = null;
    document.title = "eVote POC Login";
  });

  function handleLogin() {
    // Handle admin
    if (userID === "admin") {
      navigate("/admin");
    } 

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
    <InputSet {inputConfigs} />
    <Button label="Continue" linkTo={"/main"} onClick={handleLogin} />
  </body>
</main>

<style>
  @font-face {
    font-family: Josefin Sans;
    src: url("/fonts/JosefinSans-VariableFont_wght.ttf");
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
    font-family:
      Josefin Sans,
      sans-serif;
    font-size: 3em;
  }
</style>
