import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { plainToClass } from 'class-transformer';
import { EntityModel } from '../entities/EntityModel';
import { CreateEntityDto } from '../dto/entities/EntitiesDto';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@Injectable()
export class EntityRepository {
  constructor(
    @InjectRepository(EntityModel)
    private entityRepository: Repository<EntityModel>,
    private readonly elasticsearchService: ElasticsearchService,
  ) {}

  async save(data: DeepPartial<EntityModel>): Promise<EntityModel> {
    const entity = await this.entityRepository.save(data);
    await this.elasticsearchService.createIndexIfNotExists('entity_model', {
      properties: {
        name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    });
    await this.elasticsearchService.upsertDocument(
      'entity_model',
      entity.id.toString(),
      {
        name: entity.name,
        id: entity.id,
      },
    );
    return entity;
  }

  createModel(data: CreateEntityDto): EntityModel {
    return this.entityRepository.create(data as DeepPartial<EntityModel>);
  }

  async addEntity(name: string): Promise<EntityModel> {
    const model = this.createModel(
      plainToClass(CreateEntityDto, {
        name,
      }),
    );
    return await this.save(model);
  }

  async searchEntities(name: string): Promise<EntityModel[] | never> {
    try {
      return await this.entityRepository
        .createQueryBuilder('entity_model')
        .where('name ILIKE :searchTerm', { searchTerm: `%${name}%` })
        .getMany();
    } catch (e) {
      console.log(e);
    }
  }

  async getEntityByName(name: string): Promise<EntityModel> | never {
    try {
      return await this.entityRepository
        .createQueryBuilder('entity_model')
        .where('entity_model.name = :name', { name })
        .getOne();
    } catch (e) {
      console.log(e);
    }
  }

  async find(): Promise<EntityModel[]> | never {
    return await this.entityRepository.find();
  }
}
