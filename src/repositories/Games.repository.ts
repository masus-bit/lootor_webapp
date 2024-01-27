import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { Games } from '../entities/Games';
import { Banners } from '../entities/Banners';

@Injectable()
export class GamesRepository {
  constructor(
    @InjectRepository(Games)
    private gamesRepository: Repository<Games>,
  ) {}

  async save(data: DeepPartial<Games>): Promise<Games> {
    return await this.gamesRepository.save(data);
  }

  createModel(data: DeepPartial<Games>): Games {
    return this.gamesRepository.create(data);
  }

  async getById(id: string): Promise<Games | never> {
    return await this.gamesRepository
      .createQueryBuilder('games')
      .leftJoinAndSelect('games.pubs', 'publishers')
      .leftJoinAndSelect('games.genres', 'genres')
      .leftJoinAndSelect('games.platform', 'platforms')
      .leftJoinAndSelect('games.devs', 'developers')
      .innerJoinAndMapMany(
        'games.banners',
        Banners,
        'banners',
        'banners.games_id = games.id',
      )
      .where('games.id = :id', { id })
      .getOne();
  }

  async searchGames(name: string, limit: string): Promise<Games[] | never> {
    try {
      return await this.gamesRepository
        .createQueryBuilder('games')
        .leftJoinAndSelect('games.pubs', 'publishers')
        .leftJoinAndSelect('games.genres', 'genres')
        .leftJoinAndSelect('games.platform', 'platforms')
        .leftJoinAndSelect('games.devs', 'developers')
        .innerJoinAndMapMany(
          'games.banners',
          Banners,
          'banners',
          'banners.games_id = games.id',
        )
        .where('game_title ILIKE :searchTerm', { searchTerm: `%${name}%` })
        .take(Number(limit))
        .getMany();
    } catch (e) {
      console.log(e);
    }
  }
}
