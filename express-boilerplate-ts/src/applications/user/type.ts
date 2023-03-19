import { UserDTO } from "@domains/user/dto";

export interface UserLoginResponse {
    user: UserDTO;
    token: Token;
}

interface Token {
    type: string;
    token: string;
    refreshToken: string;
}