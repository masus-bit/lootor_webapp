import { HttpException, HttpStatus } from '@nestjs/common';

/**
 * Error 500
 */
export class HttpInternalServerError extends HttpException {
  constructor(response: string = 'Internal server error') {
    super(response, HttpStatus.INTERNAL_SERVER_ERROR);
  }
}
