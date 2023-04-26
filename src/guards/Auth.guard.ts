import {
  CanActivate,
  ExecutionContext,
  Inject,
  Injectable,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { JwtService } from '@nestjs/jwt';
import { TokenPayloadDto } from '../dto/TokenPayloadDto';
import { UserRepository } from '../repositories/User.repository';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(
    private readonly jwtService: JwtService,
    @Inject(ConfigService) private configService: ConfigService,
    @Inject(UserRepository) private userRepository: UserRepository,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    try {
      const bearer =
        context.switchToHttp().getRequest()?.headers.authorization ?? '';
      const token = bearer.split(' ')[1];
      if (!token) return false;

      const { login } = this.jwtService.verify<TokenPayloadDto>(token);
      context.switchToHttp().getRequest().user =
        await this.userRepository.getFullById(login);

      return true;
    } catch (err) {
      return false;
    }
  }
}
