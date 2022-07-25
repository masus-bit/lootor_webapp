import { HttpException, HttpStatus } from '@nestjs/common';

/**
 * Error 404
 */
export class HttpNotFoundError extends HttpException {
  constructor(response: string = 'Not found') {
    super(response, HttpStatus.NOT_FOUND);
  }
}
