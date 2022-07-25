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
    login,
    password,
  }: SignInDto): Promise<false | string | never> {
    try {
      const user = await this.usersRepository.getPasswordsByLogin(login);

      if (AuthService.verifyPassword(user, password)) {
        return this.getToken(user);
      }

      return false;
    } catch (err) {
      throw new HttpUnauthorizedError();
    }
  }

  /**
   * Регистрация нового пользователя
   */
  async signUp(signUpDto: SignUpDto): Promise<User | never> {
    console.log(signUpDto);
    try {
      if (await this.usersRepository.getById(signUpDto.login)) {
        throw new HttpBadRequestError('userId must be unique');
      }

      if (await this.usersRepository.getByEmail(signUpDto.email)) {
        throw new HttpBadRequestError('email must be unique');
      }

      const userModel = this.usersRepository.createModel(signUpDto);

      userModel.password_encrypted = AuthService.getHashPassword(
        userModel.password,
      );

      const user = await this.usersRepository.save(userModel);
      return await this.usersRepository.getByIdOrFail(user.login);
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
  public getToken(user: User): string {
    const payload = plainToClass(TokenPayloadDto, user);
    return this.jwtService.sign(JSON.parse(JSON.stringify(payload)));
  }
}
