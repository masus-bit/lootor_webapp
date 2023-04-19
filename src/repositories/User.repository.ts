import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { User } from '../entities/User';
import { Collection } from '../entities/Collection';

@Injectable()
export class UserRepository {
  constructor(
    @InjectRepository(User) private usersRepository: Repository<User>,
  ) {}

  async getPasswordsByEmail(email: string): Promise<User | never> {
    return await this.usersRepository.findOneOrFail({
      where: {
        email: email,
      },
      select: ['password', 'login', 'email', 'user_name'],
    });
  }

  async save(data: DeepPartial<User>): Promise<User> {
    return await this.usersRepository.save(data);
  }

  createModel(data: DeepPartial<User>): User {
    return this.usersRepository.create(data);
  }

  async getById(login: string): Promise<User | never> {
    return await this.usersRepository
      .createQueryBuilder('user')
      .where('user.login = :login', { login })
      .getOne();
  }

  async getByUserName(userName: string): Promise<User | never> {
    return await this.usersRepository.findOne({
      where: { user_name: userName },
    });
  }

  async getByIdOrFail(login: string): Promise<User | never> {
    return await this.usersRepository.findOneOrFail({
      where: { login: login },
    });
  }

  async getByEmail(email: string): Promise<User | never> {
    return await this.usersRepository.findOne({
      where: { email },
    });
  }

  async getFullById(login: string): Promise<User | never> {
    return await this.usersRepository.findOneOrFail({
      where: {
        login: login,
      },
    });
  }

  async updateByLogin(user: DeepPartial<User>, login: string): Promise<User> {
    const updatedUser = this.usersRepository.create({
      login,
      ...user,
    });
    return await this.usersRepository.save(updatedUser);
  }
}
