import { Body, Controller, Inject, Post, Res } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Response } from 'express';
import { SignInDto, SignInResponse } from '../dto/SignInDto';
import { SignUpDto, SignUpDtoResponse } from '../dto/SignUpDto';
import { HttpUnauthorizedError } from '../errors/HttpUnauthorizedError';
import { AuthService } from '../services/Auth.service';
import { plainToClass } from 'class-transformer';
import { ApiBody, ApiOperation, ApiResponse } from '@nestjs/swagger';

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
    const token = await this.authService.signIn(sigInDto);

    if (!token) {
      throw new HttpUnauthorizedError();
    }
    response.status(200).send({ accessToken: token });
  }

  @ApiOperation({ summary: 'Регистрация' })
  @ApiResponse({
    status: 200,
    type: SignUpDtoResponse,
  })
  @Post('/public/auth/signup')
  @ApiBody({ type: SignUpDto })
  async signUp(@Body() signUpDto: SignUpDto) {
    return await this.authService.signUp(plainToClass(SignUpDto, signUpDto));
  }
}
