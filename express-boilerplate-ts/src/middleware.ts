import ErrorHandlers from "@controllers/error-handler/handler";
import { ConflictError, RecordNotFoundError, UnauthorizedError, ForbiddenError } from "@domains/error/type";
import { NextFunction, Request, Response } from "express";


// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function errorHandlerMiddleware(err: Error, _req: Request, res: Response, _next: NextFunction): Response {
    const errorHandlers = new ErrorHandlers();

    if (err instanceof RecordNotFoundError) {
        return errorHandlers.onNotFound(res, err.message);
    }

    if (err instanceof ConflictError) {
        return errorHandlers.onConflict(res, err.message);
    }

    if (err instanceof UnauthorizedError) {
        return errorHandlers.onUnauthorized(res, err.message);
    }

    if (err instanceof ForbiddenError) {
        return errorHandlers.onForbidden(res, err.message);
    }

    // If no critearia is matched, return a 500 error
    return errorHandlers.onError(res, err?.message ?? "Internal Server Error");
}
