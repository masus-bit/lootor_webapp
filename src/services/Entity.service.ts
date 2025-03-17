import { Inject, Injectable } from '@nestjs/common';
import { EntityRepository } from '../repositories/Entity.repository';
import { GetEntitiesDto } from '../dto/entities/GetEntitiesDto';
import { CreateEntitiesDto } from '../dto/entities/EntitiesDto';

@Injectable()
export class EntityService {
  constructor(
    @Inject(EntityRepository)
    private entityRepository: EntityRepository,
  ) {}

  async saveEntity(dto: CreateEntitiesDto): Promise<GetEntitiesDto> {
    try {
      const entities = [...dto.names];

      let result = [];
      for (const item of entities) {
        const exist = await this.entityRepository.getEntityByName(item);
        if (!exist) {
          const saved = await this.entityRepository.addEntity(item);
          result.push(saved);
        }
      }

      return new GetEntitiesDto(result);
    } catch (err) {
      console.log(err);
    }
  }

  async searchEntities(name: string): Promise<GetEntitiesDto> {
    try {
      const result = await this.entityRepository.searchEntities(name);
      return new GetEntitiesDto(result);
    } catch (err) {
      console.log(err);
    }
  }
}
