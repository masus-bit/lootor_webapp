import { Inject, Injectable } from '@nestjs/common';
import { PlatformsRepository } from '../repositories/Platforms.repository';
import { GetPlatformsDto } from '../dto/platform/PlatformDto';

@Injectable()
export class PlatformsService {
  constructor(
    @Inject(PlatformsRepository)
    private platformsRepository: PlatformsRepository,
  ) {}

  async getPlatforms(): Promise<GetPlatformsDto> {
    try {
      const result = await this.platformsRepository.findAll();
      return new GetPlatformsDto(result);
    } catch (err) {
      console.log(err);
    }
  }
}
