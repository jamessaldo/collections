import bcrypt from 'bcryptjs';

class UserEntity {
    id: number;
    username: string;
    email: string;
    active: boolean;
    display_name: string;
    first_name: string;
    last_name: string;
    password: string;
    salt: string;

    constructor() {
        this.id = 0;
        this.username = '';
        this.email = '';
        this.active = false;
        this.display_name = '';
        this.first_name = '';
        this.last_name = '';
        this.password = '';
        this.salt = '';
    }

    async verifyPassword(password: string): Promise<boolean> {
        const result = await bcrypt.compare(password, this.password);
        return result;
    }
}

export { UserEntity };