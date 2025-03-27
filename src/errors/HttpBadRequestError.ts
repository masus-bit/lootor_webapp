import { HttpException, HttpStatus } from '@nestjs/common';

/**
 * Error 400
 */
export class HttpBadRequestError extends HttpException {
  constructor(
    response: string = 'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, воспользуйтесь формой обратной связи.',
  ) {
    super(response, HttpStatus.BAD_REQUEST);
  }
}
