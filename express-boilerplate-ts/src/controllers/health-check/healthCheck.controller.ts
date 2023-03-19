import { controller, httpGet, interfaces } from 'inversify-express-utils';
import { sendSuccessResponse, SuccessResponse } from '@utils/sendResponse';

export interface HealthCheckController extends interfaces.Controller {
  checkHealth(): SuccessResponse<string>;
}

@controller('/')
export class HealthCheckControllerImpl implements HealthCheckController {
  @httpGet('health')
  public checkHealth(): SuccessResponse<string> {
    return sendSuccessResponse(200, "I'm alive!");
  }
}
