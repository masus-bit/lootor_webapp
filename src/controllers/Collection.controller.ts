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
import { CreateCollectionDto } from '../dto/collections/CreateCollectionDto';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { plainToClass } from 'class-transformer';
import { Request } from '../types/base';

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

  @Get('/all')
  @UseGuards(AuthGuard)
  async getAllByUserLogin(@Query() query: { userLogin }) {
    return await this.collectionsService.getByUserId(query.userLogin);
  }

  @Get()
  @UseGuards(AuthGuard)
  async getOne(
    @Query()
    query: {
      id?: string;
      transliteration?: string;
      userLogin?: string;
    },
  ) {
    return await this.collectionsService.getOne(
      query.id,
      query.transliteration,
      query.userLogin,
    );
  }

  @Get('/like')
  @UseGuards(AuthGuard)
  async like(@Query() query: { id: string }, @Req() request: Request) {
    return await this.collectionsService.like(query.id, request.user.login);
  }
}
