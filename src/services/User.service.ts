import { Inject, Injectable } from '@nestjs/common';
import { User } from '../entities/User';
import { UserRepository } from '../repositories/User.repository';
import { GetUserByIdDto, UserByIdDto } from '../dto/user/GetUserByIdDto';
import { ChangePasswordDto } from '../dto/user/ChangePasswordDto';
import { hashSync } from 'bcrypt';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import {RatingDto} from "../dto/user/RatingDto";

@Injectable()
export class UserService {
  constructor(
    @Inject(UserRepository)
    private userRepository: UserRepository,
  ) {}

  async getByUserName(userName: string): Promise<User> {
    return await this.userRepository.getByUserName(userName);
  }

  async getByLogin(login: string): Promise<GetUserByIdDto> {
    const user = await this.userRepository.getById(login);
    return new GetUserByIdDto(new UserByIdDto(user));
  }

  async changeRating (dto: RatingDto, login: string): Promise<string> {
    try {
    const existsUser = await this.userRepository.getById(login);
    const userModel = await this.userRepository.createModel(existsUser);
    if (dto.isLike) {
      const likes = existsUser.likes + 1;
      userModel.likes = likes || 1;
      await this.userRepository.save(userModel);
      return 'Рейтинг успешно изменен';
    }
    const dislikes = existsUser.dislikes + 1;
    userModel.dislikes = dislikes || 1;
      await this.userRepository.save(userModel);
      return 'Рейтинг успешно изменен';
    }
    catch (err) {
      throw new HttpBadRequestError('Что-то пошло не так');
    }
  }

  async changePassword(
    dto: ChangePasswordDto,
    login: string,
  ): Promise<GetUserByIdDto> {
    try {
      const existsUser = await this.userRepository.getById(login);
      const userModel = await this.userRepository.createModel(existsUser);
      userModel.password = dto.password;
      userModel.password_encrypted = UserService.getHashPassword(dto.password);
      const user = await this.userRepository.save(userModel);
      return new GetUserByIdDto(new UserByIdDto(user));
    } catch (err) {
      throw new HttpBadRequestError('Что-то пошло не так');
    }
  }

  private static getHashPassword(password: string): string {
    return hashSync(password, 10);
  }
}
