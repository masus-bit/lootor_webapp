import {
  Body,
  Controller,
  Delete,
  Get,
  Inject,
  Patch,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';
import { CreateCollectionDto } from '../dto/collections/CreateCollectionDto';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { plainToClass } from 'class-transformer';
import { CollectionItemService } from '../services/CollectionItem.service';
import { CollectionItemCreateDto } from '../dto/collectionItem/CollectionItemCreateDto';
import { HttpInternalServerError } from '../errors/HttpInternalServerError';

@Controller('/secured/collection_item')
export class CollectionItemController {
  constructor(
    @Inject(CollectionItemService)
    private collectionItemsService: CollectionItemService,
  ) {}

  @Post()
  @UseGuards(AuthGuard)
  async createCollection(@Body() dto: CreateCollectionDto) {
    try {
      return await this.collectionItemsService.create(
        plainToClass(CollectionItemCreateDto, dto),
      );
    } catch (e) {
      throw new HttpBadRequestError('Bad request');
    }
  }

  @Delete()
  @UseGuards(AuthGuard)
  async deleteCollection(@Query() query: { id: string }) {
    try {
      return await this.collectionItemsService.delete(query.id);
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
      return await this.collectionItemsService.update(
        plainToClass(CollectionItemCreateDto, dto),
        query.id,
      );
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }

  @Get()
  @UseGuards(AuthGuard)
  async getById(@Query() query: { id: string }) {
    try {
      return await this.collectionItemsService.getById(query.id);
    } catch (err) {
      throw new HttpInternalServerError(err);
    }
  }
}
