import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../models/auth_dto.dart';

abstract class AuthRepository {
  Future<AuthConfigDto> getAuthConfig();

  Future<AuthResponseDto> authenticateWithLinkedIn({
    required String code,
    required String redirectUri,
    String? codeVerifier,
  });

  Future<RefreshTokenResponseDto> refreshToken(String refreshToken);

  Future<void> logout(String refreshToken);
}

class HttpAuthRepository implements AuthRepository {
  const HttpAuthRepository(this._dio);

  final Dio _dio;

  @override
  Future<AuthConfigDto> getAuthConfig() async {
    try {
      final response = await _dio.get<Map<String, dynamic>>('/v1/auth/config');
      return AuthConfigDto.fromJson(response.data!);
    } on DioException catch (e) {
      throw mapDioException(e);
    }
  }

  @override
  Future<AuthResponseDto> authenticateWithLinkedIn({
    required String code,
    required String redirectUri,
    String? codeVerifier,
  }) async {
    try {
      final response = await _dio.post<Map<String, dynamic>>(
        '/v1/auth/linkedin',
        data: {
          'code': code,
          'redirect_uri': redirectUri,
          'code_verifier': ?codeVerifier,
        },
      );
      return AuthResponseDto.fromJson(response.data!);
    } on DioException catch (e) {
      throw mapDioException(e);
    }
  }

  @override
  Future<RefreshTokenResponseDto> refreshToken(String refreshToken) async {
    try {
      final response = await _dio.post<Map<String, dynamic>>(
        '/v1/auth/refresh',
        data: {'refresh_token': refreshToken},
      );
      return RefreshTokenResponseDto.fromJson(response.data!);
    } on DioException catch (e) {
      throw mapDioException(e);
    }
  }

  @override
  Future<void> logout(String refreshToken) async {
    try {
      await _dio.post<void>(
        '/v1/auth/logout',
        data: {'refresh_token': refreshToken},
      );
    } on DioException catch (e) {
      throw mapDioException(e);
    }
  }
}

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  final dio = ref.watch(dioProvider);
  return HttpAuthRepository(dio);
});
