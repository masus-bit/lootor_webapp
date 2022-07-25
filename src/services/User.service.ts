import { Inject, Injectable } from '@nestjs/common';
import { User } from '../entities/User';
import { UserRepository } from '../repositories/User.repository';

@Injectable()
export class UserService {
  constructor(
    @Inject(UserRepository)
    private userRepository: UserRepository,
  ) {}
  async getByUserName(userName: string): Promise<User> {
    return await this.userRepository.getByUserName(userName);
  }
}
