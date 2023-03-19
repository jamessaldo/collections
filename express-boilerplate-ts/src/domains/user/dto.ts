class UserDTO {
    id: number;
    username: string;
    email: string;
    active: boolean;
    displayName: string;
    firstName: string;
    lastName: string;

    constructor(user: UserDTO) {
        this.id = user.id;
        this.username = user.username;
        this.email = user.email;
        this.active = user.active;
        this.displayName = user.displayName;
        this.firstName = user.firstName;
        this.lastName = user.lastName;
    }
}

export { UserDTO };