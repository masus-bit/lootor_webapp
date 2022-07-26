import {
  Body,
  Controller,
  Get,
  Inject,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';
import { UserService } from '../services/User.service';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { ChangePasswordDto } from '../dto/user/ChangePasswordDto';
import { AuthGuard } from '../guards/Auth.guard';

@Controller()
export class UserController {
  constructor(@Inject(UserService) private userService: UserService) {}

  @Post('/secured/user')
  @UseGuards(AuthGuard)
  async changePassword(
    @Body() changeDto: ChangePasswordDto,
    @Query() query: { login: string },
  ) {
    try {
      return await this.userService.changePassword(changeDto, query.login);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  @Get('/secured/user')
  @UseGuards(AuthGuard)
  async getByLogin(@Query() query: { login: string }) {
    try {
      return this.userService.getByLogin(query.login);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }
}
