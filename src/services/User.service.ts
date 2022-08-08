import { Inject, Injectable } from '@nestjs/common';
import { User } from '../entities/User';
import { UserRepository } from '../repositories/User.repository';
import { GetUserByIdDto, UserByIdDto } from '../dto/user/GetUserByIdDto';
import { ChangePasswordDto } from '../dto/user/ChangePasswordDto';
import { hashSync } from 'bcrypt';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';

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
