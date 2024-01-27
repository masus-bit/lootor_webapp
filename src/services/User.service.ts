import { Inject, Injectable } from '@nestjs/common';
import { User } from '../entities/User';
import { UserRepository } from '../repositories/User.repository';
import { GetUserByIdDto, UserByIdDto } from '../dto/user/GetUserByIdDto';
import { ChangePasswordDto } from '../dto/user/ChangePasswordDto';
import { hashSync } from 'bcrypt';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { RatingDto } from '../dto/user/RatingDto';
import { UpdateUserDto } from '../dto/user/UpdateUserDto';
import { EventRepository } from '../repositories/Event.repository';
import { EventActions, EventTargets } from '../types/base';
import { CollectionItemRepository } from '../repositories/CollectionItem.repository';
import { CollectionRepository } from '../repositories/Collection.repository';

@Injectable()
export class UserService {
  constructor(
    @Inject(UserRepository)
    private userRepository: UserRepository,
    @Inject(EventRepository)
    private eventRepository: EventRepository,
    @Inject(CollectionItemRepository)
    private collectionItemRepository: CollectionItemRepository,
    @Inject(CollectionRepository)
    private collectionRepository: CollectionRepository,
  ) {}

  async getByUserName(userName: string): Promise<User> {
    return await this.userRepository.getByUserName(userName);
  }

  async getByLogin(login: string, authUser?: User): Promise<GetUserByIdDto> {
    const user = await this.userRepository.getById(login);
    const collectionsItemsCount =
      await this.collectionItemRepository.getCountByUserLogin(login);
    const sum = await this.collectionItemRepository.sumByUserLogin(login);
    const collectionsCount =
      await this.collectionRepository.getCountByUserLogin(login);

    if (login === authUser?.login) {
      return new GetUserByIdDto(
        new UserByIdDto(user, collectionsItemsCount, collectionsCount, sum),
      );
    }
    let subArray;
    if (authUser) {
      subArray = authUser?.subscriptions;
    }
    const canSubscribe = authUser ? !subArray.includes(login) : false;
    return new GetUserByIdDto(
      new UserByIdDto(
        user,
        collectionsItemsCount,
        collectionsCount,
        sum,
        canSubscribe,
      ),
    );
  }

  async changeRating(dto: RatingDto, login: string): Promise<string> {
    try {
      const existsUser = await this.userRepository.getById(login);
      const userModel = await this.userRepository.createModel(existsUser);
      if (dto.isLike) {
        const likes = existsUser.likes + 1;
        userModel.likes = likes || 1;
        await this.userRepository.save(userModel);
        return 'Рейтинг успешно изменен';
      }
      const dislikes = existsUser.dislikes + 1;
      userModel.dislikes = dislikes || 1;
      await this.userRepository.save(userModel);
      return 'Рейтинг успешно изменен';
    } catch (err) {
      throw new HttpBadRequestError('Что-то пошло не так');
    }
  }

  async changePassword(
    dto: ChangePasswordDto,
    login: string,
  ): Promise<GetUserByIdDto> {
    try {
      const existsUser = await this.userRepository.getById(login);
      const userModel = await this.userRepository.createModel(existsUser);
      userModel.password = dto.password;
      userModel.password_encrypted = UserService.getHashPassword(dto.password);
      const user = await this.userRepository.save(userModel);
      return new GetUserByIdDto(new UserByIdDto(user));
    } catch (err) {
      throw new HttpBadRequestError('Что-то пошло не так');
    }
  }

  async update(dto: UpdateUserDto, login: string): Promise<GetUserByIdDto> {
    const existsUser = await this.userRepository.getById(login);
    if (existsUser) {
      const user = await this.userRepository.updateByLogin(dto, login);
      const result = await this.userRepository.getById(user.login);
      return new GetUserByIdDto(new UserByIdDto(result));
    }
  }

  async subscribe(
    subscriptionTargetUserLogin: string,
    subscriber: User,
    isSubscribe: string,
  ): Promise<string> {
    try {
      const subscriberModel = await this.userRepository.createModel(subscriber);
      const subscribeTargetUser = await this.userRepository.getById(
        subscriptionTargetUserLogin,
      );
      const subscribeTargetModel = await this.userRepository.createModel(
        subscribeTargetUser,
      );

      if (JSON.parse(isSubscribe)) {
        subscriberModel.subscriptions?.length
          ? // @ts-ignore
            (subscriberModel.subscriptions = `{ ${subscriberModel.subscriptions}, ${subscriptionTargetUserLogin} }`)
          : // @ts-ignore
            (subscriberModel.subscriptions = `{ ${subscriptionTargetUserLogin} }`);
        subscribeTargetModel.subscribers =
          subscribeTargetModel.subscribers += 1;
        await this.userRepository.save(subscribeTargetModel);
        await this.userRepository.save(subscriberModel);
        await this.eventRepository.addEvent(
          // @ts-ignore
          subscriber.login,
          EventActions.subscribe,
          EventTargets.user,
          subscriptionTargetUserLogin,
          subscriptionTargetUserLogin,
        );
        return 'Подписка оформлена';
      } else {
        const subscribers = subscriberModel.subscriptions.filter(
          (s) => s !== subscriptionTargetUserLogin,
        );
        // @ts-ignore
        subscriberModel.subscriptions = `{${subscribers}}`;
        subscribeTargetModel.subscribers =
          subscribeTargetModel.subscribers -= 1;
        await this.userRepository.save(subscribeTargetModel);
        await this.userRepository.save(subscriberModel);
        await this.eventRepository.deleteEvent(
          subscriber.login,
          subscriptionTargetUserLogin,
          EventTargets.user,
          subscriptionTargetUserLogin,
        );
        return 'Отписка оформлена :D';
      }
    } catch (err) {
      console.log(err);
      throw new HttpBadRequestError('Что-то пошло не так');
    }
  }

  private static getHashPassword(password: string): string {
    return hashSync(password, 10);
  }
}
