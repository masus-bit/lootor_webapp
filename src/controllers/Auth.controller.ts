import { Body, Controller, Inject, Post, Res } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Response } from 'express';
import { SignInDto, SignInResponse } from '../dto/SignInDto';
import { SignUpDto } from '../dto/SignUpDto';
import { HttpUnauthorizedError } from '../errors/HttpUnauthorizedError';
import { AuthService } from '../services/Auth.service';
import { plainToClass } from 'class-transformer';
import { ApiBody, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { RefreshDto } from '../dto/RefreshDto';

@Controller()
export class AuthController {
  constructor(
    @Inject(AuthService) private authService: AuthService,
    @Inject(ConfigService) private configService: ConfigService,
  ) {}

  @ApiOperation({ summary: 'Вход' })
  @ApiResponse({
    status: 200,
    type: SignInResponse,
  })
  @Post('/public/auth/signin')
  @ApiBody({ type: SignInDto })
  async signIn(@Body() sigInDto: SignInDto, @Res() response: Response) {
    const tokens = await this.authService.signIn(sigInDto);

    if (!tokens) {
      throw new HttpUnauthorizedError();
    }
    response.status(200).send(tokens);
  }

  @ApiOperation({ summary: 'Регистрация' })
  @ApiResponse({
    status: 200,
    type: SignInResponse,
  })
  @Post('/public/auth/signup')
  @ApiBody({ type: SignUpDto })
  async signUp(@Body() signUpDto: SignUpDto) {
    return await this.authService.signUp(plainToClass(SignUpDto, signUpDto));
  }

  @ApiOperation({ summary: 'Обновление токенov' })
  @ApiResponse({
    status: 200,
    type: SignInResponse,
  })
  @Post('/public/auth/refresh')
  @ApiBody({ type: RefreshDto })
  async refresh(@Body() refreshDto: RefreshDto) {
    return await this.authService.refreshTokens(refreshDto);
  }
}
