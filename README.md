<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
<div align="center">
    <img src="https://img.shields.io/badge/go-1.21.6-blue?style=for-the-badge&labelColor=white&logo=go">
    <a href="https://github.com/hyperledger/fabric">
        <img src="https://img.shields.io/badge/hyperledger_fabric-2.3-blue?style=for-the-badge&labelColor=white&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz48IS0tIFVwbG9hZGVkIHRvOiBTVkcgUmVwbywgd3d3LnN2Z3JlcG8uY29tLCBHZW5lcmF0b3I6IFNWRyBSZXBvIE1peGVyIFRvb2xzIC0tPgo8c3ZnIGZpbGw9IiMwMDAwMDAiIHdpZHRoPSI4MDBweCIgaGVpZ2h0PSI4MDBweCIgdmlld0JveD0iMCAwIDI0IDI0IiByb2xlPSJpbWciIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHBhdGggZD0iTTE2LjU5MyAxLjAwOCA1LjcyNSAyLjA2N2EuNjcxLjY3MSAwIDAgMC0uNjg4LS41ODUuNzIuNzIgMCAwIDAtLjcxLjcxYzAgLjI1Ni4xNC40NTUuMzI2LjU4TC4xMDYgMTIuNjM2bC0uMDQ3LjA0NGMtLjA2LjA2LS4wNTkuMDYtLjA1OS4xMiAwIDAgMCAuMDU4LjA1OS4wNThsNi4zMzggOC45NzRhLjcxNi43MTYgMCAwIDAtLjE3NS40NDljMCAuNDE1LjM1NS43MS43MS43MWEuNjguNjggMCAwIDAgLjcwNS0uNjY2bDEwLjY2LTEuMDQuMTMyLjA0OGEuMjI2LjIyNiAwIDAgMCAuMDUyLS4wNjVjLjAzNy0uMDA0LjA2OC0uMDE0LjA2OC0uMDU0bDQuNTUtOS44NzZjLjA2Mi4wMTguMTI1LjAzOS4xOS4wMzkuMzU2IDAgLjcxMS0uMzU1LjcxMS0uNzFhLjcyLjcyIDAgMCAwLS43MS0uNzEzLjY5My42OTMgMCAwIDAtLjI2My4wNTRMMTYuNzEgMS4wNjdjLS4wNi0uMDYtLjA2LS4wNi0uMTE4LS4wNnptLS41MjIuMjgtNC45NTIgMi43MTEtNS40MzQtMS41OThhLjk2Mi45NjIgMCAwIDAgLjA0NS0uMTV6bS40MDQuMDE2IDEuMzYgNS45MDctMy4xMTUtMS45MzNhLjk0NC45NDQgMCAwIDAgLjE1NC0uMzYuNzIuNzIgMCAwIDAtLjcxLS43MS43MS43MSAwIDAgMC0uNjY3LjQ5bC0yLjA3LS42MDh6bS4yNzIuMjQ0IDYuMDYyIDguNjA3Yy0uMDMuMDI4LS4wNTIuMDYtLjA3Ni4wOTNsLTQuNjc2LTIuOXptLTExLjA5NC45NDQgNS4zMzkgMS41NzctMy42NDQgMS45OTUtMS44NzYtMy4zMzlhLjgwMi44MDIgMCAwIDAgLjE4LS4yMzN6bS0uMzguMzYzTDcuMTE4IDYuMTRsLjAxMi4wNTktMi4wOTMgMS45MTdWMi45MDNhLjYyNS42MjUgMCAwIDAgLjIzNi0uMDQ4em0tLjQzLjAxYy4wNDQuMDE0LjA5LjAyMy4xMzUuMDI3djUuMjhMLjU2MSAxMi4yMnptNi40NTYgMS4yOTUgMi4xNzYuNjQzYy0uMDA4LjAzOS0uMDI0LjA3Ni0uMDI0LjExNSAwIC4xNjYuMDY3LjMwNC4xNTguNDJsLTIuNTMgMi4xODgtMy41NjQtMS4yNzh6bTIuMzM5IDEuMjJhLjcyLjcyIDAgMCAwIC41MjUuMjVjLjAyMiAwIC4wNDMtLjAwNS4wNjUtLjAwOGwuNjI1IDMuMjU3LTMuNzIzLTEuMzM0em0uOTk5LjAxOCAzLjI0NSAyLjAxNy42NCAyLjc3OS0zLjU5OC0xLjI5LS42Mi0zLjI5NWEuNzEuNzEgMCAwIDAgLjMzMy0uMjExek03LjQwNyA2LjRsMy41MjcgMS4yNTEtMi44MDMgMi40MjR6bS0uMjIuMDg4Ljc0OCAzLjc1Ni0yLjM4IDIuMDU2YS43MDUuNzA1IDAgMCAwLS41MTgtLjIxMVY4LjQ1N3ptMTAuOTE2IDEuMDY1IDQuNTQ1IDIuODI2Yy0uMDE1LjAzMy0uMDIyLjA2Ni0uMDMyLjFsLTMuOTE3LS4yODV6bS03LjExOC4xMTcgMy45MDIgMS4zODQuNzA3IDMuNjktMy4wNzcgMi42NC0zLjU4NC0xLjIzOC0uNzkzLTQuMDIzem0tNi4wMDcuODR2My41ODVhLjcyLjcyIDAgMCAwLS42NTIuNzA1LjcyLjcyIDAgMCAwIC43MTEuNzExYy4wMyAwIC4wNTYtLjAxMy4wODUtLjAxN2wuNDA3IDIuMTQ3LTUuMTc0LTIuOXptOS45OC41NjggMy40NzEgMS4yMzMtMi43OSAyLjM5NHptLTcuMDE0IDEuMjE0Ljc1MiAzLjc3Mi0zLjAxNS0xLjA0MWEuOTU1Ljk1NSAwIDAgMCAuMDY3LS4yMjMuNzA4LjcwOCAwIDAgMC0uMTctLjQ2NnptMTAuNzguMTQgMy44NTguMjgzYS42Ny42NyAwIDAgMCAuMTE3LjMyMWwtNC4wNjggMy42NDR6bS0uMTgxLjEzOC0uMDk1IDQuMjczLTEuNTY2IDEuNDAzYS43Mi43MiAwIDAgMC0uNDY4LS4xODZjLS4wNSAwLS4wOTMuMDE3LS4xNC4wMjVsLS41NzYtMy4wNjh6bTQuMjczLjYxOGMuMDI2LjAyNC4wNTkuMDM3LjA4OC4wNTdsLTQuNDA3IDkuNTc2LjEzLTUuOTR6bS03LjE2MiAxLjg2Ni41ODIgMy4wMzhhLjY1NC42NTQgMCAwIDAtLjQ4My40MTJsLTIuOTMyLTEuMDE0em0tMTAuMDA4LjA1NSAzLjA3IDEuMDYyLjczIDMuNjY0LTMuODM4LTIuMTUtLjQyNy0yLjIwNGEuNzEzLjcxMyAwIDAgMCAuNDY1LS4zNzJ6bS01LjEyMy4wMTIgNS4wNyAyLjg1IDEuMDc1IDUuNjUzYy0uMDMuMDEzLS4wNTQuMDM0LS4wODEuMDV6bTguNDMgMS4xMzIgMy40ODIgMS4yMDItMi43NzUgMi4zODJ6bTkuNDkxLjc4OC0uMTMzIDUuOTk3LTUuNjY5LTIuMDIgMy4xOTgtMS44NTVhLjcyLjcyIDAgMCAwIC41NzQuMzIuNzIyLjcyMiAwIDAgMCAuNzEyLS43MTMuNjc2LjY3NiAwIDAgMC0uMTYyLS40MjZ6bS01LjcwNC41MTkgMi45OTYgMS4wMzVjLS4wMTMuMDU4LS4wMzIuMTEyLS4wMzIuMTc1IDAgLjExLjAyOC4yMS4wNzEuMzAxTDEyLjUgMTguOTY3bC0yLjYzMy0uOTM3em0tNy4wNjkuNDU2IDMuNzg2IDIuMTMtMi4yMTMgMy40OTNhLjY4Ny42ODcgMCAwIDAtLjQ5Mi0uMDQzem0zLjk4NCAyLjIxNyAyLjU1My45MDMtNC43MDYgMi43MjNjLS4wMjUtLjAzNS0uMDQ3LS4wNzMtLjA4LS4xMDN6bTIuNjk3Ljk1NCA1LjUzMyAxLjk1NS0xMC4yNzQuOTU3Yy0uMDExLS4wNDctLjAyMi0uMDk0LS4wNDItLjEzN3oiLz48L3N2Zz4=">
    </a>
     <img src="https://img.shields.io/npm/v/npm.svg?style=for-the-badge&labelColor=white&logo=npm">
    <img src="https://img.shields.io/badge/svelte-3.55-blue?style=for-the-badge&labelColor=white&logo=svelte">
</div>

<br/>
<div align="center">
    <h1 align="center">Blockchain-based e-Voting System</h3>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#user-flow">User Flow</a></li>
        <li><a href="#architecture">Architecture</a></li>
      </ul>
    </li>
    <!-- <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li> -->
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

This proof-of-concept was developed as my final-year project during my internship to provide a solution to moving from paper ballots to a full electronic voting system that was:
* Secure & Tamper-proof
* Transparent
* Fast
* Cheap


The blockchain network it uses is Hyperledger Fabric due to it being a private and permissioned blockchain with fast and adequately fault tolerant consensus protocols, as well as support for common languages and ease of scalability.

### User Flow
![User Flow Diagram][flow-diagram]

### Architecture
![Architecture Diagram][architecture-diagram]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap

- [x] Set up AWS resources using Terraform
- [x] Build basic frontend with Svelte
- [x] Tear down AWS resources and replace with localstack
  - [ ] Expand to full proposed architecture
- [ ] Set up local fabric network
  - [ ] Shift to CouchDB instead of LevelDB
  - [ ] Add second organisation
  - [ ] Add local wallet storage
- [ ] Optimisations 
  - [ ] Vote Counting with goroutines
  - [ ] Batch multiple submitted votes as one transaction
- [ ] QOL
  - [ ] Changelog
  - [ ] Improve frontend ballot choices
  - [ ] Add history of ballot asset
  
<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

A big thanks to GovTech Singapore for encouraging me to go outside my comfort zone and to come up with this project.

#### Resources Used
* [README Template](https://github.com/othneildrew/Best-README-Template)
* [Img Shields](https://shields.io)

<!-- IMAGES -->
[flow-diagram]: images/eVote_flow.png
[architecture-diagram]: images/eVote_poc_arch.png