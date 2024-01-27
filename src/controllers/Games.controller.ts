// import { Controller, Get, Inject, Query, UseGuards } from '@nestjs/common';
// import { HttpBadRequestError } from '../errors/HttpBadRequestError';
// import { AuthGuard } from '../guards/Auth.guard';
// import { ApiOperation, ApiQuery, ApiResponse } from '@nestjs/swagger';
// import { GamesService } from '../services/Games.service';
// import { ReturnArrayGamesDto } from '../dto/games/GamesDto';
//
// @Controller('/secured')
// export class GamesController {
//   constructor(@Inject(GamesService) private gamesService: GamesService) {}
//
//   @ApiOperation({
//     summary: 'Поиск игор',
//   })
//   @ApiResponse({
//     status: 200,
//     type: ReturnArrayGamesDto,
//   })
//   @Get('/games')
//   @UseGuards(AuthGuard)
//   @ApiQuery({ name: 'searchString', type: String, required: true })
//   @ApiQuery({ name: 'limit', type: String, required: true })
//   async searchGames(@Query() query: { searchString: string; limit: string }) {
//     try {
//       return await this.gamesService.searchGames(
//         query.searchString,
//         query.limit,
//       );
//     } catch (err) {
//       throw new HttpBadRequestError(
//         'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
//       );
//     }
//   }
//
//   @ApiOperation({
//     summary: 'Получить игру по id',
//   })
//   @ApiResponse({
//     status: 200,
//     type: ReturnArrayGamesDto,
//   })
//   @Get('/game')
//   @UseGuards(AuthGuard)
//   @ApiQuery({ name: 'id', type: String })
//   async getOneById(@Query() query: { id: string }) {
//     try {
//       return await this.gamesService.getGameById(query.id);
//     } catch (err) {
//       throw new HttpBadRequestError(
//         'Что-то пошло не так, попробуйте еще раз, если проблема повторяется, обратитесь в техподдержку',
//       );
//     }
//   }
// }
