import { HttpException, HttpStatus } from '@nestjs/common';

/**
 * Error 403
 */
export class HttpForbiddenError extends HttpException {
  constructor(response: string = 'Forbidden') {
    super(response, HttpStatus.FORBIDDEN);
  }
}
