import { Inject, Injectable } from '@nestjs/common';
import { TagsRepository } from '../repositories/Tags.repository';
import { GetTagDto, GetTagsDto } from '../dto/tags/GetTagsDto';

@Injectable()
export class TagsService {
  constructor(
    @Inject(TagsRepository)
    private tagsRepository: TagsRepository,
  ) {}

  async saveTag(name: string): Promise<GetTagDto> {
    try {
      const saved = await this.tagsRepository.addTag(name);
      return new GetTagDto(saved);
    } catch (err) {
      console.log(err);
    }
  }

  async searchTags(name: string): Promise<GetTagsDto> {
    try {
      const result = await this.tagsRepository.searchTags(name);
      return new GetTagsDto(result);
    } catch (err) {
      console.log(err);
    }
  }
}
