# Wedigo Authorization

<!-- PROJECT LOGO -->
<div align="center">
<p>
  <a href="https://github.com/jamessaldo/wedigo/tree/main/auth">
    <img src="../assets/logo.jpeg" alt="Logo">
  </a>

  <h3 align="center">Wedigo Authorization</h3>

  <p align="center">
   Wedigo
    <br />
    <a href="https://github.com/jamessaldo/wedigo/tree/main/auth"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="mailto:ghozyghlmlaff@gmail.com">Report Bug</a>
    ·
    <a href="mailto:ghozyghlmlaff@gmail.com">Request Feature</a>
  </p>
</p>
</div>

<!-- TABLE OF CONTENTS -->

## Table of Contents

- [Wedigo Authorization](#wedigo-authorization)
  - [Table of Contents](#table-of-contents)
  - [About The Project](#about-the-project)
  - [Objectives](#objectives)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
  - [Roadmap](#roadmap)
  - [Maintainers](#maintainers)
  - [Acknowledgements](#acknowledgements)
  - [License](#license)

<!-- ABOUT THE PROJECT -->

## About The Project

This repository serves as the backend interface for the Wedigo App's Authorization Service, which includes OAuth2 Google login functionality for secure user authentication. In addition, the service enables efficient management of teams, memberships, and invitations for user collaboration. It provides a reliable way to verify user credentials and permissions, as well as manage access control policies. The repository includes a limited set of essential settings, such as token expiration time, and it is designed to be highly scalable and customizable to fit the specific needs of the Wedigo App.

## Objectives

1. Implement **Access** Control to authorize **User** access to Wedigo
   **Applications** and their **Endpoints** based on the user's **Role** within a
   **Team**.
2. Manage **Team** **Membership** and **Invitations**.

<!-- GETTING STARTED -->

## Getting Started

To get a local copy up and running follow these steps.

### Prerequisites

Install required dependencies from the official documentation or through your
system's package manager:

1. [Golang](https://go.dev/doc/install)
2. [Envoy](https://www.envoyproxy.io/docs/envoy/latest/start/install)

### Installation

1. Clone the repository

   ```sh
   > git clone https://github.com/jamessaldo/wedigo.git
   ```
2. Change directory to authorization

    ```sh
    > cd authorization
    ```
3. Install Golang dependencies

   ```sh
   > make install
   ```

<!-- USAGE EXAMPLES -->

## Usage

While developing, you will want to run a dev server locally. To run the dev server locally:

1. Setup your environment variables appropriately, by first copying
   `app.env.example` to `app.env` and filling the values according to the given
   [configuration file](./docs/config.md).
2. Run the dev server

   ```sh
   > make run
   ```
3. Run the external authorization

   ```sh
   > make authz
   ```
4. Run the proxy server

   ```sh
   > make proxy
   ```

<!-- ROADMAP -->

## Roadmap

See the [open issues](https://github.com/jamessaldo/wedigo/issues) for a list of proposed features (and known issues).

<!-- MAINTAINERS -->

## Maintainers

List of Maintainers

- [Ghozy Ghulamul Afif](mailto:ghozyghlmlaff@gmail.com)

<!-- ACKNOWLEDGEMENTS -->

## Acknowledgements

List of libraries used

- [Gitlab CI](https://docs.gitlab.com/ee/ci/)
- [Markdown](https://www.markdownguide.org/)
- [Dockerfile](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

## License

Copyright (c) 2022-, [wedigo.id](https://wedigo.id).
