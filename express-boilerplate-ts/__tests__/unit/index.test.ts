import { startServer } from '@configs/express';
process.env.port = '3001';

jest.mock('@configs/container', () => ({
  createContainer: jest.fn().mockImplementation(() => {
    return { get: jest.fn() };
  }),
}));
jest.mock('@configs/express');

describe('Index', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should work', async () => {
    await import('@server/index');
    expect(startServer).toBeCalledTimes(1);
  });
});
