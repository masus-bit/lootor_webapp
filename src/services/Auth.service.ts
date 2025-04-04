import { Inject, Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { compareSync, hashSync } from 'bcrypt';
import { plainToClass } from 'class-transformer';
import { SignInDto } from '../dto/SignInDto';
import { SignUpDto } from '../dto/SignUpDto';
import { AuthDto } from '../dto/TokenPayloadDto';
import { User } from '../entities/User';
import { HttpUnauthorizedError } from '../errors/HttpUnauthorizedError';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { UserRepository } from '../repositories/User.repository';
import { Algorithm } from 'jsonwebtoken';
import { RefreshDto } from '../dto/RefreshDto';
import { getUnixDate } from '../utils/date';
import * as CryptoJS from 'crypto-js';
import axios from 'axios';

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
    login,
  }: SignInDto): Promise<
    | false
    | { accessToken: string | false; refreshToken: string | false }
    | never
  > {
    let tokens;
    if (email) {
      try {
        const user = await this.usersRepository.getPasswordsByEmail(email);
        if (AuthService.verifyPassword(user, password)) {
          const { access, refresh } = this.getToken(user);
          tokens = { accessToken: access, refreshToken: refresh };
        } else {
          tokens = false;
        }
      } catch (err) {
        throw new HttpUnauthorizedError();
      }
    }
    if (login) {
      try {
        const user = await this.usersRepository.getPasswordsByLogin(login);
        if (AuthService.verifyPassword(user, password)) {
          const { access, refresh } = this.getToken(user);
          tokens = { accessToken: access, refreshToken: refresh };
        } else {
          tokens = false;
        }
      } catch (err) {
        throw new HttpUnauthorizedError();
      }
    }
    return tokens;
    // return false;
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
      if (await this.usersRepository.getById(signUpDto.login.toLowerCase())) {
        throw new HttpBadRequestError('Логин должен быть уникальным');
      }

      if (
        await this.usersRepository.getByEmail(signUpDto.email.toLowerCase())
      ) {
        throw new HttpBadRequestError('Email должен быть уникальным');
      }

      const userModel = this.usersRepository.createModel(signUpDto);
      userModel.created = getUnixDate();

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
    const payload = plainToClass(AuthDto, user);
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

  async vkOauth(
    code: string,
    codeVerifier: string,
    deviceId: string,
    state: string,
    codeChallenge: string,
  ) {
    const params = new URLSearchParams();
    params.append('grant_type', 'authorization_code');
    params.append('client_id', process.env.VK_CLIENT_ID);
    params.append('code', code);
    params.append('code_verifier', codeVerifier);
    params.append('redirect_uri', `${process.env.FRONTEND_URL}/sign-in`);
    params.append('device_id', deviceId);
    params.append('state', state);

    const response = await axios.post('https://id.vk.com/oauth2/auth', params, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    });
    console.log(params, 'params');
    console.log(codeVerifier, 'codeVerifier');
    console.log(response.data, 'response');
    const userInfo = await this.getUserInfo(response.data.access_token);
    return await this.findOrCreateUser(userInfo);
  }

  async getUserInfo(accessToken: string) {
    const response = await axios.get('https://api.vk.com/method/users.get', {
      params: {
        access_token: accessToken,
        v: '5.131',
        fields: 'email,photo_200',
      },
    });
    return response.data?.response?.[0];
  }

  async findOrCreateUser(userInfo: any) {
    console.log(userInfo, 'userInfo');
    let user = await this.usersRepository.findByVkId(userInfo?.id);
    if (!user) {
      const userModel = this.usersRepository.createModel({
        vk_id: userInfo?.id,
        email: userInfo?.email,
        user_name: userInfo?.first_name,
        avatar_url: userInfo?.photo_200,
        login: userInfo?.id.toString(),
      });
      userModel.created = getUnixDate();

      userModel.password = `${userInfo?.id.toString()}@${userInfo?.first_name}`;

      userModel.password_encrypted = AuthService.getHashPassword(
        userModel.password,
      );
      // await this.usersRepository.save(userModel);
      const userFinal = await this.usersRepository.save(userModel);

      const returnedUser = await this.usersRepository.getPasswordsByLogin(
        userFinal.login,
      );
      return await this.signIn({
        email: null,
        password: returnedUser.password,
        login: returnedUser.login,
      });
    }
    const returnedUser = await this.usersRepository.getPasswordsByLogin(
      userInfo.id,
    );
    return await this.signIn({
      email: null,
      password: returnedUser.password,
      login: returnedUser.login,
    });
  }

  async verifyTelegramData(
    data: any,
  ): Promise<
    | false
    | { accessToken: string | false; refreshToken: string | false }
    | never
  > {
    const { hash, ...userData } = data;
    let dataCheckArr = [];
    for (const [key, value] of Object.entries(userData)) {
      dataCheckArr.push(`${key}=${value}`);
    }
    dataCheckArr.sort();
    const dataCheckString = dataCheckArr.join('\n');
    const secretKey = CryptoJS.SHA256(process.env.TELEGRAM_BOT_TOKEN);
    const computedHash = CryptoJS.HmacSHA256(
      dataCheckString,
      secretKey,
    ).toString(CryptoJS.enc.Hex);
    if (computedHash === hash) {
      let user = await this.usersRepository.findByTgId(userData.id);

      if (!user) {
        const userModel = this.usersRepository.createModel({
          telegram_id: userData.id,
          user_name: `${userData.first_name} ${userData?.last_name}`,
          avatar_url: userData.photo_url,
          login: userData.username,
        });
        userModel.created = getUnixDate();
        userModel.password = `${userData?.id}@${userData?.first_name}`;
        userModel.password_encrypted = AuthService.getHashPassword(
          `${userData?.id}@${userData?.first_name}`,
        );
        // await this.usersRepository.save(userModel);
        const userFinal = await this.usersRepository.save(userModel);

        const returnedUser = await this.usersRepository.getPasswordsByLogin(
          userFinal.login,
        );
        return await this.signIn({
          email: null,
          password: returnedUser.password,
          login: returnedUser.login,
        });
      }
      const returnedUser = await this.usersRepository.getPasswordsByLogin(
        userData.username,
      );
      return await this.signIn({
        email: null,
        password: returnedUser.password,
        login: returnedUser.login,
      });
    }
  }
}
