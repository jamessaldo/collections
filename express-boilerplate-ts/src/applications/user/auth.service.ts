import { UserLoginResponse } from '@applications/user/type';
import { types } from '@configs/constants';
import { logGroup } from '@configs/customLogger';
import { inject, injectable } from 'inversify';
import { Logger } from 'winston';
import { UserDTO } from '@domains/user/dto';
import jwt from "jsonwebtoken";
import config from "@configs/environment";
import { UnauthorizedError } from '@server/domains/error/type';
import { UserAuthRepository } from '@server/repositories/user/auth.repository';

interface UserAuthService {
    login(email: string, password: string): Promise<UserLoginResponse>;
}

@logGroup()
@injectable()
class UserAuthServiceImpl implements UserAuthService {
    constructor(@inject(types.Logger) private logger: Logger, @inject(types.Repository.USER_AUTH) private userRepo: UserAuthRepository) { }

    public async login(email: string, password: string): Promise<UserLoginResponse> {
        this.logger.info(`Login with email: ${email}`);

        const user = await this.userRepo.findByEmail(email);
        const userDTO = new UserDTO({ id: user.id, username: user.username, email: user.email, active: user.active, displayName: user.display_name, firstName: user.first_name, lastName: user.last_name });

        if (!await user.verifyPassword(password)) {
            throw new UnauthorizedError(`Invalid password for user: ${email}`);
        }

        const token = createToken({ ...userDTO }, config.SECRET_KEY, config.TOKEN_EXPIRES_IN);
        const refreshToken = createRefreshToken({ id: user.id }, config.SECRET_KEY, config.REFRESH_TOKEN_EXPIRES_IN);

        return {
            user: userDTO,
            token: { type: 'Bearer', token, refreshToken },
        };
    }
}

// create a JWT token
const createToken = (payload: UserDTO, secretKey: string, expiresIn: string) => {
    return jwt.sign(payload, secretKey, { expiresIn });
};

// create a Refresh Token
const createRefreshToken = (payload: { id: number }, secretKey: string, expiresIn: string) => {
    return jwt.sign(payload, secretKey, { expiresIn });
};

export { UserAuthService, UserAuthServiceImpl };