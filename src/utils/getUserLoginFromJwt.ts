import { TokenPayloadDto } from '../dto/TokenPayloadDto';

export const getUserLoginFromJwt = (jwtService: any, token: string) => {
  if (!token) return undefined;
  try {
    // @ts-ignore
    const { login } = jwtService.verify<TokenPayloadDto>(token);
    return login;
  } catch (err) {
    console.log(err.message);
  }
};
