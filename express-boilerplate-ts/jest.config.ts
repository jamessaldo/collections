import path from 'path';
const rootDirectory = path.resolve(__dirname);

export default {
  clearMocks: true,
  collectCoverage: true,
  coverageDirectory: 'coverage',
  coverageProvider: 'v8',
  coverageThreshold: {
    global: {
      branches: 70,
      function: 80,
      lines: 80,
      statements: 80,
    },
  },
  globals: {
    'ts-jest': {
      tsconfig: path.resolve(__dirname, 'tsconfig.json'),
    },
  },
  moduleDirectories: ['node_modules'],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  moduleNameMapper: {
    '@tests(.*)$': `${rootDirectory}/__tests__$1`,
    '@server(.*)$': `${rootDirectory}/src$1`,
    '@configs(.*)$': `${rootDirectory}/src/configs$1`,
    '@controllers(.*)$': `${rootDirectory}/src/controllers$1`,
    '@applications(.*)$': `${rootDirectory}/src/applications$1`,
    '@domains(.*)$': `${rootDirectory}/src/domains$1`,
    '@utils(.*)$': `${rootDirectory}/src/utils$1`,
    '@infrastructures(.*)$': `${rootDirectory}/src/infrastructures$1`,
    '@repositories(.*)$': `${rootDirectory}/src/repositories$1`,
  },
  reporters: [
    'default',
    [
      path.resolve(__dirname, 'node_modules', 'jest-html-reporter'),
      {
        pageTitle: 'Demo test Report',
        outputPath: 'test-report.html',
      },
    ],
  ],
  rootDir: rootDirectory,
  roots: [rootDirectory],
  setupFilesAfterEnv: [`${rootDirectory}/__tests__/setup.ts`],
  testPathIgnorePatterns: [
    '/node_modules/',
    '<rootDir>/build',
    `${rootDirectory}/__tests__/fixtures`,
    `${rootDirectory}/__tests__/setup.ts`,
  ],
  transform: {
    '^.+\\.ts$': 'ts-jest',
  },
  testRegex: ['((/__tests__/.*)|(\\.|/)(test|spec))\\.tsx?$'],
};
