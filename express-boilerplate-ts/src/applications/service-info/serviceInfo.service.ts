import { ServiceInfoResponse } from '@applications/service-info/type';
import { types } from '@configs/constants';
import environment from '@configs/environment';
import { logGroup } from '@configs/customLogger';
import { inject, injectable } from 'inversify';
import { Logger } from 'winston';

interface ServiceInfoService {
  getServiceInfo(): Promise<ServiceInfoResponse>;
}

@logGroup()
@injectable()
class ServiceInfoServiceImpl implements ServiceInfoService {
  constructor(@inject(types.Logger) private logger: Logger) { }

  public async getServiceInfo(): Promise<ServiceInfoResponse> {
    this.logger.info(`getting service info`);

    return {
      serviceName: environment.SERVICE_NAME,
      appVersion: environment.APP_VERSION,
      timestamp: Date.now().toString(),
    };
  }
}

export { ServiceInfoService, ServiceInfoServiceImpl };
