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
import { CollectionService } from '../services/Collection.service';
import {
  CreateCollectionDto,
  ReturnCollectionsDto,
  ReturnCreateCollection,
} from '../dto/collections/CreateCollectionDto';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { plainToClass } from 'class-transformer';
import { Request } from '../types/base';
import { ApiBody, ApiOperation, ApiQuery, ApiResponse } from '@nestjs/swagger';
import { GetOneCollectionDto } from '../dto/collections/CollectionDto';

@Controller('/secured/collections')
export class CollectionController {
  constructor(
    @Inject(CollectionService) private collectionsService: CollectionService,
  ) {}

  @ApiOperation({ summary: 'Создание коллекции' })
  @ApiResponse({
    status: 200,
    type: ReturnCreateCollection,
  })
  @Post()
  @ApiBody({ type: CreateCollectionDto })
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

  @ApiOperation({ summary: 'Удаление коллекции' })
  @ApiResponse({
    status: 200,
    type: String,
  })
  @Delete()
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'id', type: String })
  async deleteCollection(@Query() query: { id: string }) {
    try {
      return await this.collectionsService.delete(query.id);
    } catch (e) {
      throw new HttpBadRequestError(e);
    }
  }

  @ApiOperation({ summary: 'Изменение коллекции' })
  @ApiResponse({
    status: 200,
    type: ReturnCreateCollection,
  })
  @Patch()
  @UseGuards(AuthGuard)
  @ApiBody({ type: CreateCollectionDto })
  @ApiQuery({ name: 'id', type: String })
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

  @ApiOperation({ summary: 'Получение коллекций по логину пользователя' })
  @ApiResponse({
    status: 200,
    type: ReturnCollectionsDto,
  })
  @Get('/all')
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'userLogin', type: String })
  async getAllByUserLogin(
    @Query() query: { userLogin },
    @Req() request: Request,
  ) {
    return await this.collectionsService.getByUserId(
      query.userLogin,
      request.user,
    );
  }

  @ApiOperation({
    summary:
      'Получение одной коллекциb по id или по транслиту имени или по shareString',
  })
  @ApiResponse({
    status: 200,
    type: GetOneCollectionDto,
  })
  @Get()
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'id', type: String, required: false })
  @ApiQuery({ name: 'transliteration', type: String, required: false })
  @ApiQuery({ name: 'userLogin', type: String, required: false })
  @ApiQuery({ name: 'shareString', type: String, required: false })
  async getOne(
    @Query()
    query: {
      id?: string;
      transliteration?: string;
      userLogin?: string;
      shareString?: string;
    },
    @Req() request: Request,
  ) {
    return await this.collectionsService.getOne(
      request.user,
      query.id,
      query.transliteration,
      query.userLogin,
      query.shareString,
    );
  }

  @ApiOperation({
    summary: 'Лайк коллекции',
  })
  @ApiResponse({
    status: 200,
    type: String,
  })
  @Get('/like')
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'id', type: String })
  async like(@Query() query: { id: string }, @Req() request: Request) {
    return await this.collectionsService.like(query.id, request.user.login);
  }

  @ApiOperation({
    summary: 'Подписка на коллекцию',
  })
  @ApiResponse({
    status: 200,
    type: String,
  })
  @Get('/subscribe')
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'id', type: String })
  @ApiQuery({ name: 'isSubscribe', type: Boolean })
  async subscribe(
    @Query() query: { id: string; isSubscribe: string },
    @Req() request: Request,
  ) {
    return await this.collectionsService.subscribe(
      query.id,
      request.user,
      query.isSubscribe,
    );
  }

  @ApiOperation({
    summary: 'Получение одной коллекциb по id или по транслиту имени',
  })
  @ApiResponse({
    status: 200,
    type: GetOneCollectionDto,
  })
  @Get('/tag')
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'tag', type: String, required: true })
  async getAllByTag(
    @Query()
    query: {
      tag: string;
    },
    @Req() request: Request,
  ) {
    return await this.collectionsService.getByTag(query.tag, request.user);
  }
}
