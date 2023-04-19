import {
  Body,
  Controller,
  Get,
  Inject,
  Patch,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';
import { UserService } from '../services/User.service';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { ChangePasswordDto } from '../dto/user/ChangePasswordDto';
import { AuthGuard } from '../guards/Auth.guard';
import { QueryUserDto } from '../dto/user/GetUserByIdDto';
import { RatingDto } from '../dto/user/RatingDto';
import { CreateCollectionDto } from '../dto/collections/CreateCollectionDto';
import { plainToClass } from 'class-transformer';
import { UpdateUserDto } from '../dto/user/UpdateUserDto';

@Controller('/secured/user')
export class UserController {
  constructor(@Inject(UserService) private userService: UserService) {}

  // api для смены пароля
  @Post()
  @UseGuards(AuthGuard)
  async changePassword(
    @Body() changeDto: ChangePasswordDto,
    @Query() query: QueryUserDto,
  ) {
    try {
      return await this.userService.changePassword(changeDto, query.login);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  //api для изменения рейтинга
  @Post('/rating')
  @UseGuards(AuthGuard)
  async changeRating(
    @Body() ratingDto: RatingDto,
    @Query() query: QueryUserDto,
  ) {
    try {
      return await this.userService.changeRating(ratingDto, query.login);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  @Get()
  @UseGuards(AuthGuard)
  async getByLogin(@Query() query: QueryUserDto) {
    try {
      return this.userService.getByLogin(query.login);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  @Patch('/update')
  @UseGuards(AuthGuard)
  async updateUser(
    @Body() dto: UpdateUserDto,
    @Query() query: { login: string },
  ) {
    try {
      return await this.userService.update(
        plainToClass(UpdateUserDto, dto),
        query.login,
      );
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }
}
