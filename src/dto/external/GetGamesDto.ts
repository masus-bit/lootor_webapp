import { Exclude, Expose } from 'class-transformer';

@Exclude()
export class GetGamesDto {
  @Expose()
  readonly genres: string;

  @Expose()
  readonly search: string;

  @Expose()
  readonly developers: string;

  @Expose()
  readonly page: string;

  @Expose({ name: 'pageSize' })
  readonly page_size: string;
}
