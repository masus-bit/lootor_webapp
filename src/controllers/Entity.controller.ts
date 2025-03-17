import {
  Body,
  Controller,
  Get,
  Inject,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { ApiBody, ApiOperation, ApiQuery, ApiResponse } from '@nestjs/swagger';
import { EntityService } from '../services/Entity.service';
import { GetEntitiesDto, GetEntityDto } from '../dto/entities/GetEntitiesDto';
import { CreateEntitiesDto } from '../dto/entities/EntitiesDto';

@Controller('/secured/entities')
export class EntityController {
  constructor(@Inject(EntityService) private entityService: EntityService) {}

  @ApiOperation({
    summary: 'Поиск по сущностям',
  })
  @ApiResponse({
    status: 200,
    type: GetEntitiesDto,
  })
  @Get()
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'name', type: String })
  async searchEntities(@Query() query: { name: string }) {
    try {
      return await this.entityService.searchEntities(query.name);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  @ApiOperation({ summary: 'Создание сущности' })
  @ApiResponse({
    status: 200,
    type: GetEntityDto,
  })
  @Post()
  @ApiBody({ type: CreateEntitiesDto })
  @UseGuards(AuthGuard)
  async createTag(@Body() dto: CreateEntitiesDto) {
    try {
      return await this.entityService.saveEntity(dto);
    } catch (e) {
      throw new HttpBadRequestError('Bad request');
    }
  }
}
