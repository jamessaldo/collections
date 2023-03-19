const types = {
  Logger: Symbol.for('Logger'),
  Controller: {
    HEALTH_CHECK: Symbol.for('HealthCheckController'),
    SERVICE_INFO: Symbol.for('ServiceInfoController'),
    USER_AUTH: Symbol.for('UserAuthController'),
  },
  Service: {
    SERVICE_INFO: Symbol.for('ServiceInfoService'),
    USER_AUTH: Symbol.for('UserAuthService'),
  },
  Repository: {
    USER_AUTH: Symbol.for('UserAuthRepository'),
  },
};

const reflectMetadataKeys = {
  CLASS_NAME: 'className',
  METHOD_NAME: 'methodName',
};

export { types, reflectMetadataKeys };
