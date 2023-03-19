export interface ErrorResponse {
    code: number;
    status: 'error';
    message: string;
}

export interface SuccessResponse<T> {
    code: number;
    status: 'success';
    data: T;
    message: string;
}

export const sendErrorResponse = (code: number, errorMessage: string): ErrorResponse => ({
    code,
    status: 'error',
    message: errorMessage,
});

export const sendSuccessResponse = <T>(code: number, data: T, message = 'Successful'): SuccessResponse<T> => ({
    code,
    status: 'success',
    data,
    message,
});
