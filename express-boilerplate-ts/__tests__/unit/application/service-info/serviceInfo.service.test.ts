import {
  ServiceInfoService,
  ServiceInfoServiceImpl,
} from '@server/applications/service-info/serviceInfo.service';
import { createChildLogger } from '@server/configs/logger';

describe('serviceInfo Service', () => {
  let service: ServiceInfoService;

  beforeEach(() => {
    service = new ServiceInfoServiceImpl(
      createChildLogger('ServiceInfoServiceImpl'),
    );
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should return service info', async () => {
    jest.spyOn(Date, 'now').mockReturnValue(123456789);

    const expectResult = {
      serviceName: 'boilerplate',
      appVersion: '1.0.1',
      timestamp: '123456789',
    };

    const info = await service.getServiceInfo();
    expect(info).toEqual(expectResult);
  });
});
