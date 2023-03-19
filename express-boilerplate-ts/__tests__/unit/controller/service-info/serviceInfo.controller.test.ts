import supertest from 'supertest';
import { createContainer } from '@configs/container';
import { InversifyExpressServer } from 'inversify-express-utils';

describe('ServiceInfo controller', () => {
  let agent: any;
  beforeAll(() => {
    const container = createContainer();
    const server = new InversifyExpressServer(container);
    const app = server.build();
    agent = supertest.agent(app);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should return service info', (done) => {
    jest.spyOn(Date, 'now').mockReturnValue(123456789);

    const expectResult = {
      code: 200,
      status: 'success',
      data: {
        timestamp: '123456789',
        serviceName: 'boilerplate',
        appVersion: '1.0.1'
      },
      message: 'Successful'
    };

    agent.get('/info').expect(expectResult, done);
  });
});
