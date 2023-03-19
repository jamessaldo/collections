import { ServiceInfoService } from '@server/applications/service-info/serviceInfo.service';
import { ServiceInfoResponse } from '@applications/service-info/type';
import { types } from '@configs/constants';
import { inject } from 'inversify';
import { controller, httpGet, interfaces } from 'inversify-express-utils';
import { sendSuccessResponse, SuccessResponse, ErrorResponse, sendErrorResponse } from '@utils/sendResponse';

export interface ServiceInfoController extends interfaces.Controller {
  getServiceInfo(): Promise<SuccessResponse<ServiceInfoResponse> | ErrorResponse>;
}

@controller('/info')
export class ServiceInfoControllerImpl implements ServiceInfoController {
  constructor(@inject(types.Service.SERVICE_INFO) private serviceInfo: ServiceInfoService) { }

  @httpGet('/')
  public async getServiceInfo(): Promise<SuccessResponse<ServiceInfoResponse> | ErrorResponse> {
    try {
      const serviceInfo = await this.serviceInfo.getServiceInfo();
      return sendSuccessResponse(200, serviceInfo);
    } catch (e) {
      return sendErrorResponse(500, 'Internal Server Error');
    }
  }
}
