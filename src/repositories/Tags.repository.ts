import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { plainToClass } from 'class-transformer';
import { Tags } from '../entities/Tags';
import { CreateTagDto } from '../dto/tags/TagsDto';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@Injectable()
export class TagsRepository {
  constructor(
    @InjectRepository(Tags)
    private tagsRepository: Repository<Tags>,
    private readonly elasticsearchService: ElasticsearchService,
  ) {}

  async save(data: DeepPartial<Tags>): Promise<Tags> {
    const tags = await this.tagsRepository.save(data);
    await this.elasticsearchService.createIndexIfNotExists('tags', {
      properties: {
        name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    });
    await this.elasticsearchService.upsertDocument('tags', tags.id.toString(), {
      name: tags.name,
      id: tags.id,
    });
    return tags;
  }

  createModel(data: CreateTagDto): Tags {
    return this.tagsRepository.create(data as DeepPartial<Tags>);
  }

  async addTag(name: string): Promise<Tags> {
    const model = this.createModel(
      plainToClass(CreateTagDto, {
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

  async getTagByName(name: string): Promise<Tags> | never {
    try {
      return await this.tagsRepository
        .createQueryBuilder('tags')
        .where('tags.name = :name', { name })
        .getOne();
    } catch (e) {
      console.log(e);
    }
  }
}
