import { Request as RequestExpress } from 'express';
import { User } from '../entities/User';

export interface Request extends RequestExpress {
  user?: User;
}

export enum EventTargets {
  user = 'targetUser',
  collection = 'targetCollection',
  collectionItem = 'targetCollectionItem',
}

export enum EventActions {
  create = 'create',
  delete = 'delete',
  update = 'update',
  like = 'like',
  subscribe = 'subscribe',
}
