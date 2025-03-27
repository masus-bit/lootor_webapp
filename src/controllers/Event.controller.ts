import { Controller, Get, Inject, Query, Req, UseGuards } from '@nestjs/common';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { Request } from '../types/base';
import { ApiOperation, ApiQuery, ApiResponse } from '@nestjs/swagger';
import { EventService } from '../services/Event.service';
import { GetEventsDto } from '../dto/events/GetEventsDto';

@Controller('/secured/events')
export class EventController {
  constructor(@Inject(EventService) private eventService: EventService) {}

  @ApiOperation({
    summary: 'Получить историю событий',
  })
  @ApiResponse({
    status: 200,
    type: GetEventsDto,
  })
  @Get()
  @UseGuards(AuthGuard)
  async getEvents(@Req() request: Request) {
    try {
      return await this.eventService.getEvents(request.user);
    } catch (err) {
      throw new HttpBadRequestError();
    }
  }

  @ApiOperation({
    summary: 'Получить историю событий по конкретной сущности',
  })
  @ApiResponse({
    status: 200,
    type: GetEventsDto,
  })
  @Get('/filter')
  @UseGuards(AuthGuard)
  @ApiQuery({
    name: 'userLogin',
    type: String,
    required: false,
  })
  @ApiQuery({
    name: 'collectionId',
    type: String,
    required: false,
  })
  @ApiQuery({
    name: 'collectionItemId',
    type: String,
    required: false,
  })
  async getFilteredEvents(
    @Query()
    query: {
      userLogin: string;
      collectionId: string;
      collectionItemId;
    },
  ) {
    try {
      return await this.eventService.getFilteredEvents(
        query?.userLogin,
        query?.collectionId,
        query?.collectionItemId,
      );
    } catch (err) {
      throw new HttpBadRequestError();
    }
  }
}
