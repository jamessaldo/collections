# xendit-callback

## Getting Started

This project is intended to be used with [Node.js][nodejs] 16.13.1 version.

### Clone repository

To clone the repository, use the following commands:

```sh
git clone git@github.com:jamessaldo/wedigo.git
cd wedigo/xendit-callback
yarn install
```

## Available Scripts

- `start` - using nodemon to watch changes
- `clean` - remove coverage data, Jest cache and transpiled files,
- `prebuild` - lint source files and tests before building,
- `build` - transpile TypeScript based on configuration in tsconfig.build.json,
- `build:watch` - interactive watch mode to automatically transpile source files,
- `lint` - lint source files and tests,
- `test` - run tests

[ts-badge]: https://img.shields.io/badge/TypeScript-4.3-blue.svg
[nodejs-badge]: https://img.shields.io/badge/Node.js->=%2014.16-blue.svg
[nodejs]: https://nodejs.org/dist/latest-v14.x/docs/api/
[typescript]: https://www.typescriptlang.org/
[typescript-4-3]: https://www.typescriptlang.org/docs/handbook/release-notes/typescript-4-3.html
[jest]: https://facebook.github.io/jest/
[eslint]: https://github.com/eslint/eslint
[prettier]: https://prettier.io
[repo-template-action]: https://github.com/Maithanhdanh/express-typescript-biolerplate.git
