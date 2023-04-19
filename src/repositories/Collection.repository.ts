import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Collection } from '../entities/Collection';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';

@Injectable()
export class CollectionRepository {
  constructor(
    @InjectRepository(Collection)
    private collectionRepository: Repository<Collection>,
  ) {}

  async save(data: DeepPartial<Collection>): Promise<Collection> {
    return await this.collectionRepository.save(data);
  }

  createModel(data: DeepPartial<Collection>): Collection {
    return this.collectionRepository.create(data);
  }

  async getById(id: string): Promise<Collection | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.id = :id', { id })
      .innerJoinAndSelect('collection.user', 'user')
      .getOne();
  }

  async getByUserId(id: string): Promise<Collection[]> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.user = :id', { id })
      .innerJoinAndSelect('collection.user', 'user')
      .getMany();
  }

  async getByIdWithoutUser(id: string): Promise<Collection | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.id = :id', { id })
      .getOne();
  }

  async updateById(
    collection: DeepPartial<Collection>,
    id: string,
  ): Promise<Collection> {
    const updatedCollection = this.collectionRepository.create({
      id,
      ...collection,
    });
    return await this.collectionRepository.save(updatedCollection);
  }

  async delete(id: string) {
    return await this.collectionRepository.delete({ id });
  }
}
