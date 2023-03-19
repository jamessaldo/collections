import { UserAuthService } from '@server/applications/user/auth.service';
import { UserLoginResponse } from '@applications/user/type';
import { types } from '@configs/constants';
import { inject } from 'inversify';
import { controller, httpPost, interfaces, next, requestBody } from 'inversify-express-utils';
import { sendSuccessResponse, SuccessResponse, ErrorResponse } from '@utils/sendResponse';
import { UserLoginRequest } from './type';
import { NextFunction } from 'express';

export interface UserAuthController extends interfaces.Controller {
    login(body: UserLoginRequest, next: NextFunction): Promise<SuccessResponse<UserLoginResponse> | ErrorResponse>;
}

@controller('/login')
export class UserAuthControllerImpl implements UserAuthController {
    constructor(@inject(types.Service.USER_AUTH) private userAuth: UserAuthService) { }

    @httpPost('/')
    public async login(@requestBody() body: UserLoginRequest, @next() next: NextFunction): Promise<SuccessResponse<UserLoginResponse> | ErrorResponse> {
        const email = body.email;
        const password = body.password;

        try {
            const data = await this.userAuth.login(email, password);
            return sendSuccessResponse(200, data, 'Login successfully');
        } catch (e) {
            next(e);
            return null;
        }
    }
}
