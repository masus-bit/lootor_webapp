import {
  Body,
  Controller,
  Delete,
  Inject,
  Patch,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';
import { CollectionService } from '../services/Collection.service';
import { CreateCollectionDto } from '../dto/collections/CreateCollectionDto';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { plainToClass } from 'class-transformer';

@Controller('/secured/collections')
export class CollectionController {
  constructor(
    @Inject(CollectionService) private collectionsService: CollectionService,
  ) {}

  @Post()
  @UseGuards(AuthGuard)
  async createCollection(@Body() dto: CreateCollectionDto) {
    try {
      return await this.collectionsService.createCollection(
        plainToClass(CreateCollectionDto, dto),
      );
    } catch (e) {
      throw new HttpBadRequestError('Bad request');
    }
  }

  @Delete()
  @UseGuards(AuthGuard)
  async deleteCollection(@Query() query: { id: string }) {
    try {
      return await this.collectionsService.delete(query.id);
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }

  @Patch()
  @UseGuards(AuthGuard)
  async updateCollections(
    @Body() dto: CreateCollectionDto,
    @Query() query: { id: string },
  ) {
    try {
      return await this.collectionsService.update(
        plainToClass(CreateCollectionDto, dto),
        query.id,
      );
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }
}
