import { Controller, Get, Inject, Req, UseGuards } from '@nestjs/common';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { Request } from '../types/base';
import { ApiOperation, ApiResponse } from '@nestjs/swagger';
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
      throw new HttpBadRequestError(
        'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
      );
    }
  }
}
