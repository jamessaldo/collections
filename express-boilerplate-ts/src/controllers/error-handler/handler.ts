import { sendErrorResponse } from '@server/utils/sendResponse'
import { Response } from 'express'
import HttpStatus from 'http-status'


class ErrorHandlers {
    public onConflict(res: Response, message: string): Response {
        return res.status(HttpStatus.CONFLICT).json(sendErrorResponse(HttpStatus.CONFLICT, message))
    }

    public onNotFound(res: Response, message: string): Response {
        return res.status(HttpStatus.NOT_FOUND).json(sendErrorResponse(HttpStatus.NOT_FOUND, message))
    }

    public onUnauthorized(res: Response, message: string): Response {
        return res.status(HttpStatus.UNAUTHORIZED).json(sendErrorResponse(HttpStatus.UNAUTHORIZED, message))
    }

    public onForbidden(res: Response, message: string): Response {
        return res.status(HttpStatus.FORBIDDEN).json(sendErrorResponse(HttpStatus.FORBIDDEN, message))
    }

    public onError(res: Response, message: string): Response {
        return res.status(HttpStatus.INTERNAL_SERVER_ERROR).json(sendErrorResponse(HttpStatus.INTERNAL_SERVER_ERROR, message))
    }
}

export default ErrorHandlers


