import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { plainToClass } from 'class-transformer';
import { Tags } from '../entities/Tags';
import { CreateTagsDto } from '../dto/tags/TagsDto';

@Injectable()
export class TagsRepository {
  constructor(
    @InjectRepository(Tags)
    private tagsRepository: Repository<Tags>,
  ) {}

  async save(data: DeepPartial<Tags>): Promise<Tags> {
    return await this.tagsRepository.save(data);
  }

  createModel(data: CreateTagsDto): Tags {
    return this.tagsRepository.create(data as DeepPartial<Tags>);
  }

  async addTag(name: string): Promise<Tags> {
    const model = this.createModel(
      plainToClass(CreateTagsDto, {
        name,
      }),
    );
    return await this.save(model);
  }

  async searchTags(name: string): Promise<Tags[] | never> {
    try {
      return await this.tagsRepository
        .createQueryBuilder('tags')
        .where('name ILIKE :searchTerm', { searchTerm: `%${name}%` })
        .getMany();
    } catch (e) {
      console.log(e);
    }
  }
}
