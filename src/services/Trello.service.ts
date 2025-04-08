import { Injectable } from '@nestjs/common';
import { ReturnDataDto } from '../dto/ReturnDataDto';
import { FeedbackAddDto } from '../dto/feedback/FeedbackAddDto';
import axios from 'axios';
import * as process from 'node:process';
import * as FormData from 'form-data';

@Injectable()
export class TrelloService {
  constructor() {}

  async addIssue(dto: FeedbackAddDto, file: any): Promise<ReturnDataDto> {
    const response = await axios.post(
      `${process.env.TRELLO_API_URL}?idList=${process.env.TRELLO_LIST_ID}&key=${process.env.TRELLO_API_KEY}&token=${process.env.TRELLO_TOKEN}`,
      new URLSearchParams({
        name: dto.title,
        desc: dto.description,
      }),
      { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } },
    );
    const cardId = response.data.id;

    if (file) {
      const formData = new FormData();
      formData.append('file', file.buffer, { filename: file.originalname });

      await axios.post(
        `${process.env.TRELLO_API_URL}/${cardId}/attachments?key=${process.env.TRELLO_API_KEY}&token=${process.env.TRELLO_TOKEN}`,
        formData,
        { headers: formData.getHeaders() },
      );
    }
    return new ReturnDataDto({ success: true });
  }
}
