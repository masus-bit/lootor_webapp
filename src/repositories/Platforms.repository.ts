import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Platforms } from '../entities/Platforms';

@Injectable()
export class PlatformsRepository {
  constructor(
    @InjectRepository(Platforms)
    private platformsRepository: Repository<Platforms>,
  ) {}

  async findAll(): Promise<Platforms[]> | never {
    return await this.platformsRepository
      .createQueryBuilder('platforms')
      .getMany();
  }
}
