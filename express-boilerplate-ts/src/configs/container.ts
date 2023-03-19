import { ServiceInfoService, ServiceInfoServiceImpl } from '@server/applications/service-info/serviceInfo.service';
import { types } from '@configs/constants';
import { CustomLogger, CustomLoggerImpl } from '@configs/customLogger';
import { logger } from '@configs/logger';
import { HealthCheckController, HealthCheckControllerImpl } from '@controllers/health-check/healthCheck.controller';
import { ServiceInfoController, ServiceInfoControllerImpl } from '@controllers/service-info/serviceInfo.controller';
import { getClassNameFromRequest } from '@utils/container';
import { Container, interfaces } from 'inversify';
import { UserAuthController, UserAuthControllerImpl } from '@server/controllers/user/auth.controller';
import { UserAuthService, UserAuthServiceImpl } from '@server/applications/user/auth.service';
import { UserAuthRepository, UserAuthRepositoryImpl } from '@server/repositories/user/auth.repository';

const createContainer = (): Container => {
  logger.debug(`[${createContainer.name}] Register service on Container`);
  const container = new Container();

  //Logger (reference: https://dev.to/maithanhdanh/enhance-logger-using-inversify-context-and-decorators-2gbe)
  container.bind<CustomLogger>(types.Logger).toDynamicValue((context: interfaces.Context) => {
    const namedMetadata = getClassNameFromRequest(context);
    const logger = new CustomLoggerImpl();
    logger.setContext(namedMetadata);
    return logger;
  });

  //Controller
  container.bind<HealthCheckController>(types.Controller.HEALTH_CHECK).to(HealthCheckControllerImpl);
  container.bind<ServiceInfoController>(types.Controller.SERVICE_INFO).to(ServiceInfoControllerImpl);
  container.bind<UserAuthController>(types.Controller.USER_AUTH).to(UserAuthControllerImpl);

  //Service
  container.bind<ServiceInfoService>(types.Service.SERVICE_INFO).to(ServiceInfoServiceImpl);
  container.bind<UserAuthService>(types.Service.USER_AUTH).to(UserAuthServiceImpl);

  //Repository
  container.bind<UserAuthRepository>(types.Repository.USER_AUTH).to(UserAuthRepositoryImpl);

  return container;
};

export { createContainer };
