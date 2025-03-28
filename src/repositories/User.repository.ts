import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { User } from '../entities/User';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@Injectable()
export class UserRepository {
  constructor(
    @InjectRepository(User) private usersRepository: Repository<User>,
    private readonly elasticsearchService: ElasticsearchService,
  ) {}

  async getPasswordsByEmail(email: string): Promise<User | never> {
    return await this.usersRepository
      .createQueryBuilder('user')
      .addSelect(['user.password'])
      .where('user.email = :email', { email })
      .getOne();
  }

  async getPasswordsByLogin(login: string): Promise<User | never> {
    return await this.usersRepository
      .createQueryBuilder('user')
      .addSelect(['user.password'])
      .where('user.login = :login', { login })
      .getOne();
    // return await this.usersRepository.findOneOrFail({
    //   where: {
    //     login: login,
    //   },
    //   select: [
    //     'password',
    //     'login',
    //     'email',
    //     'user_name',
    //     'likes',
    //     'dislikes',
    //     'subscriptions',
    //     'created',
    //     'avatar_url',
    //     'vk_id',
    //     'telegram_id'
    //   ],
    // });
  }

  async save(data: DeepPartial<User>): Promise<User> {
    const user = await this.usersRepository.save(data);
    await this.elasticsearchService.createIndexIfNotExists('collection_item', {
      properties: {
        login: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
        user_name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
        email: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    });
    await this.elasticsearchService.indexDocument('user', {
      login: user.login,
      user_name: user.user_name,
      email: user.email,
    });
    return user;
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

  async findByVkId(vkId: string): Promise<User | never> {
    return await this.usersRepository.findOne({
      where: { vk_id: vkId },
    });
  }

  async findByTgId(tgId: string): Promise<User | never> {
    return await this.usersRepository.findOne({
      where: { telegram_id: tgId },
    });
  }
}
