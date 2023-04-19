import { Controller, Get, Inject, Query, UseGuards } from '@nestjs/common';
import { GetGamesDto } from '../dto/external/GetGamesDto';
import { ExternalService } from '../services/External.service';
import { AuthGuard } from '../guards/Auth.guard';

@Controller()
export class ExternalController {
  constructor(@Inject(ExternalService) private extService: ExternalService) {}

  @Get('/secured/external/games')
  @UseGuards(AuthGuard)
  async getGames(@Query() query: GetGamesDto) {
    try {
      return await this.extService.getGames(query);
    } catch (e) {
      console.log(e);
    }
  }
}
