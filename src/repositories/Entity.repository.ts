import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { plainToClass } from 'class-transformer';
import { EntityModel } from '../entities/EntityModel';
import { CreateEntityDto } from '../dto/entities/EntitiesDto';

@Injectable()
export class EntityRepository {
  constructor(
    @InjectRepository(EntityModel)
    private entityRepository: Repository<EntityModel>,
  ) {}

  async save(data: DeepPartial<EntityModel>): Promise<EntityModel> {
    return await this.entityRepository.save(data);
  }

  createModel(data: CreateEntityDto): EntityModel {
    return this.entityRepository.create(data as DeepPartial<EntityModel>);
  }

  async addEntity(name: string, collectionItem: string): Promise<EntityModel> {
    const model = this.createModel(
      plainToClass(CreateEntityDto, {
        name,
        collectionItem,
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
}
