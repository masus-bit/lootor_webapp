import { Inject, Injectable } from '@nestjs/common';
import { GamesRepository } from '../repositories/Games.repository';
import { GamesDto, ReturnArrayGamesDto } from '../dto/games/GamesDto';

@Injectable()
export class GamesService {
  constructor(
    @Inject(GamesRepository)
    private gamesRepository: GamesRepository,
  ) {}

  async getGameById(id: string): Promise<{ data: GamesDto }> {
    try {
      const game = await this.gamesRepository.getById(id);
      return { data: new GamesDto(game) };
    } catch (err) {
      throw new Error(err);
    }
  }

  async searchGames(
    searchString: string,
    limit: string,
  ): Promise<ReturnArrayGamesDto> {
    try {
      const games = await this.gamesRepository.searchGames(searchString, limit);
      const preparedDto = games.map((item) => new GamesDto(item));
      return new ReturnArrayGamesDto(preparedDto);
    } catch (err) {
      throw new Error(err);
    }
  }
}
