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
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { plainToClass } from 'class-transformer';
import { CollectionItemService } from '../services/CollectionItem.service';
import {
  CollectionItemCreateDto,
  ReturnCreateCollectionItem,
} from '../dto/collectionItem/CollectionItemCreateDto';
import { HttpInternalServerError } from '../errors/HttpInternalServerError';
import { ApiBody, ApiOperation, ApiQuery, ApiResponse } from '@nestjs/swagger';

@Controller('/secured/collection_item')
export class CollectionItemController {
  constructor(
    @Inject(CollectionItemService)
    private collectionItemsService: CollectionItemService,
  ) {}

  @ApiOperation({
    summary: 'Создать элемент коллекции',
  })
  @ApiResponse({
    status: 200,
    type: ReturnCreateCollectionItem,
  })
  @Post()
  @UseGuards(AuthGuard)
  @ApiBody({ type: CollectionItemCreateDto })
  async create(@Body() dto: CollectionItemCreateDto) {
    try {
      return await this.collectionItemsService.create(
        plainToClass(CollectionItemCreateDto, dto),
      );
    } catch (e) {
      throw new HttpBadRequestError('Bad request');
    }
  }

  @ApiOperation({
    summary: 'Удалить элемент коллекции',
  })
  @ApiResponse({
    status: 200,
    type: String,
  })
  @Delete()
  @UseGuards(AuthGuard)
  @ApiQuery({ type: String, name: 'id' })
  async delete(@Query() query: { id: string }) {
    try {
      return await this.collectionItemsService.delete(query.id);
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }

  @ApiOperation({
    summary: 'Изменить элемент коллекции',
  })
  @ApiResponse({
    status: 200,
    type: ReturnCreateCollectionItem,
  })
  @Patch()
  @UseGuards(AuthGuard)
  @ApiBody({ type: CollectionItemCreateDto })
  @ApiQuery({ name: 'id', type: String })
  async update(
    @Body() dto: CollectionItemCreateDto,
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

  @ApiOperation({
    summary: 'Получить элемент коллекции',
  })
  @ApiResponse({
    status: 200,
    type: ReturnCreateCollectionItem,
  })
  @Get()
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'id', type: String })
  async getById(@Query() query: { id: string }) {
    try {
      return await this.collectionItemsService.getById(query.id);
    } catch (err) {
      throw new HttpInternalServerError(err);
    }
  }
}
