import mysql from '@server/infrastructures/persistent/mysql';
import { UserEntity } from '@domains/user/model';
import { RecordNotFoundError } from '@server/domains/error/type';
import { types } from '@configs/constants';
import { logGroup } from '@configs/customLogger';
import { inject, injectable } from 'inversify';
import { Logger } from 'winston';

interface UserAuthRepository {
    findByEmail(email: string): Promise<UserEntity>;
}


@logGroup()
@injectable()
class UserAuthRepositoryImpl implements UserAuthRepository {
    constructor(@inject(types.Logger) private logger: Logger) { }

    public async findByEmail(email: string): Promise<UserEntity> {
        const user = new UserEntity();
        const userKeys = Object.keys(user);

        const query = `SELECT ${userKeys.map((key) => `${key}`).join(", ")} FROM users WHERE email = ?`;
        this.logger.debug(query);

        const [results,] = await mysql.query(query, [email]);
        if (!results) {
            throw new RecordNotFoundError(`User with email ${email} is not found`);
        }

        Object.assign(user, results[0]); // populate the user object with the data from the first row
        return user;
    }
}



export { UserAuthRepository, UserAuthRepositoryImpl };