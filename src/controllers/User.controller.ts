import {
  Body,
  Controller,
  Get,
  Inject,
  Patch,
  Post,
  Query,
  Req,
  UseGuards,
} from '@nestjs/common';
import { UserService } from '../services/User.service';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { ChangePasswordDto } from '../dto/user/ChangePasswordDto';
import { AuthGuard } from '../guards/Auth.guard';
import { GetUserByIdDto, QueryUserDto } from '../dto/user/GetUserByIdDto';
import { RatingDto } from '../dto/user/RatingDto';
import { plainToClass } from 'class-transformer';
import { UpdateUserDto } from '../dto/user/UpdateUserDto';
import { Request } from '../types/base';
import { ApiBody, ApiOperation, ApiQuery, ApiResponse } from '@nestjs/swagger';

@Controller('/secured/user')
export class UserController {
  constructor(@Inject(UserService) private userService: UserService) {}

  // api для смены пароля
  @ApiOperation({
    summary: 'Отдельно сменить пароль авторизованного пользователя',
  })
  @ApiResponse({
    status: 200,
    type: GetUserByIdDto,
  })
  @Post()
  @UseGuards(AuthGuard)
  @ApiBody({ type: ChangePasswordDto })
  async changePassword(
    @Body() changeDto: ChangePasswordDto,
    @Req() request: Request,
  ) {
    try {
      return await this.userService.changePassword(
        changeDto,
        request.user.login,
      );
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  //api для изменения рейтинга
  @ApiOperation({
    summary: 'Изменить рейтинг пользователя',
  })
  @ApiResponse({
    status: 200,
    type: String,
  })
  @Post('/rating')
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'login', type: String })
  @ApiBody({ type: RatingDto })
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

  @ApiOperation({
    summary: 'Получить пользователя по логину',
  })
  @ApiResponse({
    status: 200,
    type: GetUserByIdDto,
  })
  @Get()
  @UseGuards(AuthGuard)
  async getByLogin(@Query() query: QueryUserDto, @Req() request: Request) {
    try {
      return this.userService.getByLogin(query.login, request.user);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  @ApiOperation({
    summary: 'Изменить данные авторизованного пользователя',
  })
  @ApiResponse({
    status: 200,
    type: GetUserByIdDto,
  })
  @Patch('/update')
  @UseGuards(AuthGuard)
  @ApiBody({ type: UpdateUserDto })
  async updateUser(@Body() dto: UpdateUserDto, @Req() request: Request) {
    try {
      return await this.userService.update(
        plainToClass(UpdateUserDto, dto),
        request.user.login,
      );
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }

  @ApiOperation({
    summary: 'Подписка авторизованного пользователя на другого',
  })
  @ApiResponse({
    status: 200,
    type: String,
  })
  @Get('/subscriptions')
  @UseGuards(AuthGuard)
  @ApiQuery({ type: Boolean, name: 'isSubscribe' })
  @ApiQuery({ type: QueryUserDto })
  async subscribe(
    @Req() request: Request,
    @Query() query: { login: string; isSubscribe: string },
  ) {
    try {
      return await this.userService.subscribe(
        query.login,
        request.user,
        query.isSubscribe,
      );
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }
}
