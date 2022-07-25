import { HttpException, HttpStatus } from '@nestjs/common';

/**
 * Error 401
 */
export class HttpUnauthorizedError extends HttpException {
  constructor(response: string = 'Unauthorized') {
    super(response, HttpStatus.UNAUTHORIZED);
  }
}
