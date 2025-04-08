import {
  Body,
  Controller,
  Inject,
  Post,
  UploadedFile,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { ApiBody, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { TrelloService } from '../services/Trello.service';
import { ReturnDataDto } from '../dto/ReturnDataDto';
import { FeedbackAddDto } from '../dto/feedback/FeedbackAddDto';
import { FileInterceptor } from '@nestjs/platform-express';

@Controller('/secured/feedback')
export class FeedbackController {
  constructor(@Inject(TrelloService) private trelloService: TrelloService) {}

  @ApiOperation({ summary: 'Создание обращения' })
  @ApiResponse({
    status: 200,
    type: ReturnDataDto,
  })
  @Post()
  @ApiBody({ type: FeedbackAddDto })
  @UseGuards(AuthGuard)
  @UseInterceptors(FileInterceptor('file'))
  async createTag(@UploadedFile() file: any, @Body() dto: FeedbackAddDto) {
    try {
      return await this.trelloService.addIssue(dto, file);
    } catch (e) {
      throw new HttpBadRequestError();
    }
  }
}
