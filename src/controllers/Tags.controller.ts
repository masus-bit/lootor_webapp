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
import { TagsService } from '../services/Tags.service';
import { GetTagDto, GetTagsDto } from '../dto/tags/GetTagsDto';
import { CreateTagsDto } from '../dto/tags/TagsDto';

@Controller('/secured/tags')
export class TagsController {
  constructor(@Inject(TagsService) private tagsService: TagsService) {}

  @ApiOperation({
    summary: 'Поиск тегов',
  })
  @ApiResponse({
    status: 200,
    type: GetTagsDto,
  })
  @Get()
  @UseGuards(AuthGuard)
  @ApiQuery({ name: 'name', type: String })
  async searchTags(@Query() query: { name: string }) {
    try {
      return await this.tagsService.searchTags(query.name);
    } catch (err) {
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }

  @ApiOperation({ summary: 'Создание тега' })
  @ApiResponse({
    status: 200,
    type: GetTagDto,
  })
  @Post()
  @ApiBody({ type: CreateTagsDto })
  @UseGuards(AuthGuard)
  async createCollection(@Body() dto: CreateTagsDto) {
    try {
      return await this.tagsService.saveTag(dto.name);
    } catch (e) {
      throw new HttpBadRequestError('Bad request');
    }
  }
}
