import { HttpException, HttpStatus } from '@nestjs/common';

/**
 * Error 400
 */
export class HttpBadRequestError extends HttpException {
  constructor(
    response: string = 'Что-то пошло не так, закройте еблище и попробуйте снова.',
  ) {
    super(response, HttpStatus.BAD_REQUEST);
  }
}
