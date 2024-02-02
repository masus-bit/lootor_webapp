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
import { getUserLoginFromJwt } from '../utils/getUserLoginFromJwt';
import { JwtService } from '@nestjs/jwt';

@Controller()
export class UserController {
  constructor(
    private readonly jwtService: JwtService,
    @Inject(UserService) private userService: UserService,
  ) {}

  // api для смены пароля
  @ApiOperation({
    summary: 'Отдельно сменить пароль авторизованного пользователя',
  })
  @ApiResponse({
    status: 200,
    type: GetUserByIdDto,
  })
  @Post('/secured/user')
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
  @Post('/secured/user/rating')
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
  @Get('/public/user')
  async getByLogin(@Query() query: QueryUserDto, @Req() request: Request) {
    try {
      let requestUser;
      if (request.headers?.authorization) {
        requestUser = getUserLoginFromJwt(
          this.jwtService,
          request.headers.authorization.split(' ')[1],
        );
      }
      return this.userService.getByLogin(query.login, requestUser);
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
  @Patch('/secured/user/update')
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
  @Get('/secured/user/subscriptions')
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
