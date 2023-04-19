import { Injectable } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { GetGamesDto } from '../dto/external/GetGamesDto';

require('dotenv').config();

@Injectable()
export class ExternalService {
  constructor(private readonly httpService: HttpService) {}

  async getGames(dto: GetGamesDto): Promise<any> {
    const params = {
      ...dto,
      key: process.env.API_KEY,
    };
    const result = await this.httpService
      .get(`${process.env.API_URL}/games`, {
        params,
      })
      .toPromise();
    return result.data;
  }
}
