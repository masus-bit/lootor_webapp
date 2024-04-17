import { Inject, Injectable } from '@nestjs/common';
import { TagsRepository } from '../repositories/Tags.repository';
import { GetTagsDto } from '../dto/tags/GetTagsDto';

@Injectable()
export class TagsService {
  constructor(
    @Inject(TagsRepository)
    private tagsRepository: TagsRepository,
  ) {}

  async saveTag(names: string[]): Promise<GetTagsDto> {
    try {
      const tags = [...names];

      let result = [];
      for (const item of tags) {
        const exist = await this.tagsRepository.getTagByName(item);
        if (!exist) {
          const saved = await this.tagsRepository.addTag(item);
          result.push(saved);
        }
      }

      return new GetTagsDto(result);
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
