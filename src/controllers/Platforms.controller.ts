import { Controller, Get, Inject, UseGuards } from '@nestjs/common';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { AuthGuard } from '../guards/Auth.guard';
import { ApiOperation, ApiResponse } from '@nestjs/swagger';
import { PlatformsService } from '../services/Platforms.service';
import { GetPlatformsDto } from '../dto/platform/PlatformDto';

@Controller('/secured/platforms')
export class PlatformsController {
  constructor(
    @Inject(PlatformsService) private platformsService: PlatformsService,
  ) {}

  @ApiOperation({
    summary: 'Получить список платформ',
  })
  @ApiResponse({
    status: 200,
    type: GetPlatformsDto,
  })
  @Get()
  @UseGuards(AuthGuard)
  async getAllPlatforms() {
    try {
      return await this.platformsService.getPlatforms();
    } catch (err) {
      throw new HttpBadRequestError();
    }
  }
}
