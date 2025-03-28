import {
  Body,
  Controller,
  Get,
  Inject,
  Post,
  Query,
  Res,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Response } from 'express';
import { SignInDto, SignInResponse } from '../dto/SignInDto';
import { SignUpDto } from '../dto/SignUpDto';
import { HttpUnauthorizedError } from '../errors/HttpUnauthorizedError';
import { AuthService } from '../services/Auth.service';
import { plainToClass } from 'class-transformer';
import { ApiBody, ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { RefreshDto } from '../dto/RefreshDto';
import { TelegramSignInDto } from '../dto/auth/TelegramSignInDto';

@ApiTags('Auth')
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
  async signIn(@Body() signInDto: SignInDto, @Res() response: Response) {
    const tokens = await this.authService.signIn(signInDto);

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

  @ApiOperation({ summary: 'Обновление токенов' })
  @ApiResponse({
    status: 200,
    type: SignInResponse,
  })
  @Post('/public/auth/refresh')
  @ApiBody({ type: RefreshDto })
  async refresh(@Body() refreshDto: RefreshDto) {
    return await this.authService.refreshTokens(refreshDto);
  }

  @Get('/public/auth/vk/callback')
  async vkAuthCallback(
    @Query('code') code: string,
    @Query('state') state: string,
    @Query('device_id') deviceId: string,
    @Query('code_verifier') codeVerifier: string,
    @Query('code_challenge') codeChallenge: string,
    @Res() res: Response,
  ) {
    try {
      const tokens = await this.authService.vkOauth(
        code,
        codeVerifier,
        deviceId,
        state,
        codeChallenge,
      );
      if (!tokens) {
        throw new HttpUnauthorizedError();
      }
      res.status(200).send(tokens);
    } catch (error) {
      console.error('Error:', error);
      res.status(500).send('Authorization failed');
    }
  }

  @Post('/public/auth/telegram')
  @ApiBody({ type: TelegramSignInDto })
  async telegramAuth(
    @Body() authData: TelegramSignInDto,
    @Res() res: Response,
  ) {
    try {
      const tokens = await this.authService.verifyTelegramData(authData);
      if (!tokens) {
        throw new HttpUnauthorizedError();
      }
      res.status(200).send(tokens);
    } catch (error) {
      console.error('Error:', error);
      res.status(500).send('Authorization failed');
    }
  }
}
