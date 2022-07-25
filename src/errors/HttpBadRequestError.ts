import { HttpException, HttpStatus } from '@nestjs/common';

/**
 * Error 400
 */
export class HttpBadRequestError extends HttpException {
  constructor(response: string = 'Bad request') {
    super(response, HttpStatus.BAD_REQUEST);
  }
}
