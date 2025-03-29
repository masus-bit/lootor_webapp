import {
  Body,
  Controller,
  Delete,
  Get,
  Inject,
  Patch,
  Post,
  Query,
  Req,
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
import { Request } from '../types/base';

@Controller()
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
  @Post('/secured/collection_item')
  @UseGuards(AuthGuard)
  @ApiBody({ type: CollectionItemCreateDto })
  async create(@Body() dto: CollectionItemCreateDto, @Req() request: Request) {
    try {
      return await this.collectionItemsService.create(
        plainToClass(CollectionItemCreateDto, dto),
        request.user,
      );
    } catch (e) {
      throw new HttpBadRequestError();
    }
  }

  @ApiOperation({
    summary: 'Удалить элемент коллекции',
  })
  @ApiResponse({
    status: 200,
    type: String,
  })
  @Delete('/secured/collection_item')
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
  @Patch('/secured/collection_item')
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
  @Get('/public/collection_item')
  @ApiQuery({ name: 'id', type: String })
  async getById(@Query() query: { id: string }) {
    try {
      return await this.collectionItemsService.getById(query.id);
    } catch (err) {
      throw new HttpInternalServerError(err);
    }
  }

  @ApiOperation({
    summary: 'Копировать или перенести элемент коллекции',
  })
  @ApiResponse({
    status: 200,
    type: ReturnCreateCollectionItem,
  })
  @Get('/secured/collection_item/copy')
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'id', type: String, required: true })
  @ApiQuery({ name: 'targetCollectionId', type: String, required: true })
  @ApiQuery({ name: 'sourceCollectionId', type: String })
  async copyOrMove(
    @Query()
    query: {
      id: string;
      targetCollectionId: string;
      sourceCollectionId: string;
    },
  ) {
    try {
      return await this.collectionItemsService.copyOrMove(
        query.id,
        query.targetCollectionId,
        query.sourceCollectionId,
      );
    } catch (err) {
      throw new HttpInternalServerError(err);
    }
  }
}
