# Wedigo Mailer

<!-- PROJECT LOGO -->
<div align="center">
<p>
  <a href="https://github.com/jamessaldo/wedigo/tree/main/mailer">
    <img src="../assets/logo.jpeg" alt="Logo">
  </a>

  <h3 align="center">Wedigo Mailer</h3>

  <p align="center">
   Wedigo
    <br />
    <a href="https://github.com/jamessaldo/wedigo/tree/main/mailer"><strong>Explore the docs »</strong></a>
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

- [Wedigo Mailer](#wedigo-mailer)
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

This repository serves as the backend interface for the Wedigo Mailer Service, which provides a reliable and efficient way to manage email campaigns and send messages to users. The service includes features such as email templates, scheduling, tracking, and analytics to help businesses and organizations improve their email marketing efforts. The repository includes a limited set of essential settings, such as email service provider credentials and default sender information, and it is designed to be highly scalable and customizable to fit the specific needs of the Wedigo Mailer App.

## Objectives

1. Implement a secure and reliable email delivery system for the Wedigo Mailer App.
2. Provide an intuitive and customizable interface for managing email campaigns and templates.
3. Enable users to schedule emails for future delivery and track their performance through analytics.
4. Ensure compliance with email regulations and best practices for email marketing.

<!-- GETTING STARTED -->

## Getting Started

To get a local copy up and running follow these steps.

### Prerequisites

Install required dependencies from the official documentation or through your
system's package manager:

1. [Golang](https://go.dev/doc/install)
2. [Redis](https://redis.io/docs/getting-started/installation/)

### Installation

1. Clone the repository

   ```sh
   > git clone https://github.com/jamessaldo/wedigo.git
   ```
2. Change directory to mailer

    ```sh
    > cd mailer
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
