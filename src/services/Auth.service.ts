import { Inject, Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { compareSync, hashSync } from 'bcrypt';
import { plainToClass } from 'class-transformer';
import { SignInDto } from '../dto/SignInDto';
import { SignUpDto } from '../dto/SignUpDto';
import { TokenPayloadDto } from '../dto/TokenPayloadDto';
import { User } from '../entities/User';
import { HttpUnauthorizedError } from '../errors/HttpUnauthorizedError';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { UserRepository } from '../repositories/User.repository';
import { Algorithm } from 'jsonwebtoken';
import { RefreshDto } from '../dto/RefreshDto';

@Injectable()
export class AuthService {
  constructor(
    @Inject(UserRepository) private usersRepository: UserRepository,
    private readonly jwtService: JwtService,
  ) {}

  /**
   * Авторизация поьзователя
   */
  async signIn({
    email,
    password,
  }: SignInDto): Promise<
    | false
    | { accessToken: string | false; refreshToken: string | false }
    | never
  > {
    if (email) {
      try {
        const user = await this.usersRepository.getPasswordsByEmail(email);
        if (AuthService.verifyPassword(user, password)) {
          const { access, refresh } = this.getToken(user);
          return { accessToken: access, refreshToken: refresh };
        }

        return false;
      } catch (err) {
        throw new HttpUnauthorizedError();
      }
    }
    return false;
  }

  async refreshTokens(
    refreshDto: RefreshDto,
  ): Promise<
    | false
    | { accessToken: string | false; refreshToken: string | false }
    | never
  > {
    try {
      const canActivate = await this.jwtService.verify(
        refreshDto.refreshToken,
        {
          secret: process.env.JWT_REFRESH_SECRET,
        },
      );
      const exp = canActivate.exp * 1000;
      const now = +new Date();
      if (now > exp) {
        throw new HttpUnauthorizedError('Авторизуйтесь заново.');
      }
      const returnedUser = await this.usersRepository.getPasswordsByEmail(
        canActivate.email,
      );
      return await this.signIn({
        email: returnedUser.email,
        password: returnedUser.password,
      });
    } catch (err) {
      throw new HttpUnauthorizedError();
    }
  }

  /**
   * Регистрация нового пользователя
   */
  async signUp(
    signUpDto: SignUpDto,
  ): Promise<
    | false
    | { accessToken: string | false; refreshToken: string | false }
    | never
  > {
    try {
      if (await this.usersRepository.getById(signUpDto.login)) {
        throw new HttpBadRequestError('userId must be unique');
      }

      if (await this.usersRepository.getByEmail(signUpDto.email)) {
        throw new HttpBadRequestError('email must be unique');
      }

      const userModel = this.usersRepository.createModel(signUpDto);
      userModel.created = new Date(Date.now());

      userModel.password_encrypted = AuthService.getHashPassword(
        userModel.password,
      );

      const user = await this.usersRepository.save(userModel);

      const returnedUser = await this.usersRepository.getPasswordsByEmail(
        user.email,
      );
      return await this.signIn({
        email: returnedUser.email,
        password: returnedUser.password,
      });
    } catch (err) {
      console.error(err.message);
      throw err;
    }
  }

  /**
   * Выход пользователя
   */
  async signOut(): Promise<boolean> {
    return false;
  }

  /**
   * Проверка пароля пользователя
   */
  private static verifyPassword(user: User, password?: string | null): boolean {
    const { password_encrypted } = user;

    return (
      (password_encrypted && compareSync(password, password_encrypted)) ||
      (!password_encrypted && (user.password ?? '') === (password ?? ''))
    );
  }

  /**
   * Генерация хеша пароля
   */
  private static getHashPassword(password: string): string {
    return hashSync(password, 15);
  }

  /**
   * Генерация токена
   */
  public getToken(user: User): { access: string; refresh: string } {
    const payload = plainToClass(TokenPayloadDto, user);
    let refresh = this.jwtService.sign(JSON.parse(JSON.stringify(payload)), {
      secret: process.env.JWT_REFRESH_SECRET,
      expiresIn: process.env.JWT_REFRESH_EXPIRES,
      algorithm: process.env.JWT_ALGORITHM as Algorithm,
    });
    return {
      access: this.jwtService.sign(JSON.parse(JSON.stringify(payload))),
      refresh,
    };
  }
}
