import {
  Controller,
  Inject,
  OnModuleInit,
  Post,
  Req,
  Res,
} from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Request, Response } from 'express';
import { Observable, lastValueFrom } from 'rxjs';
import { status } from '@grpc/grpc-js';
import { REFRESH_TOKEN_TIME } from 'src/utils/constant';

interface AuthService {
  Register(data: RegisterRequest): Observable<RegisterResponse>;
  Login(data: LoginRequest): Observable<AuthResponse>;
  RefreshToken(data: RefreshTokenRequest): Observable<AuthResponse>;
}

interface RegisterRequest {
  username: string;
  password: string;
}

interface RegisterResponse {
  message: string;
}

interface LoginRequest {
  username: string;
  password: string;
  device_id: string;
  user_agent: string;
}

interface RefreshTokenRequest {
  refreshToken: string;
}

interface AuthResponse {
  accessToken: string;
  refreshToken: string;
}

@Controller()
export class UserRoute implements OnModuleInit {
  private authService: AuthService;

  constructor(@Inject('AUTH_PACKAGE') private client: ClientGrpc) {}

  onModuleInit() {
    this.authService = this.client.getService<AuthService>('AuthService');
  }

  @Post('register')
  async register(@Req() req: Request, @Res() res: Response) {
    try {
      const { username, password } = req.body;

      const result = await lastValueFrom(
        this.authService.Register({ username, password }),
      );

      return res.status(200).json(result);
    } catch (err: any) {
      console.error(err);

      const statusCode =
        err.code === status.ALREADY_EXISTS
          ? 409
          : err.code === status.INVALID_ARGUMENT
            ? 400
            : 500;

      return res.status(statusCode).json({
        message: err.details || 'Registration failed',
      });
    }
  }

  @Post('login')
  async login(@Req() req: Request, @Res() res: Response) {
    try {
      const { username, password, device_id } = req.body;
      const user_agent = req.headers['user-agent'] || 'unknown';

      const { accessToken, refreshToken } = await lastValueFrom(
        this.authService.Login({ username, password, device_id, user_agent }),
      );

      res.cookie('refresh_token', refreshToken, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        maxAge: REFRESH_TOKEN_TIME,
        path: '/api/v1/refresh',
      });

      return res.status(200).json({ access_token: accessToken });
    } catch (err: any) {
      console.error(err);

      const statusCode =
        err.code === status.UNAUTHENTICATED
          ? 401
          : err.code === status.INVALID_ARGUMENT
            ? 400
            : 500;

      return res.status(statusCode).json({
        message: err.details || 'Login failed',
      });
    }
  }

  @Post('refresh')
  async refreshToken(@Req() req: Request, @Res() res: Response) {
    try {
      const refreshToken = req.cookies?.refresh_token;

      if (!refreshToken) {
        return res.status(401).json({ message: 'Refresh token missing' });
      }

      const { accessToken, refreshToken: newRefreshToken } =
        await lastValueFrom(this.authService.RefreshToken({ refreshToken }));

      res.cookie('refresh_token', newRefreshToken, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        maxAge: REFRESH_TOKEN_TIME,
        path: '/api/v1/refresh',
      });

      return res.status(200).json({ access_token: accessToken });
    } catch (err: any) {
      console.error(err);

      const statusCode =
        err.code === status.UNAUTHENTICATED
          ? 401
          : err.code === status.INVALID_ARGUMENT
            ? 400
            : 500;

      return res.status(statusCode).json({
        message: err.details || 'Could not refresh token',
      });
    }
  }
}
