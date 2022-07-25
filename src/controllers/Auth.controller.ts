import { Body, Controller, Inject, Post, Res } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Response } from 'express';
import { SignInDto } from '../dto/SignInDto';
import { SignUpDto } from '../dto/SignUpDto';
import { HttpUnauthorizedError } from '../errors/HttpUnauthorizedError';
import { AuthService } from '../services/Auth.service';
import { plainToClass } from 'class-transformer';

@Controller()
export class AuthController {
  constructor(
    @Inject(AuthService) private authService: AuthService,
    @Inject(ConfigService) private configService: ConfigService,
  ) {}

  @Post('/public/auth/signin')
  async signIn(@Body() sigInDto: SignInDto, @Res() response: Response) {
    console.log(sigInDto);
    const token = await this.authService.signIn(sigInDto);

    if (!token) {
      throw new HttpUnauthorizedError();
    }
    response.status(200).send({ accessToken: token });
  }

  @Post('/public/auth/signup')
  async signUp(@Body() signUpDto: SignUpDto) {
    return await this.authService.signUp(plainToClass(SignUpDto, signUpDto));
  }
}
