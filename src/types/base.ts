import { Request as RequestExpress } from 'express';
import { User } from '../entities/User';

export interface Request extends RequestExpress {
  user?: User;
}
